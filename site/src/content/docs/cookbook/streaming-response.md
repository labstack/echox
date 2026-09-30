---
title: Streaming Response
description: Stream data to the client as it is produced using chunked transfer encoding.
sidebar:
  order: 15
---

This recipe streams a JSON response to the client as each record is produced,
using chunked transfer encoding:

- Send data as it is produced.
- Stream a JSON response with chunked transfer encoding.

The handler encodes one record at a time and calls
`http.NewResponseController(...).Flush()` after each to push it to the client
immediately, pausing one second between records.

## Server

```go file=cookbook/streaming-response/server.go
```

## Client

```sh
curl localhost:1323
```

### Output

```js
{"Altitude":-97,"Latitude":37.819929,"Longitude":-122.478255}
{"Altitude":1899,"Latitude":39.096849,"Longitude":-120.032351}
{"Altitude":2619,"Latitude":37.865101,"Longitude":-119.538329}
{"Altitude":42,"Latitude":33.812092,"Longitude":-117.918974}
{"Altitude":15,"Latitude":37.77493,"Longitude":-122.419416}
```
