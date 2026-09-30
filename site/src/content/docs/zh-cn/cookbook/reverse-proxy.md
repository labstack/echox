---
title: 反向代理
description: 使用 Echo 作为上游应用前的反向代理和负载均衡器。
sidebar:
  order: 19
---

此示例演示如何在你的应用（如 WordPress、Node.js、Java、Python、Ruby 或 Go）前使用 Echo
作为反向代理和负载均衡器。为简单起见，这里的上游也是处理 WebSocket 的 Go 服务器。

## 1) 确定上游目标 URL

```go file=cookbook/reverse-proxy/server.go#targets
```

## 2) 使用上游目标设置代理中间件

下面的片段使用轮询负载均衡。你也可以使用 `middleware.NewRandomBalancer()`。

```go file=cookbook/reverse-proxy/server.go#middleware
```

要为子路由设置代理，请使用 `Echo#Group()`。

```go file=cookbook/reverse-proxy/server.go#grouped-proxy
```

```sh
go run . -grouped
```

## 3) 启动上游服务器

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) 启动代理服务器

```sh
go run server.go
```

访问 `http://localhost:1323`，你应该会看到网页中 HTTP 请求由 "server 1" 提供，
WebSocket 请求由 "server 2" 提供。

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## 源码

### 上游服务器

```go file=cookbook/reverse-proxy/upstream/server.go
```

### 代理服务器

```go file=cookbook/reverse-proxy/server.go
```

## Echo 安全更新 (v5.4.0 / v4.16.0)

Proxy 从 `Context#Scheme()` 设置 `X-Forwarded-Proto` 并删除旧协议标头。v5 的 `X-Real-IP` 来自 `Context#RealIP()`。如果 Echo 前面还有代理，请配置[协议](/zh-cn/guide/request-scheme/)和 [IP](/zh-cn/guide/ip-address/) 提取器，确保上游收到正确的值。
