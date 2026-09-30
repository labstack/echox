---
title: Response en streaming
description: Transmite datos al cliente a medida que se producen usando chunked transfer encoding.
sidebar:
  order: 15
---

Esta receta transmite una response JSON al cliente a medida que se produce cada registro,
usando chunked transfer encoding:

- Enviar datos a medida que se producen.
- Transmitir una response JSON con chunked transfer encoding.

El handler codifica un registro a la vez y llama a `http.NewResponseController(...).Flush()`
después de cada uno para empujarlo al cliente de inmediato, con una pausa de un segundo
entre registros.

## Servidor

```go file=cookbook/streaming-response/server.go
```

## Cliente

```sh
curl localhost:1323
```

### Salida

```js
{"Altitude":-97,"Latitude":37.819929,"Longitude":-122.478255}
{"Altitude":1899,"Latitude":39.096849,"Longitude":-120.032351}
{"Altitude":2619,"Latitude":37.865101,"Longitude":-119.538329}
{"Altitude":42,"Latitude":33.812092,"Longitude":-117.918974}
{"Altitude":15,"Latitude":37.77493,"Longitude":-122.419416}
```
