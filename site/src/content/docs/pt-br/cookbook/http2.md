---
title: Servidor HTTP/2
description: Sirva tráfego por HTTP/2 iniciando Echo com um certificado TLS.
sidebar:
  order: 9
---

HTTP/2 melhora a latência por meio de multiplexação de requests, compressão de headers e
server push. O servidor HTTP do Go negocia HTTP/2 automaticamente sobre TLS, então servir
HTTP/2 com Echo é uma questão de iniciar o servidor com um certificado.

<a id="1-generate-a-self-signed-x509-tls-certificate"></a>

## 1. Gerar um certificado TLS X.509 autoassinado

Execute o comando a seguir para gerar `cert.pem` e `key.pem`:

```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

:::note
Para fins de demonstração, usamos um certificado autoassinado. Em produção, obtenha
um certificado de uma [certificate authority](https://en.wikipedia.org/wiki/Certificate_authority).
:::

## 2. Criar um handler que ecoa informações do request

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

## 3. Iniciar o servidor TLS

Inicie o servidor com o certificado e a chave gerados:

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

Como alternativa, use um `http.Server` customizado com seu próprio `tls.Config`:

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

Inicie o servidor e acesse `https://localhost:1323/request`. Você deve ver
uma saída semelhante a:

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## Código-fonte

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
