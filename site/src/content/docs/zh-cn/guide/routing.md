---
title: 路由
description: 在 Echo 的零分配基数树上将请求 URL 匹配到处理函数。
sidebar:
  order: 3
---

Echo 优化过的路由器使用基数树（radix tree）将请求 URL 匹配到处理函数，**每个路由零动态内存分配**，并具有智能的路由优先级。

## 注册路由

在 `Echo` 实例上使用 HTTP 方法辅助函数。每个函数接受一个路径模式和一个
`HandlerFunc`（`func(c *echo.Context) error`），并可选地附带路由级中间件。

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` 为所有受支持的方法注册处理函数，`Match` 则用于指定的方法集合：

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

## 匹配类型

| 模式             | 类型     | 匹配示例              |
| ---------------- | -------- | --------------------- |
| `/users/profile` | 静态     | `/users/profile`      |
| `/users/:id`     | 参数     | `/users/42`           |
| `/static/*`      | 通配符   | `/static/css/app.css` |

:::note
优先级为 **静态 → 参数 → 通配符**，因此 `/users/profile` 始终优先于
`/users/:id`，而后者又优先于 `/users/*`。
:::

## 路径参数

使用 `c.Param()` 从上下文读取命名参数（或使用 `c.ParamOr()` 提供默认值）：

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

通配符段可通过 `*` 参数获取：

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
