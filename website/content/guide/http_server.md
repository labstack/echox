+++
title = "Starting HTTP(S) server with Echo"
description = "Serving Echo from HTTP(S) server"
[menu.main]
name = "HTTP(s) Server"
parent = "guide"
+++

Echo provides following convenience methods to start HTTP server with Echo as a request handler:

* `Echo.Start(address string)`
* `Echo.StartTLS(address string, certFile, keyFile interface{})`
* `Echo.StartAutoTLS(address string)`
* `Echo.StartH2CServer(address string, h2s *http2.Server)`
* `Echo.StartServer(s *http.Server)`

## HTTP Server

`Echo.Start` is convenience method that starts http server with Echo serving requests.
```go
func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  if err := e.Start(":8080"); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
```

Following is equivalent to `Echo.Start` previous example
```go
func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  s := http.Server{
    Addr:        ":8080",
    Handler:     e,
    //ReadTimeout: 30 * time.Second, // customize http.Server timeouts
  }
  if err := s.ListenAndServe(); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
```

## HTTPS Server

`Echo.StartTLS` is convenience method that starts HTTPS server with Echo serving requests on given address and uses 
`server.crt` and `server.key` as TLS certificate pair.
```go
func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  if err := e.StartTLS(":8443", "server.crt", "server.key"); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
```

Following is equivalent to `Echo.StartTLS` previous example
```go
func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  s := http.Server{
    Addr:    ":8443",
    Handler: e, // set Echo as handler
    TLSConfig: &tls.Config{
      //MinVersion: 1, // customize TLS configuration
    },
    //ReadTimeout: 30 * time.Second, // use custom timeouts
  }
  if err := s.ListenAndServeTLS("server.crt", "server.key"); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
```

## Auto TLS Server with Let’s Encrypt

See [Auto TLS Recipe](/cookbook/auto-tls#server)

## HTTP/2 Cleartext Server (HTTP2 over HTTP)

`Echo.StartH2CServer` is convenience method that starts a custom HTTP/2 cleartext server on given address
```go
func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  s := &http2.Server{
    MaxConcurrentStreams: 250,
    MaxReadFrameSize:     1048576,
    IdleTimeout:          10 * time.Second,
  }
  if err := e.StartH2CServer(":8080", s); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
```

Following is equivalent to `Echo.StartH2CServer` previous example
```go
func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  h2s := &http2.Server{
    MaxConcurrentStreams: 250,
    MaxReadFrameSize:     1048576,
    IdleTimeout:          10 * time.Second,
  }
  s := http.Server{
    Addr:    ":8080",
    Handler: h2c.NewHandler(e, h2s),
  }
  if err := s.ListenAndServe(); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
```
