---
title: Request
description: Recupere dados de formulário, query e path de um request, e valide-os.
sidebar:
  order: 7
---

Um handler lê dados do request por meio do `echo.Context`. Echo pode recuperar valores
individualmente por nome, fazer binding deles em structs (veja [Binding](/pt-br/guide/binding/)) e
delegar a validação para um validador que você registra.

## Recuperar dados

### Dados de formulário

Recupere um campo de formulário por nome com `Context#FormValue(name string)`:

```go
e.POST("/form", func(c *echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
})
```

Para tipos diferentes de `string`, use a função genérica `echo.FormValue[T]`:

```go
age, err := echo.FormValue[int](c, "age")
if err != nil {
	return err
}
```

Teste com:

```sh
curl -X POST http://localhost:1323/form -d 'name=Joe&age=30'
```

Para fazer binding de um tipo de dado customizado, implemente a interface `echo.BindUnmarshaler`:

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

### Parâmetros de query

Recupere um parâmetro de query por nome com `Context#QueryParam(name string)`:

```go
func(c *echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}
```

Para tipos diferentes de `string`, use a função genérica `echo.QueryParam[T]`:

```go
age, err := echo.QueryParam[int](c, "age")
if err != nil {
	return err
}
```

```sh
curl -X GET "http://localhost:1323?name=Joe&age=30"
```

### Parâmetros de caminho

Recupere um parâmetro de caminho registrado por nome com `Context#Param(name string)`:

```go
e.GET("/users/:name", func(c *echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
})
```

Para tipos diferentes de `string`, use a função genérica `echo.PathParam[T]`:

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

### Binding de dados

Echo também pode fazer binding de dados do request em structs e variáveis nativas do Go. Veja
[Binding](/pt-br/guide/binding/).

<a id="validate-data"></a>

## Validar dados

Echo não tem validação de dados embutida. Você pode registrar um validador customizado via
`Echo#Validator` e usar uma biblioteca de terceiros como
[go-playground/validator](https://github.com/go-playground/validator).

O exemplo abaixo valida uma struct vinculada:

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
