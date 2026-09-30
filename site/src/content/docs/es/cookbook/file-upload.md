---
title: Subida de archivos
description: Maneja uploads multipart de uno o varios archivos junto a campos de formulario.
sidebar:
  order: 7
---

Echo lee datos de formularios multipart mediante el contexto del request. Usa `c.FormValue`
para campos de texto, `c.FormFile` para un solo archivo y `c.MultipartForm` para acceder
a múltiples archivos bajo el mismo nombre de campo.

## Subir un solo archivo con campos

### Servidor

```go file=cookbook/file-upload/single/server.go
```

### Cliente

```html file=cookbook/file-upload/single/public/index.html
```

## Subir múltiples archivos con campos

### Servidor

```go file=cookbook/file-upload/multiple/server.go
```

### Cliente

```html file=cookbook/file-upload/multiple/public/index.html
```
