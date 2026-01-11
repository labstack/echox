---
description: Using templates
slug: /templates
sidebar_position: 12
---

# Templates

## Rendering

`Context#Render(code int, name string, data any) error` renders a template
with data and sends a text/html response with status code. Templates can be registered by setting `Echo.Renderer`, allowing us to use any template engine.

Example below shows how to use Go `html/template`:

1. Use default template renderer
```go
e.Renderer = &echo.TemplateRenderer{
	Template: template.Must(template.New("hello").Parse("Hello, {{.}}!")),
}
```

or Implement `echo.Renderer` interface

```go
type Template struct {
	templates *template.Template
}

func (t *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
 ```

2. Pre-compile templates

    `public/views/hello.html`

    ```html
    {{define "hello"}}Hello, {{.}}!{{end}}
    ```

    ```go
    t := &Template{
        templates: template.Must(template.ParseGlob("public/views/*.html")),
    }
    ```

3. Register templates

    ```go
    e := echo.New()
    e.Renderer = t
    e.GET("/hello", Hello)
    ```

4. Render a template inside your handler

    ```go
    func Hello(c *echo.Context) error {
    	return c.Render(http.StatusOK, "hello", "World")
    }
    ```

## Advanced - Calling Echo from templates

In certain situations it might be useful to generate URIs from the templates. In order to do so, you need to call `Echo#Reverse` from the templates itself. Golang's `html/template` package is not the best suited for this job, but this can be done in two ways: by providing a common method on all objects passed to templates or by passing `map[string]any` and augmenting this object in the custom renderer. Given the flexibility of the latter approach, here is a sample program:

`template.html`

```html
<html>
<body>
<h1>Hello {{index . "name"}}</h1>

<p>{{ with $x := index . "reverse" }}
   {{ call $x "foobar" }} <!--- this will call the $x with parameter "foobar" -->
   {{ end }}
</p>
</body>
</html>
```

`server.go`

```go
package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]any); isMap {
		viewContext["reverse"] = c.RouteInfo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("main/*.html")),
	}

	e.GET("/something/:name", func(c *echo.Context) error {
		return c.Render(http.StatusOK, "template.html", map[string]any{
			"name": "Dolly!",
		})
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("shutting down the server", "error", err)
	}
}
```
