---
title: 测试
description: 使用 httptest 和 echotest 辅助方法测试处理函数和中间件。
sidebar:
  order: 13
---

Echo 处理函数和中间件都是围绕 `echo.Context` 的普通函数，因此可以直接使用标准
`net/http/httptest` 包测试。`echotest` 包提供了可减少样板代码的辅助方法。

## 测试处理函数

考虑两个处理函数：

**CreateUser**：`POST /users`

- 接受 JSON 负载。
- 成功时返回 `201 Created`。
- 出错时返回 `500 Internal Server Error`。

**GetUser**：`GET /users/:email`

- 成功时返回 `200 OK`。
- 如果用户不存在，返回 `404 Not Found`；否则返回 `500 Internal Server Error`。

`handler.go`：

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

`handler_test.go`：

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

### 使用 echotest 辅助方法

`echotest.ContextConfig` 会根据请求的声明式描述构建上下文（和 recorder）：

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

### 使用表单负载

```go
// import "net/url"
f := make(url.Values)
f.Set("name", "Jon Snow")
f.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
```

使用 `echotest` 构建 multipart 表单负载：

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

### 设置路径参数

```go
c.SetPathValues(echo.PathValues{
	{Name: "id", Value: "1"},
	{Name: "email", Value: "jon@labstack.com"},
})
```

### 设置查询参数

```go
// import "net/url"
q := make(url.Values)
q.Set("email", "jon@labstack.com")
req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
```

## 测试中间件

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
更多示例请参见 Echo 源码中的
[中间件测试用例](https://github.com/labstack/echo/tree/master/middleware)。
:::

## Echo 安全更新 (v5.4.0 / v4.16.0)

`httptest.NewRequest` 将 `RemoteAddr` 设为 `192.0.2.1:1234`，默认协议提取器不信任此地址。测试设置 `X-Forwarded-Proto: https` 时，还应设 `req.RemoteAddr = "10.0.0.1:1234"` 以模拟可信私有代理。仅测试旧行为时可使用 `e.SchemeExtractor = echo.LegacySchemeExtractor()`。参见[请求协议](/zh-cn/guide/request-scheme/)。
