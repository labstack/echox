---
title: 快速开始
description: 在五分钟内构建可用于生产的 Echo API。
sidebar:
  order: 1
---

Echo 是一个高性能、极简的 Go Web 框架。本指南会让服务器在五分钟内运行起来。

## 要求

Echo 需要 **Go 1.25 或更新版本**。检查你的版本：

```bash
go version
```

## 安装

创建一个 module 并添加 Echo：

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```

## Hello, World

创建 `main.go`：

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

运行它：

```bash
go mod tidy
go run main.go
```

你的服务器已经在 `http://localhost:1323` 上运行。Echo 的路由器在每条路由分发请求时实现
**零动态内存分配**。

:::tip[Ask Echo]
遇到问题？按下右下角的 **Ask Echo** 按钮并询问
*"How do I add JWT auth?"*，答案会直接来自这些文档。
:::

## 下一步

- [路由](/zh-cn/guide/routing/)：静态、参数化和通配符路由。
- [上下文](/zh-cn/guide/context/)：每个请求的请求/响应对象。
- [绑定](/zh-cn/guide/binding/)：将请求数据解析到带类型的 struct 中。
