---
description: Handling request
slug: /request
sidebar_position: 9
---

# Request

## Retrieve Data

### Form Data

Form data can be retrieved by name using `Context#FormValue(name string)`.

```go
// Handler
func(c echo.Context) error {
  name := c.FormValue("name")
  return c.String(http.StatusOK, name)
}
```

```sh
curl -X POST http://localhost:1323 -d 'name=Joe'
```

To bind a custom data type, you can implement `Echo#BindUnmarshaler` interface.

```go
type Timestamp time.Time

func (t *Timestamp) UnmarshalParam(src string) error {
  ts, err := time.Parse(time.RFC3339, src)
  *t = Timestamp(ts)
  return err
}
```

### Query Parameters

Query parameters can be retrieved by name using `Context#QueryParam(name string)`.

```go
// Handler
func(c echo.Context) error {
  name := c.QueryParam("name")
  return c.String(http.StatusOK, name)
})
```

```sh
curl \
  -X GET \
  http://localhost:1323\?name\=Joe
```

Similar to form data, custom data type can be bind using `Context#QueryParam(name string)`.

### Path Parameters

Registered path parameters can be retrieved by name using `Context#PathParam(name string) string`.

```go
e.GET("/users/:name", func(c echo.Context) error {
  name := c.PathParam("name")
  return c.String(http.StatusOK, name)
})
```

```sh
curl http://localhost:1323/users/Joe
```

### Binding Data

Also binding of request data to native Go structs and variables is supported.
See [Binding Data](./binding.md)

## Validate Data

Echo doesn't have built-in data validation capabilities, however, you can register
a custom validator using `Echo#Validator` and leverage third-party [libraries](https://github.com/avelino/awesome-go#validation).

Example below uses https://github.com/go-playground/validator framework for validation:

```go
package main

import (
  "net/http"

  "github.com/go-playground/validator"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
)

type (
  User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
  }

  CustomValidator struct {
    validator *validator.Validate
  }
)

func (cv *CustomValidator) Validate(i interface{}) error {
  if err := cv.validator.Struct(i); err != nil {
    // Optionally, you could return the error to give each route more control over the status code
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }
  return nil
}

func main() {
  e := echo.New()
  e.Validator = &CustomValidator{validator: validator.New()}
  e.POST("/users", func(c echo.Context) (err error) {
    u := new(User)
    if err = c.Bind(u); err != nil {
      return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    if err = c.Validate(u); err != nil {
      return err
    }
    return c.JSON(http.StatusOK, u)
  })
  e.Logger.Fatal(e.Start(":1323"))
}
```

```sh
curl -X POST http://localhost:1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@invalid-domain"}'
{"message":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}
```
