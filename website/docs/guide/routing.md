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

## HTTP Methods

Echo supports all standard HTTP methods. Each method follows the signature:
`e.METHOD(path string, h HandlerFunc, middleware ...Middleware) *Route`

### Available HTTP Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `GET` | `e.GET(path, handler, ...middleware)` | Retrieve data |
| `POST` | `e.POST(path, handler, ...middleware)` | Create new resource |
| `PUT` | `e.PUT(path, handler, ...middleware)` | Update entire resource |
| `PATCH` | `e.PATCH(path, handler, ...middleware)` | Partial update of resource |
| `DELETE` | `e.DELETE(path, handler, ...middleware)` | Remove resource |
| `HEAD` | `e.HEAD(path, handler, ...middleware)` | Get headers only |
| `OPTIONS` | `e.OPTIONS(path, handler, ...middleware)` | Get allowed methods |

### Method Parameters

- **`path`** (string): The route pattern (e.g., `/users/:id`, `/api/*`)
- **`handler`** (HandlerFunc): Function that handles the request
- **`middleware`** (...Middleware): Optional middleware functions

### Handler Function Signature

Echo defines handler functions as `func(echo.Context) error` where:

- **`echo.Context`**: Provides access to:
  - Request and response objects
  - Path parameters (`c.Param("id")`)
  - Query parameters (`c.QueryParam("name")`)
  - Form data (`c.FormValue("field")`)
  - JSON binding (`c.Bind(&struct{})`)
  - Response helpers (`c.JSON()`, `c.String()`, etc.)

### Example Usage

```go
// GET request
e.GET("/users/:id", getUser)

// POST request with JSON body
e.POST("/users", createUser)

// PUT request for updates
e.PUT("/users/:id", updateUser)

// DELETE request
e.DELETE("/users/:id", deleteUser)

// PATCH request for partial updates
e.PATCH("/users/:id", patchUser)

// HEAD request (same as GET but without body)
e.HEAD("/users/:id", getUser)

// OPTIONS request for CORS
e.OPTIONS("/users", handleOptions)
```

## Route Registration Methods

You can use `Echo.Any(path string, h Handler)` to register a handler for all HTTP methods.
If you want to register it for some methods use `Echo.Match(methods []string, path string, h Handler)`.

### Echo.Any()

Registers a handler for all HTTP methods on a given path:

```go
e.Any("/api/*", func(c echo.Context) error {
    return c.String(http.StatusOK, "Handles all methods")
})
```

### Echo.Match()

Registers a handler for specific HTTP methods:

```go
// Handle both GET and POST
e.Match([]string{"GET", "POST"}, "/users", userHandler)

// Handle multiple methods
e.Match([]string{"PUT", "PATCH"}, "/users/:id", updateHandler)
```

## Route Parameters

Echo supports various types of route parameters for flexible URL patterns.

### Path Parameters

Path parameters are defined using `:paramName` syntax and can be accessed in handlers using `c.Param("paramName")`.

```go
// Single parameter
e.GET("/users/:id", func(c echo.Context) error {
    id := c.Param("id")
    return c.String(http.StatusOK, "User ID: " + id)
})

// Multiple parameters
e.GET("/users/:id/posts/:postId", func(c echo.Context) error {
    userId := c.Param("id")
    postId := c.Param("postId")
    return c.String(http.StatusOK, "User: " + userId + ", Post: " + postId)
})

// Optional parameters (using query parameters instead)
e.GET("/users", func(c echo.Context) error {
    id := c.QueryParam("id") // Optional
    return c.String(http.StatusOK, "User ID: " + id)
})
```

### Query Parameters

Query parameters are accessed using `c.QueryParam("name")` and are optional by nature:

```go
e.GET("/search", func(c echo.Context) error {
    query := c.QueryParam("q")
    limit := c.QueryParam("limit")
    return c.String(http.StatusOK, "Search: " + query + ", Limit: " + limit)
})
```

**Example URLs:**

- `/search?q=golang&limit=10`
- `/users/123?include=posts&limit=5`

### Form Parameters

For POST/PUT requests with form data, use `c.FormValue("field")`:

```go
e.POST("/users", func(c echo.Context) error {
    name := c.FormValue("name")
    email := c.FormValue("email")
    return c.String(http.StatusOK, "Name: " + name + ", Email: " + email)
})
```

## Match-any / wildcard

Matches zero or more characters in the path. For example, pattern `/users/*` will
match:

- `/users/`
- `/users/1`
- `/users/1/files/1`
- `/users/anything...`

:::caution

There can be only one effective match-any parameter in route. When route is added with multiple match-any
`/v1/*/images/*`. The router matches always the first `*` till the end of request URL i.e. it works as `/v1/*`.

:::

## Path Matching Order

- Static
- Param
- Match any

### Example

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

### Basic Group Usage

```go
// Create a group with prefix
api := e.Group("/api")

// Add routes to the group
api.GET("/users", getUsers)
api.POST("/users", createUser)
api.GET("/users/:id", getUser)
```

### Group with Middleware

```go
// Admin group with authentication
admin := e.Group("/admin")
admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
    if username == "joe" && password == "secret" {
        return true, nil
    }
    return false, nil
}))

// All routes under /admin/* will require authentication
admin.GET("/dashboard", dashboard)
admin.GET("/users", adminUsers)
admin.POST("/users", createUser)
```

### Advanced Group Examples

```go
// API v1 group with logging and CORS
apiV1 := e.Group("/api/v1")
apiV1.Use(middleware.Logger())
apiV1.Use(middleware.CORS())

apiV1.GET("/users", getUsers)
apiV1.POST("/users", createUser)

// API v2 group with different middleware
apiV2 := e.Group("/api/v2")
apiV2.Use(middleware.Logger())
apiV2.Use(middleware.CORS())
apiV2.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

apiV2.GET("/users", getUsersV2)
apiV2.POST("/users", createUserV2)

// Nested groups
admin := e.Group("/admin")
admin.Use(middleware.BasicAuth(...))

// Admin users sub-group
adminUsers := admin.Group("/users")
adminUsers.Use(middleware.RoleBasedAuth("admin"))
adminUsers.GET("/", listUsers)
adminUsers.POST("/", createUser)
```

### Group Middleware Management

```go
// Create group without initial middleware
api := e.Group("/api")

// Add middleware later
api.Use(middleware.Logger())
api.Use(middleware.Recover())

// Add routes
api.GET("/health", healthCheck)
api.GET("/users", getUsers)

// Add more middleware to specific routes
api.GET("/sensitive", sensitiveHandler, middleware.AuthRequired())
```

### Practical Example: RESTful API Structure

```go
func setupRoutes(e *echo.Echo) {
    // Public routes
    e.GET("/", homeHandler)
    e.GET("/login", loginPage)
    e.POST("/login", loginHandler)

    // API group with common middleware
    api := e.Group("/api")
    api.Use(middleware.Logger())
    api.Use(middleware.Recover())
    api.Use(middleware.CORS())

    // Public API routes
    api.GET("/health", healthCheck)
    api.POST("/register", registerUser)

    // Protected API routes
    protected := api.Group("")
    protected.Use(middleware.JWT([]byte("secret")))

    protected.GET("/users", getUsers)
    protected.POST("/users", createUser)
    protected.GET("/users/:id", getUser)
    protected.PUT("/users/:id", updateUser)
    protected.DELETE("/users/:id", deleteUser)

    // Admin API routes
    admin := protected.Group("/admin")
    admin.Use(middleware.RoleBasedAuth("admin"))

    admin.GET("/stats", getStats)
    admin.POST("/users/:id/ban", banUser)
}
```

### Group Benefits

1. **Middleware Inheritance**: Groups inherit parent middleware
2. **URL Organization**: Logical grouping of related routes
3. **Middleware Reuse**: Apply common middleware to multiple routes
4. **Nested Structure**: Create hierarchical route structures
5. **Easy Maintenance**: Modify group middleware affects all routes

## Route Naming

Each of the registration methods returns a `Route` object, which can be used to name a route after the registration. Route names are useful for generating URIs from templates, avoiding hardcoded URLs, and when you have multiple routes with the same handler.

### Basic Route Naming

```go
// Method 1: Assign name after registration
route := e.POST("/users", createUser)
route.Name = "create-user"

// Method 2: Inline syntax
e.GET("/users/:id", getUser).Name = "get-user"
e.PUT("/users/:id", updateUser).Name = "update-user"
e.DELETE("/users/:id", deleteUser).Name = "delete-user"
```

### Advanced Route Naming

```go
// Multiple parameters
e.GET("/users/:id/posts/:postId", getPost).Name = "get-user-post"

// Nested routes
e.GET("/api/v1/users/:id", getUser).Name = "api-get-user"
e.GET("/api/v2/users/:id", getUserV2).Name = "api-get-user-v2"

// Group routes with names
admin := e.Group("/admin")
admin.GET("/dashboard", dashboard).Name = "admin-dashboard"
admin.GET("/users", adminUsers).Name = "admin-users"
```

## URI Building

Echo provides two methods for generating URIs: `Echo.URI()` and `Echo.Reverse()`.

### Echo.URI() - Handler-based URI Generation

`Echo.URI(handler HandlerFunc, params ...interface{})` generates URIs based on handler functions:

```go
// Define handlers
func getUser(c echo.Context) error {
    return c.String(http.StatusOK, "User")
}

func getUserPost(c echo.Context) error {
    return c.String(http.StatusOK, "User Post")
}

// Register routes
e.GET("/users/:id", getUser)
e.GET("/users/:id/posts/:postId", getUserPost)

// Generate URIs
userURI := e.URI(getUser, 123)                    // "/users/123"
postURI := e.URI(getUserPost, 123, 456)          // "/users/123/posts/456"
```

### Echo.Reverse() - Name-based URI Generation

`Echo.Reverse(name string, params ...interface{})` generates URIs based on route names:

```go
// Register named routes
e.GET("/users/:id", getUser).Name = "get-user"
e.GET("/users/:id/posts/:postId", getUserPost).Name = "get-user-post"
e.GET("/api/v1/users/:id", getUser).Name = "api-user"

// Generate URIs using names
userURI := e.Reverse("get-user", 123)             // "/users/123"
postURI := e.Reverse("get-user-post", 123, 456)  // "/users/123/posts/456"
apiURI := e.Reverse("api-user", 123)             // "/api/v1/users/123"
```

### Practical Examples

```go
// In templates or handlers
func generateUserLinks(userID int) map[string]string {
    return map[string]string{
        "profile": e.Reverse("get-user", userID),
        "edit":    e.Reverse("update-user", userID),
        "delete":  e.Reverse("delete-user", userID),
    }
}

// In middleware for redirects
func redirectToUser(c echo.Context) error {
    userID := c.Param("id")
    return c.Redirect(http.StatusFound, e.Reverse("get-user", userID))
}
```

### Benefits of Route Naming

1. **URL Centralization**: Change URLs in one place
2. **Template Integration**: Generate links in HTML templates
3. **Refactoring Safety**: Rename handlers without breaking links
4. **API Documentation**: Generate API documentation with correct URLs

## List Routes

`Echo#Routes() []*Route` can be used to list all registered routes in the order
they are defined. Each route contains HTTP method, path and an associated handler.

### Example: Listing Routes

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
