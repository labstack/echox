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

```go file=cookbook/http2/server.go#handler
```

## 3. Iniciar o servidor TLS

Inicie o servidor com o certificado e a chave gerados:

```go file=cookbook/http2/server.go#start-tls
```

Como alternativa, use um `http.Server` customizado com seu próprio `tls.Config`:

```go file=cookbook/http2/server.go#custom-server
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

```go file=cookbook/http2/server.go
```
