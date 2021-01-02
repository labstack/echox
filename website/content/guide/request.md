+++
title = "Request"
description = "Handling HTTP request in Echo"
[menu.main]
  name = "Request"
  parent = "guide"
+++

## Bind request data with bind functions

Following functions provide handful of methods for binding to Go native types from request query or path parameters.

* `echo.QueryParamsBinder(c)` - binds query parameters (source URL)
* `echo.PathParamsBinder(c)` - binds path parameters (source URL)
* `echo.FormFieldBinder(c)` - binds form fields (source URL + body). See also [Request.ParseForm](https://golang.org/pkg/net/http/#Request.ParseForm).

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

### Methods
`FailFast` flags binder to stop binding after first bind error during binder call chain. Enabled by default.
`BindError()` returns first bind error from binder and resets errors in binder. Useful along with `FailFast()` method
to do binding and returns on first problem
`BindErrors()` returns all bind errors from binder and resets errors in binder.

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
* CustomFunc() - callback function for your custom conversion logic

For every supported type there are following methods:

* `<Type>("param", &destination)` - if parameter value exists then binds it to given destination of that type i.e `Int64(...)`.
* `Must<Type>("param", &destination)` - parameter value is required to exist, binds it to given destination of that type i.e `MustInt64(...)`.
* `<Type>s("param", &destination)` - (for slices) if parameter values exists then binds it to given destination of that type i.e `I`nt64s(...)`.
* `Must<Type>s("param", &destination)` - (for slices) parameter value is required to exist, binds it to given destination of that type i.e `MustInt64s(...)`.

for some slice types `BindWithDelimiter("param", &dest, ",")` supports slitting parameter values before type conversion is done
i.e. URL `/api/search?id=1,2,3&id=1` can be bind to `[]int64{1,2,3,1}`

## Bind request data with struct tags

Echo provides following method to bind data from different sources (path params, query params, request body) to structure using 
`Context#Bind(i interface{})` method.
The default binder supports decoding application/json, application/xml and
application/x-www-form-urlencoded data based on the Content-Type header.

Supported struct tags:
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

* For `query`, `param`, `form` **only** fields **with** tags are binded. 
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

### Example 1
Example below binds the request payload into `User` struct based on tags:

```go
// User
type User struct {
  Name  string `json:"name" form:"name" query:"name"`
  Email string `json:"email" form:"email" query:"email"`
}
```

```go
// Handler
func(c echo.Context) (err error) {
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
curl \
  -X POST \
  http://localhost:1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@labstack"}'
```

### Form Data

```sh
curl \
  -X POST \
  http://localhost:1323/users \
  -d 'name=Joe' \
  -d 'email=joe@labstack.com'
```

### Query Parameters

```sh
curl \
  -X GET \
  http://localhost:1323/users\?name\=Joe\&email\=joe@labstack.com
```

## Custom Binder

Custom binder can be registered using `Echo#Binder`.

*Example*

```go
type CustomBinder struct {}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	// You may use default binder
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	// Define your custom implementation

	return
}
```

## Retrieve Data

### Form Data

Form data can be retrieved by name using `Context#FormValue(name string)`.

*Example*

```go
// Handler
func(c echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
}
```

```sh
curl \
  -X POST \
  http://localhost:1323 \
  -d 'name=Joe'
```

To bind a custom data type, you can implement `Echo#BindUnmarshaler` interface.

*Example*

```go
type Timestamp time.Time

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}
```

### Query Parameters

Query parameters can be retrieved by name using `Context#QueryParam(name string)`.

*Example*

```go
// Handler
func(c echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
})
```

```sh
curl \
  -X GET \
  http://localhost:1323\?name\=Joe
```

Similar to form data, custom data type can be bind using `Context#QueryParam(name string)`.

### Path Parameters

Registered path parameters can be retrieved by name using `Context#Param(name string) string`.

*Example*

```go
e.GET("/users/:name", func(c echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
})
```

```sh
$ curl http://localhost:1323/users/Joe
```

## Validate Data

Echo doesn't have built-in data validation capabilities, however, you can register
a custom validator using `Echo#Validator` and leverage third-party [libraries](https://github.com/avelino/awesome-go#validation).

Example below uses https://github.com/go-playground/validator framework for validation:

```go
type (
	User struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return echo.NewHTTPError(http.StatusInternalServerError, cv.validator.Struct(i).Error())
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```

```sh
curl \
  -X POST \
  http://localhost:1323/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@invalid-domain"}'
{"message":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}
```
