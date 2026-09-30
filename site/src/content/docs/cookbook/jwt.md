---
title: JWT
description: Authenticate requests with JSON Web Tokens using the echo-jwt middleware.
sidebar:
  order: 11
---

This recipe demonstrates JWT authentication with Echo using the
[`echo-jwt`](https://github.com/labstack/echo-jwt) middleware:

- JWT authentication using the HS256 algorithm.
- The token is read from the `Authorization` request header.

See the [JWT middleware](/middleware/jwt/) page for full configuration options.

## Server

### Using custom claims

Define a claims type that embeds `jwt.RegisteredClaims`, then point the middleware
at it with `NewClaimsFunc`. Inside the restricted handler, retrieve the parsed token
from the context with the generic `echo.ContextGet`.

```go file=cookbook/jwt/custom-claims/server.go
```

### Using a user-defined KeyFunc

When tokens are signed by an external identity provider, supply a `KeyFunc` that
resolves the signing key dynamically. This example validates tokens issued by
Google Sign-In by fetching Google's public key set.

```go file=cookbook/jwt/user-defined-keyfunc/server.go
```

:::caution
Fetching the key set on every request, as shown above, is for demonstration only.
In production, cache the key set and refresh it periodically.
:::

## Client

### Login

Log in with a username and password to retrieve a token.

```sh
curl -X POST -d 'username=jon' -d 'password=shhh!' localhost:1323/login
```

Response:

```js
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
}
```

### Request

Request a restricted resource using the token in the `Authorization` request header.

```sh
curl localhost:1323/restricted -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
```

Response:

```sh
Welcome Jon Snow!
```
