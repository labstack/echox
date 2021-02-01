+++
title = "Binding Request Data"
description = "Binding request data"
[menu.main]
  name = "Binding"
  parent = "guide"
+++

## Bind using struct tags

Echo provides following method to bind data from different sources (path params, query params, request body) to structure using 
`Context#Bind(i interface{})` method.
The default binder supports decoding application/json, application/xml and
application/x-www-form-urlencoded data based on the Content-Type header.

In the struct definitions each field can be tagged to restrict binding to specific source.

* `query` - source is request query parameters.
* `param` - source is route path parameter.
* `form` - source is form. Values are taken from query and request body. Uses Go standard library form parsing.
* `json` - source is request body. Uses Go [json](https://golang.org/pkg/encoding/json/) package fo unmarshalling.
* `xml` - source is request body. Uses Go [xml](https://golang.org/pkg/encoding/xml/) package fo unmarshalling.

```go
type User struct {
  ID string `path:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

Request data is binded to the struct in given order:

1. Path parameters
2. Query parameters (only for GET/DELETE methods)
3. Request body

Notes:

* For `query`, `param`, `form` **only** fields **with** tags are bound.
* For `json` and `xml` can bind to *public* fields without tags but this is by their standard library implementation.
* Each step can overwrite binded fields from the previous step. This means if your json request has query param
  `&name=query` and body is `{"name": "body"}` then the result will be `User{Name: "body"}`.
* To avoid security flaws try to avoid passing binded structs directly to other methods if
  these structs contain fields that should not be bindable. It is advisable to have separate struct for binding and map it
  explicitly to your business struct. Consider what will happen if your binded struct has public
  field `IsAdmin bool` and request body would contain `{IsAdmin: true, Name: "hacker"}`.
* When binding forms take note that Echo implementation uses standard library form parsing which parses form data 
  from BOTH URL and BODY if content type is not MIMEMultipartForm. See documentation for [non-MIMEMultipartForm](https://golang.org/pkg/net/http/#Request.ParseForm)
  and [MIMEMultipartForm](https://golang.org/pkg/net/http/#Request.ParseMultipartForm)
* To bind data only from request body use following code
  ```go
  if err := (&DefaultBinder{}).BindBody(c, &payload); err != nil {
    return err
  }
  ```
* To bind data only from query parameters use following code
  ```go
  if err := (&DefaultBinder{}).BindQueryParams(c, &payload); err != nil {
    return err
  }
  ```
* To bind data only from path parameters use following code
  ```go
  if err := (&DefaultBinder{}).BindPathParams(c, &payload); err != nil {
    return err
  }
  ```

### Example

Example below binds the request payload into `User` struct based on tags:

```go
// User
type User struct {
  Name  string `json:"name" form:"name" query:"name"`
  Email string `json:"email" form:"email" query:"email"`
}
```

```go
e.POST("/users", func(c echo.Context) (err error) {
  u := new(User)
  if err = c.Bind(u); err != nil {
    return
  }
  // To avoid security flaws try to avoid passing binded structs directly to other methods 
  // if these structs contain fields that should not be bindable. 
  user := UserDTO{
    Name: u.Name,
    Email: u.Email,
    IsAdmin: false // because you could accidentally expose fields that should not be bind
  }
  executeSomeBusinessLogic(user)
  
  return c.JSON(http.StatusOK, u)
}
```

### JSON Data

```sh
curl -X POST http://localhost:1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@labstack"}'
```

### Form Data

```sh
curl -X POST http://localhost:1323/users \
  -d 'name=Joe' \
  -d 'email=joe@labstack.com'
```

### Query Parameters

```sh
curl -X GET http://localhost:1323/users\?name\=Joe\&email\=joe@labstack.com
```

## Fast binding with dedicated helpers

For binding data found in a request a handful of helper functions are provided. This will allow binding of query parameters, path parameters or data found in the body like forms or JSON data.

Following functions provide a handful of methods for binding to Go native types from request query or path parameters. These binders offer a fluent syntax and can be chained to configure, execute binding and handle errors. 

* `echo.QueryParamsBinder(c)` - binds query parameters (source URL)
* `echo.PathParamsBinder(c)` - binds path parameters (source URL)
* `echo.FormFieldBinder(c)` - binds form fields (source URL + body). See also [Request.ParseForm](https://golang.org/pkg/net/http/#Request.ParseForm).

A binder is usually completed by `BindError()` or `BindErrors()` which returns errors if binding fails.
With `FailFast()` the binder can be configured stop binding on the first error or continue binding for 
the binder call chain. Fail fast is enabled by default and should be disabled when using `BindErrors()`.

`BindError()` returns the first bind error from binder and resets all errors in this binder.
`BindErrors()` returns all bind errors from binder and resets errors in binder.

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

### Supported types

Types that are supported:

* bool
* float32
* float64
* int
* int8
* int16
* int32
* int64
* uint
* uint8/byte (does not support `bytes()`. Use BindUnmarshaler/CustomFunc to convert value from base64 etc to []byte{})
* uint16
* uint32
* uint64
* string
* time
* duration
* BindUnmarshaler() interface
* UnixTime() - converts unix time (integer) to time.Time
* UnixTimeNano() - converts unix time with nano second precision (integer) to time.Time
* CustomFunc() - callback function for your custom conversion logic

For every supported type there are following methods:

* `<Type>("param", &destination)` - if parameter value exists then binds it to given destination of that type i.e `Int64(...)`.
* `Must<Type>("param", &destination)` - parameter value is required to exist, binds it to given destination of that type i.e `MustInt64(...)`.
* `<Type>s("param", &destination)` - (for slices) if parameter values exists then binds it to given destination of that type i.e `Int64s(...)`.
* `Must<Type>s("param", &destination)` - (for slices) parameter value is required to exist, binds it to given destination of that type i.e `MustInt64s(...)`.

for some slice types `BindWithDelimiter("param", &dest, ",")` supports splitting parameter values before type conversion is done. For example URL `/api/search?id=1,2,3&id=1` can be bind to `[]int64{1,2,3,1}`

## Custom Binder

Custom binder can be registered using `Echo#Binder`.

```go
type CustomBinder struct {}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
  // You may use default binder
  db := new(echo.DefaultBinder)
  if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
    return
  }

  // Define your custom implementation here
  return
}
```
