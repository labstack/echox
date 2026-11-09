---
title: Request
description: Retrieve form, query, and path data from a request, and validate it.
sidebar:
  order: 7
---

A handler reads request data through the `echo.Context`. Echo can retrieve values
individually by name, bind them into structs (see [Binding](/guide/binding/)), and
hand off validation to a validator you register.

## Retrieve data

### Form data

Retrieve a form field by name with `Context#FormValue(name string)`:

```go
e.POST("/form", func(c *echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
})
```

For types other than `string`, use the generic `echo.FormValue[T]` function:

```go
age, err := echo.FormValue[int](c, "age")
if err != nil {
	return err
}
```

Test with:

```sh
curl -X POST http://localhost:1323/form -d 'name=Joe&age=30'
```

To bind a custom data type, implement the `echo.BindUnmarshaler` interface:

```go
type Timestamp time.Time

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	if err != nil {
		return err
	}
	*t = Timestamp(ts)
	return nil
}
```

### Query parameters

Retrieve a query parameter by name with `Context#QueryParam(name string)`:

```go
func(c *echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}
```

For types other than `string`, use the generic `echo.QueryParam[T]` function:

```go
age, err := echo.QueryParam[int](c, "age")
if err != nil {
	return err
}
```

```sh
curl -X GET "http://localhost:1323?name=Joe&age=30"
```

### Path parameters

Retrieve a registered path parameter by name with `Context#Param(name string)`:

```go
e.GET("/users/:name", func(c *echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
})
```

For types other than `string`, use the generic `echo.PathParam[T]` function:

```go
id, err := echo.PathParam[int](c, "id")
if err != nil {
	return err
}
```

```sh
curl http://localhost:1323/users/Joe
curl http://localhost:1323/users/123
```

### Binding data

Echo can also bind request data into native Go structs and variables. See
[Binding](/guide/binding/).

## Validate data

Echo has no built-in data validation. You can register a custom validator via
`Echo#Validator` and use a third-party library such as
[go-playground/validator](https://github.com/go-playground/validator).

The example below validates a bound struct:

```go
package main

import (
	"net/http"

	"github.com/go-playground/validator/v10" // go get github.com/go-playground/validator/v10
	"github.com/labstack/echo/v5"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally return the error to let each route control the status code.
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/users", func(c *echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if err := c.Validate(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

```sh
curl -X POST http://localhost:1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@invalid-domain"}'
{"message":"Bad Request"}
```
