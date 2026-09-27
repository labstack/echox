---
title: Servidor HTTP/2
description: Sirve tráfico sobre HTTP/2 iniciando Echo con un certificado TLS.
sidebar:
  order: 9
---

HTTP/2 mejora la latencia mediante multiplexing de requests, compresión de headers y
server push. El servidor HTTP de Go negocia HTTP/2 automáticamente sobre TLS, por lo que
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

```go
e.GET("/request", func(c *echo.Context) error {
	req := c.Request()
	format := `
		<code>
			Protocol: %s<br>
			Host: %s<br>
			Remote Address: %s<br>
			Method: %s<br>
			Path: %s<br>
		</code>
	`
	return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
})
```

## 3. Iniciar el servidor TLS

Inicia el servidor con el certificado y la key generados:

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

Alternativamente, usa un `http.Server` personalizado con tu propio `tls.Config`:

```go
s := http.Server{
  Addr:    ":8443",
  Handler: e, // set Echo as handler
  TLSConfig: &tls.Config{
    //Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
  },
  //ReadTimeout: 30 * time.Second, // use custom timeouts
}
if err := s.ListenAndServeTLS("cert.pem", "key.pem"); err != http.ErrServerClosed {
  log.Fatal(err)
}
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

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	e.GET("/request", func(c *echo.Context) error {
		req := c.Request()
		format := `
			<code>
				Protocol: %s<br>
				Host: %s<br>
				Remote Address: %s<br>
				Method: %s<br>
				Path: %s<br>
			</code>
		`
		return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
	})
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```
