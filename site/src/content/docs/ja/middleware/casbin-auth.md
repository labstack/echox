---
title: Casbin 認可
description: 小さなカスタムミドルウェアで Casbin アクセス制御ライブラリを使い、リクエストを認可します。
sidebar:
  order: 4
---

[Casbin](https://github.com/casbin/casbin) は、Go 向けの強力で効率的なオープンソース
アクセス制御ライブラリです。多くのモデルで認可を強制できます。

- ACL（Access Control List）
- スーパーユーザー付き ACL
- ユーザーなし ACL：認証やユーザーログインがないシステムに有用
- リソースなし ACL：個別リソースではなく、リソース種別（例：`write-article`、`read-log`）を対象にする
- RBAC（Role-Based Access Control）
- リソースロール付き RBAC：ユーザーとリソースの両方がロールを持てる
- ドメイン/テナント付き RBAC：ユーザーがドメイン/テナントごとに異なるロールセットを持てる
- ABAC（Attribute-Based Access Control）
- RESTful
- Deny-override：allow と deny の両方のルールをサポートし、deny が allow を上書きする

詳細は [API overview](https://casbin.org/docs/api-overview) と
[Casbin documentation](https://casbin.org/docs/) を参照してください。

## 依存関係

```bash
go get github.com/casbin/casbin/v3
```

```go
import (
	"github.com/casbin/casbin/v3"
)
```

## 実装

Echo には Casbin ミドルウェアは同梱されていません。この連携は Casbin enforcer の小さなラッパーです。

```go file=cookbook/casbin/server.go#middleware
```

## 例

Casbin モデルファイル `auth_model.conf` を作成します。

```ini file=cookbook/casbin/auth_model.conf
```

Casbin ポリシーファイル `auth_policy.csv` を作成します。

```csv file=cookbook/casbin/auth_policy.csv
```

認証と認可は別の関心事です。JWT や Basic Auth など別のミドルウェアでユーザーを認証し、
Casbin がリクエストを認可できるよう `userGetter` を渡します。

### JWT と使う

```go file=cookbook/casbin/server.go#jwt
```

次で試します。

```bash
curl -v "http://localhost:8080/dataset1/any" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
```

### Basic Auth と使う

```go
// BasicAuth middleware does authentication
e.Use(middleware.BasicAuth(func(c *echo.Context, user, password string) (bool, error) {
	return subtle.ConstantTimeCompare([]byte(user), []byte("alice")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1, nil
}))
basicAuthUser := func(c *echo.Context) (string, error) { // Basic auth user getter for Casbin authorization
	username, _, _ := c.Request().BasicAuth() // password is verified by the BasicAuth middleware above
	return username, nil
}
e.Use(NewCasbinMiddleware(ce, basicAuthUser)) // Casbin does authorization
```

次で試します。

```bash
# should pass
curl -v -u "alice:password" http://localhost:8080/dataset1/any
# should fail
curl -v -u "alice:password" http://localhost:8080/dataset2/resource2
```

### Casbin + JWT の完全な例

上記のモデルとポリシーファイルを使って、[Casbin + JWT の完全な例](https://github.com/labstack/echox/tree/master/cookbook/casbin)を実行できます。
