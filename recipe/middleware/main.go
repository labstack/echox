package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func main() {
	e := echo.New()

	// Debug mode
	e.Debug()

	//------------
	// Middleware
	//------------

	// Logger
	e.Use(middleware.Logger())

	// Recover
	e.Use(middleware.Recover())

	// Basic auth
	e.Use(middleware.BasicAuth(func(usr, pwd string) bool {
		if usr == "joe" && pwd == "secret" {
			return true
		}
		return false
	}))

	// Gzip
	e.Use(middleware.Gzip())

	// Routes
	e.Get("/", hello)

	// Start server
	e.Run(standard.New(":1323"))
}
