---
title: CRUD
description: Crea, lee, actualiza y elimina recursos con Echo y binding JSON.
sidebar:
  order: 2
---

Una API CRUD (create, read, update, delete) completa respaldada por un store en memoria.
Cada handler vincula el body JSON del request a un struct, muta el store bajo un lock y
devuelve el resultado como JSON.

## Servidor

```go
package main

import (
	"context"
	"net/http"
	"strconv"
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
	lock.Lock()
	defer lock.Unlock()
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getUser(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	existingUser, ok := users[id]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	existingUser.Name = u.Name
	return c.JSON(http.StatusOK, existingUser)
}

func deleteUser(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func getAllUsers(c *echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, users)
}

func main() {
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

	// Start server
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

## Cliente

### Crear usuario

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

### Obtener usuario

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

### Actualizar usuario

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

### Eliminar usuario

Request:

```sh
curl -X DELETE localhost:1323/users/1
```

Response: `204 No Content`.
