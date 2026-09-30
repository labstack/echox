---
title: Embed Resources
description: Serve static assets bundled into the binary with Go's embed package.
sidebar:
  order: 5
---

Go's `embed` package (Go 1.16+) lets you compile static assets directly into the
binary, so a single executable can ship with its frontend. This recipe serves the
embedded filesystem through Echo, with an optional live mode that reads from disk
during development.

## Server

```go file=cookbook/embed/server.go
```

:::tip
Run the binary with the `live` argument (`go run server.go live`) to serve assets
from the `app` directory on disk instead of the embedded copy, which is handy
during development.
:::
