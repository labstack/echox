---
title: JWT
description: 使用 echo-jwt 中间件通过 JSON Web Tokens 认证请求。
sidebar:
  order: 11
---

此示例演示如何使用 [`echo-jwt`](https://github.com/labstack/echo-jwt) 中间件在 Echo 中进行
JWT 认证：

- 使用 HS256 算法进行 JWT 认证。
- 从 `Authorization` 请求 header 读取 token。

完整配置选项请参见 [JWT 中间件](/zh-cn/middleware/jwt/)页面。

## 服务器

### 使用自定义 claims

定义一个嵌入 `jwt.RegisteredClaims` 的 claims 类型，然后通过 `NewClaimsFunc` 将中间件指向它。
在受限处理函数中，使用泛型 `echo.ContextGet` 从上下文获取已解析的 token。

```go file=cookbook/jwt/custom-claims/server.go
```

### 使用用户定义的 KeyFunc

当 token 由外部身份提供商签名时，请提供一个 `KeyFunc` 来动态解析签名 key。
此示例通过获取 Google 公钥集来验证 Google Sign-In 签发的 token。

```go file=cookbook/jwt/user-defined-keyfunc/server.go
```

:::caution
如上所示，在每个请求上获取 key set 仅用于演示。在生产环境中，请缓存 key set 并定期刷新。
:::

## 客户端

### 登录

使用用户名和密码登录以获取 token。

```sh
curl -X POST -d 'username=jon' -d 'password=shhh!' localhost:1323/login
```

响应：

```js
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
}
```

### 请求

使用 `Authorization` 请求 header 中的 token 请求受限资源。

```sh
curl localhost:1323/restricted -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
```

响应：

```sh
Welcome Jon Snow!
```
