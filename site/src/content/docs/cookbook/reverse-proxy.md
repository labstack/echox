---
title: Reverse Proxy
description: Use Echo as a reverse proxy and load balancer in front of upstream applications.
sidebar:
  order: 19
---

This recipe demonstrates how to use Echo as a reverse proxy and load balancer in
front of your applications, such as WordPress, Node.js, Java, Python, Ruby, or Go.
For simplicity, the upstreams here are Go servers that also handle WebSocket.

Proxy forwards `X-Forwarded-Proto` from `Context#Scheme()` and drops the older
`X-Forwarded-Ssl`, `X-Forwarded-Protocol`, and `X-Url-Scheme` headers. In v5 it
sets `X-Real-IP` from `Context#RealIP()`. If another proxy sits before Echo, set
the [scheme extractor](/guide/request-scheme/) and [IP extractor](/guide/ip-address/)
for that trusted proxy; otherwise the upstream may see the proxy's IP or the wrong
scheme.

## 1) Identify upstream target URL(s)

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

## 2) Set up proxy middleware with upstream targets

The snippet below uses round-robin load balancing. You can also use
`middleware.NewRandomBalancer()`.

```go
e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
```

To set up a proxy for a sub-route, use `Echo#Group()`.

```go
g := e.Group("/blog")
g.Use(middleware.Proxy(...))
```

## 3) Start upstream servers

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) Start the proxy server

```sh
go run server.go
```

Browse to `http://localhost:1323`, and you should see a webpage with an HTTP
request served from "server 1" and a WebSocket request served from "server 2".

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Source code

### Upstream server

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

### Proxy server

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
