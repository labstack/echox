---
description: Testing handler and middleware
slug: /testing
sidebar_position: 13
---

# Testing

## Testing Handler

`GET` `/users/:id`

Handler below retrieves user by id from the database. If user is not found it returns
`404` error with a message.

### CreateUser

`POST` `/users`

- Accepts JSON payload
- On success `201 - Created`
- On error `500 - Internal Server Error`

### GetUser

`GET` `/users/:email`

- On success `200 - OK`
- On error `404 - Not Found` if user is not found otherwise `500 - Internal Server Error`

`handler.go`

```go
package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

func (h *handler) createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *handler) getUser(c echo.Context) error {
	email := c.Param("email")
	user := h.db[email]
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}
```

`handler_test.go`

```go
package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	mockDB = map[string]*User{
		"jon@labstack.com": &User{"Jon Snow", "jon@labstack.com"},
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
		assert.Equal(t, userJSON, rec.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:email")
	c.SetParamNames("email")
	c.SetParamValues("jon@labstack.com")
	h := &handler{mockDB}

	// Assertions
	if assert.NoError(t, h.getUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}
```

### Using Form Payload

```go
// import "net/url"
f := make(url.Values)
f.Set("name", "Jon Snow")
f.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
```

### Setting Path Params

```go
c.SetParamNames("id", "email")
c.SetParamValues("1", "jon@labstack.com")
```

### Setting Query Params

```go
// import "net/url"
q := make(url.Values)
q.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
```

## Testing Middleware

*TBD*

For now you can look into built-in middleware [test cases](https://github.com/labstack/echo/tree/master/middleware).
