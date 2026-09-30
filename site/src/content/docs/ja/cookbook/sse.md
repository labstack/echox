---
title: Server-Sent Events (SSE)
description: Echo ハンドラから server-sent events をストリーミングします。接続ごと、または多数のクライアントへのブロードキャストに対応します。
sidebar:
  order: 14
---

[Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)
は複数の方法で使用できます。下の最初の例は、接続ごと、ハンドラごとの SSE です。
より複雑なブロードキャストロジックについては、[r3labs/sse](https://github.com/r3labs/sse)
を使う 2 つ目の例を参照してください。

:::caution
SSE 接続は長時間維持されるため、サーバーの write timeout を無効にする必要があります。
どちらの例も `BeforeServeFunc` で `s.WriteTimeout = 0` を設定します。
:::

## SSE を使う

### サーバー

ハンドラは SSE header を書き込み、その後クライアントが切断するまで毎秒イベントを送信します。
`http.NewResponseController(w).Flush()` は各イベントをすぐにクライアントへ送ります。

```go file=cookbook/sse/simple/server.go
```

### Event 構造と Marshal メソッド

```go file=cookbook/sse/simple/serversentevent.go
```

### SSE を提供する HTML

```html file=cookbook/sse/simple/index.html
```

## r3labs/sse でブロードキャストする

単一のイベントストリームを多数の購読者へブロードキャストする必要がある場合、
[r3labs/sse](https://github.com/r3labs/sse) ライブラリがストリームと購読者の管理を処理します。

### サーバー

```go file=cookbook/sse/broadcast/server.go
```

### SSE を提供する HTML

```html file=cookbook/sse/broadcast/index.html
```
