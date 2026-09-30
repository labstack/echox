---
title: 嵌入资源
description: 使用 Go 的 embed 包提供打包进二进制文件的静态资源。
sidebar:
  order: 5
---

Go 的 `embed` 包（Go 1.16+）允许你把静态资源直接编译进二进制文件，因此单个可执行文件即可携带前端。
此示例通过 Echo 提供嵌入式文件系统，并支持一个可选的 live 模式，在开发期间从磁盘读取。

## 服务器

```go file=cookbook/embed/server.go
```

:::tip
使用 `live` 参数运行二进制文件（`go run server.go live`），可以从磁盘上的 `app` 目录提供资源，
而不是使用嵌入副本，这在开发期间很方便。
:::
