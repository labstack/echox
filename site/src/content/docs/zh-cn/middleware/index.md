---
title: 按任务选择中间件
description: 从要解决的问题出发，查看可运行示例和准确的 Echo API。
sidebar:
  order: 0
---

第一次使用 Echo？请先按[快速入门](../guide/quickstart/)创建服务器。每个中间件页面都说明适用场景、用法，并展示从指定 Echo 修订版本提取的字段和函数签名。

## 你要完成什么任务？

| 任务 | 从这里开始 | 动手试试 |
| --- | --- | --- |
| 记录请求和故障 | [Request Logger](./logger/) 和 [Recover](./recover/) | 运行[日志示例](./logger/) |
| 允许浏览器访问 API | [CORS](./cors/) | 明确列出可信来源 |
| 保护网页表单 | [CSRF](./csrf/) | 选择适合客户端的令牌存储方式 |
| 为 API 验证身份 | [JWT](./jwt/) | 运行 [JWT 示例](../cookbook/jwt/) |
| 限制过量请求 | [Rate Limiter](./rate-limiter/) | 选择标识符和存储方式 |
| 提供静态文件 | [Static](./static/) | 运行[静态文件示例](./static/) |
| 转发到其他服务 | [Proxy](./proxy/) | 运行[反向代理示例](../cookbook/reverse-proxy/) |

## 从代码到可验证的请求

1. 选择上面的任务，复制最简单的用法示例。
2. 在配置部分查看指定 Echo 修订版本的字段和函数签名。改变行为前先阅读默认值与安全说明。
3. 运行完整示例或链接的教程，再用 `curl` 发送请求。

构建第一个 API 时，接着阅读[路由](../guide/routing/)、[数据绑定](../guide/binding/)、[错误处理](../guide/error-handling/)和[测试](../guide/testing/)。处理生产流量时，先加入 [Recover](./recover/) 和 [Request Logger](./logger/)，再选择客户端需要的安全中间件。

JWT 等外部集成由独立模块维护，相关页面会标明所属包。Echo 核心中间件 API 位于 [`github.com/labstack/echo/v5/middleware`](https://pkg.go.dev/github.com/labstack/echo/v5/middleware)。
