---
title: JSONP
description: Sirva responses JSONP para requests entre domínios com Context#JSONP.
sidebar:
  order: 13
---

JSONP é uma técnica que permite chamadas de servidor entre domínios a partir do navegador. Echo
serve responses JSONP com `c.JSONP()`, que envolve o payload JSON em uma chamada para
a função de callback nomeada no request.

## Servidor

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

## Cliente

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

## Atualização de segurança do Echo (v5.4.0 / v4.16.0)

O callback deve ser vazio, um identificador JavaScript ou um caminho de identificadores separados por pontos (letras ASCII, dígitos, `_` e `$`). Um valor inválido retorna HTTP 400 com `ErrInvalidJSONPCallback` e não escreve o corpo JSONP. A resposta válida inclui `X-Content-Type-Options: nosniff`. **Qualquer site pode ler JSONP com os cookies do usuário**: não sirva dados privados ou autenticados; use JSON com CORS.
