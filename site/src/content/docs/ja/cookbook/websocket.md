---
title: WebSocket
description: Echo で golang.org/x/net/websocket または gorilla/websocket を使って WebSocket 接続を処理します。
sidebar:
  order: 16
---

Echo ハンドラは、基礎となる HTTP 接続をアップグレードして WebSocket 接続を提供できます。
このレシピでは 2 つの方法を示します。標準の `golang.org/x/net/websocket` パッケージと、
広く使われている [`gorilla/websocket`](https://github.com/gorilla/websocket) ライブラリです。

## net WebSocket を使う

### サーバー

```go file=cookbook/websocket/net/server.go
```

## gorilla WebSocket を使う

### サーバー

```go file=cookbook/websocket/gorilla/server.go
```

## クライアント

```html file=cookbook/websocket/public/index.html
```

## 出力

**サーバー**

```sh
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
```

**クライアント**

```sh
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
```
