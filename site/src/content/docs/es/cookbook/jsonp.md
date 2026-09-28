---
title: JSONP
description: Sirve responses JSONP para requests cross-domain con Context#JSONP.
sidebar:
  order: 13
---

JSONP es una técnica que permite llamadas de servidor cross-domain desde el navegador. Echo
sirve responses JSONP con `c.JSONP()`, que envuelve el payload JSON en una llamada a la
función callback nombrada en el request.

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

## Actualización de seguridad de Echo (v5.4.0 / v4.16.0)

El callback debe estar vacío, ser un identificador JavaScript o una ruta de identificadores separados por puntos (letras ASCII, dígitos, `_` y `$`). Un valor inválido devuelve HTTP 400 con `ErrInvalidJSONPCallback` y no escribe JSONP. La respuesta incluye `X-Content-Type-Options: nosniff`. **Cualquier sitio puede leer JSONP con las cookies del usuario**: no sirvas datos privados o autenticados; usa JSON con CORS.
