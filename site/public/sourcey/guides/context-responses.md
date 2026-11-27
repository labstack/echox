# Context and responses {#context-and-responses}



## Context is a struct in v5 {#context-is-a-struct-in-v5}


In v4 `Context` was an interface; in v5 it is a concrete struct, and handlers take a pointer to it:

```go
// v4 — no longer valid in v5
func handler(c echo.Context) error { ... }

// v5
func handler(c *echo.Context) error { ... }
```

Because it is a struct, new methods can be added in minor versions without breaking existing code. The framework hands each request a pooled `*Context` (`AcquireContext` / `ReleaseContext`); do not hold onto it after the handler returns.


## Reading the request {#reading-the-request}


```go
c.Request()              // *http.Request
c.Param("id")            // path parameter
c.ParamOr("id", "none")  // path parameter with fallback
c.QueryParam("q")        // query value
c.QueryParamOr("q", "")  // query value with fallback
c.QueryParams()          // url.Values — all query parameters
c.QueryString()          // raw query string
c.FormValue("name")      // form value
c.FormValueOr("name", "")// form value with fallback
c.FormValues()           // all form values
c.Cookie("session")      // *http.Cookie, or ErrCookieNotFound
c.Cookies()              // all cookies
c.RealIP()               // client IP (uses e.IPExtractor)
c.Scheme()               // http or https (uses e.SchemeExtractor)
c.IsTLS()                // whether the request is TLS
c.IsWebSocket()          // whether it is a WebSocket upgrade
c.Path()                 // registered route pattern, e.g. /users/:id
c.RouteInfo()            // matched route info (method, path, handler name)
```

`c.Get(key)` / `c.Set(key, value)` store request-scoped values. The generic helpers read them type-safely:

```go
user, ok := echo.ContextGet[User](c, "user")
tenant := echo.ContextGetOr(c, "tenant", "default")
```


## Writing responses {#writing-responses}


```go
c.String(http.StatusOK, "Hello, World!")
c.JSON(http.StatusOK, map[string]string{"hello": "world"})
c.JSONPretty(http.StatusOK, data, "  ")
c.JSONP(http.StatusOK, "callback", data)
c.XML(http.StatusOK, data)
c.XMLPretty(http.StatusOK, data, "  ")
c.HTML(http.StatusOK, "<b>hi</b>")
c.HTMLBlob(http.StatusOK, []byte("<b>hi</b>"))
c.Blob(http.StatusOK, "application/octet-stream", data)
c.Stream(http.StatusOK, "text/plain", reader)
c.File("/path/to/file.pdf")
c.FileFS("file.pdf", fs)
c.Attachment("/path/to/report.pdf", "report.pdf") // Content-Disposition: attachment
c.Inline("/path/to/report.pdf", "report.pdf")     // Content-Disposition: inline
c.NoContent(http.StatusNoContent)
c.Redirect(http.StatusSeeOther, "/login")
c.Render(http.StatusOK, "hello.html", data) // needs e.Renderer set, else ErrRendererNotRegistered
```

`c.Response()` returns the `*Response` (a wrapper around `http.ResponseWriter` with `Committed` state and byte counting). Since v5.4.7, `Response.Flush` commits an uncommitted response, and `WriteHeader` passes 1xx informational codes (such as `103 Early Hints`) through without committing.


## Templates {#templates}


```go
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c *echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

e.Renderer = &Template{templates: template.Must(template.ParseGlob("views/*.html"))}
```

Set `e.Renderer` to any type implementing `Render(io.Writer, string, any, *Context) error`. Calling `c.Render` without a renderer returns `ErrRendererNotRegistered`.
