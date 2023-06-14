+++
title = "Rate Limiter Middleware"
description = "Rate limiter middleware for Echo"
[menu.main]
  name = "Rate Limiter"
  parent = "middleware"
+++

`RateLimiter` provides a Rate Limiter middleware for limiting the amount of requests to the server from a particular IP or id within a time period.

By default an in-memory store is used for keeping track of requests. The default in-memory implementation is focused on correctness and
may not be the best option for a high number of concurrent requests or a large number of different identifiers (>16k).

### Usage

To add a rate limit to your application simply add the `RateLimiter` middleware.
The example below will limit the application to 20 requests/sec using the default in-memory store:

```go
e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
```

Note: if the provided rate is a float number, Burst will be treated as the rounded down value of the rate.

## Custom Configuration

```go
config := middleware.RateLimiterConfig{
    Skipper: middleware.DefaultSkipper,
    Store: middleware.NewRateLimiterMemoryStoreWithConfig(
        middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
    ),
    IdentifierExtractor: func(ctx echo.Context) (string, error) {
        id := ctx.RealIP()
        return id, nil
    },
    ErrorHandler: func(context echo.Context, err error) error {
        return context.JSON(http.StatusForbidden, nil)
    },
    DenyHandler: func(context echo.Context, identifier string,err error) error {
        return context.JSON(http.StatusTooManyRequests, nil)
    },
}

e.Use(middleware.RateLimiterWithConfig(config))
```

Note: If you need to implement your own store, be sure to implement the RateLimiterStore interface and pass it to RateLimiterConfig and you're good to go!

## Configuration

```go
type RateLimiterConfig struct {
    Skipper    Skipper
    BeforeFunc BeforeFunc
    // IdentifierExtractor uses echo.Context to extract the identifier for a visitor
    IdentifierExtractor Extractor
    // Store defines a store for the rate limiter
    Store RateLimiterStore
    // ErrorHandler provides a handler to be called when IdentifierExtractor returns a non-nil error
    ErrorHandler func(context echo.Context, err error) error
    // DenyHandler provides a handler to be called when RateLimiter denies access
    DenyHandler func(context echo.Context, identifier string, err error) error
}
```

### Default Configuration

```go
// DefaultRateLimiterConfig defines default values for RateLimiterConfig
var DefaultRateLimiterConfig = RateLimiterConfig{
	Skipper: DefaultSkipper,
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		id := ctx.RealIP()
		return id, nil
	},
	ErrorHandler: func(context echo.Context, err error) error {
		return &echo.HTTPError{
			Code:     ErrExtractorError.Code,
			Message:  ErrExtractorError.Message,
			Internal: err,
		}
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return &echo.HTTPError{
			Code:     ErrRateLimitExceeded.Code,
			Message:  ErrRateLimitExceeded.Message,
			Internal: err,
		}
	},
}
```
### Errors

```go
var (
	// ErrRateLimitExceeded denotes an error raised when rate limit is exceeded
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	// ErrExtractorError denotes an error raised when extractor function is unsuccessful
	ErrExtractorError = echo.NewHTTPError(http.StatusForbidden, "error while extracting identifier")
)
```
