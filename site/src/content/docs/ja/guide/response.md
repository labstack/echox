---
title: レスポンス
description: 文字列、HTML、JSON、XML、ファイル、ストリーム、リダイレクト、レスポンスフックを送信します。
sidebar:
  order: 8
---

ハンドラは `echo.Context` を通じてレスポンスを書き込みます。各ヘルパーは適切な
`Content-Type` とステータスコードを設定します。

## 文字列を送信する

`Context#String(code int, s string)` は、ステータスコード付きのプレーンテキストレスポンスを送信します。

```go
func(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

## HTML を送信する

`Context#HTML(code int, html string)` は、ステータスコード付きのシンプルな HTML レスポンスを送信します。
HTML を動的に生成するには、[テンプレート](/ja/guide/templates/)を参照してください。

```go
func(c *echo.Context) error {
	return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
}
```

### HTML blob を送信する

`Context#HTMLBlob(code int, b []byte)` は、ステータスコード付きの HTML blob を送信します。
`[]byte` を出力するテンプレートエンジンと組み合わせると便利です。

```go
func handler(c *echo.Context) error {
	blob := []byte("<strong>Hello, World!</strong>")
	return c.HTMLBlob(http.StatusOK, blob)
}
```

## テンプレートをレンダリングする

[テンプレート](/ja/guide/templates/)を参照してください。

## JSON を送信する

`Context#JSON(code int, i any)` は Go の値を JSON としてエンコードし、ステータスコード付きで送信します。

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

### JSON をストリーミングする

`Context#JSON()` は内部で `json.Marshal` を使うため、大きなペイロードでは非効率な場合があります。
その場合は JSON を直接ストリーミングします。

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

`Context#JSONPretty(code int, i any, indent string)` は整形済み JSON レスポンスを送信します。
インデントにはスペースまたはタブを使えます。

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

### JSON blob

`Context#JSONBlob(code int, b []byte)` は、たとえばデータベース由来の
事前エンコード済み JSON blob を直接送信します。

```go
func(c *echo.Context) error {
	encodedJSON := []byte{} // Encoded JSON from an external source.
	return c.JSONBlob(http.StatusOK, encodedJSON)
}
```

## JSONP を送信する

`Context#JSONP(code int, callback string, i any)` は Go の値を JSON としてエンコードし、
指定した callback でラップした JSONP ペイロードとして送信します。

```go
func handler(c *echo.Context) error {
	callback := c.QueryParam("callback")
	return c.JSONP(http.StatusOK, callback, &User{Name: "Jon", Email: "jon@labstack.com"})
}
```

クロスオリジンのリクエストには [JSON と CORS](/ja/cookbook/cors/)を使用してください。

## XML を送信する

`Context#XML(code int, i any)` は Go の値を XML としてエンコードし、ステータスコード付きで送信します。

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.XML(http.StatusOK, u)
}
```

### XML をストリーミングする

`Context#XML` は内部で `xml.Marshal` を使うため、大きなペイロードでは非効率な場合があります。
その場合は XML を直接ストリーミングします。

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

`Context#XMLPretty(code int, i any, indent string)` は整形済み XML レスポンスを送信します。
インデントにはスペースまたはタブを使えます。

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
リクエスト URL のクエリ文字列に `pretty` を追加すると、`Context#XML()` でも
整形済み XML を出力できます。

```sh
curl http://localhost:1323/users/1?pretty
```
:::

### XML blob

`Context#XMLBlob(code int, b []byte)` は、たとえばデータベース由来の
事前エンコード済み XML blob を直接送信します。

```go
func(c *echo.Context) error {
	encodedXML := []byte{} // Encoded XML from an external source.
	return c.XMLBlob(http.StatusOK, encodedXML)
}
```

## ファイルを送信する

`Context#File(file string)` はファイルの内容をレスポンスとして送信します。
正しいコンテンツタイプを設定し、キャッシュも自動的に処理します。

```go
func(c *echo.Context) error {
	return c.File("<PATH_TO_YOUR_FILE>")
}
```

## 添付ファイルを送信する

`Context#Attachment(file, name string)` は `File()` と似ていますが、
`Content-Disposition: attachment` と指定した名前でファイルを送信します。

```go
func(c *echo.Context) error {
	return c.Attachment("<PATH_TO_YOUR_FILE>", "<ATTACHMENT_NAME>")
}
```

## インライン送信する

`Context#Inline(file, name string)` は `File()` と似ていますが、
`Content-Disposition: inline` と指定した名前でファイルを送信します。

```go
func(c *echo.Context) error {
	return c.Inline("<PATH_TO_YOUR_FILE>", "<INLINE_NAME>")
}
```

## blob を送信する

`Context#Blob(code int, contentType string, b []byte)` は、指定したコンテンツタイプと
ステータスコードで任意のデータを送信します。

```go
func(c *echo.Context) error {
	data := []byte(`0306703,0035866,NO_ACTION,06/19/2006
0086003,"0005866",UPDATED,06/19/2006`)
	return c.Blob(http.StatusOK, "text/csv", data)
}
```

## ストリームを送信する

`Context#Stream(code int, contentType string, r io.Reader)` は、指定したコンテンツタイプ、
`io.Reader`、ステータスコードで任意のデータストリームを送信します。

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

## ボディなしで送信する

`Context#NoContent(code int)` は、ステータスコード付きの空ボディを送信します。

```go
func(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}
```

## リクエストをリダイレクトする

`Context#Redirect(code int, url string)` は、指定した URL とステータスコードでリクエストをリダイレクトします。

```go
func(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "<URL>")
}
```

## フック

### レスポンス前

`Response#Before(func())` は、レスポンスが書き込まれる直前に実行される関数を登録します。

### レスポンス後

`Response#After(func())` は、レスポンスが書き込まれた直後に実行される関数を登録します。
`Content-Length` が不明な場合、after 関数は実行されません。

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
複数の `Before` 関数と `After` 関数を登録できます。
:::

## Echo のセキュリティ更新 (v5.4.0 / v4.16.0)

`Context#JSONP` のコールバックは空文字、JavaScript 識別子、またはドット区切りの識別子だけです。不正な値は `ErrInvalidJSONPCallback` を含む HTTP 400 を返し、正常な応答には `X-Content-Type-Options: nosniff` が付きます。JSONP は他サイトからユーザーの Cookie とともに読めるため、非公開データには使わず CORS を設定した JSON を使ってください。
