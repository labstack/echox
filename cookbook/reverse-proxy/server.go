package main

import (
	"context"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Setup proxy
	// docs:start targets
	url1, err := url.Parse("http://localhost:8081")
	if err != nil {
		e.Logger.Error("failed to parse url", "error", err)
		return
	}
	url2, err := url.Parse("http://localhost:8082")
	if err != nil {
		e.Logger.Error("failed to parse url", "error", err)
		return
	}
	targets := []*middleware.ProxyTarget{
		{URL: url1},
		{URL: url2},
	}
	// docs:end targets
	// docs:start middleware
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	// docs:end middleware

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

// configureBlogProxy limits the proxy middleware to a route group.
func configureBlogProxy(e *echo.Echo, targets []*middleware.ProxyTarget) {
	// docs:start grouped-proxy
	g := e.Group("/blog")
	g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	// docs:end grouped-proxy
}
