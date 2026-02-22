package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// 1. Start the program
// 2. Open http://localhost:8080/form in browser
// 3. Submit the form
// 4. Result has submitted token value. It is a random string if your browser is old
//    or `_echo_csrf_using_sec_fetch_site_` when you are using a browser that supports the `Sec-Fetch-Site` header

func main() {
	e := echo.New()
	e.Renderer = &echo.TemplateRenderer{
		Template: template.Must(template.New("form.html").Parse(formHTML)),
	}

	// CSRF middleware:
	// - sets a CSRF cookie
	// - expects the token in form field named `csrf`
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// IMPORTANT: if you want to submit via hidden form field, include "form:<fieldname>"
		// Here we accept token from the form field named "csrf".
		TokenLookup:    "form:csrf",
		CookieName:     "_csrf",
		CookieHTTPOnly: true,
	}))

	e.GET("/form", func(c *echo.Context) error {
		// Echo CSRF middleware stores the token in the context under the key "csrf".
		// We'll render it into the form as a hidden field.
		token, err := echo.ContextGet[string](c, "csrf")
		if err != nil {
			return err
		}
		data := map[string]any{
			"CSRFToken": token,
		}
		return c.Render(http.StatusOK, "form.html", data)
	})

	e.POST("/submit", func(c *echo.Context) error {
		// If the CSRF token is invalid/missing, middleware returns 403 and this handler is not executed.
		name := c.FormValue("name")
		formToken := c.FormValue("csrf")
		return c.String(http.StatusOK, fmt.Sprintf("OK, submitted name: '%s', token value in form: '%s'", name, formToken))
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server: ", "error", err)
	}
}

const formHTML = `
{{ define "form.html" }}
<!doctype html>
<html>
<head><meta charset="utf-8"><title>Echo CSRF Form</title></head>
<body>
  <h1>CSRF Hidden Field Example</h1>

  <form method="POST" action="/submit">
    <input type="hidden" name="csrf" value="{{ .CSRFToken }}">
    <label>Name: <input type="text" name="name" /></label>
    <button type="submit">Submit</button>
  </form>

</body>
</html>
{{ end }}
`
