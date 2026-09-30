---
title: Server-Sent Events (SSE)
description: Transmite server-sent events desde un handler Echo, por conexión o como broadcast a muchos clientes.
sidebar:
  order: 14
---

[Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)
se pueden usar de varias formas. El primer ejemplo de abajo es SSE por conexión y por handler.
Para lógica de broadcast más compleja, consulta el segundo ejemplo usando
[r3labs/sse](https://github.com/r3labs/sse).

:::caution
Las conexiones SSE son de larga duración, por lo que debe deshabilitarse el write timeout del servidor.
Ambos ejemplos establecen `s.WriteTimeout = 0` mediante `BeforeServeFunc`.
:::

## Usar SSE

### Servidor

El handler escribe los headers SSE y luego emite un evento cada segundo hasta que el cliente
se desconecta. `http.NewResponseController(w).Flush()` empuja cada evento al cliente de inmediato.

```go file=cookbook/sse/simple/server.go
```

### Estructura Event y método Marshal

```go file=cookbook/sse/simple/serversentevent.go
```

### HTML que sirve SSE

```html file=cookbook/sse/simple/index.html
```

## Broadcast con r3labs/sse

Cuando necesitas hacer broadcast de un único stream de eventos a muchos suscriptores, la
biblioteca [r3labs/sse](https://github.com/r3labs/sse) maneja por ti la gestión de streams
y suscriptores.

### Servidor

```go file=cookbook/sse/broadcast/server.go
```

### HTML que sirve SSE

```html file=cookbook/sse/broadcast/index.html
```
