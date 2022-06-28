+++
title = "Logger Middleware"
description = "Logger middleware for Echo"
[menu.main]
  name = "Logger"
  parent = "middleware"
+++

Logger middleware logs the information about each HTTP request.

*Usage*

`e.Use(middleware.Logger())`

*Sample Output*

```js
{"time":"2017-01-12T08:58:07.372015644-08:00","remote_ip":"::1","host":"localhost:1323","method":"GET","uri":"/","status":200,"error":"","latency":14743,"latency_human":"14.743µs","bytes_in":0,"bytes_out":2}
```

## Custom Configuration

*Usage*

```go
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
  Format: "method=${method}, uri=${uri}, status=${status}\n",
}))
```

Example above uses a `Format` which logs request method and request URI.

*Sample Output*

```sh
method=GET, uri=/, status=200
```

## Skipper

*Usage*

```go
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
  Skipper: func(c echo.Context) bool {
    if c.Path() == "/ping" {
      return true
    }

    return false
  },
}))
```
You can skip the log of specific route like as above example.

## Configuration

```go
// LoggerConfig defines the config for Logger middleware.
LoggerConfig struct {
  // Skipper defines a function to skip middleware.
  Skipper Skipper

  // Tags to construct the logger format.
  //
  // - time_unix
  // - time_unix_nano
  // - time_rfc3339
  // - time_rfc3339_nano
  // - time_custom
  // - id (Request ID)
  // - remote_ip
  // - uri
  // - host
  // - method
  // - path
  // - protocol
  // - referer
  // - user_agent
  // - status
  // - error
  // - latency (In nanoseconds)
  // - latency_human (Human readable)
  // - bytes_in (Bytes received)
  // - bytes_out (Bytes sent)
  // - header:<NAME>
  // - query:<NAME>
  // - form:<NAME>
  //
  // Example "${remote_ip} ${status}"
  //
  // Optional. Default value DefaultLoggerConfig.Format.
  Format string `yaml:"format"`

  // Optional. Default value DefaultLoggerConfig.CustomTimeFormat.
  CustomTimeFormat string `yaml:"custom_time_format"`

  // Output is a writer where logs in JSON format are written.
  // Optional. Default value os.Stdout.
  Output io.Writer
}
```

*Default Configuration*

```go
DefaultLoggerConfig = LoggerConfig{
  Skipper: DefaultSkipper,
  Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
    `"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
    `"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
    `,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
  CustomTimeFormat: "2006-01-02 15:04:05.00000",
}
```
