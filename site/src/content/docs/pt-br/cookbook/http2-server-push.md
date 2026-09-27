---
title: HTTP/2 Server Push
description: Envie web assets ao cliente proativamente por HTTP/2.
sidebar:
  order: 10
---

HTTP/2 server push permite que o servidor envie recursos ao cliente antes que eles sejam
solicitados, eliminando uma ida e volta para assets que a página já sabe que vai precisar. Esta
receita envia CSS, JavaScript e imagem de uma página junto da response HTML.

:::note
Server push requer uma conexão HTTP/2. Siga [Gerar um certificado TLS X.509
autoassinado](/pt-br/cookbook/http2/#1-generate-a-self-signed-x509-tls-certificate)
para criar o certificado usado abaixo.
:::

## 1. Registrar uma rota para servir web assets

```go
e.Static("/", "static")
```

## 2. Servir index.html e enviar suas dependências

Desembrulhe a response para acessar o `http.ResponseWriter` subjacente e então envie cada
asset se o writer implementar `http.Pusher`:

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
Quando `http.Pusher` tem suporte, os web assets são enviados proativamente; caso contrário,
o cliente volta a solicitá-los separadamente.
:::

## 3. Iniciar o servidor TLS

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

## Código-fonte

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
