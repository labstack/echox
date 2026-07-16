---
title: 路由
description: 在 Echo 零分配 radix tree 上将请求 URL 匹配到处理函数。
sidebar:
  order: 3
---

Echo 优化过的路由器使用 radix tree 将请求 URL 匹配到处理函数，实现
**零动态内存分配**和智能路由优先级。

## 注册路由

在 `Echo` 实例上使用 HTTP 方法辅助函数。每个函数都接受一个路径模式和一个
`HandlerFunc`（`func(c *echo.Context) error`），并可选接收路由级中间件。

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` 会为任意 HTTP 方法（包括不在 Echo 预定义列表中的方法）注册处理函数，`Match` 则用于一组指定方法：

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

:::note
为同一路径注册的特定方法路由优先于 `Any` 路由。
:::

## 匹配类型

| 模式               | 类型     | 匹配示例                   |
| ------------------ | -------- | -------------------------- |
| `/users/profile`   | 静态     | `/users/profile`           |
| `/users/:id`       | 参数     | `/users/42`                |
| `/static/*`        | 通配符   | `/static/css/app.css`      |

:::note
优先级是**静态 → 参数 → 通配符**，因此 `/users/profile` 总是优先于
`/users/:id`，而后者优先于 `/users/*`。
:::

## 路径参数

使用 `c.Param()` 从上下文读取命名参数（或使用 `c.ParamOr()` 提供默认值）：

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

通配符片段可通过 `*` 参数访问：

```go
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, c.Param("*"))
})
```

## 路由组

使用 `e.Group()` 将共享前缀和中间件的路由分组：

```go
admin := e.Group("/admin", middleware.BasicAuth(authFn))
admin.GET("/metrics", metrics)   // -> /admin/metrics
admin.GET("/users", listUsers)   // -> /admin/users
```

路由组可以嵌套，以组合更大的路由树。
