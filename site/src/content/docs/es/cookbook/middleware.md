---
title: Middleware personalizado
description: Escribe middleware Echo personalizado para recopilar estadísticas de requests y establecer headers de response.
sidebar:
  order: 12
---

Esta receta muestra cómo escribir middleware personalizado:

- Un middleware que recopila el conteo de requests, estados de response y uptime.
- Un middleware que escribe un header `Server` personalizado en cada response.

Un middleware en Echo es una función con la firma
`func(next echo.HandlerFunc) echo.HandlerFunc`. El método `Stats.Process` de abajo
satisface esa firma directamente, mientras que `ServerHeader` es una función normal.

## Servidor

```go file=cookbook/middleware/server.go
```

## Response

### Headers

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
