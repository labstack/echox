---
title: Graceful Shutdown
description: Drain in-flight requests before stopping the server on an interrupt signal.
sidebar:
  order: 8
---

A graceful shutdown lets in-flight requests finish before the process exits. The
simplest approach is to pass a cancellable context to `StartConfig.Start` and set
a `GracefulTimeout`. When the context is cancelled by an interrupt signal, Echo
stops accepting new connections and waits up to the timeout for active requests to
complete.

## Server

```go file=cookbook/graceful-shutdown/server.go#primary-server
```

## Using a custom HTTP server

If you manage the `http.Server` yourself, start it in a goroutine, wait on the
signal context, then call `Shutdown` with a timeout:

```go file=cookbook/graceful-shutdown/server.go#custom-server
```
