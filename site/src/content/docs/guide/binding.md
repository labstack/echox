---
title: Binding
description: Parse request data into typed Go structs from path, query, header, and body.
sidebar:
  order: 5
---

Parsing request data is a crucial part of a web application. In Echo this is called
_binding_, and it can read from four parts of an HTTP request:

- URL path parameters
- URL query parameters
- Headers
- Request body

## Struct tag binding

Define a struct with tags specifying the data source and key, then call `c.Bind()`
with a pointer to it. Here the query parameter `id` binds to the `ID` field:

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

### Data sources

| Tag      | Source |
| -------- | ------ |
| `query`  | Query parameter |
| `param`  | Path parameter |
| `header` | Header value |
| `form`   | Form data (query + body) |
| `json`   | Request body (`encoding/json`) |
| `xml`    | Request body (`encoding/xml`) |

Path, query, header, and form fields require an **explicit tag**. JSON and XML fall
back to the struct field name when the tag is omitted, matching the standard library.

### Slices

Repeated query, path, form, or header values bind to a slice field — the field
collects every occurrence:

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

### Body content types

When decoding the request body, the `Content-Type` header selects the decoder:

- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`

### Multiple sources & precedence

A field can declare several sources. Data is bound in this order, each step
overwriting the previous:

1. Path parameters
2. Query parameters (GET / DELETE / HEAD / QUERY only)
3. Request body

```go
type User struct {
	ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

### Direct binding from one source

```go
echo.BindBody(c, &payload)        // request body
echo.BindQueryParams(c, &payload) // query parameters
echo.BindPathValues(c, &payload)  // path parameters
echo.BindHeaders(c, &payload)     // headers
```

:::note
Headers are **not** included by `c.Bind()`. Bind them with `echo.BindHeaders` directly.
:::

:::caution[Security]
Don't bind directly into business structs. If a bound struct exposed an `IsAdmin bool`
field, a request body of `{"IsAdmin": true}` would set it. Use a dedicated DTO and map
it explicitly:
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

## Fluent binding

For explicit, type-safe binding from a single source, use the fluent binders. They
chain configuration and execute, collecting errors:

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

Available binders: `echo.QueryParamsBinder(c)`, `echo.PathValuesBinder(c)`,
`echo.FormFieldBinder(c)`. End a chain with `BindError()` (first error) or
`BindErrors()` (all errors). `FailFast(false)` runs the whole chain; it's on by default.

Each supported type offers `Type(...)`, `MustType(...)`, `Types(...)` (slices), and
`MustTypes(...)` methods — e.g. `Int64`, `MustInt64`, `Int64s`. Use
`BindWithDelimiter("id", &dest, ",")` to split comma-joined values.

## Custom binder

Register a custom binder via `Echo#Binder`:

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
