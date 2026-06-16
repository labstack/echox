---
title: クイックスタート
description: 本番運用に耐える Echo API を5分以内で構築します。
sidebar:
  order: 1
---

Echo は高性能でミニマルな Go 製 Web フレームワークです。このガイドでは、5分以内でサーバーを起動します。

## 要件

Echo には **Go 1.25 以降** が必要です。バージョンを確認します：

```bash
go version
```

## インストール

モジュールを作成して Echo を追加します：

```bash
go mod init myapp
go get github.com/labstack/echo/v5
```

## Hello, World

`main.go` を作成します：

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

実行します：

```bash
go run main.go
```

サーバーは `http://localhost:1323` で起動しています。Echo のルーターは、ルートごとに**動的なメモリ割り当てゼロ**でリクエストをディスパッチします。

:::tip[Ask Echo]
お困りですか？右下の **Ask Echo** ボタンを押して
*「JWT 認証を追加するには？」* と質問してください。回答はこのドキュメントから直接得られます。
:::

## 次のステップ

- [ルーティング](/ja/guide/routing/) —— 静的・パラメータ付き・ワイルドカードのルート。
- [コンテキスト](/ja/guide/context/) —— リクエストごとのリクエスト/レスポンスオブジェクト。
- [バインディング](/ja/guide/binding/) —— リクエストデータを型付き構造体に解析します。
