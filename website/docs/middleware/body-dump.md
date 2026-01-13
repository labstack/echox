---
description: Body dump middleware
---

# Body Dump

Body dump middleware captures the request and response payload and calls the registered handler. Generally used for debugging/logging purpose. Avoid using it if your request/response payload is huge e.g. file upload/download, but if you still need to, add an exception for your endpoints in the skipper function.

## Usage

```go
e := echo.New()
e.Use(middleware.BodyDump(func(c *echo.Context, reqBody []byte, resBody []byte, err error) {
  // logic to handle request body and response body
}))
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{}))
```

## Configuration

```go
type BodyDumpConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Handler receives request, response payloads and handler error if there are any.
	// Required.
	Handler BodyDumpHandler

	// MaxRequestBytes limits how much of the request body to dump.
	// If the request body exceeds this limit, only the first MaxRequestBytes
	// are dumped. The handler callback receives truncated data.
	// Default: 5 * MB (5,242,880 bytes)
	// Set to -1 to disable limits (not recommended in production).
	MaxRequestBytes int64

	// MaxResponseBytes limits how much of the response body to dump.
	// If the response body exceeds this limit, only the first MaxResponseBytes
	// are dumped. The handler callback receives truncated data.
	// Default: 5 * MB (5,242,880 bytes)
	// Set to -1 to disable limits (not recommended in production).
	MaxResponseBytes int64
}
```

### Default Configuration*

```go
DefaultBodyDumpConfig = BodyDumpConfig{
  Skipper: DefaultSkipper,
}
```
