---
title: Casbin 认证
description: 使用一个小型自定义中间件，通过 Casbin 访问控制库为请求授权。
sidebar:
  order: 4
---

[Casbin](https://github.com/casbin/casbin) 是面向 Go 的强大、高效开源访问控制库。
它支持在多种模型上执行授权：

- ACL（Access Control List）
- 带超级用户的 ACL
- 无用户的 ACL，适合没有认证或用户登录的系统
- 无资源的 ACL，面向一类资源（例如 `write-article`、`read-log`），而不是单个资源
- RBAC（Role-Based Access Control）
- 带资源角色的 RBAC，用户和资源都可以有角色
- 带域/租户的 RBAC，用户可以在每个域/租户拥有不同角色集
- ABAC（Attribute-Based Access Control）
- RESTful
- Deny-override，同时支持允许和拒绝规则，拒绝会覆盖允许

详情请参见 [API overview](https://casbin.org/docs/api-overview) 和
[Casbin documentation](https://casbin.org/docs/)。

## 依赖

```bash
go get github.com/casbin/casbin/v3
```

```go
import (
	"github.com/casbin/casbin/v3"
)
```

## 实现

Echo 不自带 Casbin 中间件；该集成是对 Casbin enforcer 的一层小包装：

```go file=cookbook/casbin/server.go#middleware
```

## 示例

创建 Casbin 模型文件 `auth_model.conf`：

```ini file=cookbook/casbin/auth_model.conf
```

创建 Casbin 策略文件 `auth_policy.csv`：

```csv file=cookbook/casbin/auth_policy.csv
```

认证和授权是不同的关注点。使用另一个中间件（如 JWT 或 Basic Auth）认证用户，
然后提供 `userGetter`，使 Casbin 可以为请求授权。

### 搭配 JWT

```go file=cookbook/casbin/server.go#jwt
```

用下面命令尝试：

```bash
curl -v "http://localhost:8080/dataset1/any" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
```

### 搭配 Basic Auth

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

用下面命令尝试：

```bash
# should pass
curl -v -u "alice:password" http://localhost:8080/dataset1/any
# should fail
curl -v -u "alice:password" http://localhost:8080/dataset2/resource2
```

### 完整 Casbin + JWT 示例

使用上面的模型和策略文件运行[完整的 Casbin + JWT 示例](https://github.com/labstack/echox/tree/master/cookbook/casbin)。
