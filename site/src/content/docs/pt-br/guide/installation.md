---
title: Instalação
description: Adicione o Echo v5 ao seu módulo Go.
sidebar:
  order: 2
---

O Echo é distribuído como um módulo Go: `github.com/labstack/echo/v5`.

## Requisitos

O Echo v5 requer **Go 1.25 ou mais recente**.

```bash
go version
```

## Adicionar a um projeto

Dentro de um módulo existente:

```bash
go get github.com/labstack/echo/v5
```

Ou inicie um novo módulo:

```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/labstack/echo/v5
```

Importe-o no seu código:

```go
import "github.com/labstack/echo/v5"
```

## Versões

| Versão | Caminho de importação          | Status |
| ------ | ------------------------------ | ------ |
| **v5** | `github.com/labstack/echo/v5`  | Atual |
| v4     | `github.com/labstack/echo/v4`  | Correções de bugs até 2026-12-31; correções de segurança até 2027-12-31 |

Datas e escopo: [política de segurança do Echo](https://github.com/labstack/echo/blob/master/SECURITY.md).

:::note
O Echo segue [semantic import versioning](https://go.dev/blog/v2-go-modules) — a
versão principal faz parte do caminho de importação, então v4 e v5 podem coexistir durante uma migração.
:::

## Manter atualizado

```bash
go get github.com/labstack/echo/v5
go mod tidy
```
