+++
title = "Rate Limiter middleware"
description = "Rate limiter middleware for Echo"
[menu.main]
  name = "Rate limiter"
  parent = "middleware"
+++

RateLimiter provides a Rate Limiter middleware for limiting amount of requests to the server from a particular IP or id within a time period

By default, it stores an in-memory map of users. There are quite a few notable limitations of the in-memory implementation:  
- If your app receives above 100 parallel requests at once then this might not be the best option
- Does not play well with a large number of identifiers (>= 16k) 

*Usage:*  
For example this will allow a max of 20 requests/sec (uses in-memory):
```go
e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
```

## Custom Configuration

*Usage*
```go
config := middleware.RateLimiterConfig{
    Skipper: DefaultSkipper,
    Store: middleware.NewRateLimiterMemoryStore(
        middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute}
    )
    IdentifierExtractor: func(ctx echo.Context) (string, error) {
        id := ctx.RealIP()
        return id, nil
    },
    ErrorHandler: func(context echo.Context, err error) error {
        return context.JSON(http.StatusTooManyRequests, nil)
    },
    DenyHandler: func(context echo.Context, identifier string) error {
        return context.JSON(http.StatusForbidden, nil)
    },
}

e.Use(middleware.RateLimiterWithConfig(config))
```
NB: If you need to use a different store, just create one that implements the `Allow` method and pass it to RateLimiterConfig and you're good to go!

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

*Default Configuration*

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
*Errors*

```go
// errors
var (
	// ErrRateLimitExceeded denotes an error raised when rate limit is exceeded
	ErrRateLimitExceeded = echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
	// ErrExtractorError denotes an error raised when extractor function is unsuccessful
	ErrExtractorError = echo.NewHTTPError(http.StatusForbidden, "error while extracting identifier")
)
```

## [Example](/cookbook/jwt)
