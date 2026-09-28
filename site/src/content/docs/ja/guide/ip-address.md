---
title: IP アドレス
description: プロキシの背後で実際のクライアント IP アドレスを安全に取得します。
sidebar:
  order: 14
---

IP アドレスは HTTP において基本的な役割を持ち、アクセス制御、監査、地理ベースの分析などに
使われます。Echo はそれを取得するために `Context#RealIP()` を提供しています。

_実際の_ クライアント IP を取得するのは簡単ではありません。特に L7 プロキシが
アプリケーションの前にある場合はそうです。その場合、実際の IP はプロキシから
HTTP 経由でアプリへ伝えられる必要があります。ただし HTTP header を無条件に信頼してはいけません。
だまされる可能性があります。**これはセキュリティリスクです。**

IP を信頼性高く安全に取得するには、アプリケーションがインフラ全体を把握している必要があります。
Echo では `Echo#IPExtractor` でこれを設定します。

v5 の `Context#RealIP()` は v5.1.0 以降、既定で直接接続元のアドレスを使います。v4 の既定動作は依然として任意のクライアントが送った `X-Forwarded-For` と `X-Real-IP` を信頼するため、プロキシ構成に合わせて **必ず `Echo#IPExtractor` を設定**してください。HTTP/HTTPS の判定は[リクエストのスキーム](/ja/guide/request-scheme/)で別途設定します。

適切な方法を見つけるため、まず 2 つの質問から始めます。

1. アプリケーションの前に HTTP（L7）プロキシを置いていますか？これにはクラウドの
   ロードバランサー（AWS ALB や GCP HTTP LB など）や、オープンソースプロキシ
   （Nginx、Envoy、Istio ingress gateway など）が含まれます。
2. そうであれば、プロキシはどの HTTP header を使ってクライアント IP をアプリケーションへ渡していますか？

## ケース 1：プロキシなし

プロキシがない場合（アプリがインターネットに直接向いている場合）、信頼できる唯一のアドレスは
ネットワーク層から得られるものです。クライアントが完全に制御できるため、すべての HTTP header は
信頼できません。

`echo.ExtractIPDirect()` を使います。

```go
e.IPExtractor = echo.ExtractIPDirect()
```

## ケース 2：X-Forwarded-For header を使うプロキシ

[`X-Forwarded-For` (XFF)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For)
は、クライアント IP を中継するもっとも一般的な header です。各ホップで、プロキシは
リクエスト IP を header の末尾に追加します。

```text
            ┌──────────┐            ┌──────────┐            ┌──────────┐
───────────>│ Proxy 1  │───────────>│ Proxy 2  │───────────>│ Your app │
            │ (IP: b)  │            │ (IP: c)  │            │          │
            └──────────┘            └──────────┘            └──────────┘

Case 1.
XFF:  ""                    "a"                     "a, b"
                                                    ~~~~~~
Case 2.
XFF:  "x"                   "x, a"                  "x, a, b"
                                                    ~~~~~~~~~
                                                    ↑ What your app will see
```

この場合、**右から見て最初の信頼できない IP 読み取り値**を使います。左端の最初の値は、
クライアントが制御するため絶対に使わないでください。ここで「信頼できる」とは、その IP が
自分のインフラに属していると確信できることです。上の例では `b` と `c` が信頼できるなら、
どちらのケースでもクライアント IP は `a` であり、`x` ではありません。

`ExtractIPFromXFFHeader(...TrustOption)` を使います。

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader()
```

デフォルトでは、loopback、link-local unicast、private-use、
[RFC 6890](https://datatracker.ietf.org/doc/html/rfc6890)、
[RFC 4291](https://datatracker.ietf.org/doc/html/rfc4291)、
[RFC 4193](https://datatracker.ietf.org/doc/html/rfc4193) の unique local アドレスなど、
内部 IP アドレスを信頼します。これは `TrustOption` で制御できます。

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader(
	echo.TrustLoopback(false),   // e.g. IPv4 starting with 127.
	echo.TrustLinkLocal(false),  // e.g. IPv4 starting with 169.254.
	echo.TrustPrivateNet(false), // e.g. IPv4 starting with 10. or 192.168.
	echo.TrustIPRange(lbIPRange),
)
```

## ケース 3：X-Real-IP header を使うプロキシ

`X-Real-IP` はクライアント IP を中継するもう 1 つの header ですが、XFF と異なり
単一のアドレスだけを保持します。

プロキシがこの header を設定している場合は、`ExtractIPFromRealIPHeader(...TrustOption)` を使います。

```go
e.IPExtractor = echo.ExtractIPFromRealIPHeader()
```

XFF と同様に、デフォルトでは内部 IP アドレスを信頼し、同じ `TrustOption` を受け付けます。

:::danger
インフラのエッジにある最外部プロキシを、**受信 header をそのまま通さない**ように設定することを
**絶対に忘れないでください**。そうしないと、クライアントがそれらを偽造でき、不正への道を開きます。
:::

## デフォルトの挙動

`Echo#IPExtractor` 未設定時、**v5 は直接接続元のアドレス**を使い、クライアントが送った IP ヘッダーを無視します。**v4 は旧動作のまま**で、信頼できるプロキシか確認せずに `X-Forwarded-For` や `X-Real-IP` を使用する場合があります。v4 ではプロキシがなければ `echo.ExtractIPDirect()`、プロキシがあれば適切な信頼オプションを付けたヘッダー抽出器を設定してください。
