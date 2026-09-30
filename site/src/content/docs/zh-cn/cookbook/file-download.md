---
title: 文件下载
description: 提供文件下载、内联显示，或作为命名附件提供。
sidebar:
  order: 6
---

Echo 提供三个返回文件的上下文辅助方法：`c.File` 使用浏览器默认的 content disposition 提供文件，
`c.Inline` 提示浏览器就地显示文件，`c.Attachment` 则使用给定文件名提示下载。

## 下载文件

### 服务器

```go file=cookbook/file-download/server.go
```

### 客户端

```html file=cookbook/file-download/index.html
```

## 以内联方式下载文件

使用 `c.Inline` 发送 `Content-Disposition: inline` header，使浏览器就地渲染文件，
而不是下载它。

### 服务器

```go file=cookbook/file-download/inline/server.go
```

### 客户端

```html file=cookbook/file-download/inline/index.html
```

## 以附件方式下载文件

使用 `c.Attachment` 发送 `Content-Disposition: attachment` header，提示浏览器使用提供的名称下载文件。

### 服务器

```go file=cookbook/file-download/attachment/server.go
```

### 客户端

```html file=cookbook/file-download/attachment/index.html
```
