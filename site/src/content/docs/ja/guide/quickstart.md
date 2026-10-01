---
title: クイックスタート
description: 5 分以内に本番対応の Echo API を構築します。
sidebar:
  order: 1
---

Echo は高性能でミニマルな Go Web フレームワークです。このガイドでは、5 分以内に
サーバーを起動します。

## 要件

Echo には **Go 1.25 以降**が必要です。バージョンを確認してください。

```bash
go version
```

## インストール

module を作成し、Echo を追加します。

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```

## Hello, World

`main.go` を作成します。

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

実行します。

```bash
go mod tidy
go run main.go
```

サーバーは `http://localhost:1323` で動作しています。Echo のルーターは、
ルートごとのリクエストディスパッチを**動的メモリ割り当てゼロ**で行います。

:::tip[Ask Echo]
困ったら右下の **Ask Echo** ボタンを押して
*"How do I add JWT auth?"* と質問してください。回答はこのドキュメントから直接得られます。
:::

## 次のステップ

- [ルーティング](/ja/guide/routing/)：静的、パラメーター付き、ワイルドカードのルート。
- [コンテキスト](/ja/guide/context/)：リクエストごとのリクエスト/レスポンスオブジェクト。
- [バインディング](/ja/guide/binding/)：リクエストデータを型付き struct に解析します。
