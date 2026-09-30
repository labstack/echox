---
title: Hello World
description: あいさつを返す最小構成の Echo サーバーです。
sidebar:
  order: 1
---

最小構成の Echo アプリケーションです。インスタンスを作成し、Logger と Recover ミドルウェアを登録し、
単一のルートを追加してサーバーを起動します。

## サーバー

```go file=cookbook/hello-world/server.go
```
