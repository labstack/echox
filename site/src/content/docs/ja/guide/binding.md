---
title: バインディング
description: パス、クエリ、header、リクエストボディからリクエストデータを型付き Go struct に解析します。
sidebar:
  order: 5
---

リクエストデータの解析は、Web アプリケーションの重要な要素です。Echo ではこれを
_バインディング_ と呼び、HTTP リクエストの 4 つの部分から読み取れます。

- URL パスパラメーター
- URL クエリパラメーター
- Header
- リクエストボディ

## Struct タグによるバインディング

データソースとキーを指定するタグ付きの struct を定義し、そのポインターを
`c.Bind()` に渡します。ここではクエリパラメーター `id` が `ID` フィールドに
バインドされます。

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

### データソース

| タグ     | ソース |
| -------- | ------ |
| `query`  | クエリパラメーター |
| `param`  | パスパラメーター |
| `header` | Header 値 |
| `form`   | フォームデータ（クエリ + ボディ） |
| `json`   | リクエストボディ（`encoding/json`） |
| `xml`    | リクエストボディ（`encoding/xml`） |

パス、クエリ、header、フォームのフィールドには**明示的なタグ**が必要です。
JSON と XML はタグが省略された場合、標準ライブラリと同じように struct フィールド名へ
フォールバックします。

### スライス

繰り返し指定されたクエリ・パス・フォーム・ヘッダーの値はスライスフィールドにバインドされ、
すべての出現が収集されます：

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

### ボディのコンテンツタイプ

リクエストボディをデコードするときは、`Content-Type` header によってデコーダーが選ばれます。

- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`

### 複数ソースと優先順位

1 つのフィールドで複数のソースを宣言できます。データは次の順序でバインドされ、
各ステップが前の値を上書きします。

1. パスパラメーター
2. クエリパラメーター（GET / DELETE のみ）
3. リクエストボディ

```go
type User struct {
	ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`
}
```

### 1 つのソースから直接バインドする

```go
echo.BindBody(c, &payload)        // request body
echo.BindQueryParams(c, &payload) // query parameters
echo.BindPathValues(c, &payload)  // path parameters
echo.BindHeaders(c, &payload)     // headers
```

:::note
Header は `c.Bind()` に**含まれません**。`echo.BindHeaders` で直接バインドしてください。
:::

:::caution[セキュリティ]
ビジネス用の struct に直接バインドしないでください。バインド対象の struct が
`IsAdmin bool` フィールドを公開している場合、`{"IsAdmin": true}` というリクエストボディで
その値が設定されます。専用の DTO を使い、明示的にマッピングしてください。
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

## フルーエントバインディング

単一ソースから明示的かつ型安全にバインドするには、フルーエント binder を使います。
設定をチェーンし、実行時にエラーを収集します。

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

利用できる binder は `echo.QueryParamsBinder(c)`、`echo.PathValuesBinder(c)`、
`echo.FormFieldBinder(c)` です。チェーンは `BindError()`（最初のエラー）または
`BindErrors()`（すべてのエラー）で終了します。`FailFast(false)` はチェーン全体を実行します。
デフォルトでは早期終了が有効です。

各サポート型には `Type(...)`、`MustType(...)`、`Types(...)`（スライス）、
`MustTypes(...)` メソッドがあります。例：`Int64`、`MustInt64`、`Int64s`。
カンマ区切りの値を分割するには `BindWithDelimiter("id", &dest, ",")` を使います。

<a id="custom-binder"></a>

## カスタム binder

`Echo#Binder` でカスタム binder を登録します。

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
