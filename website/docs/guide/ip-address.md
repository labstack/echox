---
description: IP address
slug: /ip-address
sidebar_position: 8
---

# IP Address

IP address plays a fundamental role in HTTP; it's used for access control, auditing, geo-based access analysis, and more.
Echo provides a handy method [`Context#RealIP()`](https://godoc.org/github.com/labstack/echo#Context) for that.

However, it is not trivial to retrieve the _real_ IP address from requests, especially when you put L7 proxies before the application.
In such situations, _real_ IP needs to be relayed on the HTTP layer from proxies to your app, however, you must not trust HTTP headers unconditionally.
Otherwise, you might give someone a chance of deceiving you. **A security risk!**

To retrieve IP address reliably/securely, you must let your application be aware of the entire architecture of your infrastructure.
In Echo, this can be done by configuring `Echo#IPExtractor` appropriately.
This guides show you why and how.

:::caution

Note: if you don't set `Echo#IPExtractor` explicitly, Echo fallback to legacy behavior, which is not a good choice.

:::


Let's start from two questions to know the right direction:

1. Do you put any HTTP (L7) proxy in front of the application?
    - It includes both cloud solutions (such as AWS ALB or GCP HTTP LB) and OSS ones (such as Nginx, Envoy or Istio ingress gateway).
2. If yes, what HTTP header do your proxies use to pass client IP to the application?

## Case 1. With no proxy

If you put no proxy (e.g.: directory facing to the internet), all you need to (and have to) see is IP address from network layer.
Any HTTP header is untrustable because the clients have full control what headers to be set.

In this case, use `echo.ExtractIPDirect()`.

```go
e.IPExtractor = echo.ExtractIPDirect()
```

## Case 2. With proxies using X-Forwarded-For header

[`X-Forwarded-For` (XFF)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For) is the popular header to relay clients' IP addresses.
At each hop on the proxies, they append the request IP address at the end of the header.

Following example diagram illustrates this behavior.

```text
            ┌──────────┐            ┌──────────┐            ┌──────────┐
───────────>│ Proxy 1  │───────────>│ Proxy 2  │───────────>│ Your app │
            │ (IP: b)  │            │ (IP: c)  │            │          │
            └──────────┘            └──────────┘            └──────────┘

Case 1.
XFF:  ""                    "a"                     "a, b"
                                                    ~~~~~~
Case 2.
XFF:  "x"                   "x, a"                  "x, a, b"
                                                    ~~~~~~~~~
                                                    ↑ What your app will see
```

In this case, use **first _untrustable_ IP reading from right**. Never use first one reading from left, as it is configurable by client. Here "trustable" means "you are sure the IP address belongs to your infrastructure". In above example, if `b` and `c` are trustable, the IP address of the client is `a` for both cases, never be `x`.

In Echo, use `ExtractIPFromXFFHeader(...TrustOption)`.

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader()
```

By default, it trusts internal IP addresses (loopback, link-local unicast, private-use and unique local address from [RFC6890](https://tools.ietf.org/html/rfc6890), [RFC4291](https://tools.ietf.org/html/rfc4291) and [RFC4193](https://tools.ietf.org/html/rfc4193)). To control this behavior, use [`TrustOption`](https://godoc.org/github.com/labstack/echo#TrustOption)s.

E.g.:

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader(
     echo.TrustLoopback(false), // e.g. ipv4 start with 127. 
     echo.TrustLinkLocal(false), // e.g. ipv4 start with 169.254
     echo.TrustPrivateNet(false), // e.g. ipv4 start with 10. or 192.168
     echo.TrustIPRange(lbIPRange),
)
```

- Ref: https://godoc.org/github.com/labstack/echo#TrustOption

## Case 3. With proxies using X-Real-IP header

`X-Real-IP` is another HTTP header to relay clients' IP addresses, but it carries only one address unlike XFF.

If your proxies set this header, use `ExtractIPFromRealIPHeader(...TrustOption)`.

```go
e.IPExtractor = echo.ExtractIPFromRealIPHeader()
```

Again, it trusts internal IP addresses by default (loopback, link-local unicast, private-use and unique local address from [RFC6890](https://tools.ietf.org/html/rfc6890), [RFC4291](https://tools.ietf.org/html/rfc4291) and [RFC4193](https://tools.ietf.org/html/rfc4193)). To control this behavior, use [`TrustOption`](https://godoc.org/github.com/labstack/echo#TrustOption)s.

- Ref: https://godoc.org/github.com/labstack/echo#TrustOption

> **Never forget** to configure the outermost proxy (i.e.; at the edge of your infrastructure) **not to pass through incoming headers**.
> Otherwise there is a chance of fraud, as it is what clients can control.

## About default behavior

In default behavior, Echo sees all of first XFF header, X-Real-IP header and IP from network layer.

As you might already notice, after reading this article, this is not good.
Sole reason this is default is just backward compatibility.
