---
description: Graceful shutdown for multiple instances recipe
---

# Multi-instance

This example will start two echo instances that can not only be shutdown gracefully
using ctrl-C, but also automatically tries a graceful shutdown of the other servers
when one (or more) instance(s) shutdown/fail.

This allows for e.g. running seperate echo instances for e.g. a private and public
api on different ports in a single binary.

Using [signal.NotifyContext()](https://pkg.go.dev/os/signal#NotifyContext)
and [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)

## Recipe

```go reference
https://github.com/labstack/echox/blob/master/cookbook/graceful-multi/server.go
```

:::note

Requires go1.16+

:::
