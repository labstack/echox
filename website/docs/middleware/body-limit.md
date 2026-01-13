---
description: Body limit middleware
---

# Body Limit

Body limit middleware sets the maximum allowed size for a request body, if the
size exceeds the configured limit, it sends "413 - Request Entity Too Large"
response. The body limit is determined based on both `Content-Length` request
header and actual content read, which makes it super secure.

Limit is specified as bytes

## Usage

```go
e := echo.New()
e.Use(middleware.BodyLimit(2_097_152)) // 2MB
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{}))
```

## Configuration

```go
type BodyLimitConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// LimitBytes is maximum allowed size in bytes for a request body
	LimitBytes int64
}
```

### Default Configuration

```go
DefaultBodyLimitConfig = BodyLimitConfig{
  Skipper: DefaultSkipper,
}
```
