---
title: 响应
description: 发送字符串、HTML、JSON、XML、文件、流、重定向和响应钩子。
sidebar:
  order: 8
---

处理函数通过 `echo.Context` 写入响应。每个辅助方法都会为你设置合适的
`Content-Type` 和状态码。

## 发送字符串

`Context#String(code int, s string)` 会发送带状态码的纯文本响应。

```go
func(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

## 发送 HTML

`Context#HTML(code int, html string)` 会发送带状态码的简单 HTML 响应。
要动态生成 HTML，请参见[模板](/zh-cn/guide/templates/)。

```go
func(c *echo.Context) error {
	return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
}
```

### 发送 HTML blob

`Context#HTMLBlob(code int, b []byte)` 会发送带状态码的 HTML blob。
它适合与输出 `[]byte` 的模板引擎一起使用。

```go
func handler(c *echo.Context) error {
	blob := []byte("<strong>Hello, World!</strong>")
	return c.HTMLBlob(http.StatusOK, blob)
}
```

## 渲染模板

参见[模板](/zh-cn/guide/templates/)。

## 发送 JSON

`Context#JSON(code int, i any)` 会把 Go 值编码为 JSON，并带状态码发送。

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

### 流式发送 JSON

`Context#JSON()` 内部使用 `json.Marshal`，对于大型负载可能效率不高。
这种情况下，请直接流式发送 JSON：

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

`Context#JSONPretty(code int, i any, indent string)` 会发送格式化后的 JSON 响应。
缩进可以是空格或制表符。

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

`Context#JSONBlob(code int, b []byte)` 会直接发送预编码的 JSON blob，
例如来自数据库的数据。

```go
func(c *echo.Context) error {
	encodedJSON := []byte{} // Encoded JSON from an external source.
	return c.JSONBlob(http.StatusOK, encodedJSON)
}
```

## 发送 JSONP

`Context#JSONP(code int, callback string, i any)` 会把 Go 值编码为 JSON，
并作为包装在给定 callback 中的 JSONP 负载发送。

```go
func handler(c *echo.Context) error {
	callback := c.QueryParam("callback")
	return c.JSONP(http.StatusOK, callback, &User{Name: "Jon", Email: "jon@labstack.com"})
}
```

参见 [JSONP cookbook](/zh-cn/cookbook/jsonp/)。

## 发送 XML

`Context#XML(code int, i any)` 会把 Go 值编码为 XML，并带状态码发送。

```go
func(c *echo.Context) error {
	u := &User{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.XML(http.StatusOK, u)
}
```

### 流式发送 XML

`Context#XML` 内部使用 `xml.Marshal`，对于大型负载可能效率不高。
这种情况下，请直接流式发送 XML：

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

`Context#XMLPretty(code int, i any, indent string)` 会发送格式化后的 XML 响应。
缩进可以是空格或制表符。

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
你也可以在请求 URL 查询字符串中追加 `pretty`，让 `Context#XML()` 输出格式化后的 XML。

```sh
curl http://localhost:1323/users/1?pretty
```
:::

### XML blob

`Context#XMLBlob(code int, b []byte)` 会直接发送预编码的 XML blob，
例如来自数据库的数据。

```go
func(c *echo.Context) error {
	encodedXML := []byte{} // Encoded XML from an external source.
	return c.XMLBlob(http.StatusOK, encodedXML)
}
```

## 发送文件

`Context#File(file string)` 会把文件内容作为响应发送。它会设置正确的内容类型，
并自动处理缓存。

```go
func(c *echo.Context) error {
	return c.File("<PATH_TO_YOUR_FILE>")
}
```

## 发送附件

`Context#Attachment(file, name string)` 类似于 `File()`，但会使用
`Content-Disposition: attachment` 和给定名称发送文件。

```go
func(c *echo.Context) error {
	return c.Attachment("<PATH_TO_YOUR_FILE>", "<ATTACHMENT_NAME>")
}
```

## 内联发送

`Context#Inline(file, name string)` 类似于 `File()`，但会使用
`Content-Disposition: inline` 和给定名称发送文件。

```go
func(c *echo.Context) error {
	return c.Inline("<PATH_TO_YOUR_FILE>", "<INLINE_NAME>")
}
```

## 发送 blob

`Context#Blob(code int, contentType string, b []byte)` 会使用给定内容类型和状态码发送任意数据。

```go
func(c *echo.Context) error {
	data := []byte(`0306703,0035866,NO_ACTION,06/19/2006
0086003,"0005866",UPDATED,06/19/2006`)
	return c.Blob(http.StatusOK, "text/csv", data)
}
```

## 发送流

`Context#Stream(code int, contentType string, r io.Reader)` 会使用给定内容类型、
`io.Reader` 和状态码发送任意数据流。

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

## 发送无内容响应

`Context#NoContent(code int)` 会发送带状态码的空响应体。

```go
func(c *echo.Context) error {
	return c.NoContent(http.StatusOK)
}
```

## 重定向请求

`Context#Redirect(code int, url string)` 会使用给定 URL 和状态码重定向请求。

```go
func(c *echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "<URL>")
}
```

## 钩子

### 响应前

`Response#Before(func())` 注册一个函数，在响应写入之前运行。

### 响应后

`Response#After(func())` 注册一个函数，在响应写入之后运行。如果 `Content-Length`
未知，则不会运行任何 after 函数。

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
你可以注册多个 `Before` 和 `After` 函数。
:::

## Echo 安全更新 (v5.4.0 / v4.16.0)

`Context#JSONP` 的回调只能为空、JavaScript 标识符或点分隔的标识符。无效回调返回包含 `ErrInvalidJSONPCallback` 的 HTTP 400；正常响应带有 `X-Content-Type-Options: nosniff`。任何网站都可携带用户 Cookie 读取 JSONP，所以私密数据应使用带 CORS 的 JSON。
