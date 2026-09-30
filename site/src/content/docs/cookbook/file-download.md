---
title: File Download
description: Serve files for download, inline display, or as named attachments.
sidebar:
  order: 6
---

Echo provides three context helpers for returning files: `c.File` serves a file
using the browser's default content disposition, `c.Inline` hints the browser to
display the file in place, and `c.Attachment` prompts a download with a given
filename.

## Download file

### Server

```go file=cookbook/file-download/server.go
```

### Client

```html file=cookbook/file-download/index.html
```

## Download file as inline

Use `c.Inline` to send a `Content-Disposition: inline` header so the browser
renders the file in place rather than downloading it.

### Server

```go file=cookbook/file-download/inline/server.go
```

### Client

```html file=cookbook/file-download/inline/index.html
```

## Download file as attachment

Use `c.Attachment` to send a `Content-Disposition: attachment` header, prompting
the browser to download the file under the supplied name.

### Server

```go file=cookbook/file-download/attachment/server.go
```

### Client

```html file=cookbook/file-download/attachment/index.html
```
