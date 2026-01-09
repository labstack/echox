package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// Hosts
	vHosts := make(map[string]*echo.Echo)

	//-----
	// API
	//-----

	api := echo.New()
	api.Use(middleware.RequestLogger())
	api.Use(middleware.Recover())

	vHosts["api.localhost:1323"] = api

	api.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "API")
	})

	//------
	// Blog
	//------

	blog := echo.New()
	blog.Use(middleware.RequestLogger())
	blog.Use(middleware.Recover())

	vHosts["blog.localhost:1323"] = blog

	blog.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Blog")
	})

	//---------
	// Website
	//---------

	site := echo.New()
	site.Use(middleware.RequestLogger())
	site.Use(middleware.Recover())

	vHosts["localhost:1323"] = site

	site.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Website")
	})

	e := echo.NewVirtualHostHandler(vHosts)

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
