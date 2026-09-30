---
title: Reverse Proxy
description: Use Echo como reverse proxy e load balancer na frente de aplicações upstream.
sidebar:
  order: 19
---

Esta receita demonstra como usar Echo como reverse proxy e load balancer na
frente das suas aplicações, como WordPress, Node.js, Java, Python, Ruby ou Go.
Por simplicidade, os upstreams aqui são servidores Go que também tratam WebSocket.

## 1) Identificar URL(s) de destino upstream

```go file=cookbook/reverse-proxy/server.go#targets
```

## 2) Configurar middleware de proxy com destinos upstream

O snippet abaixo usa load balancing round-robin. Você também pode usar
`middleware.NewRandomBalancer()`.

```go file=cookbook/reverse-proxy/server.go#middleware
```

Para configurar um proxy para uma sub-rota, use `Echo#Group()`.

```go file=cookbook/reverse-proxy/server.go#grouped-proxy
```

## 3) Iniciar servidores upstream

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) Iniciar o servidor proxy

```sh
go run server.go
```

Acesse `http://localhost:1323`, e você deverá ver uma página com um request HTTP
servido pelo "server 1" e um request WebSocket servido pelo "server 2".

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Código-fonte

### Servidor upstream

```go file=cookbook/reverse-proxy/upstream/server.go
```

### Servidor proxy

```go file=cookbook/reverse-proxy/server.go
```

## Atualização de segurança do Echo (v5.4.0 / v4.16.0)

Proxy define `X-Forwarded-Proto` a partir de `Context#Scheme()` e remove os outros cabeçalhos de esquema. No v5, `X-Real-IP` vem de `Context#RealIP()`. Se houver outro proxy antes do Echo, configure os extratores de [esquema](/pt-br/guide/request-scheme/) e [IP](/pt-br/guide/ip-address/) para transmitir valores corretos ao upstream.
