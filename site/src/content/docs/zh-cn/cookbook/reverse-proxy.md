---
title: 反向代理
description: 使用 Echo 作为上游应用前的反向代理和负载均衡器。
sidebar:
  order: 19
---

此示例演示如何在你的应用（如 WordPress、Node.js、Java、Python、Ruby 或 Go）前使用 Echo
作为反向代理和负载均衡器。为简单起见，这里的上游也是处理 WebSocket 的 Go 服务器。

## 1) 确定上游目标 URL

```go
url1, err := url.Parse("http://localhost:8081")
if err != nil {
  e.Logger.Error("failed parse url", "error", err)
}
url2, err := url.Parse("http://localhost:8082")
if err != nil {
  e.Logger.Error("failed parse url", "error", err)
}
targets := []*middleware.ProxyTarget{
  {
    URL: url1,
  },
  {
    URL: url2,
  },
}
```

## 2) 使用上游目标设置代理中间件

下面的片段使用轮询负载均衡。你也可以使用 `middleware.NewRandomBalancer()`。

```go
e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
```

要为子路由设置代理，请使用 `Echo#Group()`。

```go
g := e.Group("/blog")
g.Use(middleware.Proxy(...))
```

## 3) 启动上游服务器

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) 启动代理服务器

```sh
go run server.go
```

访问 `http://localhost:1323`，你应该会看到网页中 HTTP 请求由 "server 1" 提供，
WebSocket 请求由 "server 2" 提供。

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## 源码

### 上游服务器

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"golang.org/x/net/websocket"
)

var index = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Upstream Server</title>
		<style>
			h1, p {
				font-weight: 300;
			}
		</style>
	</head>
	<body>
		<h1>HTTP</h1>
		<p>
			Hello from upstream server %s
		</p>
		<h1>WebSocket</h1>
		<p id="output"></p>
		<script>
			var ws = new WebSocket('ws://localhost:1323/ws')

			ws.onmessage = function(evt) {
				var out = document.getElementById('output');
				out.innerHTML += evt.data + '<br>';
			}
		</script>
	</body>
	</html>
`

func main() {
	name := os.Args[1]
	port := os.Args[2]
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(index, name))
	})

	// WebSocket handler
	e.GET("/ws", func(c *echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				// Write
				err := websocket.Message.Send(ws, fmt.Sprintf("Hello from upstream server %s!", name))
				if err != nil {
					e.Logger.Error("failed to send message", "error", err)
				}
				select {
				case <-ws.Request().Context().Done():
					return
				case <-time.After(1 * time.Second):
					continue
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	sc := echo.StartConfig{Address: port}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

### 代理服务器

```go
package main

import (
	"context"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Setup proxy
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")
	targets := []*middleware.ProxyTarget{
		{URL: url1},
		{URL: url2},
	}
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

## Echo 安全更新 (v5.4.0 / v4.16.0)

Proxy 从 `Context#Scheme()` 设置 `X-Forwarded-Proto` 并删除旧协议标头。v5 的 `X-Real-IP` 来自 `Context#RealIP()`。如果 Echo 前面还有代理，请配置[协议](/zh-cn/guide/request-scheme/)和 [IP](/zh-cn/guide/ip-address/) 提取器，确保上游收到正确的值。
