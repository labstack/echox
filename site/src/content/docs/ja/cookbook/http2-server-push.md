---
title: HTTP/2 Server Push
description: HTTP/2 で Web アセットをクライアントへ先行プッシュします。
sidebar:
  order: 10
---

HTTP/2 server push により、サーバーはリソースがリクエストされる前にクライアントへ送信できます。
ページが必要とすることが分かっているアセットの往復通信をなくせます。このレシピでは、HTML レスポンスと一緒に
ページの CSS、JavaScript、画像をプッシュします。

:::note
Server push には HTTP/2 接続が必要です。下で使う証明書を作成するには、
[自己署名 X.509 TLS 証明書を生成する](/ja/cookbook/http2/#1-generate-a-self-signed-x509-tls-certificate)
に従ってください。
:::

## 1. Web アセットを配信するルートを登録する

```go
e.Static("/", "static")
```

## 2. index.html を配信し、その依存関係をプッシュする

レスポンスを unwrap して基底の `http.ResponseWriter` にアクセスし、writer が `http.Pusher`
を実装している場合に各アセットをプッシュします。

```go
e.GET("/", func(c *echo.Context) (err error) {
	rw, err := echo.UnwrapResponse(c.Response())
	if err != nil {
		return
	}
	if pusher, ok := rw.ResponseWriter.(http.Pusher); ok {
		if err = pusher.Push("/app.css", nil); err != nil {
			return
		}
		if err = pusher.Push("/app.js", nil); err != nil {
			return
		}
		if err = pusher.Push("/echo.png", nil); err != nil {
			return
		}
	}
	return c.File("index.html")
})
```

:::tip
`http.Pusher` がサポートされている場合、Web アセットは先行してプッシュされます。
そうでない場合、クライアントはそれらを個別にリクエストします。
:::

## 3. TLS サーバーを起動する

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

## ソースコード

### index.html

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>HTTP/2 Server Push</title>
  <link rel="stylesheet" href="/app.css">
  <script src="/app.js"></script>
</head>
<body>
  <img class="echo" src="/echo.png" alt="Echo">
  <h2>The following static files are served via HTTP/2 server push</h2>
  <ul>
    <li><code>/app.css</code></li>
    <li><code>/app.js</code></li>
    <li><code>/echo.png</code></li>
  </ul>
</body>
</html>
```

### server.go

```go
package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	e.Static("/", "static")
	e.GET("/", func(c *echo.Context) (err error) {
		rw, err := echo.UnwrapResponse(c.Response())
		if err != nil {
			return
		}
		if pusher, ok := rw.ResponseWriter.(http.Pusher); ok {
			if err = pusher.Push("/app.css", nil); err != nil {
				return
			}
			if err = pusher.Push("/app.js", nil); err != nil {
				return
			}
			if err = pusher.Push("/echo.png", nil); err != nil {
				return
			}
		}
		return c.File("index.html")
	})

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```
