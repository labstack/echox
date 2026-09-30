---
title: Apagado graceful
description: Drena requests en curso antes de detener el servidor ante una señal de interrupción.
sidebar:
  order: 8
---

Un apagado graceful permite que los requests en curso terminen antes de que el proceso salga.
El enfoque más simple es pasar un contexto cancelable a `StartConfig.Start` y establecer
un `GracefulTimeout`. Cuando el contexto se cancela por una señal de interrupción, Echo
deja de aceptar nuevas conexiones y espera hasta el timeout para que los requests activos
se completen.

## Servidor

```go file=cookbook/graceful-shutdown/server.go#primary-server
```

## Usar un servidor HTTP personalizado

Si administras el `http.Server` por tu cuenta, inícialo en una goroutine, espera el
contexto de señal y luego llama a `Shutdown` con un timeout:

```go file=cookbook/graceful-shutdown/server.go#custom-server
```
