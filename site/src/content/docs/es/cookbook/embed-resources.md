---
title: Recursos embebidos
description: Sirve assets estáticos empaquetados en el binario con el paquete embed de Go.
sidebar:
  order: 5
---

El paquete `embed` de Go (Go 1.16+) te permite compilar assets estáticos directamente en el
binario, por lo que un único ejecutable puede incluir su frontend. Esta receta sirve el
filesystem embebido mediante Echo, con un modo live opcional que lee desde disco durante
el desarrollo.

## Servidor

```go file=cookbook/embed/server.go
```

:::tip
Ejecuta el binario con el argumento `live` (`go run server.go live`) para servir assets
desde el directorio `app` en disco en lugar de la copia embebida, lo que resulta útil
durante el desarrollo.
:::
