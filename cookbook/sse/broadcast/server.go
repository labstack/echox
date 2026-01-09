package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/r3labs/sse/v2"
)

func main() {
	e := echo.New()

	server := sse.New()             // create SSE broadcaster server
	server.AutoReplay = false       // do not replay messages for each new subscriber that connects
	_ = server.CreateStream("time") // EventSource in "index.html" connecting to stream named "time"

	go func(s *sse.Server) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Publish("time", &sse.Event{
					Data: []byte("time: " + time.Now().Format(time.RFC3339Nano)),
				})
			}
		}
	}(server)

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.File("/", "./index.html")

	//e.GET("/sse", echo.WrapHandler(server))

	e.GET("/sse", func(c *echo.Context) error { // longer variant with disconnect logic
		e.Logger.Info("New client connected", "ip", c.RealIP())
		go func() {
			<-c.Request().Context().Done() // Received Browser Disconnection
			e.Logger.Info("Client disconnected", "ip", c.RealIP())
			return
		}()

		server.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Error("shutting down the server", "error", err)
	}
}
