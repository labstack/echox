---
title: Timeout
description: Apply a request timeout to handlers with the ContextTimeout middleware.
sidebar:
  order: 18
---

The [`ContextTimeout`](/middleware/context-timeout/) middleware sets a deadline on the
request's `context.Context`. When the deadline passes, the context is cancelled,
and handlers that watch `c.Request().Context().Done()` can return promptly instead
of running to completion.

In the example below the middleware imposes a 5-second timeout while the handler
would otherwise take 10 seconds, so the request returns a `408 Request Timeout`.

## Server

```go file=cookbook/timeout/server.go
```
