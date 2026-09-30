---
title: Server-Sent Events (SSE)
description: 从 Echo 处理函数流式发送 server-sent events，可按连接发送或广播给多个客户端。
sidebar:
  order: 14
---

[Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)
可以有多种使用方式。下面第一个示例是每个连接、每个处理函数的 SSE。更复杂的广播逻辑请参见第二个使用
[r3labs/sse](https://github.com/r3labs/sse) 的示例。

:::caution
SSE 连接是长连接，因此必须禁用服务器写超时。两个示例都通过 `BeforeServeFunc`
设置 `s.WriteTimeout = 0`。
:::

## 使用 SSE

### 服务器

处理函数写入 SSE header，然后每秒发送一个事件，直到客户端断开连接。
`http.NewResponseController(w).Flush()` 会立即把每个事件推送到客户端。

```go file=cookbook/sse/simple/server.go
```

### Event 结构和 Marshal 方法

```go file=cookbook/sse/simple/serversentevent.go
```

### 提供 SSE 的 HTML

```html file=cookbook/sse/simple/index.html
```

## 使用 r3labs/sse 广播

当你需要把单个事件流广播给多个订阅者时，[r3labs/sse](https://github.com/r3labs/sse)
库会为你处理流和订阅者管理。

### 服务器

```go file=cookbook/sse/broadcast/server.go
```

### 提供 SSE 的 HTML

```html file=cookbook/sse/broadcast/index.html
```
