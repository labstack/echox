---
title: ストリーミングレスポンス
description: chunked transfer encoding を使い、生成されたデータをクライアントへストリーミングします。
sidebar:
  order: 15
---

このレシピでは、chunked transfer encoding を使って、各レコードが生成されるたびに JSON レスポンスを
クライアントへストリーミングします。

- 生成されたデータをその都度送信します。
- chunked transfer encoding で JSON レスポンスをストリーミングします。

ハンドラは 1 レコードずつエンコードし、各レコードの後に
`http.NewResponseController(...).Flush()` を呼び出してすぐにクライアントへ送り、
レコード間で 1 秒待機します。

## サーバー

```go file=cookbook/streaming-response/server.go
```

## クライアント

```sh
curl localhost:1323
```

### 出力

```js
{"Altitude":-97,"Latitude":37.819929,"Longitude":-122.478255}
{"Altitude":1899,"Latitude":39.096849,"Longitude":-120.032351}
{"Altitude":2619,"Latitude":37.865101,"Longitude":-119.538329}
{"Altitude":42,"Latitude":33.812092,"Longitude":-117.918974}
{"Altitude":15,"Latitude":37.77493,"Longitude":-122.419416}
```
