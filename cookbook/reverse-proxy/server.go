package main

import (
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Setup proxy
	url1, err := url.Parse("http://localhost:8081")
	if err != nil {
		e.Logger.Fatal(err)
	}
	url2, err := url.Parse("http://localhost:8082")
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Use(middleware.Proxy(&middleware.RoundRobinBalancer{
		Targets: []*middleware.ProxyTarget{
			&middleware.ProxyTarget{
				URL: url1,
			},
			&middleware.ProxyTarget{
				URL: url2,
			},
		},
	}))

	e.Logger.Fatal(e.Start(":1323"))
}
