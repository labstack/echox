---
title: 文件上传
description: 处理单个和多个 multipart 文件上传，并同时处理表单字段。
sidebar:
  order: 7
---

Echo 通过请求上下文读取 multipart 表单数据。文本字段使用 `c.FormValue`，单个文件使用
`c.FormFile`，要访问同一字段名下的多个文件则使用 `c.MultipartForm`。

## 上传带字段的单个文件

### 服务器

```go file=cookbook/file-upload/single/server.go
```

### 客户端

```html file=cookbook/file-upload/single/public/index.html
```

## 上传带字段的多个文件

### 服务器

```go file=cookbook/file-upload/multiple/server.go
```

### 客户端

```html file=cookbook/file-upload/multiple/public/index.html
```
