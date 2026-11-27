# Quickstart {#quickstart}


Echo v5 requires **Go 1.25 or higher**. If you need to install or upgrade Go, visit the [official Go download page](https://go.dev/dl/).


## Install {#install}


Create a new directory for your project, initialize Go modules, then install Echo:

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```


## Hello, World! {#hello-world}


```go
package main

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger()) // log requests with log/slog
	e.Use(middleware.Recover())       // recover from panics and return an error

	e.GET("/", hello)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

Run it and send a request:

```console
$ go run .
$ curl http://localhost:8080/
Hello, World!
```

`e.Start` shuts the server down gracefully on <kbd>Ctrl</kbd>+<kbd>C</kbd>.


## Core concepts {#core-concepts}


This simple server introduces the core concepts of Echo:

- **Echo instance** — `echo.New()` creates the application. `echo.NewWithConfig` takes a `Config` to tune it (see `Config` in the [Go API](https://echo.labstack.com/sourcey/api.md)).
- **Routes** — `e.GET`, `e.POST`, `e.PUT`, `e.DELETE`, and friends register handlers for a method + path. Paths support parameters (`:id`) and wildcards (`*`) — see [Routing and context](https://echo.labstack.com/sourcey/guides/routing-context.md).
- **Handlers** — in v5 a handler takes `*echo.Context` (a pointer to the concrete `Context` struct) and returns an `error`. `c.String`, `c.JSON`, `c.HTML`, and `c.File` write the response — see [Context and responses](https://echo.labstack.com/sourcey/guides/context-responses.md).
- **Middleware** — `e.Use` registers handlers that run after routing; `e.Pre` registers handlers that run before routing — see [Middleware](https://echo.labstack.com/sourcey/guides/middleware.md).
- **Start** — `e.Start(":8080")` starts the server. It returns only errors that are not `http.ErrServerClosed`, so a graceful shutdown on Ctrl+C is not an error.


## Next steps {#next-steps}


- [Routing and context](https://echo.labstack.com/sourcey/guides/routing-context.md) — parameters, wildcards, groups, route info.
- [Middleware](https://echo.labstack.com/sourcey/guides/middleware.md) — the `Pre`/`Use` split, RequestLogger, Recover.
- [Binding and error handling](https://echo.labstack.com/sourcey/guides/binding-error-handling.md) — binding request data, custom error responses.
