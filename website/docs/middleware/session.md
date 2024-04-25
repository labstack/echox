---
description: Session middleware
---

# Session

Session middleware facilitates HTTP session management backed by [gorilla sessions](https://github.com/gorilla/sessions). The default implementation provides cookie and
filesystem based session store; however, you can take advantage of [community maintained
implementation](https://github.com/gorilla/sessions#store-implementations) for various backends.

:::note

Echo community contribution

:::

## Dependencies

```go
import (
  "github.com/gorilla/sessions"
  "github.com/labstack/echo-contrib/session"
)
```

## Usage

This example exposes two endpoints: `/create-session` creates new session and `/read-session` read value from 
session if request contains session id.

```go
import (
    "errors"
    "fmt"
    "github.com/gorilla/sessions"
    "github.com/labstack/echo-contrib/session"
    "github.com/labstack/echo/v4"
    "log"
    "net/http"
)

func main() {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/create-session", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["foo"] = "bar"
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	e.GET("/read-session", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, fmt.Sprintf("foo=%v\n", sess.Values["foo"]))
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```

### Example usage

Requesting `/read-session` without providing session it will output nil as `foo` value
```bash
x@x:~/$ curl -v http://localhost:8080/read-session
* processing: http://localhost:8080/read-session
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /read-session HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.2.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=UTF-8
< Date: Thu, 25 Apr 2024 09:15:14 GMT
< Content-Length: 10
< 
foo=<nil>
```

Requesting `/create-session` creates new session
```bash
x@x:~/$ curl -v -c cookies.txt http://localhost:8080/create-session
* processing: http://localhost:8080/create-session
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /create-session HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.2.1
> Accept: */*
> 
< HTTP/1.1 200 OK
* Added cookie session="MTcxNDAzNjYyMHxEWDhFQVFMX2dBQUJFQUVRQUFBZ180QUFBUVp6ZEhKcGJtY01CUUFEWm05dkJuTjBjbWx1Wnd3RkFBTmlZWEk9fHJQxR5fJDUEV-6iHSWuyVzjYX2f9F5tVaMGV6pjIE1Y" for domain localhost, path /, expire 1714641420
< Set-Cookie: session=MTcxNDAzNjYyMHxEWDhFQVFMX2dBQUJFQUVRQUFBZ180QUFBUVp6ZEhKcGJtY01CUUFEWm05dkJuTjBjbWx1Wnd3RkFBTmlZWEk9fHJQxR5fJDUEV-6iHSWuyVzjYX2f9F5tVaMGV6pjIE1Y; Path=/; Expires=Thu, 02 May 2024 09:17:00 GMT; Max-Age=604800; HttpOnly
< Date: Thu, 25 Apr 2024 09:17:00 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

Using session cookie from previous response and requesting `/read-session` will output `foo` value from session.
```bash
x@x:~/$ curl -v -b cookies.txt http://localhost:8080/read-session
* processing: http://localhost:8080/read-session
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /read-session HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.2.1
> Accept: */*
> Cookie: session=MTcxNDAzNjYyMHxEWDhFQVFMX2dBQUJFQUVRQUFBZ180QUFBUVp6ZEhKcGJtY01CUUFEWm05dkJuTjBjbWx1Wnd3RkFBTmlZWEk9fHJQxR5fJDUEV-6iHSWuyVzjYX2f9F5tVaMGV6pjIE1Y
> 
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=UTF-8
< Date: Thu, 25 Apr 2024 09:18:56 GMT
< Content-Length: 8
< 
foo=bar
* Connection #0 to host localhost left intact
```

## Custom Configuration

### Usage

```go
e := echo.New()
e.Use(session.MiddlewareWithConfig(session.Config{}))
```

## Configuration

```go
Config struct {
  // Skipper defines a function to skip middleware.
  Skipper middleware.Skipper

  // Session store.
  // Required.
  Store sessions.Store
}
```

### Default Configuration

```go
DefaultConfig = Config{
  Skipper: DefaultSkipper,
}
```
