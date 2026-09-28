---
title: Testes
description: Teste handlers e middleware com httptest e os helpers echotest.
sidebar:
  order: 13
---

Handlers e middleware do Echo são funções simples sobre um `echo.Context`, então são
diretos de testar com o pacote padrão `net/http/httptest`. O pacote
`echotest` fornece helpers que reduzem boilerplate.

## Testando um handler

Considere dois handlers:

**CreateUser** — `POST /users`

- Aceita um payload JSON.
- Retorna `201 Created` em caso de sucesso.
- Retorna `500 Internal Server Error` em caso de erro.

**GetUser** — `GET /users/:email`

- Retorna `200 OK` em caso de sucesso.
- Retorna `404 Not Found` se o usuário não existir, caso contrário `500 Internal Server Error`.

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

### Usando os helpers echotest

`echotest.ContextConfig` constrói um contexto (e recorder) a partir de uma descrição
declarativa do request:

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

### Usando um payload de formulário

```go
// import "net/url"
f := make(url.Values)
f.Set("name", "Jon Snow")
f.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
```

Um payload de formulário multipart com `echotest`:

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

### Definindo parâmetros de caminho

```go
c.SetPathValues(echo.PathValues{
	{Name: "id", Value: "1"},
	{Name: "email", Value: "jon@labstack.com"},
})
```

### Definindo parâmetros de query

```go
// import "net/url"
q := make(url.Values)
q.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
```

## Testando middleware

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
Para mais exemplos, veja os [casos de teste de middleware](https://github.com/labstack/echo/tree/master/middleware)
no código-fonte do Echo.
:::

## Atualização de segurança do Echo (v5.4.0 / v4.16.0)

`httptest.NewRequest` define `RemoteAddr` como `192.0.2.1:1234`, endereço que o extrator de esquema padrão não considera confiável. Se o teste define `X-Forwarded-Proto: https`, defina também `req.RemoteAddr = "10.0.0.1:1234"` para simular um proxy privado confiável. Use `e.SchemeExtractor = echo.LegacySchemeExtractor()` apenas para testar o comportamento antigo. Veja [Esquema da requisição](/pt-br/guide/request-scheme/).
