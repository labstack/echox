---
title: Descarga de archivos
description: Sirve archivos para descarga, visualización inline o como attachments con nombre.
sidebar:
  order: 6
---

Echo proporciona tres helpers de contexto para devolver archivos: `c.File` sirve un archivo
usando la content disposition por defecto del navegador, `c.Inline` sugiere al navegador
mostrar el archivo en su lugar, y `c.Attachment` solicita una descarga con un filename dado.

## Descargar archivo

### Servidor

```go file=cookbook/file-download/server.go
```

### Cliente

```html file=cookbook/file-download/index.html
```

## Descargar archivo como inline

Usa `c.Inline` para enviar un header `Content-Disposition: inline`, de modo que el navegador
renderice el archivo en su lugar en vez de descargarlo.

### Servidor

```go file=cookbook/file-download/inline/server.go
```

### Cliente

```html file=cookbook/file-download/inline/index.html
```

## Descargar archivo como attachment

Usa `c.Attachment` para enviar un header `Content-Disposition: attachment`, solicitando
al navegador descargar el archivo con el nombre proporcionado.

### Servidor

```go file=cookbook/file-download/attachment/server.go
```

### Cliente

```html file=cookbook/file-download/attachment/index.html
```
