---
title: Streaming Response
description: Envie dados ao cliente conforme são produzidos usando chunked transfer encoding.
sidebar:
  order: 15
---

Esta receita envia uma response JSON ao cliente conforme cada registro é produzido,
usando chunked transfer encoding:

- Envie dados conforme são produzidos.
- Faça streaming de uma response JSON com chunked transfer encoding.

O handler codifica um registro por vez e chama
`http.NewResponseController(...).Flush()` depois de cada um para enviá-lo ao cliente
imediatamente, pausando um segundo entre registros.

## Servidor

```go file=cookbook/streaming-response/server.go
```

## Cliente

```sh
curl localhost:1323
```

### Saída

```js
{"Altitude":-97,"Latitude":37.819929,"Longitude":-122.478255}
{"Altitude":1899,"Latitude":39.096849,"Longitude":-120.032351}
{"Altitude":2619,"Latitude":37.865101,"Longitude":-119.538329}
{"Altitude":42,"Latitude":33.812092,"Longitude":-117.918974}
{"Altitude":15,"Latitude":37.77493,"Longitude":-122.419416}
```
