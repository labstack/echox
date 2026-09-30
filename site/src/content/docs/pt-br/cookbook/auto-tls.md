---
title: Auto TLS
description: Obtenha e renove automaticamente certificados TLS da Let's Encrypt.
sidebar:
  order: 3
---

Esta receita obtém certificados TLS para um domínio automaticamente da Let's Encrypt.
Configure um `StartConfig` com o `TLSConfig` do gerenciador autocert e escute na
porta `443`.

Acesse `https://<DOMAIN>`. Se tudo estiver configurado corretamente, você deverá ver
uma mensagem de boas-vindas servida por TLS.

:::tip
- Para segurança adicional, especifique uma política de host no gerenciador autocert.
- Faça cache de certificados para evitar atingir os [limites de taxa da Let's Encrypt](https://letsencrypt.org/docs/rate-limits).
- Para redirecionar tráfego HTTP para HTTPS, use o [middleware de redirect](/pt-br/middleware/redirect/#https-redirect).
:::

## Servidor

```go file=cookbook/auto-tls/server.go
```

## Usando um servidor HTTP customizado

Se você precisar de controle total sobre o `http.Server`, conecte o gerenciador autocert a um
`tls.Config` customizado:

```go file=cookbook/auto-tls/custom-server/server.go
```

```sh
go run ./custom-server
```
