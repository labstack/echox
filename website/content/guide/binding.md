+++
title = "Binding Request Data"
description = "Binding request data"
[menu.main]
  name = "Binding"
  parent = "guide"
+++

Parsing request data is a crucial part of a web application.  In Echo this is done with a process called _binding_. This is done with information passed by the client in the following parts of an HTTP request:

* URL Path parameter
* URL Query parameter
* Header
* Request body

## Binding with Struct Tags

The default binder is accessed through the `Context#Bind(i interface{})` method. It uses Go struct tags to infer from which source(s) data should be loaded.

In this example a struct type `User` tells the binder to bind the query string parameter `id` to its string field `ID`:

```go
type User struct {
  ID string `query:"id"`
}
```

### Data Sources

Echo supports the following tags specifying data sources:

* `query` - query parameter
* `param` - path parameter (also called route)
* `header` - header parameter
* `json` - request body. Uses builtin Go [json](https://golang.org/pkg/encoding/json/) package for unmarshalling.
* `xml` - request body. Uses builtin Go [xml](https://golang.org/pkg/encoding/xml/) package for unmarshalling.
* `form` - form data. Values are taken from query and request body. Uses Go standard library form parsing.

### Data Types

When decoding the request body, the following data types are supported as specified by the `Content-Type` header:

* `application/json`
* `application/xml`
* `application/x-www-form-urlencoded`

When binding path parameter, query parameter, header, or form data, tags must be explicitly set on each struct field. However, JSON and XML binding is done on the struct field name if the tag is omitted. This is according to the behaviour of [Go's json package](https://pkg.go.dev/encoding/json#Unmarshal).

For form data, Echo uses Go standard library form parsing. This parses form data from both the request URL and body if content type is not `MIMEMultipartForm`. See documentation for [non-MIMEMultipartForm](https://golang.org/pkg/net/http/#Request.ParseForm)and [MIMEMultipartForm](https://golang.org/pkg/net/http/#Request.ParseMultipartForm)

### Multiple Sources

It is possible to specify multiple sources on the same field. In this case request data is bound in this order:

1. Path parameters
2. Query parameters (only for GET/DELETE methods)
3. Request body

```go
type User struct {
  ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

Note that binding at each stage will overwrite data bound in a previous stage. This means if your JSON request contains the query param `name=query` and body `{"name": "body"}` then the result will be `User{Name: "body"}`.

### Direct Source

It is also possible to bind data directly from a specific source:
  
Request body:
```go
err := (&DefaultBinder{}).BindBody(c, &payload)
```

Query parameters:
```go
err := (&DefaultBinder{}).BindQueryParams(c, &payload)
```

Path parameters:
```go
err := (&DefaultBinder{}).BindPathParams(c, &payload)
```

Header parameters:
```go
err := (&DefaultBinder{}).BindHeaders(c, &payload)
```

Note that headers is not one of the included sources with `Context#Bind`. The only way to bind header data is by calling `BindHeaders` directly.

### Security

To keep your application secure, avoid passing bound structs directly to other methods if these structs contain fields that should not be bindable. It is advisable to have a separate struct for binding and map it explicitly to your business struct.

Consider what will happen if your bound struct has an Exported field `IsAdmin bool` and the request body contains `{IsAdmin: true, Name: "hacker"}`.

### Example

In this example we define a `User` struct type with field tags to bind from `json`, `form`, or `query` request data:

```go
type User struct {
  Name  string `json:"name" form:"name" query:"name"`
  Email string `json:"email" form:"email" query:"email"`
}

type UserDTO struct {
  Name  string
  Email string
  IsAdmin: bool
}
```

And a handler at the POST `/users` route binds request data to the struct:

```go
e.POST("/users", func(c echo.Context) (err error) {
  u := new(User)
  if err = c.Bind(u); err != nil {
    return
  }

  // Load into separate struct for security
  user := UserDTO{
    Name: u.Name,
    Email: u.Email,
    IsAdmin: false // avoids exposing field that should not be bound
  }

  executeSomeBusinessLogic(user)
  
  return c.JSON(http.StatusOK, u)
}
```

#### JSON Data

```sh
curl -X POST http://localhost:1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@labstack"}'
```

#### Form Data

```sh
curl -X POST http://localhost:1323/users \
  -d 'name=Joe' \
  -d 'email=joe@labstack.com'
```

#### Query Parameters

```sh
curl -X GET 'http://localhost:1323/users?name=Joe&email=joe@labstack.com'
```

## Fast Binding with Dedicated Helpers

Echo provides a handful of helper functions for binding request data. Binding of query parameters, path parameters, and form data found in the body are supported.

The following functions provide a handful of methods for binding to Go data type. These binders offer a fluent syntax and can be chained to configure & execute binding, and handle errors. 

* `echo.QueryParamsBinder(c)` - binds query parameters (source URL)
* `echo.PathParamsBinder(c)` - binds path parameters (source URL)
* `echo.FormFieldBinder(c)` - binds form fields (source URL + body). See also [Request.ParseForm](https://golang.org/pkg/net/http/#Request.ParseForm).

A binder is usually completed by `BindError()` or `BindErrors()` which returns errors if binding has failed.
With `FailFast()` the binder can be configured to stop binding on the first error or continue exsecuting 
the binder call chain. Fail fast is enabled by default and should be disabled when using `BindErrors()`.

`BindError()` returns the first bind error encountered and resets all errors in its binder.
`BindErrors()` returns all bind errors and resets errors in its binder.

```go
// url =  "/api/search?active=true&id=1&id=2&id=3&length=25"
var opts struct {
  IDs []int64
  Active bool
}
length := int64(50) // default length is 50

// creates query params binder that stops binding at first error
err := echo.QueryParamsBinder(c).
  Int64("length", &length).
  Int64s("ids", &opts.IDs).
  Bool("active", &opts.Active).
  BindError() // returns first binding error
```

### Supported Data Types

| Data Type           | Notes |
| ------------------- | ----- |
| `bool`              |       |
| `float32`           |       |
| `float64`           |       |
| `int`               |       |
| `int8`              |       |
| `int16`             |       |
| `int32`             |       |
| `int64`             |       |
| `uint`              |       |
| `uint8/byte`        | Does not support `bytes()`. Use `BindUnmarshaler`/`CustomFunc` to convert value from base64 etc to `[]byte{}`.|
| `uint16`            |       |
| `uint32`            |       |
| `uint64`            |       |
| `string`            |       |
| `time`              |       |
| `duration`          |       |
| `BindUnmarshaler()` | interface |
| `UnixTime()`        | converts unix time (integer) to `time.Time` |
| `UnixTimeNano()`    | converts unix time with nano second precision (integer) to `time.Time` |
| `CustomFunc()`      | callback function for your custom conversion logic |

Each supported type has the following methods:

* `<Type>("param", &destination)` - if parameter value exists then binds it to given destination of that type i.e `Int64(...)`.
* `Must<Type>("param", &destination)` - parameter value is required to exist, binds it to given destination of that type i.e `MustInt64(...)`.
* `<Type>s("param", &destination)` - (for slices) if parameter values exists then binds it to given destination of that type i.e `Int64s(...)`.
* `Must<Type>s("param", &destination)` - (for slices) parameter value is required to exist, binds it to given destination of that type i.e `MustInt64s(...)`.

For certain slice types `BindWithDelimiter("param", &dest, ",")` supports splitting parameter values before type conversion is done. For example binding an integer slice from the URL `/api/search?id=1,2,3&id=1` will result in `[]int64{1,2,3,1}`.

## Custom Binder

A custom binder can be registered using `Echo#Binder`.

```go
type CustomBinder struct {}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
  // You may use default binder
  db := new(echo.DefaultBinder)
  if err := db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
    return
  }

  // Define your custom implementation here
  return
}
```
