---
title: WebSocket
description: Maneja conexiones WebSocket en Echo usando golang.org/x/net/websocket o gorilla/websocket.
sidebar:
  order: 16
---

Los handlers de Echo pueden servir conexiones WebSocket actualizando la conexión HTTP
subyacente. Esta receta muestra dos enfoques: el paquete estándar
`golang.org/x/net/websocket` y la biblioteca popular
[`gorilla/websocket`](https://github.com/gorilla/websocket).

## Usar net WebSocket

### Servidor

```go file=cookbook/websocket/net/server.go
```

## Usar gorilla WebSocket

### Servidor

```go file=cookbook/websocket/gorilla/server.go
```

## Cliente

```html file=cookbook/websocket/public/index.html
```

## Salida

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
