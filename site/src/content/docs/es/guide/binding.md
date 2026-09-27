---
title: Binding
description: Analiza datos de request en structs Go tipados desde path, query, header y body.
sidebar:
  order: 5
---

Analizar datos de request es una parte crucial de una aplicación web. En Echo esto se llama
_binding_, y puede leer desde cuatro partes de un request HTTP:

- Parámetros de path de URL
- Parámetros de query de URL
- Headers
- Body del request

## Binding con tags de struct

Define un struct con tags que especifican el origen de datos y la clave, y luego llama a `c.Bind()`
con un puntero a él. Aquí el parámetro de query `id` se vincula al campo `ID`:

```go
type User struct {
	ID string `query:"id"`
}

// handler for /users?id=<userID>
var user User
if err := c.Bind(&user); err != nil {
	return c.String(http.StatusBadRequest, "bad request")
}
```

### Orígenes de datos

| Tag      | Origen |
| -------- | ------ |
| `query`  | Parámetro de query |
| `param`  | Parámetro de path |
| `header` | Valor de header |
| `form`   | Datos de formulario (query + body) |
| `json`   | Body del request (`encoding/json`) |
| `xml`    | Body del request (`encoding/xml`) |

Los campos de path, query, header y form requieren un **tag explícito**. JSON y XML usan
el nombre del campo del struct cuando se omite el tag, igual que la biblioteca estándar.

### Slices

Los valores repetidos de query, path, formulario o header se enlazan a un campo
slice, que recoge cada aparición:

```go
// GET /search?tag=go&tag=web&tag=api
type Filter struct {
	Tags []string `query:"tag"`
}

var f Filter
if err := c.Bind(&f); err != nil {
	return c.String(http.StatusBadRequest, "bad request")
}
// f.Tags == []string{"go", "web", "api"}
```

### Content types del body

Al decodificar el body del request, el header `Content-Type` selecciona el decoder:

- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`

### Múltiples orígenes y precedencia

Un campo puede declarar varios orígenes. Los datos se vinculan en este orden, y cada paso
sobrescribe el anterior:

1. Parámetros de path
2. Parámetros de query (solo GET / DELETE)
3. Body del request

```go
type User struct {
	ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

### Binding directo desde un origen

```go
echo.BindBody(c, &payload)        // request body
echo.BindQueryParams(c, &payload) // query parameters
echo.BindPathValues(c, &payload)  // path parameters
echo.BindHeaders(c, &payload)     // headers
```

:::note
Los headers **no** están incluidos por `c.Bind()`. Vincúlalos directamente con `echo.BindHeaders`.
:::

:::caution[Seguridad]
No vincules directamente a structs de negocio. Si un struct vinculado expone un campo `IsAdmin bool`,
un body de request `{"IsAdmin": true}` lo establecería. Usa un DTO dedicado y mapéalo explícitamente:
:::

```go
type UserDTO struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

e.POST("/users", func(c *echo.Context) error {
	var dto UserDTO
	if err := c.Bind(&dto); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	user := User{Name: dto.Name, Email: dto.Email, IsAdmin: false}
	executeSomeBusinessLogic(user)
	return c.JSON(http.StatusOK, user)
})
```

## Binding fluido

Para binding explícito y type-safe desde un único origen, usa los binders fluidos. Encadenan
la configuración y la ejecutan, recopilando errores:

```go
// /api/search?active=true&id=1&id=2&id=3&length=25
var opts struct {
	IDs    []int64
	Active bool
}
length := int64(50)

err := echo.QueryParamsBinder(c).
	Int64("length", &length).
	Int64s("id", &opts.IDs).
	Bool("active", &opts.Active).
	BindError() // first error, if any
```

Binders disponibles: `echo.QueryParamsBinder(c)`, `echo.PathValuesBinder(c)`,
`echo.FormFieldBinder(c)`. Termina una cadena con `BindError()` (primer error) o
`BindErrors()` (todos los errores). `FailFast(false)` ejecuta toda la cadena; está activado por defecto.

Cada tipo soportado ofrece métodos `Type(...)`, `MustType(...)`, `Types(...)` (slices) y
`MustTypes(...)`, por ejemplo `Int64`, `MustInt64`, `Int64s`. Usa
`BindWithDelimiter("id", &dest, ",")` para separar valores unidos por comas.

<a id="custom-binder"></a>

## Binder personalizado

Registra un binder personalizado mediante `Echo#Binder`:

```go
type CustomBinder struct{}

func (cb *CustomBinder) Bind(c *echo.Context, i any) error {
	db := new(echo.DefaultBinder)
	if err := db.Bind(c, i); err != echo.ErrUnsupportedMediaType {
		return err
	}
	// custom logic here
	return nil
}

e.Binder = &CustomBinder{}
```
