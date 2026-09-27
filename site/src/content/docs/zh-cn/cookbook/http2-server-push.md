---
title: HTTP/2 Server Push
description: 通过 HTTP/2 主动向客户端推送 Web 资源。
sidebar:
  order: 10
---

HTTP/2 server push 允许服务器在客户端请求资源之前发送资源，消除页面已知所需资源的一次往返。
此示例会在 HTML 响应旁推送页面的 CSS、JavaScript 和图片。

:::note
Server push 需要 HTTP/2 连接。请按照[生成自签名 X.509 TLS 证书](/zh-cn/cookbook/http2/#1-generate-a-self-signed-x509-tls-certificate)
创建下面使用的证书。
:::

## 1. 注册提供 Web 资源的路由

```go
e.Static("/", "static")
```

## 2. 提供 index.html 并推送依赖

解包响应以访问底层 `http.ResponseWriter`，然后在 writer 实现 `http.Pusher` 时推送每个资源：

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
当支持 `http.Pusher` 时，Web 资源会被主动推送；否则客户端会回退为分别请求它们。
:::

## 3. 启动 TLS 服务器

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

## 源码

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
