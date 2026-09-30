---
title: CRUD
description: Create, read, update, and delete resources with Echo and JSON binding.
sidebar:
  order: 2
---

A complete CRUD (create, read, update, delete) API backed by an in-memory store.
Each handler binds the JSON request body into a struct, mutates the store under a
lock, and returns the result as JSON.

## Server

```go
package main

import (
	"context"
	"net/http"
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
	users = map[int]*user{}
	seq   = 1
	lock  = sync.Mutex{}
)

//----------
// Handlers
//----------

func createUser(c *echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	// Assign the ID after binding so an "id" in the request body cannot overwrite another user.
	u.ID = seq
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getUser(c *echo.Context) error {
	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return err // 400 Bad Request for a non-numeric id
	}
	lock.Lock()
	defer lock.Unlock()
	u, ok := users[id]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, u)
}

func updateUser(c *echo.Context) error {
	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return err
	}
	// Bind before taking the lock so a slow request body does not block other requests.
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	existingUser, ok := users[id]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	existingUser.Name = u.Name
	return c.JSON(http.StatusOK, existingUser)
}

func deleteUser(c *echo.Context) error {
	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := users[id]; !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, users)
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
```

## Client

### Create user

Request:

```sh
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe Smith"}' \
  localhost:1323/users
```

Response:

```json
{
  "id": 1,
  "name": "Joe Smith"
}
```

### Get user

Request:

```sh
curl localhost:1323/users/1
```

Response:

```json
{
  "id": 1,
  "name": "Joe Smith"
}
```

### Update user

Request:

```sh
curl -X PUT \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe"}' \
  localhost:1323/users/1
```

Response:

```json
{
  "id": 1,
  "name": "Joe"
}
```

### Delete user

Request:

```sh
curl -X DELETE localhost:1323/users/1
```

Response: `204 No Content`.
