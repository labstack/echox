---
title: Customization
menu:
  side:
    parent: guide
    weight: 3
---

### HTTP Error Handler

`Echo#SetHTTPErrorHandler(h HTTPErrorHandler)` registers a custom `Echo#HTTPErrorHandler`.

Default HTTP error handler rules:

- If error is of type `Echo#HTTPError` it sends HTTP response with status code `HTTPError.Code`
and message `HTTPError.Message`.
- Else it sends `500 - Internal Server Error`.
- If debug mode is enabled, it uses `error.Error()` as status message.

### Debug

`Echo#SetDebug(on bool)` enable/disable debug mode.

### Log Prefix

`Echo#SetLogPrefix(prefix string)` sets the prefix for the logger. Default value
is `echo`.

### Log Output

`Echo#SetLogOutput(w io.Writer)` sets the output destination for the logger. Default
value is `os.Stdout`

To completely disable logs use `Echo#SetLogOutput(io.Discard)`

### Log Level

`Echo#SetLogLevel(l log.Level)`

SetLogLevel sets the log level for the logger. Default value is `log.ERROR`.

### HTTP Engine

Echo currently supports standard and [fasthttp](https://github.com/valyala/fasthttp)
server engines. Echo utilizes interfaces to abstract the internal implementation
of these servers so you can seamlessly switch from one engine to another based on
your preference.

#### Running a standard HTTP server

`e.Run(standard.New(":1323"))`

#### Running a fasthttp server

`e.Run(fasthttp.New(":1323"))`

#### Running a server from TLS configuration

`e.Run(<engine>.NewFromTLS(":1323", "<certfile>", "<keyfile>"))`

#### Running a server from engine configuration

`e.Run(<engine>.NewFromConfig(<config>))`

#### Configuration

```go
Config struct {
  Address      string        // TCP address to listen on.
  Listener     net.Listener  // Custom `net.Listener`. If set, server accepts connections on it.
  TLSCertfile  string        // TLS certificate file path.
  TLSKeyfile   string        // TLS key file path.
  ReadTimeout  time.Duration // Maximum duration before timing out read of the request.
  WriteTimeout time.Duration // Maximum duration before timing out write of the response.
}
```
