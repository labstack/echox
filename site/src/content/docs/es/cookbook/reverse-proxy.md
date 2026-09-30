---
title: Reverse Proxy
description: Usa Echo como reverse proxy y load balancer delante de aplicaciones upstream.
sidebar:
  order: 19
---

Esta receta demuestra cómo usar Echo como reverse proxy y load balancer delante de tus
aplicaciones, como WordPress, Node.js, Java, Python, Ruby o Go. Para simplificar, aquí
los upstreams son servidores Go que también manejan WebSocket.

## 1) Identificar URL(s) de destino upstream

```go file=cookbook/reverse-proxy/server.go#targets
```

## 2) Configurar middleware proxy con destinos upstream

El snippet de abajo usa load balancing round-robin. También puedes usar
`middleware.NewRandomBalancer()`.

```go file=cookbook/reverse-proxy/server.go#middleware
```

Para configurar un proxy para una subruta, usa `Echo#Group()`.

```go file=cookbook/reverse-proxy/server.go#grouped-proxy
```

```sh
go run . -grouped
```

## 3) Iniciar servidores upstream

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) Iniciar el servidor proxy

```sh
go run server.go
```

Abre `http://localhost:1323`, y deberías ver una página web con un request HTTP
servido desde "server 1" y un request WebSocket servido desde "server 2".

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Código fuente

### Servidor upstream

```go file=cookbook/reverse-proxy/upstream/server.go
```

### Servidor proxy

```go file=cookbook/reverse-proxy/server.go
```

## Actualización de seguridad de Echo (v5.4.0 / v4.16.0)

Proxy reenvía `X-Forwarded-Proto` desde `Context#Scheme()` y elimina las otras cabeceras de esquema. En v5, `X-Real-IP` procede de `Context#RealIP()`. Si hay otro proxy delante de Echo, configura los extractores de [esquema](/es/guide/request-scheme/) e [IP](/es/guide/ip-address/) para que el upstream reciba valores correctos.
