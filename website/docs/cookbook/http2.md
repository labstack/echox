---
description: HTTP/2 server recipe
---

# HTTP/2 Server

## 1) Generate a self-signed X.509 TLS certificate 

Run the following command to generate `cert.pem` and `key.pem` files:

```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

:::note

For demo purpose, we are using a self-signed certificate. Ideally, you should obtain
a certificate from [CA](https://en.wikipedia.org/wiki/Certificate_authority).

:::

## 2) Create a handler which simply outputs the request information to the client

```go
e.GET("/request", func(c *echo.Context) error {
  req := c.Request()
  format := `
    <code>
      Protocol: %s<br>
      Host: %s<br>
      Remote Address: %s<br>
      Method: %s<br>
      Path: %s<br>
    </code>
  `
  return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
})
```

## 3) Start TLS server using cert.pem and key.pem

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
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

## 4) Start the server and browse to https://localhost:1323/request to see the following output

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## Source Code

```go reference
https://github.com/labstack/echox/blob/master/cookbook/http2/server.go
```
