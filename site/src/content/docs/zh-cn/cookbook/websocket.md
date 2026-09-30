---
title: WebSocket
description: 在 Echo 中使用 golang.org/x/net/websocket 或 gorilla/websocket 处理 WebSocket 连接。
sidebar:
  order: 16
---

Echo 处理函数可以通过升级底层 HTTP 连接来提供 WebSocket 连接。此示例展示两种方式：
标准的 `golang.org/x/net/websocket` 包，以及流行的
[`gorilla/websocket`](https://github.com/gorilla/websocket) 库。

## 使用 net WebSocket

### 服务器

```go file=cookbook/websocket/net/server.go
```

## 使用 gorilla WebSocket

### 服务器

```go file=cookbook/websocket/gorilla/server.go
```

## 客户端

```html file=cookbook/websocket/public/index.html
```

## 输出

**服务器**

```sh
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
```

**客户端**

```sh
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
```
