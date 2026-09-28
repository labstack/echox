---
title: JSONP
description: 使用 Context#JSONP 为跨域请求提供 JSONP 响应。
sidebar:
  order: 13
---

JSONP 是一种允许浏览器发起跨域服务器调用的技术。Echo 使用 `c.JSONP()` 提供 JSONP 响应，
它会把 JSON 负载包装到请求中命名的 callback 函数调用里。

## 服务器

```go
package main

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Static("/", "public")

	// JSONP
	e.GET("/jsonp", func(c *echo.Context) error {
		callback := c.QueryParam("callback")
		var content struct {
			Response  string    `json:"response"`
			Timestamp time.Time `json:"timestamp"`
			Random    int       `json:"random"`
		}
		content.Response = "Sent via JSONP"
		content.Timestamp = time.Now().UTC()
		content.Random = rand.Intn(1000)
		return c.JSONP(http.StatusOK, callback, &content)
	})

	// Start server
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

## 客户端

```html
<!DOCTYPE html>
<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css">
    <title>JSONP</title>
    <script type="text/javascript" src="https://code.jquery.com/jquery-3.7.1.min.js"></script>
    <script type="text/javascript">
        var host_prefix = 'http://localhost:1323';
        $(function() {
            // JSONP version - add 'callback=?' to the URL - fetch the JSONP response to the request
            $("#jsonp-button").on("click", function(e) {
                e.preventDefault();
                // The only difference on the client end is the addition of 'callback=?' to the URL
                var url = host_prefix + '/jsonp?callback=?';
                $.getJSON(url, function(jsonp) {
                    console.log(jsonp);
                    $("#jsonp-response").html(JSON.stringify(jsonp, null, 2));
                });
            });
        });
    </script>

</head>

<body>
    <div class="container" style="margin-top: 50px;">
        <input type="button" class="btn btn-primary btn-lg" id="jsonp-button" value="Get JSONP response">
        <p>
            <pre id="jsonp-response"></pre>
        </p>
    </div>
</body>

</html>
```

## Echo 安全更新 (v5.4.0 / v4.16.0)

回调必须为空、JavaScript 标识符，或由点分隔的标识符路径（ASCII 字母、数字、`_`、`$`）。无效值返回包含 `ErrInvalidJSONPCallback` 的 HTTP 400，且不写入 JSONP 正文。正常响应带有 `X-Content-Type-Options: nosniff`。**任何网站都能携带用户 Cookie 读取 JSONP**；不要用它返回需要认证或私密的数据，应使用带 CORS 的 JSON。
