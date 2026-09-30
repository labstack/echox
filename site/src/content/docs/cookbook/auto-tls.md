---
title: Auto TLS
description: Automatically obtain and renew TLS certificates from Let's Encrypt.
sidebar:
  order: 3
---

This recipe obtains TLS certificates for a domain automatically from Let's Encrypt.
Configure a `StartConfig` with the autocert manager's `TLSConfig` and listen on
port `443`.

Browse to `https://<DOMAIN>`. If everything is configured correctly, you should see
a welcome message served over TLS.

:::tip
- For added security, specify a host policy in the autocert manager.
- Cache certificates to avoid hitting [Let's Encrypt rate limits](https://letsencrypt.org/docs/rate-limits).
- To redirect HTTP traffic to HTTPS, use the [redirect middleware](/middleware/redirect/#https-redirect).
:::

## Server

```go file=cookbook/auto-tls/server.go#primary-server
```

## Using a custom HTTP server

If you need full control over the `http.Server`, wire the autocert manager into a
custom `tls.Config` instead:

```go file=cookbook/auto-tls/server.go#custom-server
```
