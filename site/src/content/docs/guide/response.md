---
title: Response
description: Send strings, HTML, JSON, XML, files, streams, redirects, and response hooks.
sidebar:
  order: 8
---

A handler writes its response through the `echo.Context`. Each helper sets the
appropriate `Content-Type` and status code for you.

## Send string

`Context#String(code int, s string)` sends a plain text response with a status code.

```go
func(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

## Send HTML

`Context#HTML(code int, html string)` sends a simple HTML response with a status
code. To generate HTML dynamically, see [Templates](/guide/templates/).

```go
func(c *echo.Context) error {
	return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
}
```

### Send HTML blob

`Context#HTMLBlob(code int, b []byte)` sends an HTML blob with a status code. It is
handy with a template engine that outputs `[]byte`.

```go
func handler(c *echo.Context) error {
	blob := []byte("<strong>Hello, World!</strong>")
	return c.HTMLBlob(http.StatusOK, blob)
}
```

## Render template

See [Templates](/guide/templates/).

## Send JSON

`Context#JSON(code int, i any)` encodes a Go value as JSON and sends it with a
status code.

```go
type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.JSON(http.StatusOK, u)
}
```

### Stream JSON

`Context#JSON()` uses `json.Marshal` internally, which may be inefficient for large
payloads. In that case, stream the JSON directly:

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(u)
}
```

### JSON pretty

`Context#JSONPretty(code int, i any, indent string)` sends a pretty-printed JSON
response. The indent can be spaces or tabs.

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}
```

```json
{
  "email": "jon@labstack.com",
  "name": "Jon"
}
```

### JSON blob

`Context#JSONBlob(code int, b []byte)` sends a pre-encoded JSON blob directly, for
example from a database.

```go
func(c *echo.Context) error {
	encodedJSON := []byte{} // Encoded JSON from an external source.
	return c.JSONBlob(http.StatusOK, encodedJSON)
}
```

## Send JSONP

`Context#JSONP(code int, callback string, i any)` encodes a Go value as JSON and
sends it as a JSONP payload wrapped in the given callback.

```go
func handler(c *echo.Context) error {
	callback := c.QueryParam("callback")
	return c.JSONP(http.StatusOK, callback, &User{Name: "Jon", Email: "jon@labstack.com"})
}
```

For cross-origin requests, use [JSON with CORS](/cookbook/cors/).

The callback must be empty, a JavaScript identifier, or a dot-separated path of
identifiers. Invalid callbacks return HTTP 400 with an error wrapping
`ErrInvalidJSONPCallback`; JSONP responses include `X-Content-Type-Options: nosniff`.
Because any site can read JSONP with the user's cookies, do not return private or
authenticated data through JSONP. Use JSON with CORS for that case.

## Send XML

`Context#XML(code int, i any)` encodes a Go value as XML and sends it with a status
code.

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.XML(http.StatusOK, u)
}
```

### Stream XML

`Context#XML` uses `xml.Marshal` internally, which may be inefficient for large
payloads. In that case, stream the XML directly:

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return xml.NewEncoder(c.Response()).Encode(u)
}
```

### XML pretty

`Context#XMLPretty(code int, i any, indent string)` sends a pretty-printed XML
response. The indent can be spaces or tabs.

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.XMLPretty(http.StatusOK, u, "  ")
}
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<User>
  <Name>Jon</Name>
  <Email>jon@labstack.com</Email>
</User>
```

:::tip
You can also make `Context#XML()` output pretty-printed XML by appending `pretty`
to the request URL query string.

```sh
curl http://localhost:1323/users/1?pretty
```
:::

### XML blob

`Context#XMLBlob(code int, b []byte)` sends a pre-encoded XML blob directly, for
example from a database.

```go
func(c *echo.Context) error {
	encodedXML := []byte{} // Encoded XML from an external source.
	return c.XMLBlob(http.StatusOK, encodedXML)
}
```

## Send file

`Context#File(file string)` sends the contents of a file as the response. It sets
the correct content type and handles caching automatically.

```go
func(c *echo.Context) error {
	return c.File("<PATH_TO_YOUR_FILE>")
}
```

## Send attachment

`Context#Attachment(file, name string)` is like `File()` but sends the file with
`Content-Disposition: attachment` and the given name.

```go
func(c *echo.Context) error {
	return c.Attachment("<PATH_TO_YOUR_FILE>", "<ATTACHMENT_NAME>")
}
```

## Send inline

`Context#Inline(file, name string)` is like `File()` but sends the file with
`Content-Disposition: inline` and the given name.

```go
func(c *echo.Context) error {
	return c.Inline("<PATH_TO_YOUR_FILE>", "<INLINE_NAME>")
}
```

## Send blob

`Context#Blob(code int, contentType string, b []byte)` sends arbitrary data with a
given content type and status code.

```go
func(c *echo.Context) error {
	data := []byte(`0306703,0035866,NO_ACTION,06/19/2006
0086003,"0005866",UPDATED,06/19/2006`)
	return c.Blob(http.StatusOK, "text/csv", data)
}
```

## Send stream

`Context#Stream(code int, contentType string, r io.Reader)` sends an arbitrary data
stream with a given content type, `io.Reader`, and status code.

```go
func(c *echo.Context) error {
	f, err := os.Open("<PATH_TO_IMAGE>")
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Stream(http.StatusOK, "image/png", f)
}
```

## Send no content

`Context#NoContent(code int)` sends an empty body with a status code.

```go
func(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}
```

## Redirect request

`Context#Redirect(code int, url string)` redirects the request to the given URL with
a status code.

```go
func(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "<URL>")
}
```

## Hooks

### Before response

`Response#Before(func())` registers a function that runs just before the response is
written.

### After response

`Response#After(func())` registers a function that runs just after the response is
written. If the `Content-Length` is unknown, no after functions run.

```go
e.GET("/hooks", func(c *echo.Context) error {
	resp, err := echo.UnwrapResponse(c.Response())
	if err != nil {
		return err
	}
	resp.Before(func() {
		println("before response")
	})
	resp.After(func() {
		println("after response")
	})
	return c.String(http.StatusOK, "Hello, World!")
})
```

:::tip
You can register multiple `Before` and `After` functions.
:::
