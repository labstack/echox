---
title: ルーティング
description: Echo のゼロ割り当て radix tree でリクエスト URL をハンドラにマッチさせます。
sidebar:
  order: 3
---

Echo の最適化されたルーターは、radix tree を使ってリクエスト URL をハンドラへマッチさせ、
**動的メモリ割り当てゼロ**とスマートなルート優先順位を実現します。

## ルートを登録する

`Echo` インスタンス上の HTTP メソッドヘルパーを使います。各ヘルパーはパスパターンと
`HandlerFunc`（`func(c *echo.Context) error`）を受け取り、任意でルートレベルの
ミドルウェアも受け取れます。

```go
e := echo.New()

e.GET("/users/:id", getUser)      // named parameter
e.POST("/users", createUser)
e.PUT("/users/:id", updateUser)
e.DELETE("/users/:id", deleteUser)
e.GET("/static/*", serveFiles)    // wildcard
```

`Any` は任意の HTTP メソッド（Echo にあらかじめ定義された一覧にないものも含む）にハンドラを
登録し、`Match` は指定した集合に登録します。

```go
e.Any("/ping", pong)
e.Match([]string{http.MethodGet, http.MethodPost}, "/form", handleForm)
```

:::note
同じパスに登録されたメソッド固有のルートは、`Any` ルートより優先されます。
:::

## マッチ種別

| パターン           | 種別         | マッチ例                   |
| ------------------ | ------------ | -------------------------- |
| `/users/profile`   | 静的         | `/users/profile`           |
| `/users/:id`       | パラメーター | `/users/42`                |
| `/static/*`        | ワイルドカード | `/static/css/app.css`    |

:::note
優先順位は**静的 → パラメーター → ワイルドカード**です。そのため
`/users/profile` は常に `/users/:id` に勝ち、`/users/:id` は `/users/*` に勝ちます。
:::

## パスパラメーター

`c.Param()` でコンテキストから名前付きパラメーターを読み取ります
（デフォルト値には `c.ParamOr()` を使います）。

```go
func getUser(c *echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
```

ワイルドカードセグメントは `*` パラメーターとして取得できます。

```go
e.GET("/files/*", func(c *echo.Context) error {
	return c.String(http.StatusOK, c.Param("*"))
})
```

## グループ

`e.Group()` を使い、プレフィックスとミドルウェアを共有するルートをグループ化します。

```go
admin := e.Group("/admin", middleware.BasicAuth(authFn))
admin.GET("/metrics", metrics)   // -> /admin/metrics
admin.GET("/users", listUsers)   // -> /admin/users
```

グループはネストでき、より大きなルートツリーを構成できます。
