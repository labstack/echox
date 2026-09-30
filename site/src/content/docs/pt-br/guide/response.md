---
title: Response
description: Envie strings, HTML, JSON, XML, arquivos, streams, redirects e hooks de response.
sidebar:
  order: 8
---

Um handler escreve sua response por meio do `echo.Context`. Cada helper define o
`Content-Type` e o código de status apropriados para você.

## Enviar string

`Context#String(code int, s string)` envia uma response de texto simples com um código de status.

```go
func(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

## Enviar HTML

`Context#HTML(code int, html string)` envia uma response HTML simples com um código de
status. Para gerar HTML dinamicamente, veja [Templates](/pt-br/guide/templates/).

```go
func(c *echo.Context) error {
	return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
}
```

### Enviar blob HTML

`Context#HTMLBlob(code int, b []byte)` envia um blob HTML com um código de status. Ele é
útil com um mecanismo de templates que produz `[]byte`.

```go
func handler(c *echo.Context) error {
	blob := []byte("<strong>Hello, World!</strong>")
	return c.HTMLBlob(http.StatusOK, blob)
}
```

## Renderizar template

Veja [Templates](/pt-br/guide/templates/).

## Enviar JSON

`Context#JSON(code int, i any)` codifica um valor Go como JSON e o envia com um
código de status.

```go
type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.JSON(http.StatusOK, u)
}
```

### Stream JSON

`Context#JSON()` usa `json.Marshal` internamente, o que pode ser ineficiente para payloads grandes.
Nesse caso, envie o JSON diretamente como stream:

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(u)
}
```

### JSON pretty

`Context#JSONPretty(code int, i any, indent string)` envia uma response JSON formatada.
O indent pode ser espaços ou tabs.

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}
```

```json
{
  "email": "jon@labstack.com",
  "name": "Jon"
}
```

### Blob JSON

`Context#JSONBlob(code int, b []byte)` envia um blob JSON pré-codificado diretamente, por
exemplo vindo de um banco de dados.

```go
func(c *echo.Context) error {
	encodedJSON := []byte{} // Encoded JSON from an external source.
	return c.JSONBlob(http.StatusOK, encodedJSON)
}
```

## Enviar JSONP

`Context#JSONP(code int, callback string, i any)` codifica um valor Go como JSON e
o envia como um payload JSONP envolvido no callback informado.

```go
func handler(c *echo.Context) error {
	callback := c.QueryParam("callback")
	return c.JSONP(http.StatusOK, callback, &User{Name: "Jon", Email: "jon@labstack.com"})
}
```

Para requisições entre origens, use [JSON com CORS](/pt-br/cookbook/cors/).

## Enviar XML

`Context#XML(code int, i any)` codifica um valor Go como XML e o envia com um código
de status.

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.XML(http.StatusOK, u)
}
```

### Stream XML

`Context#XML` usa `xml.Marshal` internamente, o que pode ser ineficiente para payloads grandes.
Nesse caso, envie o XML diretamente como stream:

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return xml.NewEncoder(c.Response()).Encode(u)
}
```

### XML pretty

`Context#XMLPretty(code int, i any, indent string)` envia uma response XML formatada.
O indent pode ser espaços ou tabs.

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.XMLPretty(http.StatusOK, u, "  ")
}
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<User>
  <Name>Jon</Name>
  <Email>jon@labstack.com</Email>
</User>
```

:::tip
Você também pode fazer `Context#XML()` gerar XML formatado adicionando `pretty`
à query string da URL do request.

```sh
curl http://localhost:1323/users/1?pretty
```
:::

### Blob XML

`Context#XMLBlob(code int, b []byte)` envia um blob XML pré-codificado diretamente, por
exemplo vindo de um banco de dados.

```go
func(c *echo.Context) error {
	encodedXML := []byte{} // Encoded XML from an external source.
	return c.XMLBlob(http.StatusOK, encodedXML)
}
```

## Enviar arquivo

`Context#File(file string)` envia o conteúdo de um arquivo como a response. Ele define
o tipo de conteúdo correto e lida com cache automaticamente.

```go
func(c *echo.Context) error {
	return c.File("<PATH_TO_YOUR_FILE>")
}
```

## Enviar attachment

`Context#Attachment(file, name string)` é como `File()`, mas envia o arquivo com
`Content-Disposition: attachment` e o nome informado.

```go
func(c *echo.Context) error {
	return c.Attachment("<PATH_TO_YOUR_FILE>", "<ATTACHMENT_NAME>")
}
```

## Enviar inline

`Context#Inline(file, name string)` é como `File()`, mas envia o arquivo com
`Content-Disposition: inline` e o nome informado.

```go
func(c *echo.Context) error {
	return c.Inline("<PATH_TO_YOUR_FILE>", "<INLINE_NAME>")
}
```

## Enviar blob

`Context#Blob(code int, contentType string, b []byte)` envia dados arbitrários com um
tipo de conteúdo e código de status informados.

```go
func(c *echo.Context) error {
	data := []byte(`0306703,0035866,NO_ACTION,06/19/2006
0086003,"0005866",UPDATED,06/19/2006`)
	return c.Blob(http.StatusOK, "text/csv", data)
}
```

## Enviar stream

`Context#Stream(code int, contentType string, r io.Reader)` envia um stream de dados
arbitrário com um tipo de conteúdo, `io.Reader` e código de status informados.

```go
func(c *echo.Context) error {
	f, err := os.Open("<PATH_TO_IMAGE>")
	if err != nil {
		return err
	}
	defer f.Close()
	return c.Stream(http.StatusOK, "image/png", f)
}
```

## Enviar sem conteúdo

`Context#NoContent(code int)` envia um body vazio com um código de status.

```go
func(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}
```

## Redirecionar request

`Context#Redirect(code int, url string)` redireciona o request para a URL informada com
um código de status.

```go
func(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "<URL>")
}
```

## Hooks

### Antes da response

`Response#Before(func())` registra uma função que roda pouco antes de a response ser
escrita.

### Depois da response

`Response#After(func())` registra uma função que roda logo depois de a response ser
escrita. Se o `Content-Length` for desconhecido, nenhuma função after roda.

```go
e.GET("/hooks", func(c *echo.Context) error {
	resp, err := echo.UnwrapResponse(c.Response())
	if err != nil {
		return err
	}
	resp.Before(func() {
		println("before response")
	})
	resp.After(func() {
		println("after response")
	})
	return c.String(http.StatusOK, "Hello, World!")
})
```

:::tip
Você pode registrar várias funções `Before` e `After`.
:::

## Atualização de segurança do Echo (v5.4.0 / v4.16.0)

`Context#JSONP` aceita callback vazio, identificador JavaScript ou caminho de identificadores separados por pontos. Os demais valores retornam HTTP 400 com `ErrInvalidJSONPCallback`; uma resposta válida inclui `X-Content-Type-Options: nosniff`. Não use JSONP para dados privados: qualquer site pode ler a resposta com os cookies do usuário. Use JSON com CORS.
