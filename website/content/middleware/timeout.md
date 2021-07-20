+++
title = "Timeout Middleware"
description = "Timeout middleware for Echo"
[menu.main]
  name = "Timeout"
  parent = "middleware"
+++

Timeout middleware is used to timeout at a long running operation within a predefined period.

When timeout occurs, and the client receives timeout response the handler keeps running its code and keeps using resources until it finishes and returns!

> Timeout middleware is a serious performance hit as it buffers all responses from wrapped handler. Do not set it in front of file downloads or responses you want to stream to the client.

Timeout middleware is not a magic wand to hide slow handlers from clients. Consider designing/implementing asynchronous 
request/response API if (extremely) fast responses are to be expected and actual work can be done in background
Prefer handling timeouts in handler functions explicitly

*Usage*

`e.GET("/", handlerFunc, middleware.Timeout())`

## Custom Configuration

*Usage*

```go
// Echo instance
e := echo.New()

handlerFunc := func(c echo.Context) error {
    time.Sleep(10 * time.Second)
    return c.String(http.StatusOK, "Hello, World!\n")
}
middlewareFunc := middleware.TimeoutWithConfig(middleware.TimeoutConfig{
    Timeout: 30 * time.Second,
    ErrorMessage: "my custom error message",
})
// Handler with timeout middleware
e.GET("/", handlerFunc, middlewareFunc)
```

## Configuration

```go
// TimeoutConfig defines the config for Timeout middleware.
type TimeoutConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// ErrorMessage is written to response on timeout in addition to http.StatusServiceUnavailable (503) status code
	// It can be used to define a custom timeout error message
	ErrorMessage string

	// OnTimeoutRouteErrorHandler is an error handler that is executed for error that was returned from wrapped route after
	// request timeouted and we already had sent the error code (503) and message response to the client.
	// NB: do not write headers/body inside this handler. The response has already been sent to the client and response writer
	// will not accept anything no more. If you want to know what actual route middleware timeouted use `c.Path()`
	OnTimeoutRouteErrorHandler func(err error, c echo.Context)

	// Timeout configures a timeout for the middleware, defaults to 0 for no timeout
	// NOTE: when difference between timeout duration and handler execution time is almost the same (in range of 100microseconds)
	// the result of timeout does not seem to be reliable - could respond timeout, could respond handler output
	// difference over 500microseconds (0.5millisecond) response seems to be reliable
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

## Alternatively handle timeouts in handlers

```go
func main() {
    e := echo.New()

    doBusinessLogic := func(ctx context.Context, UID string) error {
        // NB: Do not use echo.JSON() or any other method that writes data/headers to client here. This function is executed
        // in different coroutine that should not access echo.Context and response writer

        log.Printf("uid: %v\n", UID)
        //res, err := slowDatabaseCon.ExecContext(ctx, query, args)
        time.Sleep(10 * time.Second) // simulate slow execution
        log.Print("doBusinessLogic done\n")
        return nil
    }

    handlerFunc := func(c echo.Context) error {
        defer log.Print("handlerFunc done\n")

        // extract and validate needed data from request and pass it to business function
        UID := c.QueryParam("uid")

        ctx, cancel := context.WithTimeout(c.Request().Context(), 5 * time.Second)
        defer cancel()
        result := make(chan error)
        go func() { // run actual business logic in separate coroutine
            defer func() { // unhandled panic in coroutine will crash the whole application
                if err := recover(); err != nil {
                    result <- fmt.Errorf("panic: %v", err)
                }
            }()
            result <- doBusinessLogic(ctx, UID)
        }()

        select { // wait until doBusinessLogic finishes or we timeout while waiting for the result
        case <-ctx.Done():
            err := ctx.Err()
            if err == context.DeadlineExceeded {
                return echo.NewHTTPError(http.StatusServiceUnavailable, "doBusinessLogic timeout")
            }
            return err // probably client closed the connection
        case err := <-result: // doBusinessLogic finishes
            if err != nil {
                return err
            }
        }
        return c.NoContent(http.StatusAccepted)
    }
    e.GET("/", handlerFunc)

    s := http.Server{Addr: ":8080", Handler: e}
    if err := s.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
```
