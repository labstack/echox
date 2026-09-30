---
title: サブドメイン
description: 仮想ホストハンドラを使い、host ごとにリクエストを異なる Echo インスタンスへルーティングします。
sidebar:
  order: 17
---

このレシピでは、リクエスト host に基づいてリクエストを別々の `Echo` インスタンスへルーティングします。
これにより、各サブドメインが独自のルーティングとミドルウェアを持てます。
これらのインスタンスは `echo.NewVirtualHostHandler` で結合され、host name ごとに振り分けられます。

## サーバー

```go file=cookbook/subdomain/server.go
```
