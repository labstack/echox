---
title: Quickstart
description: Build a production-ready Echo API in under five minutes.
sidebar:
  order: 1
---

Echo is a high performance, minimalist Go web framework. This guide gets a server
running in under five minutes.

## Requirements

Echo requires **Go 1.25 or newer**. Check your version:

```bash
go version
```

## Install

Create a module and add Echo:

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```

## Hello, World

Create `main.go`:

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

Run it:

```bash
go mod tidy
go run main.go
```

Your server is live at `http://localhost:1323`. Echo's router dispatches requests
with **zero dynamic memory allocation** per route.

:::tip[Ask Echo]
Stuck? Press the **Ask Echo** button (bottom-right) and ask
*"How do I add JWT auth?"* — answers come straight from these docs.
:::

## Next steps

- [Routing](/guide/routing/) — static, parameterized, and wildcard routes.
- [Context](/guide/context/) — the per-request request/response object.
- [Binding](/guide/binding/) — parse request data into typed structs.
