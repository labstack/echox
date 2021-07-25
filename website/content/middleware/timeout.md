+++
title = "Timeout Middleware"
description = "Timeout middleware for Echo"
[menu.main]
  name = "Timeout"
  parent = "middleware"
+++

Timeout middleware is used to timeout at a long running operation within a predefined period.

*Usage*

`e.Use(middleware.Timeout())`

## Custom Configuration

*Usage*

```go
e := echo.New()
e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
  Skipper: Skipper,
  OnTimeoutRouteErrorHandler: func(err error, e echo.Context) error {
      // you can handle your error here, the returning error will be 
      // passed down the middleware chain
      return err
  },
  Timeout: 30*time.Second,
}))
```

## Configuration

```go
// TimeoutConfig defines the config for Timeout middleware.
TimeoutConfig struct {
    // Skipper defines a function to skip middleware.
    Skipper Skipper
    // ErrorMessage is written to response on timeout in addition to http.StatusServiceUnavailable (503) status code
    // It can be used to define a custom timeout error message
    ErrorMessage string
    // OnTimeoutRouteErrorHandler is an error handler that is executed for error that was returned
    // from wrapped route after OnTimeoutRouteErrorHandler func(err error, c echo.Context)
    OnTimeoutRouteErrorHandler func(err error, c echo.Context)
    // Timeout configures a timeout for the middleware, defaults to 0 for no timeout
    Timeout time.Duration
}
```

*Default Configuration*

```go
DefaultTimeoutConfig = TimeoutConfig{
    Skipper:      DefaultSkipper,
    Timeout:      0,
    ErrorMessage: "",
}
```
