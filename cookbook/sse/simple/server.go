package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.File("/", "./index.html")

	e.GET("/sse", func(c *echo.Context) error {
		log.Printf("SSE client connected, ip: %v", c.RealIP())

		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		count := uint64(0)
		for {
			select {
			case <-c.Request().Context().Done():
				log.Printf("SSE client disconnected, ip: %v", c.RealIP())
				return nil
			case <-ticker.C:
				count++
				event := Event{
					Data: []byte(fmt.Sprintf("count: %d, time: %s\n\n", count, time.Now().Format(time.RFC3339Nano))),
				}
				if err := event.MarshalTo(w); err != nil {
					return err
				}
				if err := http.NewResponseController(w).Flush(); err != nil {
					return err
				}
			}
		}
	})

	sc := echo.StartConfig{
		Address: ":8080",
		BeforeServeFunc: func(s *http.Server) error {
			s.WriteTimeout = 0 // IMPORTANT: disable for SSE
			return nil
		},
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // start shutdown process on ctrl+c
	defer cancel()

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
