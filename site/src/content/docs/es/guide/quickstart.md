---
title: Inicio rápido
description: Crea una API Echo lista para producción en menos de cinco minutos.
sidebar:
  order: 1
---

Echo es un framework web Go minimalista y de alto rendimiento. Esta guía pone un servidor
en marcha en menos de cinco minutos.

## Requisitos

Echo requiere **Go 1.25 o posterior**. Comprueba tu versión:

```bash
go version
```

## Instalar

Crea un módulo y agrega Echo:

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```

## Hello, World

Crea `main.go`:

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

Ejecútalo:

```bash
go mod tidy
go run main.go
```

Tu servidor está disponible en `http://localhost:1323`. El router de Echo despacha requests
con **cero asignaciones dinámicas de memoria** por ruta.

:::tip[Ask Echo]
¿Atascado? Presiona el botón **Ask Echo** (abajo a la derecha) y pregunta
*"How do I add JWT auth?"*: las respuestas salen directamente de esta documentación.
:::

## Próximos pasos

- [Routing](/es/guide/routing/): rutas estáticas, parametrizadas y wildcard.
- [Context](/es/guide/context/): el objeto request/response por request.
- [Binding](/es/guide/binding/): analiza datos de request en structs tipados.
