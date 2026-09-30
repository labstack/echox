---
title: Custom Middleware
description: Write custom Echo middleware to collect request statistics and set response headers.
sidebar:
  order: 12
---

This recipe shows how to write custom middleware:

- A middleware that collects the request count, response statuses, and uptime.
- A middleware that writes a custom `Server` header to every response.

A middleware in Echo is a function with the signature
`func(next echo.HandlerFunc) echo.HandlerFunc`. The `Stats.Process` method below
satisfies that signature directly, while `ServerHeader` is a plain function.

## Server

```go file=cookbook/middleware/server.go
```

## Response

### Headers

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
