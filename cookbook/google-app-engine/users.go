package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	usersMap struct {
		sync.RWMutex
		users map[string]user
	}
)

var (
	usersStore usersMap
)

func init() {
	usersStore.users = map[string]user{
		"1": user{
			ID:   "1",
			Name: "Wreck-It Ralph",
		},
	}

	// hook into the echo instance to create an endpoint group
	// and add specific middleware to it plus handlers
	g := e.Group("/users")
	g.Use(middleware.CORS())

	g.POST("", createUser)
	g.GET("", getUsers)
	g.GET("/:id", getUser)
}

func createUser(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	usersStore.Lock()
	defer usersStore.Unlock()
	usersStore.users[u.ID] = *u
	return c.JSON(http.StatusCreated, u)
}

func getUsers(c echo.Context) error {
	usersStore.RLock()
	defer usersStore.RUnlock()
	return c.JSON(http.StatusOK, usersStore.users)
}

func getUser(c echo.Context) error {
	usersStore.RLock()
	defer usersStore.RUnlock()
	return c.JSON(http.StatusOK, usersStore.users[c.Param("id")])
}
