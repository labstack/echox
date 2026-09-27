---
title: 自定义
description: 自定义 Echo 的 logger、validator、binder、renderer、serializer 和错误处理。
sidebar:
  order: 12
---

Echo 在 `Echo` 实例上暴露了一组字段，让你可以用自己的实现替换内置行为。

## 日志

`Echo#Logger` 写入结构化日志。默认处理函数会把 JSON 输出到 `os.Stdout`。

### 自定义 logger

logger 是一个 `*slog.Logger`，因此你可以注册任何 `slog` 处理函数：

```go
e.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

## Validator

`Echo#Validator` 注册用于请求负载验证的 validator。

[了解更多](/zh-cn/guide/request/#validate-data)

## 自定义 binder

`Echo#Binder` 注册用于绑定请求负载的自定义 binder。

[了解更多](/zh-cn/guide/binding/#custom-binder)

## 自定义 JSON serializer

`Echo#JSONSerializer` 注册自定义 JSON serializer。参见
[json.go](https://github.com/labstack/echo/blob/master/json.go) 中的
`DefaultJSONSerializer`。

## Renderer

`Echo#Renderer` 注册用于模板渲染的 renderer。

[了解更多](/zh-cn/guide/templates/)

<a id="http-error-handler"></a>

## HTTP 错误处理函数

`Echo#HTTPErrorHandler` 注册自定义 HTTP 错误处理函数。

[了解更多](/zh-cn/guide/error-handling/)

## 路由回调

`Echo#OnAddRoute` 注册一个回调，每当新路由被添加到路由器时都会调用。

## IP 提取器

`Echo#IPExtractor` 控制如何确定真实客户端 IP 地址。要可靠且安全地获取它，
你的应用必须了解整个基础设施。

[了解更多](/zh-cn/guide/ip-address/)
