---
title: CORS
description: Habilite Cross-Origin Resource Sharing com uma allow list ou função de origem customizada.
sidebar:
  order: 4
---

O [middleware CORS](/pt-br/middleware/cors/) controla quais origens podem acessar sua API.
Você pode passar uma lista fixa de origens permitidas ou fornecer uma função que decide
por request.

## Allow list de origens

Passe as origens permitidas diretamente para `middleware.CORS`.

```go file=cookbook/cors/origin-list/server.go
```

## Função de origem customizada

Para políticas dinâmicas, use `CORSWithConfig` com `UnsafeAllowOriginFunc`. A
função recebe o contexto do request e a origem, e retorna a origem a ser ecoada de volta,
se o request é permitido, e um erro opcional.

```go file=cookbook/cors/origin-func/server.go
```
