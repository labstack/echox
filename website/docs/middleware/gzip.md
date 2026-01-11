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
GzipConfig struct {
  // Skipper defines a function to skip middleware.
  Skipper Skipper

  // Gzip compression level.
  // Optional. Default value -1.
  Level int `json:"level"`
}
```

### Default Configuration

```go
DefaultGzipConfig = GzipConfig{
  Skipper: DefaultSkipper,
  Level:   -1,
}
```
