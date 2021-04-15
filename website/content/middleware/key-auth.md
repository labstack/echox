+++
title = "Key Auth Middleware"
description = "Key auth middleware for Echo"
[menu.main]
  name = "Key Auth"
  parent = "middleware"
+++

Key auth middleware provides a key based authentication.

- For valid key it calls the next handler.
- For invalid key, it sends "401 - Unauthorized" response.
- For missing key, it sends "400 - Bad Request" response.

*Usage*

```go
e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
  return key == "valid-key", nil
}))
```

## Custom Configuration

*Usage*

```go
e := echo.New()
e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
  KeyLookup: "query:api-key",
  Validator: func(key string, c echo.Context) (bool, error) {
			return key == "valid-key", nil
		},
}))
```

## Configuration

```go
KeyAuthConfig struct {
  // Skipper defines a function to skip middleware.
  Skipper Skipper

  // KeyLookup is a string in the form of "<source>:<name>" that is used
  // to extract key from the request.
  // Optional. Default value "header:Authorization".
  // Possible values:
  // - "header:<name>"
  // - "query:<name>"
  // - "form:<name>"
  KeyLookup string `yaml:"key_lookup"`

  // AuthScheme to be used in the Authorization header.
  // Optional. Default value "Bearer".
  AuthScheme string

  // Validator is a function to validate key.
  // Required.
  Validator KeyAuthValidator

  // ErrorHandler defines a function which is executed for an invalid key.
  // It may be used to define a custom error.
  ErrorHandler KeyAuthErrorHandler
}
```

*Default Configuration*

```go
DefaultKeyAuthConfig = KeyAuthConfig{
  Skipper:    DefaultSkipper,
  KeyLookup:  "header:" + echo.HeaderAuthorization,
  AuthScheme: "Bearer",
}
```
