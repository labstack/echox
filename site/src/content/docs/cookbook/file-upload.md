---
title: File Upload
description: Handle single and multiple multipart file uploads alongside form fields.
sidebar:
  order: 7
---

Echo reads multipart form data through the request context. Use `c.FormValue` for
text fields, `c.FormFile` for a single file, and `c.MultipartForm` to access
multiple files under the same field name.

## Upload a single file with fields

### Server

```go file=cookbook/file-upload/single/server.go
```

### Client

```html file=cookbook/file-upload/single/public/index.html
```

## Upload multiple files with fields

### Server

```go file=cookbook/file-upload/multiple/server.go
```

### Client

```html file=cookbook/file-upload/multiple/public/index.html
```
