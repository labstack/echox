---
title: 请求
description: 从请求中获取表单、查询和路径数据，并验证它。
sidebar:
  order: 7
---

处理函数通过 `echo.Context` 读取请求数据。Echo 可以按名称分别获取值，
把它们绑定到 struct 中（参见[绑定](/zh-cn/guide/binding/)），并把验证交给你注册的 validator。

## 获取数据

### 表单数据

使用 `Context#FormValue(name string)` 按名称获取表单字段：

```go
e.POST("/form", func(c *echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
})
```

对于 `string` 以外的类型，请使用泛型函数 `echo.FormValue[T]`：

```go
age, err := echo.FormValue[int](c, "age")
if err != nil {
	return err
}
```

测试：

```sh
curl -X POST http://localhost:1323/form -d 'name=Joe&age=30'
```

要绑定自定义数据类型，请实现 `echo.BindUnmarshaler` 接口：

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

### 查询参数

使用 `Context#QueryParam(name string)` 按名称获取查询参数：

```go
func(c *echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}
```

对于 `string` 以外的类型，请使用泛型函数 `echo.QueryParam[T]`：

```go
age, err := echo.QueryParam[int](c, "age")
if err != nil {
	return err
}
```

```sh
curl -X GET "http://localhost:1323?name=Joe&age=30"
```

### 路径参数

使用 `Context#Param(name string)` 按名称获取已注册的路径参数：

```go
e.GET("/users/:name", func(c *echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
})
```

对于 `string` 以外的类型，请使用泛型函数 `echo.PathParam[T]`：

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

### 绑定数据

Echo 也可以把请求数据绑定到原生 Go struct 和变量中。参见
[绑定](/zh-cn/guide/binding/)。

<a id="validate-data"></a>

## 验证数据

Echo 没有内置数据验证。你可以通过 `Echo#Validator` 注册自定义 validator，并使用
[go-playground/validator](https://github.com/go-playground/validator) 这样的第三方库。

下面的示例会验证一个已绑定的 struct：

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
