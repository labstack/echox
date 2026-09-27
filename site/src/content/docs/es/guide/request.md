---
title: Request
description: Obtén datos de form, query y path desde un request, y valídalos.
sidebar:
  order: 7
---

Un handler lee datos del request a través de `echo.Context`. Echo puede obtener valores
individualmente por nombre, vincularlos a structs (consulta [Binding](/es/guide/binding/)) y
delegar la validación a un validator que registres.

## Obtener datos

### Datos de formulario

Obtén un campo de formulario por nombre con `Context#FormValue(name string)`:

```go
e.POST("/form", func(c *echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
})
```

Para tipos distintos de `string`, usa la función genérica `echo.FormValue[T]`:

```go
age, err := echo.FormValue[int](c, "age")
if err != nil {
	return err
}
```

Prueba con:

```sh
curl -X POST http://localhost:1323/form -d 'name=Joe&age=30'
```

Para vincular un tipo de datos personalizado, implementa la interfaz `echo.BindUnmarshaler`:

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

### Parámetros de query

Obtén un parámetro de query por nombre con `Context#QueryParam(name string)`:

```go
func(c *echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}
```

Para tipos distintos de `string`, usa la función genérica `echo.QueryParam[T]`:

```go
age, err := echo.QueryParam[int](c, "age")
if err != nil {
	return err
}
```

```sh
curl -X GET "http://localhost:1323?name=Joe&age=30"
```

### Parámetros de path

Obtén un parámetro de path registrado por nombre con `Context#Param(name string)`:

```go
e.GET("/users/:name", func(c *echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
})
```

Para tipos distintos de `string`, usa la función genérica `echo.PathParam[T]`:

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

### Binding de datos

Echo también puede vincular datos de request a structs y variables nativas de Go. Consulta
[Binding](/es/guide/binding/).

<a id="validate-data"></a>

## Validar datos

Echo no incluye validación de datos integrada. Puedes registrar un validator personalizado mediante
`Echo#Validator` y usar una biblioteca de terceros como
[go-playground/validator](https://github.com/go-playground/validator).

El ejemplo siguiente valida un struct vinculado:

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
{"message":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}
```
