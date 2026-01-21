---
description: Error handling
slug: /error-handling
sidebar_position: 6
---

# Error Handling

Echo advocates for centralized HTTP error handling by returning error from middleware
and handlers. Centralized error handler allows us to log errors to external services
from a unified location and send a customized HTTP response to the client.

You can return a standard `error` or `echo.*HTTPError`.

For example, when basic auth middleware finds invalid credentials it returns
401 - Unauthorized error, aborting the current HTTP request.

```go
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c *echo.Context) error {
    // Extract the credentials from HTTP request header and perform a security
    // check

    // For invalid credentials
    return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")

    // For valid credentials call next
    // return next(c)
  }
})
```

You can also use `echo.NewHTTPError()` without a message, in that case status text is used
as an error message. For example, "Unauthorized".

## Default HTTP Error Handler

Echo provides a default HTTP error handler which sends error in a JSON format.

```js
{
  "message": "error connecting to redis"
}
```

For a standard `error`, response is sent as `500 - Internal Server Error`; however,
if you are running in a debug mode, the original error message is sent. If error
is `*HTTPError`, response is sent with the provided status code and message.
If logging is on, the error message is also logged.

## Custom HTTP Error Handler

Custom HTTP error handler can be set via `e.HTTPErrorHandler`

For most cases default error HTTP handler should be sufficient; however, a custom HTTP
error handler can come handy if you want to capture different type of errors and
take action accordingly e.g. send notification email or log error to a centralized
system. You can also send customized response to the client e.g. error page or
just a JSON response.

To check if the response has been already sent to the client ("commited") you can use `echo.UnwrapResponse()`,
```go
if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
    if resp.Committed {
        return // response has been already sent to the client by handler or some middleware
    }
}
```

To find error in an error chain that implements `echo.HTTPStatusCoder`,
```go
code := http.StatusInternalServerError
var sc echo.HTTPStatusCoder
if errors.As(err, &sc) { // find error in an error chain that implements HTTPStatusCoder
    if tmp := sc.StatusCode(); tmp != 0 {
        code = tmp
    }
}
```

### Error Pages

The following custom HTTP error handler shows how to display error pages for different
type of errors and logs the error. The name of the error page should be like `<CODE>.html` e.g. `500.html`. You can look into this project
https://github.com/AndiDittrich/HttpErrorPages for pre-built error pages.

```go
func customHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return // response has been already sent to the client by handler or some middleware
		}
	}

	code := http.StatusInternalServerError
	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) { // find error in an error chain that implements HTTPStatusCoder
		if tmp := sc.StatusCode(); tmp != 0 {
			code = tmp
		}
	}

	var cErr error
	if c.Request().Method == http.MethodHead {
		cErr = c.NoContent(code)
	} else {
		errorPage := fmt.Sprintf("%d.html", code)
		cErr = c.File(errorPage)
	}
	if cErr != nil {
		c.Logger().Error("failed to send error page to client", "error", errors.Join(err, cErr))
	}
}

e.HTTPErrorHandler = customHTTPErrorHandler
```

:::tip

Instead of writing logs to the logger, you can also write them to an external service like Elasticsearch or Splunk.

:::
