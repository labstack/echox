package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"golang.org/x/net/websocket"
)

func hello(c *echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			if err := websocket.Message.Send(ws, "Hello, Client!"); err != nil {
				c.Logger().Error("failed to write WS message", "error", err)
			}

			// Read
			msg := ""
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				c.Logger().Error("failed to write WS message", "error", err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Static("/", "../public")
	e.GET("/ws", hello)

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
