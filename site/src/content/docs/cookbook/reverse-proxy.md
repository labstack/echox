---
title: Reverse Proxy
description: Use Echo as a reverse proxy and load balancer in front of upstream applications.
sidebar:
  order: 19
---

This recipe demonstrates how to use Echo as a reverse proxy and load balancer in
front of your applications, such as WordPress, Node.js, Java, Python, Ruby, or Go.
For simplicity, the upstreams here are Go servers that also handle WebSocket.

Proxy forwards `X-Forwarded-Proto` from `Context#Scheme()` and drops the older
`X-Forwarded-Ssl`, `X-Forwarded-Protocol`, and `X-Url-Scheme` headers. In v5 it
sets `X-Real-IP` from `Context#RealIP()`. If another proxy sits before Echo, set
the [scheme extractor](/guide/request-scheme/) and [IP extractor](/guide/ip-address/)
for that trusted proxy; otherwise the upstream may see the proxy's IP or the wrong
scheme.

## 1) Identify upstream target URL(s)

```go file=cookbook/reverse-proxy/server.go#targets
```

## 2) Set up proxy middleware with upstream targets

The snippet below uses round-robin load balancing. You can also use
`middleware.NewRandomBalancer()`.

```go file=cookbook/reverse-proxy/server.go#middleware
```

To set up a proxy for a sub-route, use `Echo#Group()`.

```go file=cookbook/reverse-proxy/server.go#grouped-proxy
```

```sh
go run . -grouped
```

## 3) Start upstream servers

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) Start the proxy server

```sh
go run server.go
```

Browse to `http://localhost:1323`, and you should see a webpage with an HTTP
request served from "server 1" and a WebSocket request served from "server 2".

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Source code

### Upstream server

```go file=cookbook/reverse-proxy/upstream/server.go
```

### Proxy server

```go file=cookbook/reverse-proxy/server.go
```
