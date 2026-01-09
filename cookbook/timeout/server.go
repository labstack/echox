package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.ContextTimeout(5 * time.Second))

	// Route => handler
	e.GET("/", func(c *echo.Context) error {
		select {
		case <-c.Request().Context().Done():
			return echo.NewHTTPError(http.StatusRequestTimeout, "Request timed out")
		case <-time.After(10 * time.Second):
			return c.String(http.StatusOK, "Hello, World!\n")
		}
	})

	// Start server
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
