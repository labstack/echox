---
title: JWT
description: Autentica requests con JSON Web Tokens usando el middleware echo-jwt.
sidebar:
  order: 11
---

Esta receta demuestra autenticación JWT con Echo usando el middleware
[`echo-jwt`](https://github.com/labstack/echo-jwt):

- Autenticación JWT usando el algoritmo HS256.
- El token se lee desde el header de request `Authorization`.

Consulta la página del [middleware JWT](/es/middleware/jwt/) para ver todas las opciones de configuración.

## Servidor

### Usar claims personalizados

Define un tipo de claims que embebe `jwt.RegisteredClaims`, y luego apunta el middleware a él
con `NewClaimsFunc`. Dentro del handler restringido, recupera el token parseado desde el
contexto con el genérico `echo.ContextGet`.

```go file=cookbook/jwt/custom-claims/server.go
```

### Usar un KeyFunc definido por el usuario

Cuando los tokens están firmados por un proveedor de identidad externo, proporciona un `KeyFunc`
que resuelva dinámicamente la signing key. Este ejemplo valida tokens emitidos por Google Sign-In
obteniendo el conjunto de claves públicas de Google.

```go file=cookbook/jwt/user-defined-keyfunc/server.go
```

:::caution
Obtener el conjunto de claves en cada request, como se muestra arriba, es solo para demostración.
En producción, cachea el conjunto de claves y actualízalo periódicamente.
:::

## Cliente

### Login

Inicia sesión con un username y password para obtener un token.

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

Solicita un recurso restringido usando el token en el header de request `Authorization`.

```sh
curl localhost:1323/restricted -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"
```

Response:

```sh
Welcome Jon Snow!
```
