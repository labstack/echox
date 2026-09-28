---
title: リクエストのスキームと信頼できるプロキシ
description: 信頼できるプロキシの背後で HTTP と HTTPS を判定する方法を設定します。
sidebar:
  order: 15
---

`Context#Scheme()` はリクエストの HTTP/HTTPS を返します。HTTPS リダイレクト、
HSTS、Proxy ミドルウェアがこの値を使用します。Echo v5.4.0 と v4.16.0 以降、
`X-Forwarded-Proto`、`X-Forwarded-Protocol`、`X-Forwarded-Ssl`、
`X-Url-Scheme` が使われるのは、**直接接続した相手**がループバック、リンクローカル、
プライベートアドレス、または Unix ソケットの場合だけです。それ以外では接続自体から
スキームを判定します。公開クライアントが HTTP で `X-Forwarded-Proto: https` を
送って HTTPS リダイレクトを回避することはできません。

## 抽出方法を選ぶ

`Echo#SchemeExtractor` で判定方法を設定します。v5 では
`Config.SchemeExtractor` も使えます。既定の
`echo.ExtractSchemeFromHeaders()` は上記の直接接続元だけを信頼します。
`echo.ExtractSchemeDirect()` は転送ヘッダーを無視します。
`echo.LegacySchemeExtractor()` は旧動作に戻しますが、信頼できないクライアントが
アプリに到達できる場合や、プロキシがクライアントのヘッダーをそのまま渡す場合は危険です。

`X-Forwarded-Proto` がある場合は**最後の値だけ**を使い、小文字のスキームを
返します。不正な値は `http` となり、別の転送スキームヘッダーには切り替わりません。
信頼するプロキシはクライアントから届いた値を渡さず、観測したスキームで
`X-Forwarded-Proto` を**上書き**してください。nginx では
`proxy_set_header X-Forwarded-Proto $scheme;` を設定します。

## 公開アドレスのプロキシ

直接接続するプロキシが公開アドレスまたは `100.64.0.0/10` を使う場合は、
そのアドレス範囲を明示的に信頼してください。設定しないと転送スキームが無視され、
HTTPS リダイレクトがループしたり、Secure ミドルウェアが HSTS を送らなくなったり
します。Cloudflare、CloudFront、Azure Front Door、GCP 外部 HTTP(S)
ロードバランサー、GKE Ingress などが該当します。実際に使う範囲だけを信頼します。
GCP の範囲を指定する例:

```go
_, gclb1, _ := net.ParseCIDR("35.191.0.0/16")
_, gclb2, _ := net.ParseCIDR("130.211.0.0/22")
e.SchemeExtractor = echo.ExtractSchemeFromHeaders(
	echo.TrustIPRange(gclb1),
	echo.TrustIPRange(gclb2),
)
```

この v5 の例では `net` と `github.com/labstack/echo/v5` をインポートします。
v4 では `github.com/labstack/echo/v4` を使います。同一ホスト、プライベート
ネットワーク、Unix ソケットのプロキシは既定の設定で動作します。アダプターが
`RemoteAddr` をクライアントのアドレスに置き換える場合は、実際の構成に合わせて
抽出方法を設定してください。クライアント IP は
[IP アドレス](/ja/guide/ip-address/)で別途設定します。
