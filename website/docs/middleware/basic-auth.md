---
description: Basic auth middleware
---

# Basic Auth

Basic auth middleware provides an HTTP basic authentication.

- For valid credentials it calls the next handler.
- For missing or invalid credentials, it sends "401 - Unauthorized" response.

## Usage

```go
e.Use(middleware.BasicAuth(func(username, password string, c *echo.Context) (bool, error) {
	// Be careful to use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
		return true, nil
	}
	return false, nil
}))
```

## Custom Configuration

### Usage

```go
e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{}))
```

## Configuration

```go
type BasicAuthConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Validator is a function to validate BasicAuthWithConfig credentials. Note: if request contains multiple basic auth headers
	// this function would be called once for each header until first valid result is returned
	// Required.
	Validator BasicAuthValidator

	// Realm is a string to define realm attribute of BasicAuthWithConfig.
	// Default value "Restricted".
	Realm string

	// AllowedCheckLimit set how many headers are allowed to be checked. This is useful
	// environments like corporate test environments with application proxies restricting
	// access to environment with their own auth scheme.
	// Defaults to 1.
	AllowedCheckLimit uint
}
```

### Default Configuration

```go
DefaultBasicAuthConfig = BasicAuthConfig{
	Skipper: DefaultSkipper,
}
```
