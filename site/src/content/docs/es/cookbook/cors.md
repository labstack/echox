---
title: CORS
description: Habilita Cross-Origin Resource Sharing con una allow list o una función de origin personalizada.
sidebar:
  order: 4
---

El [middleware CORS](/es/middleware/cors/) controla qué origins pueden acceder a tu API.
Puedes pasar una lista fija de origins permitidos o proporcionar una función que decida
por request.

## Allow list de origins

Pasa los origins permitidos directamente a `middleware.CORS`.

```go file=cookbook/cors/origin-list/server.go
```

## Función de origin personalizada

Para policies dinámicas, usa `CORSWithConfig` con `UnsafeAllowOriginFunc`. La
función recibe el contexto del request y el origin, y devuelve el origin que se debe
reflejar, si el request está permitido y un error opcional.

```go file=cookbook/cors/origin-func/server.go
```
