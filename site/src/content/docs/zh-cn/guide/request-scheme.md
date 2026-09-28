---
title: 请求协议与可信代理
description: 配置 Echo 如何在可信代理后判断 HTTP 或 HTTPS。
sidebar:
  order: 15
---

`Context#Scheme()` 返回请求使用的 HTTP 或 HTTPS。HTTPS 重定向、HSTS 和 Proxy
中间件都依赖它。从 Echo v5.4.0 和 v4.16.0 开始，只有**直接连接的对端**是
环回、链路本地或私有地址，或通过 Unix 套接字连接时，Echo 才使用
`X-Forwarded-Proto`、`X-Forwarded-Protocol`、`X-Forwarded-Ssl` 和
`X-Url-Scheme`。否则由连接本身决定协议。这样，公网客户端无法通过 HTTP 发送
`X-Forwarded-Proto: https` 来绕过 HTTPS 重定向。

## 选择协议提取器

使用 `Echo#SchemeExtractor` 配置判断方式；v5 还支持
`Config.SchemeExtractor`。默认的 `echo.ExtractSchemeFromHeaders()` 仅信任上述
直接对端。`echo.ExtractSchemeDirect()` 忽略转发标头。
`echo.LegacySchemeExtractor()` 恢复旧行为；如果不可信客户端能直接访问应用，
或代理透传客户端提供的标头，旧行为并不安全。

存在 `X-Forwarded-Proto` 时，Echo **只使用最后一个值**，并以小写形式返回协议。
无效值会得到 `http`，不会回退到其他转发协议标头。可信代理必须用自己观察到的
协议**覆盖**该标头，不能透传客户端的值。nginx 可配置
`proxy_set_header X-Forwarded-Proto $scheme;`。

## 公网地址的代理

如果直接连接的代理使用公网地址或 `100.64.0.0/10`，需要明确设置信任的地址段。
否则 Echo 会忽略其 `X-Forwarded-Proto`，导致 HTTPS 重定向循环或 Secure
中间件不发送 HSTS。Cloudflare、CloudFront、Azure Front Door、GCP 外部
HTTP(S) 负载均衡器和 GKE Ingress 都可能属于这种情况。只信任部署实际使用的
地址段。以下是 GCP 地址段示例：

```go
_, gclb1, _ := net.ParseCIDR("35.191.0.0/16")
_, gclb2, _ := net.ParseCIDR("130.211.0.0/22")
e.SchemeExtractor = echo.ExtractSchemeFromHeaders(
	echo.TrustIPRange(gclb1),
	echo.TrustIPRange(gclb2),
)
```

这个 v5 示例需要导入 `net` 和 `github.com/labstack/echo/v5`；v4 使用
`github.com/labstack/echo/v4`。同一主机、私有网络和 Unix 套接字上的代理可沿用
默认设置。如果适配器把 `RemoteAddr` 改为客户端地址，请按实际拓扑配置提取器。
客户端 IP 需要另外配置，参见 [IP 地址](/zh-cn/guide/ip-address/)。
