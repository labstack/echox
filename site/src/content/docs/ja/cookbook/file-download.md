---
title: ファイルダウンロード
description: ダウンロード、インライン表示、名前付き添付ファイルとしてファイルを配信します。
sidebar:
  order: 6
---

Echo はファイルを返すための 3 つのコンテキストヘルパーを提供します。`c.File` はブラウザーの
デフォルト content disposition でファイルを配信し、`c.Inline` はブラウザーにその場で表示するよう促し、
`c.Attachment` は指定したファイル名でダウンロードを促します。

## ファイルをダウンロードする

### サーバー

```go file=cookbook/file-download/server.go
```

### クライアント

```html file=cookbook/file-download/index.html
```

## ファイルをインラインでダウンロードする

`c.Inline` を使って `Content-Disposition: inline` header を送信し、ブラウザーがファイルを
ダウンロードするのではなくその場でレンダリングするようにします。

### サーバー

```go file=cookbook/file-download/inline/server.go
```

### クライアント

```html file=cookbook/file-download/inline/index.html
```

## ファイルを添付ファイルとしてダウンロードする

`c.Attachment` を使って `Content-Disposition: attachment` header を送信し、
指定した名前でファイルをダウンロードするようブラウザーに促します。

### サーバー

```go file=cookbook/file-download/attachment/server.go
```

### クライアント

```html file=cookbook/file-download/attachment/index.html
```
