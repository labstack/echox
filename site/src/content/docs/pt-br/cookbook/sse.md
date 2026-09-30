---
title: Server-Sent Events (SSE)
description: Faça streaming de server-sent events a partir de um handler Echo, por conexão ou em broadcast para muitos clientes.
sidebar:
  order: 14
---

[Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format)
podem ser usados de várias formas. O primeiro exemplo abaixo é SSE por conexão e por handler.
Para lógica de broadcast mais complexa, veja o segundo exemplo usando
[r3labs/sse](https://github.com/r3labs/sse).

:::caution
Conexões SSE são long-lived, então o write timeout do servidor precisa ser desabilitado.
Os dois exemplos definem `s.WriteTimeout = 0` via `BeforeServeFunc`.
:::

## Usando SSE

### Servidor

O handler escreve os headers SSE e então emite um evento a cada segundo até que o
cliente desconecte. `http.NewResponseController(w).Flush()` envia cada evento ao
cliente imediatamente.

```go file=cookbook/sse/simple/server.go
```

### Estrutura Event e método Marshal

```go file=cookbook/sse/simple/serversentevent.go
```

### HTML servindo SSE

```html file=cookbook/sse/simple/index.html
```

## Broadcast com r3labs/sse

Quando você precisa transmitir um único stream de eventos para muitos subscribers, a
biblioteca [r3labs/sse](https://github.com/r3labs/sse) cuida do gerenciamento de stream e subscribers
para você.

### Servidor

```go file=cookbook/sse/broadcast/server.go
```

### HTML servindo SSE

```html file=cookbook/sse/broadcast/index.html
```
