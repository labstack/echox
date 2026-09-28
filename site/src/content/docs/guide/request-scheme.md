---
title: Request Scheme and Trusted Proxies
description: Configure how Echo determines HTTP or HTTPS behind a trusted proxy.
sidebar:
  order: 15
---

`Context#Scheme()` tells Echo whether a request used HTTP or HTTPS. HTTPS redirects,
HSTS, and the Proxy middleware rely on it. Since Echo v5.4.0 and v4.16.0, Echo only
uses `X-Forwarded-Proto`, `X-Forwarded-Protocol`, `X-Forwarded-Ssl`, or
`X-Url-Scheme` when the **direct peer** is a loopback, link-local, or private address,
or a Unix socket. Otherwise, the connection itself determines the scheme. This
prevents a public client from sending `X-Forwarded-Proto: https` over HTTP to bypass
an HTTPS redirect. See the [v5.4.0 release notes](https://github.com/labstack/echo/releases/tag/v5.4.0).

## Choose a scheme extractor

`Echo#SchemeExtractor` controls this behavior. In v5, `Config.SchemeExtractor` can
set it when creating Echo. The default is `echo.ExtractSchemeFromHeaders()` with the
trusted direct-peer addresses above. `echo.ExtractSchemeDirect()` ignores forwarding
headers. `echo.LegacySchemeExtractor()` restores the older behavior, which is unsafe
if an untrusted client can reach the app or a proxy forwards a client-supplied scheme
header.

When `X-Forwarded-Proto` is present, Echo uses **only its last value** and returns
the scheme in lowercase. An invalid value becomes `http`; Echo does not fall back
to another forwarded-scheme header. Your trusted proxy must **overwrite**
`X-Forwarded-Proto` with the scheme it observed. It must not pass through a value
sent by the client. For example, configure nginx with
`proxy_set_header X-Forwarded-Proto $scheme;`.

## Proxies with public addresses

If the direct proxy has a public address, or uses `100.64.0.0/10`, explicitly trust
its address ranges. Otherwise Echo ignores its `X-Forwarded-Proto`: HTTPS redirects
can loop and Secure middleware will not send HSTS. This applies to Cloudflare,
CloudFront, and Azure Front Door connecting to a public origin, and to GCP external
HTTP(S) load balancers or GKE Ingress. Trust only the ranges your deployment uses.

For the GCP load balancer ranges:

```go
_, gclb1, _ := net.ParseCIDR("35.191.0.0/16")
_, gclb2, _ := net.ParseCIDR("130.211.0.0/22")
e.SchemeExtractor = echo.ExtractSchemeFromHeaders(
	echo.TrustIPRange(gclb1),
	echo.TrustIPRange(gclb2),
)
```

Import `net` and `github.com/labstack/echo/v5` for this v5 example. For v4, use
`github.com/labstack/echo/v4` instead. Proxies on the same host or a private
network, and Unix sockets, keep working with the default extractor. If an adapter
replaces `RemoteAddr` with the client's address, configure the extractor for the
actual deployment rather than assuming Echo sees the proxy address.

The scheme and client IP are separate decisions. Configure `Echo#IPExtractor` for
client IPs as described in [IP Address](/guide/ip-address/).
