---
title: CORS
description: allow list またはカスタム origin 関数で Cross-Origin Resource Sharing を有効にします。
sidebar:
  order: 4
---

[CORS ミドルウェア](/ja/middleware/cors/)は、どの origin が API にアクセスできるかを制御します。
許可する origin の固定リストを渡すことも、リクエストごとに判断する関数を指定することもできます。

## origin の allow list

許可する origin を `middleware.CORS` に直接渡します。

```go file=cookbook/cors/origin-list/server.go
```

## カスタム origin 関数

動的ポリシーには、`UnsafeAllowOriginFunc` を指定した `CORSWithConfig` を使います。
この関数はリクエストコンテキストと origin を受け取り、返す origin、リクエストを許可するか、
任意のエラーを返します。

```go file=cookbook/cors/origin-func/server.go
```
