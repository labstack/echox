---
title: Servidor HTTP/2
description: Sirve tráfico sobre HTTP/2 iniciando Echo con un certificado TLS.
sidebar:
  order: 9
---

HTTP/2 mejora la latencia mediante multiplexing de requests y compresión de headers. El servidor HTTP de Go negocia HTTP/2 automáticamente sobre TLS, por lo que
servir HTTP/2 con Echo consiste en iniciar el servidor con un certificado.

<a id="1-generate-a-self-signed-x509-tls-certificate"></a>

## 1. Generar un certificado TLS X.509 autofirmado

Ejecuta el siguiente comando para generar `cert.pem` y `key.pem`:

```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

:::note
Con fines de demostración usamos un certificado autofirmado. En producción, obtén
un certificado de una [certificate authority](https://en.wikipedia.org/wiki/Certificate_authority).
:::

## 2. Crear un handler que refleje información del request

```go file=cookbook/http2/server.go#handler
```

## 3. Iniciar el servidor TLS

Inicia el servidor con el certificado y la key generados:

```go file=cookbook/http2/server.go#start-tls
```

Alternativamente, usa un `http.Server` personalizado con tu propio `tls.Config`:

```go file=cookbook/http2/server.go#custom-server
```

```sh
go run . -custom-server
```

## 4. Verificar

Inicia el servidor y abre `https://localhost:1323/request`. Deberías ver una salida
similar a:

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## Código fuente

```go file=cookbook/http2/server.go
```
