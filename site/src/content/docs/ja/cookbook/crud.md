---
title: CRUD
description: Echo と JSON バインディングでリソースを作成、読み取り、更新、削除します。
sidebar:
  order: 2
---

インメモリ store を使った完全な CRUD（create, read, update, delete）API です。
ハンドラは JSON リクエストボディを struct にバインドし、ロック下で store にアクセスして、
JSON を返します。存在しないユーザーには `404` を、不正な id や空の name には `400` を返します。

## サーバー

```go file=cookbook/crud/server.go title="cookbook/crud/server.go"
```

## クライアント

### ユーザーを作成する

リクエスト：

```sh
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe Smith"}' \
  localhost:1323/users
```

レスポンス：

```json
{
  "id": 1,
  "name": "Joe Smith"
}
```

### ユーザーを取得する

リクエスト：

```sh
curl localhost:1323/users/1
```

レスポンス：

```json
{
  "id": 1,
  "name": "Joe Smith"
}
```

### ユーザー一覧を取得する

リクエスト：

```sh
curl localhost:1323/users
```

レスポンス：

```json
[
  {
    "id": 1,
    "name": "Joe Smith"
  }
]
```

### ユーザーを更新する

リクエスト：

```sh
curl -X PUT \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe"}' \
  localhost:1323/users/1
```

レスポンス：

```json
{
  "id": 1,
  "name": "Joe"
}
```

### ユーザーを削除する

リクエスト：

```sh
curl -X DELETE localhost:1323/users/1
```

レスポンス：`204 No Content`。

### 存在しないユーザー

リクエスト：

```sh
curl localhost:1323/users/9999
```

レスポンス（`404 Not Found`）：

```json
{
  "message": "user not found"
}
```
