# Migration v4 → v5 {#migration-v4--v5}


Echo `v5` (current release line, first released 2026-01-18; this site documents v5.4.7) is a maintenance release with major breaking changes. The full list is in [API_CHANGES_V5.md](https://github.com/labstack/echo/blob/master/API_CHANGES_V5.md) in the repository. The changes that touch nearly every application are below.


## Context: interface → concrete struct {#context-interface--concrete-struct}


The single biggest change. In v4 `Context` was an interface; in v5 it is a struct, and **every handler must take a pointer**:

```go
// v4
func handler(c echo.Context) error {
	return c.JSON(200, map[string]string{"hello": "world"})
}

// v5
func handler(c *echo.Context) error {
	return c.JSON(200, map[string]string{"hello": "world"})
}
```

Middleware signatures change the same way: `echo.HandlerFunc` is now `func(c *echo.Context) error`, and `echo.MiddlewareFunc` wraps that. Custom types that implemented the v4 `Context` interface no longer compile against v5 — use the struct directly, or wrap it.

The payoff: because `Context` is now a struct, the maintainers can add methods in minor versions without breaking your code.


## Logging: the old interface is gone, slog is in {#logging-the-old-interface-is-gone-slog-is-in}


v5 drops the old logging interface and uses the standard library's `log/slog`. The `middleware.Logger` middleware is gone; use `middleware.RequestLogger` (or `RequestLoggerWithConfig`) instead, and reach the logger from handlers with `c.Logger()`. See [Middleware](https://echo.labstack.com/sourcey/guides/middleware.md).


## Routing: a Router interface {#routing-a-router-interface}


v5 adds a `Router` interface (`e.Router`, `RouterConfig`, `DefaultRouter`) so routing implementations can be swapped. The default radix-tree behavior is unchanged: static routes match first, then parameters and wildcards. `c.RouteInfo()` is new — it returns the matched route's method, path pattern, and handler name.


## Generic helpers take *Context {#generic-helpers-take-context}


The generic helpers (`ContextGet`, `ContextGetOr`) now take `*Context`, and the form helpers were renamed to the `FormValue*` family (`FormValue`, `FormValueOr`, `FormValues`). Update call sites that used the old names.


## New test helpers: echotest {#new-test-helpers-echotest}


v5 adds the `echotest` package with test helpers for handlers and middleware. It is a separate package in the module, not covered by this reference (root package only).


## Version support {#version-support}


| Version | Bug fixes | Security fixes |
| --- | --- | --- |
| 5.x | Yes | Yes |
| 4.x, latest release | Until 2026-12-31 | Until 2027-12-31 |
| 4.x, older releases | No | No — upgrade to the latest 4.x |

Routine v4 bug fixes end on 2026-12-31. Free, public, security-only fixes for v4 continue until 2027-12-31, on patch releases of the latest v4 version. Echo v4 requires Go 1.26 or later. Start new projects on v5, and plan v4 upgrades with `API_CHANGES_V5.md`.


## Checklist {#checklist}


1. Change every handler and middleware signature from `echo.Context` to `*echo.Context`.
2. Replace `middleware.Logger` with `middleware.RequestLogger`; move logger configuration to `RequestLoggerConfig` / `log/slog`.
3. Update the module path: `go get github.com/labstack/echo/v5`.
4. Replace removed/renamed helpers (form helpers → `FormValue*`; check `API_CHANGES_V5.md` for the full list).
5. Re-run your tests — the `echotest` package can replace hand-rolled test harnesses.
