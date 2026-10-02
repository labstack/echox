package main

import (
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"golang.org/x/net/websocket"
)

func TestHelloReturnsAfterClientDisconnect(t *testing.T) {
	e := echo.New()
	e.Logger = slog.New(slog.DiscardHandler)
	handlerReturned := make(chan struct{})
	e.GET("/ws", func(c *echo.Context) error {
		err := hello(c)
		close(handlerReturned)
		return err
	})

	server := httptest.NewServer(e)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, err := websocket.Dial(wsURL, "", server.URL)
	if err != nil {
		server.Close()
		t.Fatalf("failed to connect websocket client: %v", err)
	}
	if err := ws.Close(); err != nil {
		server.CloseClientConnections()
		t.Fatalf("failed to close websocket client: %v", err)
	}

	select {
	case <-handlerReturned:
		server.Close()
	case <-time.After(2 * time.Second):
		t.Fatal("websocket handler did not return after client disconnected")
	}
}
