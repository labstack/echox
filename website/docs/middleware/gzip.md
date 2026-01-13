---
description: Gzip middleware
---

# Gzip

Gzip middleware compresses HTTP response using gzip compression scheme.

## Usage

`e.Use(middleware.Gzip())`

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
  Level: 5,
}))
```

:::tip

A middleware skipper can be passed to avoid gzip to certain URL(s).

:::

#### Example

```go
e := echo.New()
e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
  Skipper: func(c *echo.Context) bool {
    return strings.Contains(c.Path(), "metrics") // Change "metrics" for your own path
  },
}))
```

## Configuration

```go
type GzipConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// Gzip compression level.
	// Optional. Default value -1.
	Level int

	// Length threshold before gzip compression is applied.
	// Optional. Default value 0.
	//
	// Most of the time you will not need to change the default. Compressing
	// a short response might increase the transmitted data because of the
	// gzip format overhead. Compressing the response will also consume CPU
	// and time on the server and the client (for decompressing). Depending on
	// your use case such a threshold might be useful.
	//
	// See also:
	// https://webmasters.stackexchange.com/questions/31750/what-is-recommended-minimum-object-size-for-gzip-performance-benefits
	MinLength int
}
```

