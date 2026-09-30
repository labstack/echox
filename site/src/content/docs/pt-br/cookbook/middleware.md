---
title: Middleware customizado
description: Escreva middleware Echo customizado para coletar estatísticas de request e definir headers de response.
sidebar:
  order: 12
---

Esta receita mostra como escrever middleware customizado:

- Um middleware que coleta a contagem de requests, status de response e uptime.
- Um middleware que escreve um header `Server` customizado em toda response.

Um middleware no Echo é uma função com a assinatura
`func(next echo.HandlerFunc) echo.HandlerFunc`. O método `Stats.Process` abaixo
satisfaz essa assinatura diretamente, enquanto `ServerHeader` é uma função simples.

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
