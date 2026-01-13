---
description: Reverse proxy recipe
---

# Reverse Proxy

This recipe demonstrates how you can use Echo as a reverse proxy server and load balancer in front of your favorite applications like WordPress, Node.js, Java, Python, Ruby or even Go. For simplicity, I will use Go upstream servers with WebSocket.

## 1) Identify upstream target URL(s)

```go
url1, err := url.Parse("http://localhost:8081")
if err != nil {
  e.Logger.Error("failed parse url", "error", err)
}
url2, err := url.Parse("http://localhost:8082")
if err != nil {
  e.Logger.Error("failed parse url", "error", err)
}
targets := []*middleware.ProxyTarget{
  {
    URL: url1,
  },
  {
    URL: url2,
  },
}
```

## 2) Setup proxy middleware with upstream targets

In the following code snippet we are using round-robin load balancing technique. You may also use `middleware.NewRandomBalancer()`.

```go
e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
```

To setup proxy for a sub-route use `Echo#Group()`.

```go
g := e.Group("/blog")
g.Use(middleware.Proxy(...))
```

## 3) Start upstream servers

- `cd upstream`
- `go run server.go server1 :8081`
- `go run server.go server2 :8082`

## 4) Start the proxy server

```sh
go run server.go
```

Browse to http://localhost:1323, and you should see a webpage with an HTTP request being served from "server 1" and a WebSocket request being served from "server 2."

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## Source Code

```go reference
https://github.com/labstack/echox/blob/master/cookbook/reverse-proxy/upstream/server.go
```

```go reference
https://github.com/labstack/echox/blob/master/cookbook/reverse-proxy/server.go
```
