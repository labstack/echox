---
description: HTTP/2 server push recipe
---

# HTTP/2 Server Push

:::note

Requires go1.8+

:::

## Send web assets using HTTP/2 server push

### [Generate a self-signed X.509 TLS certificate](http2#step-1-generate-a-self-signed-x-509-tls-certificate)

### 1) Register a route to serve web assets

```go
e.Static("/", "static")
```

### 2) Create a handler to serve index.html and push it's dependencies

```go
e.GET("/", func(c *echo.Context) (err error) {
  pusher, ok := c.Response().Writer.(http.Pusher)
  if ok {
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

:::info

If `http.Pusher` is supported, web assets are pushed; otherwise, client makes separate requests to get them.

:::

### 3) Start TLS server using cert.pem and key.pem

```go
if err := e.StartTLS(":1323", "cert.pem", "key.pem"); err != http.ErrServerClosed {
  log.Fatal(err)
}
```
or use customized HTTP server with your own TLSConfig
```go
s := http.Server{
  Addr:    ":8443",
  Handler: e, // set Echo as handler
  TLSConfig: &tls.Config{
    //Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
  },
  //ReadTimeout: 30 * time.Second, // use custom timeouts
}
if err := s.ListenAndServeTLS("cert.pem", "key.pem"); err != http.ErrServerClosed {
  log.Fatal(err)
}
```

### 4) Start the server and browse to https://localhost:1323

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## Source Code

```html reference
https://github.com/labstack/echox/blob/master/cookbook/http2-server-push/index.html
```

```go reference
https://github.com/labstack/echox/blob/master/cookbook/http2-server-push/server.go
```
