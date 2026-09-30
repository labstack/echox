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

## 3) Start upstream servers

From `cookbook/reverse-proxy`, start each upstream in its own terminal:

```sh
go run ./upstream server1 :8081
```

```sh
go run ./upstream server2 :8082
```

## 4) Start the proxy server

In a third terminal, also in `cookbook/reverse-proxy`, start the proxy:

```sh
go run .
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

To use grouped mode instead, stop the default proxy and run:

```sh
go run . -grouped
```

Browse to `http://localhost:1323/blog/`. The proxy strips `/blog` before forwarding requests; the page’s relative WebSocket URL stays under `/blog/`.

## Source code

### Upstream server

```go file=cookbook/reverse-proxy/upstream/server.go
```

### Proxy server

```go file=cookbook/reverse-proxy/server.go
```
