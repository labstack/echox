---
title: JWT
description: echo-jwt ミドルウェアで JSON Web Tokens を使ってリクエストを認証します。
sidebar:
  order: 11
---

このレシピでは、[`echo-jwt`](https://github.com/labstack/echo-jwt) ミドルウェアを使った
Echo の JWT 認証を示します。

- HS256 アルゴリズムを使った JWT 認証。
- token は `Authorization` リクエスト header から読み取ります。

完全な設定オプションは [JWT ミドルウェア](/ja/middleware/jwt/)ページを参照してください。

## サーバー

### カスタム claims を使う

`jwt.RegisteredClaims` を埋め込む claims 型を定義し、`NewClaimsFunc` でミドルウェアに指定します。
制限付きハンドラ内では、ジェネリック `echo.ContextGet` を使ってコンテキストから解析済み token を取得します。

```go file=cookbook/jwt/custom-claims/server.go
```

### ユーザー定義 KeyFunc を使う

token が外部 ID プロバイダーによって署名されている場合は、署名 key を動的に解決する
`KeyFunc` を指定します。この例では、Google の公開 key set を取得して Google Sign-In が発行した
token を検証します。

```go file=cookbook/jwt/user-defined-keyfunc/server.go
```

:::caution
上のように各リクエストで key set を取得するのはデモ目的のみです。本番では key set をキャッシュし、
定期的に更新してください。
:::

## クライアント

### ログイン

ユーザー名とパスワードでログインし、token を取得します。

```sh
curl -X POST -d 'username=jon' -d 'password=shhh!' localhost:1323/login
```

レスポンス：

```js
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
}
```

### リクエスト

`Authorization` リクエスト header の token を使って制限付きリソースをリクエストします。

```sh
curl localhost:1323/restricted -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
```

レスポンス：

```sh
Welcome Jon Snow!
```
