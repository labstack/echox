package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Handler with timeout middleware
	handlerFunc := func(c echo.Context) error {
		time.Sleep(10 * time.Second)
		return c.String(http.StatusOK, "Hello, World!\n")
	}
	middlewareFunc := middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      30 * time.Second,
		ErrorMessage: "my custom error message",
	})
	e.GET("/", handlerFunc, middlewareFunc)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
