package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGracefulShutdown(t *testing.T) {
	// Skip in short mode as this test takes time
	if testing.Short() {
		t.Skip("Skipping graceful shutdown test in short mode")
	}

	// Setup Echo server
	e := echo.New()
	e.HideBanner = true

	// Add the test endpoint that sleeps
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := e.Start(":0"); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait a moment for server to start
	time.Sleep(100 * time.Millisecond)

	// Get the actual port the server is listening on
	address := e.Listener.Addr().String()

	// Start a request that will take 5 seconds
	requestDone := make(chan bool, 1)
	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(fmt.Sprintf("http://%s/", address))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				requestDone <- true
				return
			}
		}
		requestDone <- false
	}()

	// Wait a bit then initiate shutdown
	time.Sleep(100 * time.Millisecond)

	// Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownErr := make(chan error, 1)
	go func() {
		shutdownErr <- e.Shutdown(ctx)
	}()

	// Wait for either server error or shutdown completion
	select {
	case err := <-serverErr:
		t.Fatalf("Server failed to start: %v", err)
	case err := <-shutdownErr:
		require.NoError(t, err, "Server shutdown should not error")
	case <-time.After(15 * time.Second):
		t.Fatal("Test timeout - shutdown took too long")
	}

	// Verify the long-running request completed successfully
	select {
	case success := <-requestDone:
		assert.True(t, success, "Long-running request should complete successfully during graceful shutdown")
	case <-time.After(2 * time.Second):
		t.Error("Long-running request did not complete in time")
	}
}

func TestServerBasicFunctionality(t *testing.T) {
	// Setup Echo server
	e := echo.New()
	e.HideBanner = true

	// Add the endpoint from the main example
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	// Start server
	go func() {
		e.Start(":0")
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test basic functionality (without waiting for full response)
	address := e.Listener.Addr().String()
	client := &http.Client{Timeout: 1 * time.Second}

	// This will timeout, but we're just testing that the connection is established
	_, err := client.Get(fmt.Sprintf("http://%s/", address))

	// We expect a timeout error since the handler sleeps for 5 seconds
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}

func TestShutdownWithoutActiveRequests(t *testing.T) {
	// Setup Echo server
	e := echo.New()
	e.HideBanner = true

	e.GET("/quick", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Quick response")
	})

	// Start server
	go func() {
		e.Start(":0")
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Make a quick request
	address := e.Listener.Addr().String()
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s/quick", address))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Shutdown should be quick with no active requests
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = e.Shutdown(ctx)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Less(t, duration, 1*time.Second, "Shutdown without active requests should be quick")
}