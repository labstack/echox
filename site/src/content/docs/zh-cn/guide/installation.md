---
title: 安装
description: 将 Echo v5 添加到你的 Go module。
sidebar:
  order: 2
---

Echo 以 Go module 的形式发布：`github.com/labstack/echo/v5`。

## 要求

Echo v5 需要 **Go 1.25 或更新版本**。

```bash
go version
```

## 添加到项目

在现有 module 中：

```bash
go get github.com/labstack/echo/v5
```

或者启动一个新 module：

```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/labstack/echo/v5
```

在代码中导入它：

```go
import "github.com/labstack/echo/v5"
```

## 版本

| 版本 | 导入路径                      | 状态 |
| ---- | ----------------------------- | ---- |
| **v5** | `github.com/labstack/echo/v5` | 当前 |
| v4   | `github.com/labstack/echo/v4` | 错误修复至 2026-12-31，安全修复至 2027-12-31 |

日期和范围见 [Echo 安全策略](https://github.com/labstack/echo/blob/master/SECURITY.md)。

:::note
Echo 遵循[语义化导入版本](https://go.dev/blog/v2-go-modules)：
主版本是导入路径的一部分，因此 v4 和 v5 可以在迁移期间共存。
:::

## 保持最新

```bash
go get github.com/labstack/echo/v5
go mod tidy
```
