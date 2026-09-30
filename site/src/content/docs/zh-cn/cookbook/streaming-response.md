---
title: 流式响应
description: 使用 chunked transfer encoding 在数据生成时流式发送给客户端。
sidebar:
  order: 15
---

此示例使用 chunked transfer encoding，在每条记录生成时把 JSON 响应流式发送给客户端：

- 在数据生成时发送数据。
- 使用 chunked transfer encoding 流式发送 JSON 响应。

处理函数一次编码一条记录，并在每条记录后调用 `http.NewResponseController(...).Flush()`，
立即将其推送给客户端，记录之间暂停一秒。

## 服务器

```go file=cookbook/streaming-response/server.go
```

## 客户端

```sh
curl localhost:1323
```

### 输出

```js
{"Altitude":-97,"Latitude":37.819929,"Longitude":-122.478255}
{"Altitude":1899,"Latitude":39.096849,"Longitude":-120.032351}
{"Altitude":2619,"Latitude":37.865101,"Longitude":-119.538329}
{"Altitude":42,"Latitude":33.812092,"Longitude":-117.918974}
{"Altitude":15,"Latitude":37.77493,"Longitude":-122.419416}
```
