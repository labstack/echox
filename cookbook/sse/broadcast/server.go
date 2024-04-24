package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"
	"log"
	"net/http"
	"time"
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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.File("/", "./index.html")

	//e.GET("/sse", echo.WrapHandler(server))

	e.GET("/sse", func(c echo.Context) error { // longer variant with disconnect logic
		log.Printf("The client is connected: %v\n", c.RealIP())
		go func() {
			<-c.Request().Context().Done() // Received Browser Disconnection
			log.Printf("The client is disconnected: %v\n", c.RealIP())
			return
		}()

		server.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
