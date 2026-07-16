---
title: Routing
description: Hace coincidir URLs de request con handlers en el árbol radix sin asignaciones de Echo.
sidebar:
  order: 3
---

El router optimizado de Echo hace coincidir URLs de request con handlers usando un árbol radix con
**cero asignaciones dinámicas de memoria** y priorización inteligente de rutas.

## Registrar rutas

Usa los helpers de métodos HTTP en la instancia `Echo`. Cada uno recibe un patrón de path y una
`HandlerFunc` (`func(c *echo.Context) error`), con middleware opcional a nivel de ruta.

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` registra un handler para cualquier método HTTP — incluidos los que no están en
la lista predefinida de Echo — y `Match` para un conjunto específico:

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

:::note
Una ruta de método específico registrada para el mismo path tiene precedencia sobre la ruta `Any`.
:::

## Tipos de coincidencia

| Patrón             | Tipo     | Ejemplo de coincidencia    |
| ------------------ | -------- | -------------------------- |
| `/users/profile`   | Static   | `/users/profile`           |
| `/users/:id`       | Param    | `/users/42`                |
| `/static/*`        | Wildcard | `/static/css/app.css`      |

:::note
La prioridad es **static → param → wildcard**, por lo que `/users/profile` siempre gana sobre
`/users/:id`, que a su vez gana sobre `/users/*`.
:::

## Parámetros de path

Lee parámetros nombrados desde el contexto con `c.Param()` (o `c.ParamOr()` para un valor por defecto):

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

El segmento wildcard está disponible como el parámetro `*`:

```go
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, c.Param("*"))
})
```

## Grupos

Agrupa rutas que comparten un prefijo y middleware con `e.Group()`:

```go
admin := e.Group("/admin", middleware.BasicAuth(authFn))
admin.GET("/metrics", metrics)   // -> /admin/metrics
admin.GET("/users", listUsers)   // -> /admin/users
```

Los grupos pueden anidarse para componer árboles de rutas más grandes.
