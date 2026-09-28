---
title: Testing
description: Prueba handlers y middleware con httptest y los helpers de echotest.
sidebar:
  order: 13
---

Los handlers y middleware de Echo son funciones simples sobre un `echo.Context`, por lo que son
fáciles de probar con el paquete estándar `net/http/httptest`. El paquete `echotest` ofrece
helpers que reducen el código repetitivo.

## Probar un handler

Considera dos handlers:

**CreateUser**: `POST /users`

- Acepta un payload JSON.
- Devuelve `201 Created` si tiene éxito.
- Devuelve `500 Internal Server Error` si hay error.

**GetUser**: `GET /users/:email`

- Devuelve `200 OK` si tiene éxito.
- Devuelve `404 Not Found` si el usuario no existe; de lo contrario, `500 Internal Server Error`.

`handler.go`:

```go
package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type (
	User struct {
		Name  string `json:"name" form:"name"`
		Email string `json:"email" form:"email"`
	}
	handler struct {
		db map[string]*User
	}
)

func (h *handler) createUser(c *echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *handler) getUser(c *echo.Context) error {
	email := c.Param("email")
	user := h.db[email]
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}
```

`handler_test.go`:

```go
package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

var (
	mockDB = map[string]*User{
		"jon@labstack.com": {Name: "Jon Snow", Email: "jon@labstack.com"},
	}
	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com"}`
)

func TestCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handler{mockDB}

	// Assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, userJSON+"\n", rec.Body.String())
	}
}
```

### Usar los helpers de echotest

`echotest.ContextConfig` crea un contexto (y recorder) a partir de una descripción
declarativa del request:

```go
// Same test as above, using echotest.
func TestCreateUserWithEchoTest(t *testing.T) {
	c, rec := echotest.ContextConfig{
		Headers: map[string][]string{
			echo.HeaderContentType: {echo.MIMEApplicationJSON},
		},
		JSONBody: []byte(userJSON),
	}.ToContextRecorder(t)

	h := &handler{mockDB}

	// Assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, userJSON+"\n", rec.Body.String())
	}
}

// Even shorter, using ServeWithHandler.
func TestCreateUserWithServeHandler(t *testing.T) {
	h := &handler{mockDB}

	rec := echotest.ContextConfig{
		Headers: map[string][]string{
			echo.HeaderContentType: {echo.MIMEApplicationJSON},
		},
		JSONBody: []byte(userJSON),
	}.ServeWithHandler(t, h.createUser)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, userJSON+"\n", rec.Body.String())
}

func TestGetUser(t *testing.T) {
	c, rec := echotest.ContextConfig{
		PathValues: echo.PathValues{
			{Name: "email", Value: "jon@labstack.com"},
		},
		Headers: map[string][]string{
			echo.HeaderContentType: {echo.MIMEApplicationJSON},
		},
	}.ToContextRecorder(t)

	h := &handler{mockDB}

	// Assertions
	if assert.NoError(t, h.getUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON+"\n", rec.Body.String())
	}
}
```

### Usar un payload de formulario

```go
// import "net/url"
f := make(url.Values)
f.Set("name", "Jon Snow")
f.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
```

Un payload de formulario multipart con `echotest`:

```go
func TestContext_MultipartForm(t *testing.T) {
	testConf := echotest.ContextConfig{
		MultipartForm: &echotest.MultipartForm{
			Fields: map[string]string{
				"key": "value",
			},
			Files: []echotest.MultipartFormFile{
				{
					Fieldname: "file",
					Filename:  "test.json",
					Content:   echotest.LoadBytes(t, "testdata/test.json"),
				},
			},
		},
	}
	c := testConf.ToContext(t)

	assert.Equal(t, "value", c.FormValue("key"))
	assert.Equal(t, http.MethodPost, c.Request().Method)
	assert.Equal(t, true, strings.HasPrefix(c.Request().Header.Get(echo.HeaderContentType), "multipart/form-data; boundary="))

	fv, err := c.FormFile("file")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test.json", fv.Filename)
}
```

### Establecer parámetros de path

```go
c.SetPathValues(echo.PathValues{
	{Name: "id", Value: "1"},
	{Name: "email", Value: "jon@labstack.com"},
})
```

### Establecer parámetros de query

```go
// import "net/url"
q := make(url.Values)
q.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
```

## Probar middleware

```go
func TestMiddleware(t *testing.T) {
	handler := func(c *echo.Context) error {
		return c.JSON(http.StatusTeapot, fmt.Sprintf("email: %s", c.Param("email")))
	}
	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("user_id", int64(1234))
			return next(c)
		}
	}

	c, rec := echotest.ContextConfig{
		PathValues: echo.PathValues{{Name: "email", Value: "jon@labstack.com"}},
	}.ToContextRecorder(t)

	if err := middleware(handler)(c); err != nil {
		t.Fatal(err)
	}

	// Check that the middleware set the value.
	userID, err := echo.ContextGet[int64](c, "user_id")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), userID)

	// Check that the handler returned the correct response.
	assert.Equal(t, http.StatusTeapot, rec.Code)
}
```

:::tip
Para más ejemplos, consulta los [casos de prueba de middleware](https://github.com/labstack/echo/tree/master/middleware)
en el código fuente de Echo.
:::

## Actualización de seguridad de Echo (v5.4.0 / v4.16.0)

`httptest.NewRequest` establece `RemoteAddr` en `192.0.2.1:1234`, una dirección que el extractor de esquema predeterminado no considera de confianza. Si una prueba establece `X-Forwarded-Proto: https`, usa también `req.RemoteAddr = "10.0.0.1:1234"` para representar un proxy privado de confianza, o configura `e.SchemeExtractor = echo.LegacySchemeExtractor()` solo para probar el comportamiento antiguo. Consulta [Esquema de la solicitud](/es/guide/request-scheme/).
