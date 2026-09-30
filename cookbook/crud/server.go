package main

import (
	"cmp"
	"context"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	users = map[int]user{}
	seq   = 1
	lock  = sync.Mutex{}
)

//----------
// Handlers
//----------

// The store holds user values, so handlers copy a user under the lock and write the response after unlocking.

func createUser(c *echo.Context) error {
	u, err := bindUser(c)
	if err != nil {
		return err
	}
	lock.Lock()
	// Assign the ID after binding so an "id" in the request body cannot overwrite another user.
	u.ID = seq
	users[u.ID] = u
	seq++
	lock.Unlock()
	return c.JSON(http.StatusCreated, u)
}

func getUser(c *echo.Context) error {
	id, err := userID(c)
	if err != nil {
		return err
	}
	lock.Lock()
	u, ok := users[id]
	lock.Unlock()
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, u)
}

func updateUser(c *echo.Context) error {
	id, err := userID(c)
	if err != nil {
		return err
	}
	u, err := bindUser(c)
	if err != nil {
		return err
	}
	lock.Lock()
	existing, ok := users[id]
	if ok {
		existing.Name = u.Name
		users[id] = existing
	}
	lock.Unlock()
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, existing)
}

func deleteUser(c *echo.Context) error {
	id, err := userID(c)
	if err != nil {
		return err
	}
	lock.Lock()
	_, ok := users[id]
	delete(users, id)
	lock.Unlock()
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c *echo.Context) error {
	lock.Lock()
	list := make([]user, 0, len(users)) // non-nil, so an empty store returns [] rather than null
	for _, u := range users {
		list = append(list, u)
	}
	lock.Unlock()
	slices.SortFunc(list, func(a, b user) int { return cmp.Compare(a.ID, b.ID) })
	return c.JSON(http.StatusOK, list)
}

// bindUser binds the JSON request body and requires a non-blank name.
func bindUser(c *echo.Context) (user, error) {
	var u user
	if err := c.Bind(&u); err != nil {
		return u, err
	}
	if strings.TrimSpace(u.Name) == "" {
		return u, echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	return u, nil
}

// userID returns the :id path parameter, or a 400 error when it is not a number.
func userID(c *echo.Context) (int, error) {
	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}
	return id, nil
}

func newServer() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	return e
}

func main() {
	e := newServer()

	// Start server
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
