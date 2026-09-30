---
title: WebSocket
description: Trate conexões WebSocket no Echo usando golang.org/x/net/websocket ou gorilla/websocket.
sidebar:
  order: 16
---

Handlers Echo podem servir conexões WebSocket fazendo upgrade da conexão
HTTP subjacente. Esta receita mostra duas abordagens: o pacote padrão
`golang.org/x/net/websocket` e a biblioteca popular
[`gorilla/websocket`](https://github.com/gorilla/websocket).

## Usando net WebSocket

### Servidor

```go file=cookbook/websocket/net/server.go
```

## Usando gorilla WebSocket

### Servidor

```go file=cookbook/websocket/gorilla/server.go
```

## Cliente

```html file=cookbook/websocket/public/index.html
```

## Saída

**Servidor**

```sh
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
Hello, Server!
```

**Cliente**

```sh
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
Hello, Client!
```
