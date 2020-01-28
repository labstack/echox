+++
title = "Customization"
description = "Customizing Echo"
[menu.main]
  name = "Customization"
  parent = "guide"
weight = 2
+++

## Debug

`Echo#Debug` can be used to enable / disable debug mode. Debug mode sets the log level
to `DEBUG`.

## Logging

The default format for logging is JSON, which can be changed by modifying the header.

### Log Header

`Echo#Logger.SetHeader(io.Writer)` can be used to set the header for
the logger. Default value:

```json
{"time":"${time_rfc3339_nano}","level":"${level}","prefix":"${prefix}","file":"${short_file}","line":"${line}"}
```

*Example*
```go
import "github.com/labstack/gommon/log"

/* ... */

if l, ok := e.Logger.(*log.Logger); ok {
  l.SetHeader("${time_rfc3339} ${level}")
}
```

```sh
2018-05-08T20:30:06-07:00 INFO info
```

#### Available Tags

- `time_rfc3339`
- `time_rfc3339_nano`
- `level`
- `prefix`
- `long_file`
- `short_file`
- `line`

### Log Output

`Echo#Logger.SetOutput(io.Writer)` can be used to set the output destination for
the logger. Default value is `os.Stdout`

To completely disable logs use `Echo#Logger.SetOutput(ioutil.Discard)` or `Echo#Logger.SetLevel(log.OFF)`

### Log Level

`Echo#Logger.SetLevel(log.Lvl)` can be used to set the log level for the logger.
Default value is `ERROR`. Possible values:

- `DEBUG`
- `INFO`
- `WARN`
- `ERROR`
- `OFF`

### Custom Logger

Logging is implemented using `echo.Logger` interface which allows you to register
a custom logger using `Echo#Logger`.

## Custom Server

`Echo#StartServer()` can be used to run a custom server.

*Example*

```go
s := &http.Server{
  Addr:         ":1323",
  ReadTimeout:  20 * time.Minute,
  WriteTimeout: 20 * time.Minute,
}
e.Logger.Fatal(e.StartServer(s))
```

## Custom HTTP/2 Cleartext Server

`Echo#StartH2CServer()` can be used to run a custom HTTP/2 cleartext server.

*Example*

```go
import "golang.org/x/net/http2"

s := &http2.Server{
  MaxConcurrentStreams: 250,
  MaxReadFrameSize:     1048576,
  IdleTimeout:          10 * time.Second,
}
e.Logger.Fatal(e.StartH2CServer(":1323", s))
```

## Startup Banner

`Echo#HideBanner` can be used to hide the startup banner.

## Custom Listener

`Echo#*Listener` can be used to run a custom listener.

*Example*

```go
l, err := net.Listen("tcp", ":1323")
if err != nil {
  e.Logger.Fatal(l)
}
e.Listener = l
e.Logger.Fatal(e.Start(""))
```

## Disable HTTP/2

`Echo#DisableHTTP2` can be used disable HTTP/2 protocol.

## Read Timeout

`Echo#*Server#ReadTimeout` can be used to set the maximum duration before timing out read
of the request.

## Write Timeout

`Echo#*Server#WriteTimeout` can be used to set the maximum duration before timing out write
of the response.

## Validator

`Echo#Validator` can be used to register a validator for performing data validation
on request payload.

[Learn more](/guide/request#validate-data)

## Custom Binder

`Echo#Binder` can be used to register a custom binder for binding request payload.

[Learn more](/guide/request/#custom-binder)

## Renderer

`Echo#Renderer` can be used to register a renderer for template rendering.

[Learn more](/guide/templates)

## HTTP Error Handler

`Echo#HTTPErrorHandler` can be used to register a custom http error handler.

[Learn more](/guide/error-handling)
