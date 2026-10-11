---
title: Error Handling
description: Centralized HTTP error handling by returning errors from handlers and middleware.
sidebar:
  order: 6
---

Echo advocates **centralized** error handling: handlers and middleware return an
`error`, and a single error handler turns it into an HTTP response. This keeps logging
and response formatting in one place.

Return a plain `error` or an `*echo.HTTPError`:

```go
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if !authenticated(c) {
			// invalid credentials → abort with 401
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}
		return next(c)
	}
})
```

`echo.NewHTTPError(code)` without a message uses the status text (e.g. `"Unauthorized"`).
Echo also ships sentinel errors like `echo.ErrBadRequest`, `echo.ErrNotFound`, and
`echo.ErrUnauthorized`.

## Default error handler

Echo's default handler responds in JSON:

```json
{ "message": "error connecting to redis" }
```

A plain `error` becomes `500 Internal Server Error` (the original message is included
when running with errors exposed). An `*HTTPError` uses its status code and message.

## Custom error handler

Set your own via `e.HTTPErrorHandler` — useful for error pages, notifications, or
sending errors to a centralized system.

Check whether the response was already sent with `echo.UnwrapResponse()`, and find a
status code in the error chain via `echo.HTTPStatusCoder`:

```go
func customHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return // already sent by a handler/middleware
		}
	}

	code := http.StatusInternalServerError
	if tmp := echo.StatusCode(err); tmp != 0 {
		code = tmp
	}

	var cErr error
	if c.Request().Method == http.MethodHead {
		cErr = c.NoContent(code)
	} else {
		cErr = c.File(fmt.Sprintf("%d.html", code)) // e.g. 404.html, 500.html
	}
	if cErr != nil {
		c.Logger().Error("failed to send error page", "error", errors.Join(err, cErr))
	}
}

e.HTTPErrorHandler = customHTTPErrorHandler
```

:::tip
Instead of (or in addition to) the logger, forward errors to an external service like
Sentry, Elasticsearch, or Splunk from the central handler.
:::
