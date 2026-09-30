---
title: CORS
description: Enable Cross-Origin Resource Sharing with an allow list or a custom origin function.
sidebar:
  order: 4
---

The [CORS middleware](/middleware/cors/) controls which origins may access your API.
You can pass a fixed list of allowed origins, or supply a function that decides
per request.

## Allow list of origins

Pass the allowed origins directly to `middleware.CORS`.

```go file=cookbook/cors/origin-list/server.go
```

## Custom origin function

For dynamic policies, use `CORSWithConfig` with `UnsafeAllowOriginFunc`. The
function receives the request context and origin and returns the origin to echo
back, whether the request is allowed, and an optional error.

```go file=cookbook/cors/origin-func/server.go
```
