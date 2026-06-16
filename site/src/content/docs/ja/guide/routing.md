---
title: ルーティング
description: Echo のゼロアロケーション基数木で、リクエスト URL をハンドラにマッチさせます。
sidebar:
  order: 3
---

Echo の最適化されたルーターは、基数木（radix tree）を用いてリクエスト URL をハンドラにマッチさせます。**ルートごとに動的なメモリ割り当てはゼロ**で、賢いルート優先順位付けを備えています。

## ルートの登録

`Echo` インスタンスの HTTP メソッドヘルパーを使います。各ヘルパーはパスパターンと
`HandlerFunc`（`func(c *echo.Context) error`）を受け取り、任意でルートレベルのミドルウェアを指定できます。

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` はサポートされるすべてのメソッドに、`Match` は指定したメソッド集合にハンドラを登録します：

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

## マッチの種類

| パターン         | 種類           | マッチ例              |
| ---------------- | -------------- | --------------------- |
| `/users/profile` | 静的           | `/users/profile`      |
| `/users/:id`     | パラメータ     | `/users/42`           |
| `/static/*`      | ワイルドカード | `/static/css/app.css` |

:::note
優先順位は **静的 → パラメータ → ワイルドカード** です。したがって `/users/profile` は常に
`/users/:id` より優先され、さらにそれは `/users/*` より優先されます。
:::

## パスパラメータ

`c.Param()` でコンテキストから名前付きパラメータを読み取ります（デフォルト値が必要な場合は `c.ParamOr()`）：

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

ワイルドカードのセグメントは `*` パラメータとして取得できます：

```go
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, c.Param("*"))
})
```

## グループ

接頭辞とミドルウェアを共有するルートを `e.Group()` でまとめます：

```go
admin := e.Group("/admin", middleware.BasicAuth(authFn))
admin.GET("/metrics", metrics)   // -> /admin/metrics
admin.GET("/users", listUsers)   // -> /admin/users
```

グループはネストして、より大きなルートツリーを構成できます。
