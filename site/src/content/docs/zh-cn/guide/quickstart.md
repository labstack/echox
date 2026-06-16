---
title: 快速开始
description: 在五分钟内构建一个可用于生产的 Echo API。
sidebar:
  order: 1
---

Echo 是一个高性能、极简的 Go Web 框架。本指南将在五分钟内让服务器运行起来。

## 环境要求

Echo 需要 **Go 1.25 或更高版本**。检查你的版本：

```bash
go version
```

## 安装

创建一个模块并添加 Echo：

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

运行：

```bash
go run main.go
```

服务器已在 `http://localhost:1323` 上运行。Echo 的路由器在分发请求时**每个路由零动态内存分配**。

:::tip[Ask Echo]
遇到问题？点击右下角的 **Ask Echo** 按钮，提问
*“如何添加 JWT 认证？”*——答案直接来自这些文档。
:::

## 下一步

- [路由](/zh-cn/guide/routing/) —— 静态、带参数和通配符路由。
- [上下文](/zh-cn/guide/context/) —— 每个请求的请求/响应对象。
- [绑定](/zh-cn/guide/binding/) —— 将请求数据解析为类型化的结构体。
