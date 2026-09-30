---
title: カスタムミドルウェア
description: リクエスト統計を収集し、レスポンス header を設定するカスタム Echo ミドルウェアを書きます。
sidebar:
  order: 12
---

このレシピでは、カスタムミドルウェアの書き方を示します。

- リクエスト数、レスポンスステータス、稼働時間を収集するミドルウェア。
- すべてのレスポンスにカスタム `Server` header を書き込むミドルウェア。

Echo のミドルウェアは `func(next echo.HandlerFunc) echo.HandlerFunc` というシグネチャの関数です。
下の `Stats.Process` メソッドはそのシグネチャを直接満たし、`ServerHeader` は通常の関数です。

## サーバー

```go file=cookbook/middleware/server.go
```

## レスポンス

### Header

```sh
Content-Length:122
Content-Type:application/json; charset=utf-8
Date:Thu, 14 Apr 2016 20:31:46 GMT
Server:Echo/5.0
```

### Body

```js
{
  "uptime": "2016-04-14T13:28:48.486548936-07:00",
  "requestCount": 5,
  "statuses": {
    "200": 4,
    "404": 1
  }
}
```
