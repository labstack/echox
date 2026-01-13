---
description: Context Timeout middleware
---

# Context Timeout

Timeout middleware is used to timeout request context within a predefined period so context aware methods could return
early.

## Usage

```go
e.Use(middleware.ContextTimeout(60 * time.Second))
```

## Custom Configuration

```go
e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
	// Skipper defines a function to skip middleware.
	Skipper: nil,
	// ErrorHandler is a function when error aries in middleware execution.
	ErrorHandler: nil,
	// Timeout configures a timeout for the middleware, defaults to 0 for no timeout
	Timeout: 60 * time.Second,
}))
```


