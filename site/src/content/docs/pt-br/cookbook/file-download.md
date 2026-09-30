---
title: Download de arquivos
description: Sirva arquivos para download, exibição inline ou como attachments nomeados.
sidebar:
  order: 6
---

Echo fornece três helpers de contexto para retornar arquivos: `c.File` serve um arquivo
usando a content disposition padrão do navegador, `c.Inline` sugere que o navegador
exiba o arquivo no local, e `c.Attachment` solicita um download com um
nome de arquivo informado.

## Baixar arquivo

### Servidor

```go file=cookbook/file-download/server.go
```

### Cliente

```html file=cookbook/file-download/index.html
```

## Baixar arquivo como inline

Use `c.Inline` para enviar um header `Content-Disposition: inline`, de modo que o navegador
renderize o arquivo no local em vez de baixá-lo.

### Servidor

```go file=cookbook/file-download/inline/server.go
```

### Cliente

```html file=cookbook/file-download/inline/index.html
```

## Baixar arquivo como attachment

Use `c.Attachment` para enviar um header `Content-Disposition: attachment`, solicitando
que o navegador baixe o arquivo com o nome fornecido.

### Servidor

```go file=cookbook/file-download/attachment/server.go
```

### Cliente

```html file=cookbook/file-download/attachment/index.html
```
