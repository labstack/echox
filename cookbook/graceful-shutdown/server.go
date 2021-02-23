package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	// start background processes or anything that needs to be opened/started out of Echo
	// and before Echo starts serving requests

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	go mustGracefullyShutdown(e)

	// Start blocks until Echo finishes serving ongoing requests and shuts down or fails to start with an error
	err := e.Start(":1323")
	if err != nil && err != http.ErrServerClosed { // http.ErrServerClosed indicates normal Echo shutdown by e.Shutdown
		e.Logger.Errorf("Shutting down the server due error: %v", err)
	}

	// stop background processes or anything that was opened/started out of Echo
}

func mustGracefullyShutdown(e *echo.Echo) {
	var shutdownSig = make(chan os.Signal, 1)
	signal.Notify(shutdownSig, os.Interrupt) // catch interrupt signal (ctrl+C)
	<-shutdownSig

	e.Logger.Info("Received shutdown signal. Starting gracefully to shut down Echo server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err == context.DeadlineExceeded {
		// Fatal() terminates program with os.Exit(1)
		// NB: this does not close anything cleanly
		e.Logger.Fatal("Graceful shutdown timeout exceeded. Exiting program forcefully!")
	}
}
