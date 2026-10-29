---
title: インストール
description: Echo v5 を Go module に追加します。
sidebar:
  order: 2
---

Echo は Go module として配布されています：`github.com/labstack/echo/v5`。

## 要件

Echo v5 には **Go 1.25 以降**が必要です。

```bash
go version
```

## プロジェクトに追加する

既存の module 内では次を実行します。

```bash
go get github.com/labstack/echo/v5
```

または新しい module を開始します。

```bash
mkdir myapp && cd myapp
go mod init myapp
go get github.com/labstack/echo/v5
```

コードで import します。

```go
import "github.com/labstack/echo/v5"
```

## バージョン

| バージョン | インポートパス                | ステータス |
| ---------- | ----------------------------- | ---------- |
| **v5**     | `github.com/labstack/echo/v5` | 現行 |
| v4         | `github.com/labstack/echo/v4` | バグ修正は 2026-12-31 まで、セキュリティ修正は 2027-12-31 まで |

日付と対象範囲は [Echo のセキュリティポリシー](https://github.com/labstack/echo/blob/master/SECURITY.md) を参照してください。

:::note
Echo は[セマンティックインポートバージョニング](https://go.dev/blog/v2-go-modules)に従います。
メジャーバージョンはインポートパスの一部なので、移行中に v4 と v5 を共存させられます。
:::

## 最新に保つ

```bash
go get github.com/labstack/echo/v5
go mod tidy
```
