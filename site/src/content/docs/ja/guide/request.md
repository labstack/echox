---
title: リクエスト
description: リクエストからフォーム、クエリ、パスのデータを取得し、検証します。
sidebar:
  order: 7
---

ハンドラは `echo.Context` を通じてリクエストデータを読み取ります。Echo は値を名前で個別に取得し、
struct にバインドし（[バインディング](/ja/guide/binding/)を参照）、登録した validator に
検証を委ねられます。

## データを取得する

### フォームデータ

`Context#FormValue(name string)` でフォームフィールドを名前から取得します。

```go
e.POST("/form", func(c *echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
})
```

`string` 以外の型には、ジェネリック関数 `echo.FormValue[T]` を使います。

```go
age, err := echo.FormValue[int](c, "age")
if err != nil {
	return err
}
```

次でテストします。

```sh
curl -X POST http://localhost:1323/form -d 'name=Joe&age=30'
```

カスタムデータ型をバインドするには、`echo.BindUnmarshaler` インターフェイスを実装します。

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

### クエリパラメーター

`Context#QueryParam(name string)` でクエリパラメーターを名前から取得します。

```go
func(c *echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}
```

`string` 以外の型には、ジェネリック関数 `echo.QueryParam[T]` を使います。

```go
age, err := echo.QueryParam[int](c, "age")
if err != nil {
	return err
}
```

```sh
curl -X GET "http://localhost:1323?name=Joe&age=30"
```

### パスパラメーター

`Context#Param(name string)` で登録済みのパスパラメーターを名前から取得します。

```go
e.GET("/users/:name", func(c *echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
})
```

`string` 以外の型には、ジェネリック関数 `echo.PathParam[T]` を使います。

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

### データをバインドする

Echo はリクエストデータをネイティブの Go struct や変数にもバインドできます。
[バインディング](/ja/guide/binding/)を参照してください。

<a id="validate-data"></a>

## データを検証する

Echo には組み込みのデータ検証はありません。`Echo#Validator` でカスタム validator を登録し、
[go-playground/validator](https://github.com/go-playground/validator) などのサードパーティライブラリを
使えます。

次の例は、バインド済み struct を検証します。

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
