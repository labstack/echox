# Binding and error handling {#binding-and-error-handling}



## Binding request data {#binding-request-data}


`c.Bind` decodes the request body (JSON, XML, or form, by `Content-Type`) plus query and path parameters into a struct:

```go
type User struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func createUser(c *echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// use u ...
	return c.JSON(http.StatusCreated, u)
}
```

Binding reads, in order: path parameters, query parameters, then the body. A binding failure returns a `BindingError` — since v5.4.7, binding into a map whose value type is not `string`, `[]string`, or `any` returns an error instead of silently skipping the values, and named string types no longer panic.

The fine-grained helpers bind one source each: `echo.BindBody`, `echo.BindQueryParams`, `echo.BindPathValues`, and `echo.BindHeaders`. The `Binder` interface (implemented by `DefaultBinder`, the default on `e.Binder`) lets you replace the whole binding step; `e.Binder` is the field to set.


## Validation {#validation}


```go
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

e.Validator = &CustomValidator{validator: validator.New()}
```

Set `e.Validator` to any type implementing `Validate(any) error`. Then `c.Validate(i)` runs it; `c.ValidateCtx(ctx, i)` runs it with a request context. Without a validator registered, `Validate` returns `ErrValidatorNotRegistered`.


## Errors: HTTPError {#errors-httperror}


Handlers signal failures by returning an error. The `HTTPError` type carries a code and a message:

```go
return echo.ErrNotFound                       // 404, predefined
return echo.NewHTTPError(http.StatusBadRequest, "name is required")
return echo.NewHTTPError(http.StatusBadRequest, "name is required", map[string]string{"field": "name"})
```

Predefined errors cover the common codes: `ErrBadRequest`, `ErrUnauthorized`, `ErrForbidden`, `ErrNotFound`, `ErrMethodNotAllowed`, `ErrRequestTimeout`, `ErrStatusRequestEntityTooLarge`, `ErrUnsupportedMediaType`, `ErrTooManyRequests`, `ErrInternalServerError`, `ErrBadGateway`, `ErrServiceUnavailable`, plus `ErrValidatorNotRegistered`, `ErrRendererNotRegistered`, and `ErrCookieNotFound`. `HTTPError` has `Code` and `Message` fields; `NewBindingError` wraps binding failures.


## The error handler {#the-error-handler}


Every error a handler or middleware returns goes to one place: `e.HTTPErrorHandler`. The default is `DefaultHTTPErrorHandler`, which writes a JSON (or HTML, by `Accept`) error body with the code and message. Replace it to control the format:

```go
e.HTTPErrorHandler = func(err error, c *echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	// log it, then write your own body
	c.Logger().Error("request failed", "error", err)
	if !c.Response().Committed {
		c.JSON(code, map[string]any{"error": err.Error()})
	}
}
```

`ProblemDetailsHTTPErrorHandler` is a built-in alternative that writes RFC 9457 problem-details responses. Two details matter when writing a custom handler: check `c.Response().Committed` before writing (the response may already be partially sent), and remember the built-in handlers strip `Content-Encoding`, `Content-Length`, and `Content-Disposition` before writing the error body (since v5.4.7) — a custom handler that skips that step can send a stale `Content-Length` with the error.
