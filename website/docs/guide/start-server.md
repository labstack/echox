---
description: Starting server
slug: /start-server
sidebar_position: 7
---

# Start Server

Echo provides following `Echo.Start(address string)` convenience method to start the server. Which uses the default configuration for graceful shutdown.


## HTTP Server

`Echo.Start` is convenience method that starts http server with Echo serving requests.
```go
func main() {
	e := echo.New()
	// add middleware and routes
	// ...

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

same functionality using server configuration `echo.StartConfig`
```go
func main() {
	e := echo.New()
	// add middleware and routes
	// ...
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // start shutdown process on signal
	defer cancel()

	sc := echo.StartConfig{
		Address: ":1323",
		GracefulTimeout: 5 * time.Second, // defaults to 10 seconds
	}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

Following is server using `http.Server`
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
    e.Logger.Error("failed to start server", "error", err)
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

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.StartTLS(context.Background(), e, "server.crt", "server.key"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
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

See [Auto TLS Recipe](../cookbook/auto-tls.md#server)

## HTTP/2 Cleartext Server (HTTP2 over HTTP)

`Echo.StartH2CServer` is convenience method that starts a custom HTTP/2 cleartext server on given address
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
	h2Handler := h2c.NewHandler(e, h2s)

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), h2Handler); err != nil {
		e.Logger.Error("failed to start server", "error", err)
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