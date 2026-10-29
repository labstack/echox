---
title: Instalación
description: Agrega Echo v5 a tu módulo Go.
sidebar:
  order: 2
---

Echo se distribuye como un módulo Go: `github.com/labstack/echo/v5`.

## Requisitos

Echo v5 requiere **Go 1.25 o posterior**.

```bash
go version
```

## Agregar a un proyecto

Dentro de un módulo existente:

```bash
go get github.com/labstack/echo/v5
```

O inicia un módulo nuevo:

```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/labstack/echo/v5
```

Impórtalo en tu código:

```go
import "github.com/labstack/echo/v5"
```

## Versiones

| Versión | Import path                     | Estado |
| ------- | ------------------------------- | ------ |
| **v5**  | `github.com/labstack/echo/v5`   | Actual |
| v4      | `github.com/labstack/echo/v4`   | Correcciones de errores hasta 2026-12-31; correcciones de seguridad hasta 2027-12-31 |

Fechas y alcance: [política de seguridad de Echo](https://github.com/labstack/echo/blob/master/SECURITY.md).

:::note
Echo sigue [semantic import versioning](https://go.dev/blog/v2-go-modules): la
versión major forma parte del import path, por lo que v4 y v5 pueden coexistir durante una migración.
:::

## Mantenerse actualizado

```bash
go get github.com/labstack/echo/v5
go mod tidy
```
