---
title: Subdomínio
description: Roteie requests para diferentes instâncias Echo por host usando um handler de virtual host.
sidebar:
  order: 17
---

Esta receita roteia requests para instâncias `Echo` separadas com base no host do request,
para que cada subdomínio tenha suas próprias rotas e middleware. As instâncias são
combinadas com `echo.NewVirtualHostHandler`, que despacha pelo nome do host.

## Servidor

```go file=cookbook/subdomain/server.go
```
