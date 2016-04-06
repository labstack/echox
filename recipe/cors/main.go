package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/rs/cors"
)

var (
	users = []string{"Joe", "Veer", "Zion"}
)

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(standard.WrapMiddleware(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost"},
	}).Handler))

	e.Get("/api/users", getUsers)
	e.Run(standard.New(":1323"))
}
