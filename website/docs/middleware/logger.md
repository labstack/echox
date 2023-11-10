---
description: Logger middleware
---

# Logger

Logger middleware logs the information about each HTTP request. 

Echo has 2 different logger middlewares:

- Older string template based logger [`Logger`](https://github.com/labstack/echo/blob/master/middleware/logger.go) - easy to start with but has limited capabilities
- Newer customizable function based logger [`RequestLogger`](https://github.com/labstack/echo/blob/master/middleware/request_logger.go) - allows developer fully to customize what is logged and how it is logged. Suitable for usage with 3rd party logger libraries.

## Old Logger middleware

## Usage

```go
e.Use(middleware.Logger())
```

*Sample output*

```js
{"time":"2017-01-12T08:58:07.372015644-08:00","remote_ip":"::1","host":"localhost:1323","method":"GET","uri":"/","status":200,"error":"","latency":14743,"latency_human":"14.743µs","bytes_in":0,"bytes_out":2}
```

## Custom Configuration

### Usage

```go
e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
  Format: "method=${method}, uri=${uri}, status=${status}\n",
}))
```

Example above uses a `Format` which logs request method and request URI.

*Sample output*

```sh
method=GET, uri=/, status=200
```

## Configuration

```go
// LoggerConfig defines the config for Logger middleware.
LoggerConfig struct {
  // Skipper defines a function to skip middleware.
  Skipper Skipper

  // Tags to construct the logger format.
  //
  // - time_unix
  // - time_unix_milli
  // - time_unix_micro
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

### Default Configuration

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

# New RequestLogger middleware

RequestLogger middleware allows developer fully to customize what is logged and how it is logged and is more suitable
for usage with 3rd party (structured logging) libraries.

See [`RequestLoggerConfig`](https://github.com/labstack/echo/blob/master/middleware/request_logger.go) structure fields for values that logger knows to extract.


### Examples

Example for naive `fmt.Printf`
```go
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogStatus: true,
	LogURI:    true,
	BeforeNextFunc: func(c echo.Context) {
		c.Set("customValueFromContext", 42)
	},
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		value, _ := c.Get("customValueFromContext").(int)
		fmt.Printf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
		return nil
	},
}))
```

Example for slog (https://pkg.go.dev/log/slog)
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
    LogStatus:   true,
    LogURI:      true,
    LogError:    true,
    HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
    LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
        if v.Error == nil {
            logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
                slog.String("uri", v.URI),
                slog.Int("status", v.Status),
            )
        } else {
            logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
                slog.String("uri", v.URI),
                slog.Int("status", v.Status),
                slog.String("err", v.Error.Error()),
            )
        }
        return nil
    },
}))
```

Example for Zerolog (https://github.com/rs/zerolog)
```go
logger := zerolog.New(os.Stdout)
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		logger.Info().
			Str("URI", v.URI).
			Int("status", v.Status).
			Msg("request")

		return nil
	},
}))
```

Example for Zap (https://github.com/uber-go/zap)
```go
logger, _ := zap.NewProduction()
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		logger.Info("request",
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
		)

		return nil
	},
}))
```

Example for Logrus (https://github.com/sirupsen/logrus)
```go
log := logrus.New()
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
		log.WithFields(logrus.Fields{
			"URI":   values.URI,
			"status": values.Status,
		}).Info("request")

		return nil
	},
}))
```
