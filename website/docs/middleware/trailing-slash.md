---
description: Trailing slash middleware
---

# Trailing Slash

## Add Trailing Slash  

Add trailing slash middleware adds a trailing slash to the request URI.

### Usage

```go
e := echo.New()
e.Pre(middleware.AddTrailingSlash())
```

## Remove Trailing Slash

Remove trailing slash middleware removes a trailing slash from the request URI.

### Usage

```go
e := echo.New()
e.Pre(middleware.RemoveTrailingSlash())
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(middleware.AddTrailingSlashWithConfig(middleware.AddTrailingSlashConfig{
  RedirectCode: http.StatusMovedPermanently,
}))
```

Example above will add a trailing slash to the request URI and redirect with `301 - StatusMovedPermanently`.

## Configuration

```go
type AddTrailingSlashConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Status code to be used when redirecting the request.
	// Optional, but when provided the request is redirected using this code.
	// Valid status codes: [300...308]
	RedirectCode int
}
```

```go
type RemoveTrailingSlashConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Status code to be used when redirecting the request.
	// Optional, but when provided the request is redirected using this code.
	RedirectCode int
}
```

