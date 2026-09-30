---
title: Subdominio
description: Enruta requests a distintas instancias Echo por host usando un handler de virtual host.
sidebar:
  order: 17
---

Esta receta enruta requests a instancias `Echo` separadas según el host del request,
de modo que cada subdominio tenga sus propias rutas y middleware. Las instancias se
combinan con `echo.NewVirtualHostHandler`, que despacha por host name.

## Servidor

```go file=cookbook/subdomain/server.go
```
