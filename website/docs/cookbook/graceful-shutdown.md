---
description: Graceful shutdown recipe
---

# Graceful Shutdown

This example demonstrates how to implement graceful shutdown using Go's `signal.NotifyContext` and Echo's `Shutdown` method to handle interrupt signals properly.

## Using [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext) and [Echo.Shutdown](https://pkg.go.dev/github.com/labstack/echo/v4#Echo.Shutdown)

```go reference
https://github.com/labstack/echox/blob/master/cookbook/graceful-shutdown/server.go
```

### Key Features

- Uses `signal.NotifyContext` to listen for OS interrupt signals (Ctrl+C)
- Starts the server in a goroutine to allow non-blocking execution
- Implements a 10-second timeout for graceful shutdown
- Properly handles server shutdown errors

:::note

Requires Go 1.16+ for `signal.NotifyContext`. For earlier versions, use `signal.Notify` with a channel.

:::
