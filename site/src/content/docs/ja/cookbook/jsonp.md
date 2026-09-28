---
title: JSONP
description: Context#JSONP でクロスドメインリクエスト向け JSONP レスポンスを配信します。
sidebar:
  order: 13
---

JSONP は、ブラウザーからクロスドメインのサーバー呼び出しを可能にする技術です。
Echo は `c.JSONP()` で JSONP レスポンスを配信し、リクエストで指定された callback 関数呼び出しで
JSON ペイロードをラップします。

## サーバー

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

## クライアント

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

## Echo のセキュリティ更新 (v5.4.0 / v4.16.0)

コールバックは空文字、JavaScript 識別子、または識別子をドットでつないだ名前に限られます（ASCII 英字、数字、`_`、`$`）。不正な値では `ErrInvalidJSONPCallback` を含む HTTP 400 を返し、JSONP 本文は書き込みません。正常な応答には `X-Content-Type-Options: nosniff` が付きます。**どのサイトもユーザーの Cookie とともに JSONP を読み取れる**ため、認証が必要なデータには使わず、CORS を設定した JSON を使用してください。
