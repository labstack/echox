---
title: Timeout
description: Aplica un timeout de request a handlers con el middleware ContextTimeout.
sidebar:
  order: 18
---

El middleware [`ContextTimeout`](/es/middleware/context-timeout/) establece un deadline en el
`context.Context` del request. Cuando el deadline vence, el contexto se cancela, y los handlers
que observan `c.Request().Context().Done()` pueden retornar rápidamente en vez de ejecutarse
hasta completar.

En el ejemplo de abajo, el middleware impone un timeout de 5 segundos, mientras que el handler
tomaría 10 segundos de otro modo, por lo que el request devuelve `408 Request Timeout`.

## Servidor

```go file=cookbook/timeout/server.go
```
