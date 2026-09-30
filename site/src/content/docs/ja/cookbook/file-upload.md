---
title: ファイルアップロード
description: フォームフィールドと一緒に単一または複数の multipart ファイルアップロードを処理します。
sidebar:
  order: 7
---

Echo はリクエストコンテキストを通じて multipart フォームデータを読み取ります。テキストフィールドには
`c.FormValue`、単一ファイルには `c.FormFile`、同じフィールド名の複数ファイルへアクセスするには
`c.MultipartForm` を使います。

## フィールド付きで単一ファイルをアップロードする

### サーバー

```go file=cookbook/file-upload/single/server.go
```

### クライアント

```html file=cookbook/file-upload/single/public/index.html
```

## フィールド付きで複数ファイルをアップロードする

### サーバー

```go file=cookbook/file-upload/multiple/server.go
```

### クライアント

```html file=cookbook/file-upload/multiple/public/index.html
```
