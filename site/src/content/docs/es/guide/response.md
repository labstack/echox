---
title: Response
description: Envía strings, HTML, JSON, XML, archivos, streams, redirects y hooks de response.
sidebar:
  order: 8
---

Un handler escribe su response a través de `echo.Context`. Cada helper establece por ti el
`Content-Type` adecuado y el código de estado.

## Enviar string

`Context#String(code int, s string)` envía una response de texto plano con un código de estado.

```go
func(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

## Enviar HTML

`Context#HTML(code int, html string)` envía una response HTML simple con un código de estado.
Para generar HTML dinámicamente, consulta [Templates](/es/guide/templates/).

```go
func(c *echo.Context) error {
	return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
}
```

### Enviar blob HTML

`Context#HTMLBlob(code int, b []byte)` envía un blob HTML con un código de estado. Es útil
con un template engine que produce `[]byte`.

```go
func handler(c *echo.Context) error {
	blob := []byte("<strong>Hello, World!</strong>")
	return c.HTMLBlob(http.StatusOK, blob)
}
```

## Renderizar template

Consulta [Templates](/es/guide/templates/).

## Enviar JSON

`Context#JSON(code int, i any)` codifica un valor Go como JSON y lo envía con un código de estado.

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

`Context#JSON()` usa `json.Marshal` internamente, lo que puede ser ineficiente para payloads
grandes. En ese caso, transmite el JSON directamente:

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

`Context#JSONPretty(code int, i any, indent string)` envía una response JSON con formato bonito.
La indentación puede ser espacios o tabs.

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

`Context#JSONBlob(code int, b []byte)` envía directamente un blob JSON precodificado, por
ejemplo desde una base de datos.

```go
func(c *echo.Context) error {
	encodedJSON := []byte{} // Encoded JSON from an external source.
	return c.JSONBlob(http.StatusOK, encodedJSON)
}
```

## Enviar JSONP

`Context#JSONP(code int, callback string, i any)` codifica un valor Go como JSON y lo envía
como payload JSONP envuelto en el callback dado.

```go
func handler(c *echo.Context) error {
	callback := c.QueryParam("callback")
	return c.JSONP(http.StatusOK, callback, &User{Name: "Jon", Email: "jon@labstack.com"})
}
```

Para solicitudes entre orígenes, usa [JSON con CORS](/es/cookbook/cors/).

## Enviar XML

`Context#XML(code int, i any)` codifica un valor Go como XML y lo envía con un código de estado.

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

`Context#XML` usa `xml.Marshal` internamente, lo que puede ser ineficiente para payloads
grandes. En ese caso, transmite el XML directamente:

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

`Context#XMLPretty(code int, i any, indent string)` envía una response XML con formato bonito.
La indentación puede ser espacios o tabs.

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
También puedes hacer que `Context#XML()` produzca XML con formato bonito agregando `pretty`
al query string de la URL del request.

```sh
curl http://localhost:1323/users/1?pretty
```
:::

### Blob XML

`Context#XMLBlob(code int, b []byte)` envía directamente un blob XML precodificado, por
ejemplo desde una base de datos.

```go
func(c *echo.Context) error {
	encodedXML := []byte{} // Encoded XML from an external source.
	return c.XMLBlob(http.StatusOK, encodedXML)
}
```

## Enviar archivo

`Context#File(file string)` envía el contenido de un archivo como response. Establece
el content type correcto y maneja el caching automáticamente.

```go
func(c *echo.Context) error {
	return c.File("<PATH_TO_YOUR_FILE>")
}
```

## Enviar attachment

`Context#Attachment(file, name string)` es como `File()`, pero envía el archivo con
`Content-Disposition: attachment` y el nombre dado.

```go
func(c *echo.Context) error {
	return c.Attachment("<PATH_TO_YOUR_FILE>", "<ATTACHMENT_NAME>")
}
```

## Enviar inline

`Context#Inline(file, name string)` es como `File()`, pero envía el archivo con
`Content-Disposition: inline` y el nombre dado.

```go
func(c *echo.Context) error {
	return c.Inline("<PATH_TO_YOUR_FILE>", "<INLINE_NAME>")
}
```

## Enviar blob

`Context#Blob(code int, contentType string, b []byte)` envía datos arbitrarios con un
content type y código de estado dados.

```go
func(c *echo.Context) error {
	data := []byte(`0306703,0035866,NO_ACTION,06/19/2006
0086003,"0005866",UPDATED,06/19/2006`)
	return c.Blob(http.StatusOK, "text/csv", data)
}
```

## Enviar stream

`Context#Stream(code int, contentType string, r io.Reader)` envía un stream de datos
arbitrario con un content type, `io.Reader` y código de estado dados.

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

## Enviar sin contenido

`Context#NoContent(code int)` envía un body vacío con un código de estado.

```go
func(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}
```

## Redirigir request

`Context#Redirect(code int, url string)` redirige el request a la URL dada con un código de estado.

```go
func(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "<URL>")
}
```

## Hooks

### Antes de la response

`Response#Before(func())` registra una función que se ejecuta justo antes de escribir la response.

### Después de la response

`Response#After(func())` registra una función que se ejecuta justo después de escribir la response.
Si `Content-Length` es desconocido, no se ejecuta ninguna función after.

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
Puedes registrar varias funciones `Before` y `After`.
:::

## Actualización de seguridad de Echo (v5.4.0 / v4.16.0)

`Context#JSONP` acepta un callback vacío, un identificador JavaScript o una ruta de identificadores separados por puntos. Los demás valores devuelven HTTP 400 con `ErrInvalidJSONPCallback`; la respuesta válida incluye `X-Content-Type-Options: nosniff`. No uses JSONP para datos privados: cualquier sitio puede leer la respuesta con las cookies del usuario. Usa JSON con CORS.
