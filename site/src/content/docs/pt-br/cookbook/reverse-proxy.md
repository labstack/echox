---
title: Reverse Proxy
description: Use Echo como reverse proxy e load balancer na frente de aplicações upstream.
sidebar:
  order: 19
---

Esta receita demonstra como usar Echo como reverse proxy e load balancer na
frente das suas aplicações, como WordPress, Node.js, Java, Python, Ruby ou Go.
Por simplicidade, os upstreams aqui são servidores Go que também tratam WebSocket.

## 1) Identificar URL(s) de destino upstream

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

## 2) Configurar middleware de proxy com destinos upstream

O snippet abaixo usa load balancing round-robin. Você também pode usar
`middleware.NewRandomBalancer()`.

```go
e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
```

Para configurar um proxy para uma sub-rota, use `Echo#Group()`.

```go
g := e.Group("/blog")
g.Use(middleware.Proxy(...))
```

## 3) Iniciar servidores upstream

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) Iniciar o servidor proxy

```sh
go run server.go
```

Acesse `http://localhost:1323`, e você deverá ver uma página com um request HTTP
servido pelo "server 1" e um request WebSocket servido pelo "server 2".

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Código-fonte

### Servidor upstream

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

### Servidor proxy

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

## Atualização de segurança do Echo (v5.4.0 / v4.16.0)

Proxy define `X-Forwarded-Proto` a partir de `Context#Scheme()` e remove os outros cabeçalhos de esquema. No v5, `X-Real-IP` vem de `Context#RealIP()`. Se houver outro proxy antes do Echo, configure os extratores de [esquema](/pt-br/guide/request-scheme/) e [IP](/pt-br/guide/ip-address/) para transmitir valores corretos ao upstream.
