---
title: Installation
description: Add Echo v5 to your Go module.
sidebar:
  order: 2
---

Echo is distributed as a Go module: `github.com/labstack/echo/v5`.

## Requirements

Echo v5 requires **Go 1.25 or newer**.

```bash
go version
```

## Add to a project

Inside an existing module:

```bash
go get github.com/labstack/echo/v5
```

Or start a new module:

```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/labstack/echo/v5
```

Import it in your code:

```go
import "github.com/labstack/echo/v5"
```

## Versions

| Version | Import path                     | Status |
| ------- | ------------------------------- | ------ |
| **v5**  | `github.com/labstack/echo/v5`   | Current |
| v4      | `github.com/labstack/echo/v4`   | Bug fixes until 2026-12-31; security fixes until 2027-12-31 |

Dates and scope: [Echo security policy](https://github.com/labstack/echo/blob/master/SECURITY.md).

:::note
Echo follows [semantic import versioning](https://go.dev/blog/v2-go-modules) — the
major version is part of the import path, so v4 and v5 can coexist during a migration.
:::

## Keeping up to date

```bash
go get github.com/labstack/echo/v5
go mod tidy
```
