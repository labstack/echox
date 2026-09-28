---
title: リバースプロキシ
description: 上流アプリケーションの前段で Echo をリバースプロキシ兼ロードバランサーとして使います。
sidebar:
  order: 19
---

このレシピでは、WordPress、Node.js、Java、Python、Ruby、Go などのアプリケーションの前段で、
Echo をリバースプロキシ兼ロードバランサーとして使う方法を示します。簡単のため、ここでの上流は
WebSocket も処理する Go サーバーです。

## 1) 上流ターゲット URL を特定する

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

## 2) 上流ターゲットでプロキシミドルウェアを設定する

下のスニペットは round-robin ロードバランシングを使います。
`middleware.NewRandomBalancer()` も使えます。

```go
e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
```

サブルートにプロキシを設定するには、`Echo#Group()` を使います。

```go
g := e.Group("/blog")
g.Use(middleware.Proxy(...))
```

## 3) 上流サーバーを起動する

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) プロキシサーバーを起動する

```sh
go run server.go
```

`http://localhost:1323` にアクセスすると、HTTP リクエストは "server 1" から、
WebSocket リクエストは "server 2" から配信される Web ページが表示されます。

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## ソースコード

### 上流サーバー

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

### プロキシサーバー

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

## Echo のセキュリティ更新 (v5.4.0 / v4.16.0)

Proxy は `Context#Scheme()` から `X-Forwarded-Proto` を設定し、古いスキームヘッダーを削除します。v5 の `X-Real-IP` は `Context#RealIP()` から設定されます。Echo の前段に別のプロキシがある場合は[スキーム](/ja/guide/request-scheme/)と [IP](/ja/guide/ip-address/) の抽出器を設定し、上流へ正しい値を渡してください。
