---
title: 绑定
description: 从路径、查询、header 和请求体将请求数据解析到带类型的 Go struct 中。
sidebar:
  order: 5
---

解析请求数据是 Web 应用的重要组成部分。在 Echo 中这称为
_绑定_，它可以从 HTTP 请求的四个部分读取数据：

- URL 路径参数
- URL 查询参数
- Header
- 请求体

## Struct 标签绑定

定义一个带标签的 struct 来指定数据来源和键，然后把它的指针传给 `c.Bind()`。
这里查询参数 `id` 会绑定到 `ID` 字段：

```go
type User struct {
	ID string `query:"id"`
}

// handler for /users?id=<userID>
var user User
if err := c.Bind(&user); err != nil {
	return c.String(http.StatusBadRequest, "bad request")
}
```

### 数据来源

| 标签     | 来源 |
| -------- | ---- |
| `query`  | 查询参数 |
| `param`  | 路径参数 |
| `header` | Header 值 |
| `form`   | 表单数据（查询 + 请求体） |
| `json`   | 请求体（`encoding/json`） |
| `xml`    | 请求体（`encoding/xml`） |

路径、查询、header 和表单字段都需要**显式标签**。JSON 和 XML 在省略标签时会回退到
struct 字段名，这与标准库的行为一致。

### 切片

重复的查询、路径、表单或请求头值会绑定到切片字段——该字段会收集每一个值：

```go
// GET /search?tag=go&tag=web&tag=api
type Filter struct {
	Tags []string `query:"tag"`
}

var f Filter
if err := c.Bind(&f); err != nil {
	return c.String(http.StatusBadRequest, "bad request")
}
// f.Tags == []string{"go", "web", "api"}
```

### 请求体内容类型

解码请求体时，`Content-Type` header 会选择对应的解码器：

- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`

### 多个来源和优先级

一个字段可以声明多个来源。数据会按下面的顺序绑定，每一步都会覆盖前一步的值：

1. 路径参数
2. 查询参数（仅 GET / DELETE）
3. 请求体

```go
type User struct {
	ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

### 从单一来源直接绑定

```go
echo.BindBody(c, &payload)        // request body
echo.BindQueryParams(c, &payload) // query parameters
echo.BindPathValues(c, &payload)  // path parameters
echo.BindHeaders(c, &payload)     // headers
```

:::note
Header **不会**包含在 `c.Bind()` 中。请直接使用 `echo.BindHeaders` 绑定它们。
:::

:::caution[安全]
不要直接绑定到业务 struct。如果被绑定的 struct 暴露了 `IsAdmin bool` 字段，
请求体 `{"IsAdmin": true}` 就会设置它。请使用专用 DTO，并显式映射：
:::

```go
type UserDTO struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

e.POST("/users", func(c *echo.Context) error {
	var dto UserDTO
	if err := c.Bind(&dto); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	user := User{Name: dto.Name, Email: dto.Email, IsAdmin: false}
	executeSomeBusinessLogic(user)
	return c.JSON(http.StatusOK, user)
})
```

## 流式绑定

如果要从单一来源进行显式且类型安全的绑定，请使用流式 binder。它们会链式配置并执行，
同时收集错误：

```go
// /api/search?active=true&id=1&id=2&id=3&length=25
var opts struct {
	IDs    []int64
	Active bool
}
length := int64(50)

err := echo.QueryParamsBinder(c).
	Int64("length", &length).
	Int64s("id", &opts.IDs).
	Bool("active", &opts.Active).
	BindError() // first error, if any
```

可用的 binder：`echo.QueryParamsBinder(c)`、`echo.PathValuesBinder(c)`、
`echo.FormFieldBinder(c)`。用 `BindError()`（第一个错误）或 `BindErrors()`（所有错误）
结束链式调用。`FailFast(false)` 会运行完整链；默认启用快速失败。

每个支持的类型都提供 `Type(...)`、`MustType(...)`、`Types(...)`（切片）和
`MustTypes(...)` 方法，例如 `Int64`、`MustInt64`、`Int64s`。使用
`BindWithDelimiter("id", &dest, ",")` 可拆分用逗号连接的值。

<a id="custom-binder"></a>

## 自定义 binder

通过 `Echo#Binder` 注册自定义 binder：

```go
type CustomBinder struct{}

func (cb *CustomBinder) Bind(c *echo.Context, i any) error {
	db := new(echo.DefaultBinder)
	if err := db.Bind(c, i); err != echo.ErrUnsupportedMediaType {
		return err
	}
	// custom logic here
	return nil
}

e.Binder = &CustomBinder{}
```
