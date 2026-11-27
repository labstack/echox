# Echo — High performance, minimalist Go web framework

# Echo {#echo}


Echo is a **high performance, extensible, minimalist Go web framework** built on Go's standard `net/http`. It interoperates with the standard library through `echo.WrapHandler` and `echo.WrapMiddleware`, and adds the parts the standard library leaves to you: a fast radix-tree router, request binding with pluggable validation, official and third-party middleware, and centralized error handling. Echo is actively maintained, and `v5` is the current release line.

Echo requires **Go 1.25 or higher** (module `go 1.25.0` at the documented commit).


## What you get {#what-you-get}


- **Radix-tree routing** — static, parameter and wildcard routes; static routes match first — see [Routing and context](https://echo.labstack.com/sourcey/guides/routing-context.md).
- **Middleware, before and after routing** — `e.Pre` runs before routing, `e.Use` after it — see [Middleware](https://echo.labstack.com/sourcey/guides/middleware.md).
- **Data binding** — JSON, XML, form, query and path parameters, with pluggable validation — see [Binding and error handling](https://echo.labstack.com/sourcey/guides/binding-error-handling.md).
- **Centralized error handling** — one `HTTPErrorHandler` for every error a handler or middleware returns — see [Binding and error handling](https://echo.labstack.com/sourcey/guides/binding-error-handling.md).
- **A concrete `Context`** — in v5 `Context` is a struct, not an interface, and handlers take `*echo.Context` — see [Context and responses](https://echo.labstack.com/sourcey/guides/context-responses.md).
- **Migration from v4** — the breaking changes, and the v4 support dates — see [Migration v4 → v5](https://echo.labstack.com/sourcey/guides/migration-v4-v5.md).


## Where to start {#where-to-start}


- New to Echo? Read [Quickstart](https://echo.labstack.com/sourcey/quickstart.md) — install, your first "Hello, World!" server, and the core concepts.
- Already routing? [Routing and context](https://echo.labstack.com/sourcey/guides/routing-context.md) and [Middleware](https://echo.labstack.com/sourcey/guides/middleware.md) cover the daily API.
- Shipping to production? [Binding and error handling](https://echo.labstack.com/sourcey/guides/binding-error-handling.md) and [Context and responses](https://echo.labstack.com/sourcey/guides/context-responses.md) cover the sharp edges (binding failures, error responses, response helpers).


## API reference {#api-reference}


The [Go API](https://echo.labstack.com/sourcey/api.md) tab is the full reference for the `github.com/labstack/echo/v5` package, generated from the package source: 38 types, 55 package-level functions, 232 methods, 78 constants, and 21 variables — including `Echo`, `Context`, `Config`, `Router`, `Group`, `RouteInfo`, `DefaultBinder`, `HTTPError`, and the error helpers.


## Documented version {#documented-version}


This site documents **labstack/echo @ 59aacdd** (commit `59aacddc4a13904b8a80a76f1ec838c495b5feca`, 20 Nov 2026, release v5.4.7), MIT licensed. Source links in the API reference point at that exact commit on GitHub.
