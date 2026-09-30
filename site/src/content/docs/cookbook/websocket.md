---
title: WebSocket
description: Handle WebSocket connections in Echo using golang.org/x/net/websocket or gorilla/websocket.
sidebar:
  order: 16
---

Echo handlers can serve WebSocket connections by upgrading the underlying
HTTP connection. This recipe shows two approaches: the standard
`golang.org/x/net/websocket` package and the popular
[`gorilla/websocket`](https://github.com/gorilla/websocket) library.

## Using net WebSocket

### Server

```go file=cookbook/websocket/net/server.go
```

## Using gorilla WebSocket

### Server

```go file=cookbook/websocket/gorilla/server.go
```

## Client

```html file=cookbook/websocket/public/index.html
```

## Output

**Server**

```sh
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
```

**Client**

```sh
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
```
