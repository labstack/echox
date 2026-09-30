---
title: Upload de arquivos
description: Trate uploads multipart de um ou vários arquivos junto de campos de formulário.
sidebar:
  order: 7
---

Echo lê dados de formulário multipart por meio do contexto do request. Use `c.FormValue` para
campos de texto, `c.FormFile` para um único arquivo e `c.MultipartForm` para acessar
vários arquivos sob o mesmo nome de campo.

## Fazer upload de um único arquivo com campos

### Servidor

```go file=cookbook/file-upload/single/server.go
```

### Cliente

```html file=cookbook/file-upload/single/public/index.html
```

## Fazer upload de vários arquivos com campos

### Servidor

```go file=cookbook/file-upload/multiple/server.go
```

### Cliente

```html file=cookbook/file-upload/multiple/public/index.html
```
