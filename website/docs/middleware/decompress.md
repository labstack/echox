---
description: Decompress middleware
---

# Decompress

Decompress middleware decompresses HTTP request if Content-Encoding header is set to gzip.

:::note

The body will be decompressed in memory and consume it for the lifetime of the request (and garbage collection). 

:::

## Usage

```go
e.Use(middleware.Decompress())
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
  Skipper: Skipper
}))
```

## Configuration

```go
type DecompressConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// GzipDecompressPool defines an interface to provide the sync.Pool used to create/store Gzip readers
	GzipDecompressPool Decompressor

	// MaxDecompressedSize limits the maximum size of decompressed request body in bytes.
	// If the decompressed body exceeds this limit, the middleware returns HTTP 413 error.
	// This prevents zip bomb attacks where small compressed payloads decompress to huge sizes.
	// Default: 100 * MB (104,857,600 bytes)
	// Set to -1 to disable limits (not recommended in production).
	MaxDecompressedSize int64
}
```

### Default Configuration

```go
DefaultDecompressConfig = DecompressConfig{
  Skipper: DefaultSkipper,
}
```
