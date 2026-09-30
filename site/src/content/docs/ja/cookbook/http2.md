---
title: HTTP/2 サーバー
description: TLS 証明書で Echo を起動し、HTTP/2 でトラフィックを配信します。
sidebar:
  order: 9
---

HTTP/2 はリクエスト多重化、header 圧縮、server push によりレイテンシーを改善します。
Go の HTTP サーバーは TLS 上で HTTP/2 を自動的にネゴシエートするため、Echo で HTTP/2 を配信するには
証明書付きでサーバーを起動すれば済みます。

<a id="1-generate-a-self-signed-x509-tls-certificate"></a>

## 1. 自己署名 X.509 TLS 証明書を生成する

次のコマンドで `cert.pem` と `key.pem` を生成します。

```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

:::note
デモ目的のため、自己署名証明書を使います。本番では
[certificate authority](https://en.wikipedia.org/wiki/Certificate_authority) から証明書を取得してください。
:::

## 2. リクエスト情報を echo するハンドラを作成する

```go file=cookbook/http2/server.go#handler
```

## 3. TLS サーバーを起動する

生成した証明書と key でサーバーを起動します。

```go file=cookbook/http2/server.go#start-tls
```

または、独自の `tls.Config` を持つカスタム `http.Server` を使います。

```go file=cookbook/http2/server.go#custom-server
```

## 4. 検証する

サーバーを起動して `https://localhost:1323/request` にアクセスします。次のような出力が表示されます。

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## ソースコード

```go file=cookbook/http2/server.go
```
