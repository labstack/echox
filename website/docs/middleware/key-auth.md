---
description: Key auth middleware
---

# Key Auth

Key auth middleware provides a key based authentication.

- For valid key it calls the next handler.
- For invalid key, it sends "401 - Unauthorized" response.
- For missing key, it sends "400 - Bad Request" response.

## Usage

```go
e.Use(middleware.KeyAuth(func(c *echo.Context, key string, source middleware.ExtractorSource) (bool, error) {
  return key == "valid-key", nil
}))
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
  KeyLookup: "query:api-key",
  Validator: func(c *echo.Context, key string, source middleware.ExtractorSource) (bool, error) {
			return key == "valid-key", nil
		},
}))
```

## Configuration

```go
type KeyAuthConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// KeyLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract key from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// 			`<cut-prefix>` is argument value to cut/trim prefix of the extracted value. This is useful if header
	//			value has static prefix like `Authorization: <auth-scheme> <authorisation-parameters>` where part that we
	//			want to cut is `<auth-scheme> ` note the space at the end.
	//			In case of basic authentication `Authorization: Basic <credentials>` prefix we want to remove is `Basic `.
	// - "query:<name>"
	// - "form:<name>"
	// - "cookie:<name>"
	// Multiple sources example:
	// - "header:Authorization,header:X-Api-Key"
	KeyLookup string

	// AllowedCheckLimit set how many KeyLookup values are allowed to be checked. This is
	// useful environments like corporate test environments with application proxies restricting
	// access to environment with their own auth scheme.
	AllowedCheckLimit uint

	// Validator is a function to validate key.
	// Required.
	Validator KeyAuthValidator

	// ErrorHandler defines a function which is executed when all lookups have been done and none of them passed Validator
	// function. ErrorHandler is executed with last missing (ErrExtractionValueMissing) or an invalid key.
	// It may be used to define a custom error.
	//
	// Note: when error handler swallows the error (returns nil) middleware continues handler chain execution towards handler.
	// This is useful in cases when portion of your site/api is publicly accessible and has extra features for authorized users
	// In that case you can use ErrorHandler to set default public auth value to request and continue with handler chain.
	ErrorHandler KeyAuthErrorHandler

	// ContinueOnIgnoredError allows the next middleware/handler to be called when ErrorHandler decides to
	// ignore the error (by returning `nil`).
	// This is useful when parts of your site/api allow public access and some authorized routes provide extra functionality.
	// In that case you can use ErrorHandler to set a default public key auth value in the request context
	// and continue. Some logic down the remaining execution chain needs to check that (public) key auth value then.
	ContinueOnIgnoredError bool
}
```

### Default Configuration

```go
DefaultKeyAuthConfig = KeyAuthConfig{
  Skipper:    DefaultSkipper,
  KeyLookup:  "header:" + echo.HeaderAuthorization,
  AuthScheme: "Bearer",
}
```
