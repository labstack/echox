---
title: 自动 TLS
description: 自动从 Let's Encrypt 获取并续期 TLS 证书。
sidebar:
  order: 3
---

这个示例会自动从 Let's Encrypt 为域名获取 TLS 证书。使用 autocert manager 的
`TLSConfig` 配置 `StartConfig`，并监听 `443` 端口。

访问 `https://<DOMAIN>`。如果一切配置正确，你应该能看到通过 TLS 提供的欢迎消息。

:::tip
- 为增强安全性，请在 autocert manager 中指定 host policy。
- 缓存证书以避免触发 [Let's Encrypt rate limits](https://letsencrypt.org/docs/rate-limits)。
- 要将 HTTP 流量重定向到 HTTPS，请使用[重定向中间件](/zh-cn/middleware/redirect/#https-redirect)。
:::

## 服务器

```go file=cookbook/auto-tls/server.go
```

## 使用自定义 HTTP 服务器

如果你需要完全控制 `http.Server`，请改为把 autocert manager 接入自定义 `tls.Config`：

```go file=cookbook/auto-tls/custom-server/server.go
```

```sh
go run ./custom-server
```
