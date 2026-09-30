---
title: タイムアウト
description: ContextTimeout ミドルウェアでハンドラにリクエストタイムアウトを適用します。
sidebar:
  order: 18
---

[`ContextTimeout`](/ja/middleware/context-timeout/) ミドルウェアは、リクエストの `context.Context`
に deadline を設定します。deadline を過ぎるとコンテキストがキャンセルされ、
`c.Request().Context().Done()` を監視しているハンドラは、最後まで実行されるのを待たずにすぐ戻れます。

下の例では、ミドルウェアが 5 秒のタイムアウトを適用します。一方、ハンドラは通常なら 10 秒かかるため、
リクエストは `408 Request Timeout` を返します。

## サーバー

```go file=cookbook/timeout/server.go
```
