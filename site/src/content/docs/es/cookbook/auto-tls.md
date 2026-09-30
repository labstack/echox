---
title: Auto TLS
description: Obtén y renueva automáticamente certificados TLS desde Let's Encrypt.
sidebar:
  order: 3
---

Esta receta obtiene automáticamente certificados TLS para un dominio desde Let's Encrypt.
Configura un `StartConfig` con el `TLSConfig` del manager autocert y escucha en el
puerto `443`.

Abre `https://<DOMAIN>`. Si todo está configurado correctamente, deberías ver
un mensaje de bienvenida servido sobre TLS.

:::tip
- Para mayor seguridad, especifica una host policy en el manager autocert.
- Cachea certificados para evitar alcanzar los [rate limits de Let's Encrypt](https://letsencrypt.org/docs/rate-limits).
- Para redirigir tráfico HTTP a HTTPS, usa el [middleware redirect](/es/middleware/redirect/#https-redirect).
:::

## Servidor

```go file=cookbook/auto-tls/server.go#primary-server
```

## Usar un servidor HTTP personalizado

Si necesitas control total sobre `http.Server`, conecta el manager autocert a un
`tls.Config` personalizado:

```go file=cookbook/auto-tls/server.go#custom-server
```

```sh
go run . -custom-server
```
