---
title: Subdomain
description: Route requests to different Echo instances per host using a virtual host handler.
sidebar:
  order: 17
---

This recipe routes requests to separate `Echo` instances based on the request
host, so each subdomain has its own routes and middleware. The instances are
combined with `echo.NewVirtualHostHandler`, which dispatches by host name.

## Server

```go file=cookbook/subdomain/server.go
```
