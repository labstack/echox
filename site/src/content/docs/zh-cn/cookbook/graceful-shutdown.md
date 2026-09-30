---
title: 优雅关闭
description: 在中断信号停止服务器之前排空正在处理的请求。
sidebar:
  order: 8
---

优雅关闭会让正在处理的请求在进程退出前完成。最简单的方法是把可取消的上下文传给
`StartConfig.Start`，并设置 `GracefulTimeout`。当上下文被中断信号取消时，Echo 会停止接受新连接，
并最多等待到超时时间，让活跃请求完成。

## 服务器

```go file=cookbook/graceful-shutdown/server.go#primary-server
```

## 使用自定义 HTTP 服务器

如果你自己管理 `http.Server`，请在 goroutine 中启动它，等待 signal context，
然后带超时调用 `Shutdown`：

```go file=cookbook/graceful-shutdown/server.go#custom-server
```

```sh
go run . -custom-server
```
