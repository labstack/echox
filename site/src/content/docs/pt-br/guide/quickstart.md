---
title: Início rápido
description: Crie uma API Echo pronta para produção em menos de cinco minutos.
sidebar:
  order: 1
---

Echo é um framework web Go minimalista e de alta performance. Este guia coloca um
servidor em execução em menos de cinco minutos.

## Requisitos

Echo requer **Go 1.25 ou mais recente**. Verifique sua versão:

```bash
go version
```

## Instalar

Crie um módulo e adicione o Echo:

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```

## Hello, World

Crie `main.go`:

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

Execute:

```bash
go mod tidy
go run main.go
```

Seu servidor está ativo em `http://localhost:1323`. O router do Echo despacha requests
com **zero alocação dinâmica de memória** por rota.

:::tip[Ask Echo]
Travou? Pressione o botão **Ask Echo** (canto inferior direito) e pergunte
*"How do I add JWT auth?"* — as respostas vêm diretamente destes docs.
:::

## Próximos passos

- [Routing](/pt-br/guide/routing/) — rotas estáticas, parametrizadas e com wildcard.
- [Context](/pt-br/guide/context/) — o objeto request/response por request.
- [Binding](/pt-br/guide/binding/) — analise dados do request em structs tipadas.
