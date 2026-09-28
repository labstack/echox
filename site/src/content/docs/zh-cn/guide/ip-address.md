---
title: IP 地址
description: 在代理之后安全地获取真实客户端 IP 地址。
sidebar:
  order: 14
---

IP 地址在 HTTP 中扮演基础角色，它用于访问控制、审计、基于地理位置的分析等。
Echo 暴露 `Context#RealIP()` 来获取它。

获取_真实_客户端 IP 并不简单，尤其是当 L7 代理位于应用前方时。在这种情况下，
真实 IP 必须由代理通过 HTTP 转发给你的应用，但你不能无条件信任 HTTP header，
否则就可能被欺骗。**这是一个安全风险。**

要可靠且安全地获取 IP，你的应用必须了解整个基础设施。在 Echo 中，你通过
`Echo#IPExtractor` 配置这一点。

从 v5.1.0 开始，v5 的 `Context#RealIP()` 默认使用直接对端地址。v4 仍会默认信任任意客户端发送的 `X-Forwarded-For` 和 `X-Real-IP`，因此请按代理拓扑**始终配置 `Echo#IPExtractor`**。HTTP/HTTPS 的判断需另行配置，参见[请求协议与可信代理](/zh-cn/guide/request-scheme/)。

从两个问题开始，找到正确方法：

1. 你的应用前面是否有任何 HTTP（L7）代理？这包括云负载均衡器（如 AWS ALB 或 GCP HTTP LB）
   和开源代理（如 Nginx、Envoy 或 Istio ingress gateway）。
2. 如果有，这些代理使用哪个 HTTP header 将客户端 IP 传递给应用？

## 情况 1：没有代理

如果没有代理（应用直接面向互联网），唯一可以信任的地址来自网络层。每个 HTTP header 都不可信，
因为客户端可以完全控制它们。

使用 `echo.ExtractIPDirect()`：

```go
e.IPExtractor = echo.ExtractIPDirect()
```

## 情况 2：使用 X-Forwarded-For header 的代理

[`X-Forwarded-For` (XFF)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For)
是转发客户端 IP 最常见的 header。每经过一跳，代理都会把请求 IP 追加到 header 末尾。

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

在这种情况下，应从右侧开始取得**第一个不可信的 IP 读数**。永远不要取最左边的第一个，
因为客户端控制它。这里“可信”表示你确定该 IP 属于你的基础设施。在上面的例子中，
如果 `b` 和 `c` 可信，则两种情况下客户端 IP 都是 `a`，绝不是 `x`。

使用 `ExtractIPFromXFFHeader(...TrustOption)`：

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader()
```

默认情况下，它信任内部 IP 地址，包括 loopback、link-local unicast、private-use，
以及来自 [RFC 6890](https://datatracker.ietf.org/doc/html/rfc6890)、
[RFC 4291](https://datatracker.ietf.org/doc/html/rfc4291) 和
[RFC 4193](https://datatracker.ietf.org/doc/html/rfc4193) 的 unique local 地址。
可用 `TrustOption` 控制这一点：

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader(
	echo.TrustLoopback(false),   // e.g. IPv4 starting with 127.
	echo.TrustLinkLocal(false),  // e.g. IPv4 starting with 169.254.
	echo.TrustPrivateNet(false), // e.g. IPv4 starting with 10. or 192.168.
	echo.TrustIPRange(lbIPRange),
)
```

## 情况 3：使用 X-Real-IP header 的代理

`X-Real-IP` 是另一个用于转发客户端 IP 的 header，但与 XFF 不同，它只携带单个地址。

如果你的代理设置了这个 header，请使用 `ExtractIPFromRealIPHeader(...TrustOption)`：

```go
e.IPExtractor = echo.ExtractIPFromRealIPHeader()
```

与 XFF 一样，它默认信任内部 IP 地址，并接受相同的 `TrustOption`。

:::danger
**永远不要忘记**配置最外层代理（基础设施边缘的代理）**不要透传传入的 header**。
否则客户端可以伪造它们，从而打开欺诈之门。
:::

## 默认行为

未配置 `Echo#IPExtractor` 时，**v5 使用直接对端地址**，忽略客户端提供的 IP 标头。**v4 保留旧行为**，可能在不验证代理的情况下使用 `X-Forwarded-For` 或 `X-Real-IP`。v4 在无代理时使用 `echo.ExtractIPDirect()`；有代理时使用带适当信任选项的标头提取器。
