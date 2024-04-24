---
description: SSE recipe
---

# Server-Sent-Events (SSE)

[Server-sent events](
https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format) can be
used in different ways. This example here is per connection - per handler SSE. If your requirements need more complex
broadcasting logic see https://github.com/r3labs/sse library.

## Using SSE

### Server

```go reference
https://github.com/labstack/echox/blob/master/cookbook/sse/simple/server.go
```

### Event structure and Marshal method

```go reference
https://github.com/labstack/echox/blob/master/cookbook/sse/simple/serversentevent.go
```

### HTML serving SSE

```html reference
https://github.com/labstack/echox/blob/master/cookbook/sse/simple/index.html
```

## Using 3rd party library [r3labs/sse](https://github.com/r3labs/sse) to broadcast events

### Server

```go reference
https://github.com/labstack/echox/blob/master/cookbook/sse/broadcast/server.go
```

### HTML serving SSE

```html reference
https://github.com/labstack/echox/blob/master/cookbook/sse/broadcast/index.html
```