---
title: Binding
description: Analise dados de request em structs Go tipadas a partir de path, query, header e body.
sidebar:
  order: 5
---

Analisar dados de request é uma parte crucial de uma aplicação web. No Echo isso se chama
_binding_, e ele pode ler de quatro partes de um request HTTP:

- Parâmetros de caminho da URL
- Parâmetros de query da URL
- Headers
- Body do request

## Binding com tags de struct

Defina uma struct com tags que especificam a fonte e a chave dos dados, então chame `c.Bind()`
com um ponteiro para ela. Aqui o parâmetro de query `id` faz binding para o campo `ID`:

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

### Fontes de dados

| Tag      | Fonte |
| -------- | ----- |
| `query`  | Parâmetro de query |
| `param`  | Parâmetro de caminho |
| `header` | Valor de header |
| `form`   | Dados de formulário (query + body) |
| `json`   | Body do request (`encoding/json`) |
| `xml`    | Body do request (`encoding/xml`) |

Campos de path, query, header e form exigem uma **tag explícita**. JSON e XML usam
o nome do campo da struct quando a tag é omitida, igual à biblioteca padrão.

### Slices

Valores repetidos de query, path, formulário ou header são vinculados a um campo
slice, que coleta cada ocorrência:

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

### Tipos de conteúdo do body

Ao decodificar o body do request, o header `Content-Type` seleciona o decoder:

- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`

### Múltiplas fontes e precedência

Um campo pode declarar várias fontes. Os dados fazem binding nesta ordem, cada etapa
sobrescrevendo a anterior:

1. Parâmetros de caminho
2. Parâmetros de query (somente GET / DELETE)
3. Body do request

```go
type User struct {
	ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

### Binding direto de uma fonte

```go
echo.BindBody(c, &payload)        // request body
echo.BindQueryParams(c, &payload) // query parameters
echo.BindPathValues(c, &payload)  // path parameters
echo.BindHeaders(c, &payload)     // headers
```

:::note
Headers **não** são incluídos por `c.Bind()`. Faça binding deles diretamente com `echo.BindHeaders`.
:::

:::caution[Segurança]
Não faça binding diretamente em structs de negócio. Se uma struct vinculada expuser um campo `IsAdmin bool`,
um body de request `{"IsAdmin": true}` o definiria. Use um DTO dedicado e faça o mapeamento
explicitamente:
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

## Binding fluente

Para binding explícito e type-safe de uma única fonte, use os binders fluentes. Eles
encadeiam configuração e execução, coletando erros:

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

Binders disponíveis: `echo.QueryParamsBinder(c)`, `echo.PathValuesBinder(c)`,
`echo.FormFieldBinder(c)`. Termine uma cadeia com `BindError()` (primeiro erro) ou
`BindErrors()` (todos os erros). `FailFast(false)` executa a cadeia inteira; ele vem ativado por padrão.

Cada tipo suportado oferece métodos `Type(...)`, `MustType(...)`, `Types(...)` (slices) e
`MustTypes(...)` — por exemplo, `Int64`, `MustInt64`, `Int64s`. Use
`BindWithDelimiter("id", &dest, ",")` para separar valores unidos por vírgula.

<a id="custom-binder"></a>

## Binder customizado

Registre um binder customizado via `Echo#Binder`:

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
