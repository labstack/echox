---
title: 子域名
description: 使用虚拟主机处理函数，按 host 将请求路由到不同 Echo 实例。
sidebar:
  order: 17
---

此示例会根据请求 host 将请求路由到独立的 `Echo` 实例，使每个子域名拥有自己的路由和中间件。
这些实例通过 `echo.NewVirtualHostHandler` 组合，并按 host name 分发。

## 服务器

```go file=cookbook/subdomain/server.go
```
