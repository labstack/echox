---
description: Method override middleware
---

# Method Override

Method override middleware checks for the overridden method from the request and
uses it instead of the original method.

:::info

For security reasons, only `POST` method can be overridden.

:::

## Usage

```go
e.Pre(middleware.MethodOverride())
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
  Getter: middleware.MethodFromForm("_method"),
}))
```

## Configuration

```go
type MethodOverrideConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Getter is a function that gets overridden method from the request.
	// Optional. Default values MethodFromHeader(echo.HeaderXHTTPMethodOverride).
	Getter MethodOverrideGetter
}
```

### Default Configuration

```go
DefaultMethodOverrideConfig = MethodOverrideConfig{
  Skipper: DefaultSkipper,
  Getter:  MethodFromHeader(echo.HeaderXHTTPMethodOverride),
}
```
