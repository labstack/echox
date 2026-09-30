---
title: 自定义中间件
description: 编写自定义 Echo 中间件以收集请求统计并设置响应 header。
sidebar:
  order: 12
---

此示例展示如何编写自定义中间件：

- 一个收集请求数量、响应状态和运行时间的中间件。
- 一个为每个响应写入自定义 `Server` header 的中间件。

Echo 中的中间件是签名为 `func(next echo.HandlerFunc) echo.HandlerFunc` 的函数。下面的
`Stats.Process` 方法直接满足该签名，而 `ServerHeader` 是一个普通函数。

## 服务器

```go file=cookbook/middleware/server.go
```

## 响应

### Header

```sh
Content-Length:122
Content-Type:application/json; charset=utf-8
Date:Thu, 14 Apr 2016 20:31:46 GMT
Server:Echo/5.0
```

### Body

```js
{
  "uptime": "2016-04-14T13:28:48.486548936-07:00",
  "requestCount": 5,
  "statuses": {
    "200": 4,
    "404": 1
  }
}
```
