---
title: HTTP/2 服务器
description: 使用 TLS 证书启动 Echo，通过 HTTP/2 提供流量。
sidebar:
  order: 9
---

HTTP/2 通过请求复用、header 压缩和 server push 改善延迟。Go 的 HTTP 服务器会在 TLS 上自动协商
HTTP/2，因此使用 Echo 提供 HTTP/2 只需用证书启动服务器。

<a id="1-generate-a-self-signed-x509-tls-certificate"></a>

## 1. 生成自签名 X.509 TLS 证书

运行以下命令生成 `cert.pem` 和 `key.pem`：

```sh
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

:::note
出于演示目的，我们使用自签名证书。在生产环境中，请从
[certificate authority](https://en.wikipedia.org/wiki/Certificate_authority) 获取证书。
:::

## 2. 创建回显请求信息的处理函数

```go
e.GET("/request", func(c *echo.Context) error {
	req := c.Request()
	format := `
		<code>
			Protocol: %s<br>
			Host: %s<br>
			Remote Address: %s<br>
			Method: %s<br>
			Path: %s<br>
		</code>
	`
	return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
})
```

## 3. 启动 TLS 服务器

使用生成的证书和 key 启动服务器：

```go
sc := echo.StartConfig{Address: ":1323"}
if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
	e.Logger.Error("failed to start server", "error", err)
}
```

或者使用带自定义 `tls.Config` 的自定义 `http.Server`：

```go
s := http.Server{
  Addr:    ":8443",
  Handler: e, // set Echo as handler
  TLSConfig: &tls.Config{
    //Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
  },
  //ReadTimeout: 30 * time.Second, // use custom timeouts
}
if err := s.ListenAndServeTLS("cert.pem", "key.pem"); err != http.ErrServerClosed {
  log.Fatal(err)
}
```

## 4. 验证

启动服务器并访问 `https://localhost:1323/request`。你应该会看到类似下面的输出：

```sh
Protocol: HTTP/2.0
Host: localhost:1323
Remote Address: [::1]:60288
Method: GET
Path: /
```

## 源码

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()
	e.GET("/request", func(c *echo.Context) error {
		req := c.Request()
		format := `
			<code>
				Protocol: %s<br>
				Host: %s<br>
				Remote Address: %s<br>
				Method: %s<br>
				Path: %s<br>
			</code>
		`
		return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
	})
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```
