---
description: Logger middleware
---

# Logger

[`RequestLogger`](https://github.com/labstack/echo/blob/master/middleware/request_logger.go) middleware logs the information about each HTTP request. 

## RequestLogger middleware

RequestLogger middleware allows developer fully to customize what is logged and how it is logged and is more suitable
for usage with 3rd party (structured logging) libraries.

You can quickly acquaint yourself with the values that the logger knows to extract by referring to the fields of the [`RequestLoggerConfig`](https://github.com/labstack/echo/blob/master/middleware/request_logger.go) structure below. Or click the link to view the most up-to-date details.
```go
type RequestLoggerConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// BeforeNextFunc defines a function that is called before next middleware or handler is called in chain.
	BeforeNextFunc func(c *echo.Context)
	// LogValuesFunc defines a function that is called with values extracted by logger from request/response.
	// Mandatory.
	LogValuesFunc func(c *echo.Context, v RequestLoggerValues) error

	// HandleError instructs logger to call global error handler when next middleware/handler returns an error.
	// This is useful when you have custom error handler that can decide to use different status codes.
	//
	// A side-effect of calling global error handler is that now Response has been committed and sent to the client
	// and middlewares up in chain can not change Response status code or response body.
	HandleError bool

	// LogLatency instructs logger to record duration it took to execute rest of the handler chain (next(c) call).
	LogLatency bool
	// LogProtocol instructs logger to extract request protocol (i.e. `HTTP/1.1` or `HTTP/2`)
	LogProtocol bool
	// LogRemoteIP instructs logger to extract request remote IP. See `echo.Context.RealIP()` for implementation details.
	LogRemoteIP bool
	// LogHost instructs logger to extract request host value (i.e. `example.com`)
	LogHost bool
	// LogMethod instructs logger to extract request method value (i.e. `GET` etc)
	LogMethod bool
	// LogURI instructs logger to extract request URI (i.e. `/list?lang=en&page=1`)
	LogURI bool
	// LogURIPath instructs logger to extract request URI path part (i.e. `/list`)
	LogURIPath bool
	// LogRoutePath instructs logger to extract route path part to which request was matched to (i.e. `/user/:id`)
	LogRoutePath bool
	// LogRequestID instructs logger to extract request ID from request `X-Request-ID` header or response if request did not have value.
	LogRequestID bool
	// LogReferer instructs logger to extract request referer values.
	LogReferer bool
	// LogUserAgent instructs logger to extract request user agent values.
	LogUserAgent bool
	// LogStatus instructs logger to extract response status code. If handler chain returns an echo.HTTPError,
	// the status code is extracted from the echo.HTTPError returned
	LogStatus bool
	// LogError instructs logger to extract error returned from executed handler chain.
	LogError bool
	// LogContentLength instructs logger to extract content length header value. Note: this value could be different from
	// actual request body size as it could be spoofed etc.
	LogContentLength bool
	// LogResponseSize instructs logger to extract response content length value. Note: when used with Gzip middleware
	// this value may not be always correct.
	LogResponseSize bool
	// LogHeaders instructs logger to extract given list of headers from request. Note: request can contain more than
	// one header with same value so slice of values is been logger for each given header.
	//
	// Note: header values are converted to canonical form with http.CanonicalHeaderKey as this how request parser converts header
	// names to. For example, the canonical key for "accept-encoding" is "Accept-Encoding".
	LogHeaders []string
	// LogQueryParams instructs logger to extract given list of query parameters from request URI. Note: request can
	// contain more than one query parameter with same name so slice of values is been logger for each given query param name.
	LogQueryParams []string
	// LogFormValues instructs logger to extract given list of form values from request body+URI. Note: request can
	// contain more than one form value with same name so slice of values is been logger for each given form value name.
	LogFormValues []string

}
```


### Examples

Example for naive `fmt.Printf`
```go
skipper := func(c *echo.Context) bool {
	// Skip health check endpoint
    return c.Request().URL.Path == "/health"
}
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogStatus: true,
	LogURI:    true,
	Skipper: skipper,
	BeforeNextFunc: func(c *echo.Context) {
		c.Set("customValueFromContext", 42)
	},
	LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
		value, _ := c.Get("customValueFromContext").(int)
		fmt.Printf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
		return nil
	},
}))
```
*Sample output*

```exec
REQUEST: uri: /hello, status: 200, custom-value: 42
```


Example for slog (https://pkg.go.dev/log/slog)
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
    LogStatus:   true,
    LogURI:      true,
    LogError:    true,
    HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
    LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
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
*Sample output*
```exec
{"time":"2024-12-30T20:55:46.2399999+08:00","level":"INFO","msg":"REQUEST","uri":"/hello","status":200}
```

Example for Zerolog (https://github.com/rs/zerolog)
```go
logger := zerolog.New(os.Stdout)
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
		logger.Info().
			Str("URI", v.URI).
			Int("status", v.Status).
			Msg("request")

		return nil
	},
}))
```
*Sample output*
```exec
{"level":"info","URI":"/hello","status":200,"message":"request"}
```

Example for Zap (https://github.com/uber-go/zap)
```go
logger, _ := zap.NewProduction()
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
		logger.Info("request",
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
		)

		return nil
	},
}))
```
*Sample output*
```exec
{"level":"info","ts":1735564026.3197417,"caller":"cmd/main.go:20","msg":"request","URI":"/hello","status":200}
```

Example for Logrus (https://github.com/sirupsen/logrus)
```go
log := logrus.New()
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:    true,
	LogStatus: true,
	LogValuesFunc: func(c *echo.Context, values middleware.RequestLoggerValues) error {
		log.WithFields(logrus.Fields{
			"URI":   values.URI,
			"status": values.Status,
		}).Info("request")

		return nil
	},
}))
```
*Sample output*
```exec
time="2024-12-30T21:08:49+08:00" level=info msg=request URI=/hello status=200
```

### Troubleshooting Tips

#### 1. Solution for "panic: missing LogValuesFunc callback function for request logger middleware"
This panic arises when the `LogValuesFunc` callback function, which is mandatory for the request logger middleware configuration, is left unset.

To address this, you must define a suitable function that adheres to the `LogValuesFunc` specifications and then assign it within the middleware configuration. Consider the following straightforward illustration:

```go
func logValues(c *echo.Context, v middleware.RequestLoggerValues) error {
    fmt.Printf("Request Method: %s, URI: %s\n", v.Method, v.URI)
    return nil
}

e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
    LogValuesFunc: logValues,
}))
```

#### 2. If Parameters in Logs Are Empty
When investigating logging-related glitches, if you notice that certain parameters like `v.URI` and `v.Status` within the `LogValuesFunc` function produce empty outputs, your focus should shift to validating the relevant configuration elements. Specifically, check whether the corresponding items (such as `LogStatus`, `LogURI`, etc.) in `e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{...}))` have been erroneously set to `false` or failed to activate properly due to miscellaneous factors. Ensure these configuration particulars are accurately configured so that the pertinent request and response data can be precisely logged. 