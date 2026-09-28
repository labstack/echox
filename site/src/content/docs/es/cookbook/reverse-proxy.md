---
title: Reverse Proxy
description: Usa Echo como reverse proxy y load balancer delante de aplicaciones upstream.
sidebar:
  order: 19
---

Esta receta demuestra cómo usar Echo como reverse proxy y load balancer delante de tus
aplicaciones, como WordPress, Node.js, Java, Python, Ruby o Go. Para simplificar, aquí
los upstreams son servidores Go que también manejan WebSocket.

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

## 2) Configurar middleware proxy con destinos upstream

El snippet de abajo usa load balancing round-robin. También puedes usar
`middleware.NewRandomBalancer()`.

```go
e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
```

Para configurar un proxy para una subruta, usa `Echo#Group()`.

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

## 4) Iniciar el servidor proxy

```sh
go run server.go
```

Abre `http://localhost:1323`, y deberías ver una página web con un request HTTP
servido desde "server 1" y un request WebSocket servido desde "server 2".

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Código fuente

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

## Actualización de seguridad de Echo (v5.4.0 / v4.16.0)

Proxy reenvía `X-Forwarded-Proto` desde `Context#Scheme()` y elimina las otras cabeceras de esquema. En v5, `X-Real-IP` procede de `Context#RealIP()`. Si hay otro proxy delante de Echo, configura los extractores de [esquema](/es/guide/request-scheme/) e [IP](/es/guide/ip-address/) para que el upstream reciba valores correctos.
