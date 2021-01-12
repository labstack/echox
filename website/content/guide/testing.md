+++
title = "Testing"
description = "Testing handler and middleware in Echo"
[menu.main]
  name = "Testing"
  parent = "guide"
+++

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

### CheckEmail

- On invalid email `400 - BadRequest`
- On valid email `nil`

`handler.go`

```go
package handler

import (
	"net/http"
	"strings"

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

func (h *handler) checkEmail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.Param("email")
		if !strings.Contains(email, "@") {
			echo.NewHTTPError(http.StatusBadRequest, "invalid email address")
		}
		c.Set("validEmail", true)
		return nil
	}
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

func TestCheckEmail(t *testing.T) {
    var tcs = []struct{
        email string
        err string
    }{
        {
            email: "valid@email.com",
            err: "",
        },
        {
            email: "invalid",
            err: "invalid email address",
        },
    }

    for _, tc := range tcs {
        e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("email")
        c.SetParamValues(tc.email)
        h := &handler{mockDB}

        // Check that "validEmail" was properly set.
        var next = h.checkEmail(func(ec echo.Context) error {
            var e = ec.Get("validEmail")
            var email, ok = e.(bool)
            assert.True(t, ok)
            assert.Equal(t, email, tc.email)

            return nil
        })

        err := next(c)
        if tc.err == "" {
			assert.EqualError(t, err, tc.err)
        } else {
            assert.Nil(t, err)
        }
    }
}
```

For additional examples, check the built-in middleware [test cases](https://github.com/labstack/echo/tree/master/middleware).
