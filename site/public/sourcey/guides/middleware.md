# Middleware {#middleware}


Middleware is a function wrapping a handler. Echo splits it in two places:

- **`e.Pre`** runs **before routing**, so it can rewrite the path or method. Use it for path rewrites, trailing-slash handling, and method overrides.
- **`e.Use`** runs **after routing**, around the matched handler. Use it for logging, recovery, auth, CORS, and everything else.

```go
e := echo.New()

e.Pre(middleware.RemoveTrailingSlash())
e.Use(middleware.RequestLogger()) // log requests with log/slog
e.Use(middleware.Recover())       // recover from panics and return an error
e.Use(middleware.CORS())          // CORS headers
```

`e.Middlewares()` returns the registered `Use` middleware; `e.PreMiddlewares()` returns the `Pre` middleware. Both run in registration order.


## RequestLogger (v5 logging) {#requestlogger-v5-logging}


v5 drops the old logging interface and uses the standard library's `log/slog` instead. There is no `middleware.Logger` anymore — the logger is `middleware.RequestLogger`:

```go
e.Use(middleware.RequestLogger())
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogMethod:     true,
	LogURI:        true,
	LogStatus:     true,
	LogLatency:    true,
	LogRemoteIP:   true,
	LogError:      true,
	LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
		c.Logger().Info("request", "method", v.Method, "uri", v.URI, "status", v.Status)
		return nil
	},
}))
```

`RequestLoggerWithConfig` takes a `RequestLoggerConfig`; the no-arg `RequestLogger()` uses the defaults. Handlers reach the same `*slog.Logger` through `c.Logger()`, and `c.SetLogger` replaces it per request.


## Recover {#recover}


```go
e.Use(middleware.Recover())
```

`Recover` catches panics in handlers and middleware and converts them into an error passed to the `HTTPErrorHandler` (a 500 by default). Register it early — middleware registered before it cannot be recovered by it.


## Official middleware {#official-middleware}


Maintained by the Echo team, in the `middleware` package of this same module:

| Middleware | What it does |
| --- | --- |
| `RequestLogger` / `RequestLoggerWithConfig` | Structured request logging with `log/slog` |
| `Recover` / `RecoverWithConfig` | Panic recovery |
| `CORS` / `CORSWithConfig` | Cross-Origin Resource Sharing headers |
| `CSRF` / `CSRFWithConfig` | CSRF token protection |
| `Gzip` / `GzipWithConfig` | Response compression (honors `Accept-Encoding` quality values since v5.4.7) |
| `BodyLimit` / `BodyLimitWithConfig` | Request body size limit |
| `BasicAuth` / `BasicAuthWithConfig` | HTTP Basic authentication |
| `BodyDump` / `BodyDumpWithConfig` | Request/response body logging |
| `Compress` | Response compression (alternative to Gzip) |
| `Decompress` | Request body decompression |
| `ContextTimeout` | Per-request context timeout |

JWT authentication, Casbin/sessions/pprof, OpenTelemetry, and Prometheus live in separate official repositories: [labstack/echo-jwt](https://github.com/labstack/echo-jwt), [labstack/echo-contrib](https://github.com/labstack/echo-contrib), [labstack/echo-otel](https://github.com/labstack/echo-otel), and [labstack/echo-prometheus](https://github.com/labstack/echo-prometheus). Check each one for the Echo version (`v4` or `v5`) it supports before adopting it.


## Writing middleware {#writing-middleware}


A middleware takes the next handler and returns a handler:

```go
func onlyJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if c.Request().Header.Get("Accept") != "application/json" {
			return echo.NewHTTPError(http.StatusNotAcceptable)
		}
		return next(c)
	}
}

e.Use(onlyJSON)
```

`echo.HandlerFunc` is the handler type in v5: `func(c *echo.Context) error`. Returning an error sends it to the `HTTPErrorHandler` — see [Binding and error handling](https://echo.labstack.com/sourcey/guides/binding-error-handling.md).
