---
title: Recursos incorporados
description: Sirva assets estáticos empacotados no binário com o pacote embed do Go.
sidebar:
  order: 5
---

O pacote `embed` do Go (Go 1.16+) permite compilar assets estáticos diretamente no
binário, para que um único executável possa incluir seu frontend. Esta receita serve o
filesystem incorporado por meio do Echo, com um modo live opcional que lê do disco
durante o desenvolvimento.

## Servidor

```go file=cookbook/embed/server.go
```

:::tip
Execute o binário com o argumento `live` (`go run server.go live`) para servir assets
a partir do diretório `app` no disco em vez da cópia incorporada, o que é útil
durante o desenvolvimento.
:::
