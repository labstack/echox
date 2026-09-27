---
title: HTTP/2 Server Push
description: Envía web assets al cliente proactivamente sobre HTTP/2.
sidebar:
  order: 10
---

HTTP/2 server push permite que el servidor envíe recursos al cliente antes de que se
soliciten, eliminando un viaje de ida y vuelta para assets que la página sabe que necesitará.
Esta receta envía el CSS, JavaScript e imagen de una página junto con la response HTML.

:::note
Server push requiere una conexión HTTP/2. Sigue [Generar un certificado TLS X.509
autofirmado](/es/cookbook/http2/#1-generate-a-self-signed-x509-tls-certificate)
para crear el certificado usado abajo.
:::

## 1. Registrar una ruta para servir web assets

```go
e.Static("/", "static")
```

## 2. Servir index.html y hacer push de sus dependencias

Desenvuelve la response para acceder al `http.ResponseWriter` subyacente, y luego haz push
de cada asset si el writer implementa `http.Pusher`:

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
Cuando `http.Pusher` está soportado, los web assets se envían proactivamente; de lo contrario,
el cliente vuelve a solicitarlos por separado.
:::

## 3. Iniciar el servidor TLS

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

## Código fuente

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
