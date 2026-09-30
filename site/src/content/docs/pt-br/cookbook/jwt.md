---
title: JWT
description: Autentique requests com JSON Web Tokens usando o middleware echo-jwt.
sidebar:
  order: 11
---

Esta receita demonstra autenticação JWT com Echo usando o middleware
[`echo-jwt`](https://github.com/labstack/echo-jwt):

- Autenticação JWT usando o algoritmo HS256.
- O token é lido do header de request `Authorization`.

Veja a página do [middleware JWT](/pt-br/middleware/jwt/) para opções completas de configuração.

## Servidor

### Usando claims customizadas

Defina um tipo de claims que incorpora `jwt.RegisteredClaims`, então aponte o middleware
para ele com `NewClaimsFunc`. Dentro do handler restrito, recupere o token analisado
do contexto com o `echo.ContextGet` genérico.

```go file=cookbook/jwt/custom-claims/server.go
```

### Usando uma KeyFunc definida pelo usuário

Quando tokens são assinados por um provedor de identidade externo, forneça uma `KeyFunc` que
resolva a chave de assinatura dinamicamente. Este exemplo valida tokens emitidos pelo
Google Sign-In buscando o conjunto de chaves públicas do Google.

```go file=cookbook/jwt/user-defined-keyfunc/server.go
```

:::caution
Buscar o conjunto de chaves a cada request, como mostrado acima, é apenas para demonstração.
Em produção, faça cache do conjunto de chaves e atualize-o periodicamente.
:::

## Cliente

### Login

Faça login com usuário e senha para recuperar um token.

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

Solicite um recurso restrito usando o token no header de request `Authorization`.

```sh
curl localhost:1323/restricted -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
```

Response:

```sh
Welcome Jon Snow!
```
