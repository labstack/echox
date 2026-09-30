---
title: HTTP/2 Server
description: Serve traffic over HTTP/2 by starting Echo with a TLS certificate.
sidebar:
  order: 9
---

HTTP/2 improves latency through request multiplexing, header compression, and
server push. Go's HTTP server negotiates HTTP/2 automatically over TLS, so serving
HTTP/2 with Echo is a matter of starting the server with a certificate.

## 1. Generate a self-signed X.509 TLS certificate

Run the following command to generate `cert.pem` and `key.pem`:

```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

:::note
For demonstration purposes we use a self-signed certificate. In production, obtain
a certificate from a [certificate authority](https://en.wikipedia.org/wiki/Certificate_authority).
:::

## 2. Create a handler that echoes request information

```go file=cookbook/http2/server.go#handler
```

## 3. Start the TLS server

Start the server with the generated certificate and key:

```go file=cookbook/http2/server.go#start-tls
```

Alternatively, use a custom `http.Server` with your own `tls.Config`:

```go file=cookbook/http2/server.go#custom-server
```

## 4. Verify

Start the server and browse to `https://localhost:1323/request`. You should see
output similar to:

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## Source code

```go file=cookbook/http2/server.go
```
