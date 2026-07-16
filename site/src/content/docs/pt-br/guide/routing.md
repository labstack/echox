---
title: Routing
description: Associe URLs de request a handlers na árvore radix de alocação zero do Echo.
sidebar:
  order: 3
---

O router otimizado do Echo associa URLs de request a handlers usando uma árvore radix com
**zero alocação dinâmica de memória** e priorização inteligente de rotas.

## Registrar rotas

Use os helpers de método HTTP na instância `Echo`. Cada um recebe um padrão de caminho e um
`HandlerFunc` (`func(c *echo.Context) error`), com middleware opcional no nível da rota.

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` registra um handler para qualquer método HTTP — incluindo os que não estão na
lista predefinida do Echo — e `Match` para um conjunto específico:

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

:::note
Uma rota de método específico registrada para o mesmo caminho tem precedência sobre a rota `Any`.
:::

## Tipos de correspondência

| Padrão            | Tipo     | Exemplo correspondente     |
| ----------------- | -------- | -------------------------- |
| `/users/profile`  | Static   | `/users/profile`           |
| `/users/:id`      | Param    | `/users/42`                |
| `/static/*`       | Wildcard | `/static/css/app.css`      |

:::note
A prioridade é **static → param → wildcard**, então `/users/profile` sempre vence
`/users/:id`, que vence `/users/*`.
:::

## Parâmetros de caminho

Leia parâmetros nomeados do contexto com `c.Param()` (ou `c.ParamOr()` para um valor padrão):

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

O segmento wildcard fica disponível como o parâmetro `*`:

```go
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, c.Param("*"))
})
```

## Grupos

Agrupe rotas que compartilham prefixo e middleware com `e.Group()`:

```go
admin := e.Group("/admin", middleware.BasicAuth(authFn))
admin.GET("/metrics", metrics)   // -> /admin/metrics
admin.GET("/users", listUsers)   // -> /admin/users
```

Grupos podem ser aninhados para compor árvores de rotas maiores.
