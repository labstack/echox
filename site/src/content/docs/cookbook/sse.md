---
title: Server-Sent Events (SSE)
description: Stream server-sent events from an Echo handler, either per connection or broadcast to many clients.
sidebar:
  order: 14
---

[Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)
can be used in several ways. The first example below is per-connection, per-handler
SSE. For more complex broadcasting logic, see the second example using
[r3labs/sse](https://github.com/r3labs/sse).

:::caution
SSE connections are long-lived, so the server's write timeout must be disabled.
Both examples set `s.WriteTimeout = 0` via `BeforeServeFunc`.
:::

## Using SSE

### Server

The handler writes the SSE headers, then emits an event every second until the
client disconnects. `http.NewResponseController(w).Flush()` pushes each event to
the client immediately.

```go file=cookbook/sse/simple/server.go
```

### Event structure and Marshal method

```go file=cookbook/sse/simple/serversentevent.go
```

### HTML serving SSE

```html file=cookbook/sse/simple/index.html
```

## Broadcasting with r3labs/sse

When you need to broadcast a single stream of events to many subscribers, the
[r3labs/sse](https://github.com/r3labs/sse) library handles stream and subscriber
management for you.

### Server

```go file=cookbook/sse/broadcast/server.go
```

### HTML serving SSE

```html file=cookbook/sse/broadcast/index.html
```
