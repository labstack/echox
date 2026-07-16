---
title: Routing
description: Match request URLs to handlers on Echo's zero-allocation radix tree.
sidebar:
  order: 3
---

Echo's optimized router matches request URLs to handlers using a radix tree with
**zero dynamic memory allocation** and smart route prioritization.

## Registering routes

Use the HTTP-method helpers on the `Echo` instance. Each takes a path pattern and a
`HandlerFunc` (`func(c *echo.Context) error`), with optional route-level middleware.

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` registers a handler for any HTTP method — including ones not in Echo's predefined
list — and `Match` for a specific set:

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

:::note
A method-specific route registered for the same path takes precedence over the `Any` route.
:::

## Match types

| Pattern            | Type     | Example match              |
| ------------------ | -------- | -------------------------- |
| `/users/profile`   | Static   | `/users/profile`           |
| `/users/:id`       | Param    | `/users/42`                |
| `/static/*`        | Wildcard | `/static/css/app.css`      |

:::note
Priority is **static → param → wildcard**, so `/users/profile` always wins over
`/users/:id`, which wins over `/users/*`.
:::

## Path parameters

Read named parameters from the context with `c.Param()` (or `c.ParamOr()` for a default):

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

The wildcard segment is available as the `*` parameter:

```go
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, c.Param("*"))
})
```

## Groups

Group routes that share a prefix and middleware with `e.Group()`:

```go
admin := e.Group("/admin", middleware.BasicAuth(authFn))
admin.GET("/metrics", metrics)   // -> /admin/metrics
admin.GET("/users", listUsers)   // -> /admin/users
```

Groups can be nested to compose larger route trees.
