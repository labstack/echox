---
title: CORS
description: 使用允许列表或自定义 origin 函数启用 Cross-Origin Resource Sharing。
sidebar:
  order: 4
---

[CORS 中间件](/zh-cn/middleware/cors/)控制哪些 origin 可以访问你的 API。你可以传入固定的允许
origin 列表，也可以提供一个按请求决定的函数。

## origin 允许列表

把允许的 origin 直接传给 `middleware.CORS`。

```go file=cookbook/cors/origin-list/server.go
```

## 自定义 origin 函数

对于动态策略，请使用带 `UnsafeAllowOriginFunc` 的 `CORSWithConfig`。该函数接收请求上下文和
origin，并返回要回显的 origin、请求是否允许，以及可选错误。

```go file=cookbook/cors/origin-func/server.go
```
