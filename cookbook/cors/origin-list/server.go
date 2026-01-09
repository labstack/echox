package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var (
	users = []string{"Joe", "Veer", "Zion"}
)

func getUsers(c *echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	// e.Use(middleware.CORS("*"))

	// CORS restricted
	// Allows requests from any `https://labstack.com` or `https://labstack.net` origin
	e.Use(middleware.CORS("https://labstack.com", "https://labstack.net"))

	e.GET("/api/users", getUsers)

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
