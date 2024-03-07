package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
)

// TaskFunc is a errorgroup.Go compatible function signature
type TaskFunc func() error

func main() {
	// create a context that automatically cancels thenever one of the
	// configured os signals is received
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	// HTTP servers
	e1 := NewHTTPServer()
	e2 := NewHTTPServer()

	// an errorgroup runs multiple goroutines, and allows all goroutines to be
	// stopped whenever one of the goroutines fails/errors
	group, gctx := errgroup.WithContext(ctx)
	group.Go(NewStartServerTask(e1, ":8888"))
	group.Go(NewStartServerTask(e2, ":8889"))
	group.Go(NewWatchDogTask(gctx, stop, e1, e2, 10*time.Second))

	// simply wait for all tasks in the group to be done
	if err := group.Wait(); err != nil && err != context.Canceled && err != http.ErrServerClosed {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println("exit.")
}

// NewHTTPServer returns a dummy echo instance
func NewHTTPServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		time.Sleep(3 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})
	// please note we can't use Logger.Fatal, because echo has implemented that
	// with an os.Exit(1), which prevents any graceful handling from code
	e.GET("/panic", func(c echo.Context) error {
		e.Logger.Error("something went horribly wrong..")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		return e.Shutdown(ctx)
	})
	return e
}

// NewStartServerTask returns a TaskFunc that starts/runs a HTTP server
func NewStartServerTask(e *echo.Echo, listenAddr string) TaskFunc {
	return func() error {
		return e.Start(listenAddr)
	}
}

// NewWatchDogTask returns the watchdog task that is responsible for doing
// a graceful shutdown of the HTTP servers when an os interrupt is received
func NewWatchDogTask(ctx context.Context, stop context.CancelFunc, e1, e2 *echo.Echo, timeout time.Duration) TaskFunc {
	return func() error {
		// waitgroup for parallellising shutdown of multiple servers
		var wg sync.WaitGroup

		// wait for context cancellation (== os interrupt, e.g. ctrl-C)
		<-ctx.Done()
		stop()
		fmt.Println("signal received: shutting down servers")

		wg.Add(1)
		go TryGracefulShutdown(&wg, e1, timeout)
		wg.Add(1)
		go TryGracefulShutdown(&wg, e2, timeout)

		// wait for shutdown routines to finish
		wg.Wait()
		fmt.Println("shutdown complete")
		return nil
	}
}

// TryGracefulShutdown tries to shutdown a server, giving it a maximum timeout
// for finishing the shutdown. Any shutdown errors are logged.
func TryGracefulShutdown(wg *sync.WaitGroup, e *echo.Echo, timeout time.Duration) {
	// decrement waitgroup to signal this gorouting has finished when we return
	defer wg.Done()

	// set context with a timeout to use for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// try graceful shutdown
	e.Logger.Warn("received shutdown signal")
	if err := e.Shutdown(ctx); err != nil {
		fmt.Printf("error shutting down server: %v\n", err)
	}
}
