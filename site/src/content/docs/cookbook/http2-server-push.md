---
title: HTTP/2 Server Push
description: Push web assets to the client proactively over HTTP/2.
sidebar:
  order: 10
---

HTTP/2 server push lets the server send resources to the client before they are
requested, eliminating a round trip for assets the page is known to need. This
recipe pushes a page's CSS, JavaScript, and image alongside the HTML response.

:::note
Server push requires an HTTP/2 connection. Follow [Generate a self-signed X.509
TLS certificate](/cookbook/http2/#1-generate-a-self-signed-x509-tls-certificate)
to create the certificate used below.
:::

## 1. Register a route to serve web assets

```go
e.Static("/", "static")
```

## 2. Serve index.html and push its dependencies

Unwrap the response to access the underlying `http.ResponseWriter`, then push each
asset if the writer implements `http.Pusher`:

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
When `http.Pusher` is supported, the web assets are pushed proactively; otherwise
the client falls back to requesting them separately.
:::

## 3. Start the TLS server

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

## Source code

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
