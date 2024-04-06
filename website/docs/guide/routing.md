---
description: Routing requests
slug: /routing
sidebar_position: 10
---

# Routing

Echo's router is based on [radix tree](http://en.wikipedia.org/wiki/Radix_tree), making
route lookup really fast. It leverages [sync pool](https://golang.org/pkg/sync/#Pool)
to reuse memory and achieve zero dynamic memory allocation with no GC overhead.

Routes can be registered by specifying HTTP method, path and a matching handler.
For example, code below registers a route for method `GET`, path `/hello` and a
handler which sends `Hello, World!` HTTP response.

```go
// Handler
func hello(c echo.Context) error {
  	return c.String(http.StatusOK, "Hello, World!")
}

// Route
e.GET("/hello", hello)
```

You can use `Echo.Any(path string, h Handler)` to register a handler for all HTTP methods.
If you want to register it for some methods use `Echo.Match(methods []string, path string, h Handler)`.

Echo defines handler function as `func(echo.Context) error` where `echo.Context` primarily
holds HTTP request and response interfaces.

## Match-any / wildcard

Matches zero or more characters in the path. For example, pattern `/users/*` will
match:

- `/users/`
- `/users/1`
- `/users/1/files/1`
- `/users/anything...`

:::caution

There can be only one effective match-any parameter in route. When route is added with multiple match-any
`/v1/*/images/*`. The router matches always the first `*` till the end of requst URL i.e. it works as `/v1/*`.

:::

## Path Matching Order

- Static
- Param
- Match any

*Example*

```go
e.GET("/users/:id", func(c echo.Context) error {
	return c.String(http.StatusOK, "/users/:id")
})

e.GET("/users/new", func(c echo.Context) error {
	return c.String(http.StatusOK, "/users/new")
})

e.GET("/users/1/files/*", func(c echo.Context) error {
	return c.String(http.StatusOK, "/users/1/files/*")
})
```

Above routes would resolve in the following order:

- `/users/new`
- `/users/:id`
- `/users/1/files/*`

:::tip

Routes can be written in any order.

:::

## Group

`Echo#Group(prefix string, m ...Middleware) *Group`

Routes with common prefix can be grouped to define a new sub-router with optional
middleware. In addition to specified middleware group also inherits parent middleware.
To add middleware later in the group you can use `Group.Use(m ...Middleware)`.
Groups can also be nested.

In the code below, we create an admin group which requires basic HTTP authentication
for routes `/admin/*`.

```go
g := e.Group("/admin")
g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}))
```

## Route Naming

Each of the registration methods returns a `Route` object, which can be used to name a route after the registration. For example:

```go
route := e.POST("/users", func(c echo.Context) error {
})
route.Name = "create-user"

// or using the inline syntax
e.GET("/users/:id", func(c echo.Context) error {
}).Name = "get-user"
```

Route names can be very useful when generating URIs from the templates, where you can't access the handler references or when you have multiple routes with the same handler.

## URI Building

`Echo#URI(handler HandlerFunc, params ...interface{})` can be used to generate URI for any handler with specified path parameters. It's helpful to centralize all your
URI patterns which ease in refactoring your application.


For example, `e.URI(h, 1)` will generate `/users/1` for the route registered below:

```go
// Handler
h := func(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// Route
e.GET("/users/:id", h)
```

In addition to `Echo#URI`, there is also `Echo#Reverse(name string, params ...interface{})` which is used to generate URIs based on the route name. For example a call to `Echo#Reverse("foobar", 1234)` would generate the URI `/users/1234` if the `foobar` route is registered like below:

```go
// Handler
h := func(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// Route
e.GET("/users/:id", h).Name = "foobar"
```

## List Routes

`Echo#Routes() []*Route` can be used to list all registered routes in the order
they are defined. Each route contains HTTP method, path and an associated handler.

*Example*

```go
// Handlers
func createUser(c echo.Context) error {
}

func findUser(c echo.Context) error {
}

func updateUser(c echo.Context) error {
}

func deleteUser(c echo.Context) error {
}

// Routes
e.POST("/users", createUser)
e.GET("/users", findUser)
e.PUT("/users", updateUser)
e.DELETE("/users", deleteUser)
```

Using the following code you can output all the routes to a JSON file:

```go
data, err := json.MarshalIndent(e.Routes(), "", "  ")
if err != nil {
	return err
}
os.WriteFile("routes.json", data, 0644)
```

`routes.json`

```js
[
  {
    "method": "POST",
    "path": "/users",
    "name": "main.createUser"
  },
  {
    "method": "GET",
    "path": "/users",
    "name": "main.findUser"
  },
  {
    "method": "PUT",
    "path": "/users",
    "name": "main.updateUser"
  },
  {
    "method": "DELETE",
    "path": "/users",
    "name": "main.deleteUser"
  }
]
```
