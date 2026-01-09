package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	e.Static("/", "static")
	e.GET("/", func(c *echo.Context) (err error) {
		rw, err := echo.UnwrapResponse(c.Response())
		if err != nil {
			return
		}
		if pusher, ok := rw.ResponseWriter.(http.Pusher); ok {
			if err = pusher.Push("/app.css", nil); err != nil {
				return
			}
			if err = pusher.Push("/app.js", nil); err != nil {
				return
			}
			if err = pusher.Push("/echo.png", nil); err != nil {
				return
			}
		}
		return c.File("index.html")
	})

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
