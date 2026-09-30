---
title: Timeout
description: Aplique um timeout de request aos handlers com o middleware ContextTimeout.
sidebar:
  order: 18
---

O middleware [`ContextTimeout`](/pt-br/middleware/context-timeout/) define um deadline no
`context.Context` do request. Quando o deadline passa, o contexto é cancelado,
e handlers que observam `c.Request().Context().Done()` podem retornar rapidamente em vez
de continuar até o fim.

No exemplo abaixo, o middleware impõe um timeout de 5 segundos enquanto o handler
levaria 10 segundos; por isso, o request retorna `408 Request Timeout`.

## Servidor

```go file=cookbook/timeout/server.go
```
