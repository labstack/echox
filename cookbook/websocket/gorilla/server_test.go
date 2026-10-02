package main

import (
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
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
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect websocket client: %v", err)
	}
	if err := ws.Close(); err != nil {
		t.Fatalf("failed to close websocket client: %v", err)
	}

	select {
	case <-handlerReturned:
	case <-time.After(2 * time.Second):
		t.Fatal("websocket handler did not return after client disconnected")
	}
}
