# Routing and context {#routing-and-context}


Echo's router is a radix tree. Static routes match first, then parameter and wildcard routes.


## Basic routing {#basic-routing}


```go
e := echo.New()

e.GET("/", func(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
})
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
```

Every HTTP method has a helper: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`, `CONNECT`, `TRACE`, plus `Any` (all methods) and `Match` (an explicit method list). `Add` registers a route for any method string, and `AddRoute` registers a pre-built `Route`.


## Parameters and wildcards {#parameters-and-wildcards}


```go
// GET /users/42 — c.Param("id") == "42"
e.GET("/users/:id", func(c *echo.Context) error {
	return c.String(http.StatusOK, "user "+c.Param("id"))
})

// GET /files/docs/guide.pdf — c.Param("*") == "docs/guide.pdf"
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, "file "+c.Param("*"))
})
```

`c.Param` reads a captured value by name (`"*"` for the wildcard). `c.ParamOr("id", "default")` returns a fallback when the parameter is absent, and `c.Path()` returns the registered route pattern (for example `/users/:id`), not the concrete request path. `c.PathValues()` returns all captured values at once.


## Groups {#groups}


Groups share a prefix and middleware:

```go
api := e.Group("/api")
api.Use(middleware.JWTWithConfig(jwtConfig)) // runs for every route in the group

api.GET("/users", listUsers)      // GET /api/users
api.POST("/users", createUser)    // POST /api/users

admin := api.Group("/admin")      // GET /api/admin/...
admin.GET("/stats", stats)
```

`Group` supports the same route helpers as `Echo` (`GET`, `POST`, `Group` for nesting, `Static`, `File`, `Use`). Middleware registered on a group runs only for that group's routes.


## Static files {#static-files}


```go
e.Static("/assets", "public")      // serves ./public at /assets/*
e.StaticFS("/assets", fs)          // serves an fs.FS at /assets/*
e.File("/favicon.ico", "favicon.ico")
```

Since v5.4.7, `Static` with a prefix without a trailing slash registers `prefix/*` plus `prefix` for the redirect to `prefix/` — so `/assetsadmin/private.txt` can no longer escape a guard on `/assets/admin` (GHSA-76qj-xf66-r429). Prefer `Group.Static` over `middleware.Static` attached with `Group.Use` when files sit behind route guards: the middleware variant is not protected by guards on routes inside the group (GHSA-5hp7-4fcm-57f9).


## The Router interface (v5) {#the-router-interface-v5}


v5 adds a `Router` interface so routing implementations can be swapped. `e.Router` returns the current router, and `RouterConfig` configures it. The default is `DefaultRouter`, created with `NewDefaultRouter` or `NewConcurrentRouter`. `e.RouteNotFound` registers a handler for requests that match no route; without it, unmatched requests get the default 404 from `ErrNotFound`.

`c.RouteInfo()` returns the matched route's info (`RouteInfo`): method, path pattern, and handler name. It is the reliable way to log or meter by route pattern instead of by concrete path.


## Context pool {#context-pool}


`Echo` reuses `Context` values through a pool: `AcquireContext` takes one, `ReleaseContext` returns it. Handlers never need to call these — the framework does it per request — but custom servers built on `ServeHTTP` must not hold a `*Context` after the handler returns.
