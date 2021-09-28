+++
title = "Prometheus Middleware"
description = "Prometheus metrics middleware for Echo"
[menu.main]
  name = "Prometheus"
  parent = "middleware"
+++

Prometheus middleware generates metrics for HTTP requests.

*Usage*

```go
package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/prometheus"
)
func main() {
    e := echo.New()
    // Enable metrics middleware
    p := prometheus.NewPrometheus("echo", nil)
    p.Use(e)

    e.Logger.Fatal(e.Start(":1323"))
}
```

*Sample Output*

```bash
curl http://localhost:1323/metrics

# HELP echo_request_duration_seconds The HTTP request latencies in seconds.
# TYPE echo_request_duration_seconds summary
echo_request_duration_seconds_sum 0.41086482
echo_request_duration_seconds_count 1
# HELP echo_request_size_bytes The HTTP request sizes in bytes.
# TYPE echo_request_size_bytes summary
echo_request_size_bytes_sum 56
echo_request_size_bytes_count 1
# HELP echo_requests_total How many HTTP requests processed, partitioned by status code and HTTP method.
# TYPE echo_requests_total counter
echo_requests_total{code="200",host="localhost:8080",method="GET",url="/"} 1
# HELP echo_response_size_bytes The HTTP response sizes in bytes.
# TYPE echo_response_size_bytes summary
echo_response_size_bytes_sum 61
echo_response_size_bytes_count 1
...
```

## Custom Configuration

*Usage*

A middleware skipper can be passed to avoid generating metrics to certain URLs:

```go
package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/prometheus"
)

// urlSkipper ignores metrics route on some middleware
func urlSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/testurl") {
		return true
	}
	return false
}

func main() {
    e := echo.New()
    // Enable metrics middleware
    p := prometheus.NewPrometheus("echo", urlSkipper)
    p.Use(e)

    e.Logger.Fatal(e.Start(":1323"))
}
```
## Complex scenarios
Serve /metrics endpoint separately from main server
```go
package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Setup Main Server
	echoMainServer := echo.New()
	echoMainServer.HideBanner = true
	echoMainServer.Use(middleware.Logger())
	echoMainServer.GET("/", hello)

	// Create Prometheus server and Middleware
	echoPrometheus := echo.New()
	echoPrometheus.HideBanner = true
	prom := prometheus.NewPrometheus("echo", nil)

	// Scrape metrics from Main Server
	echoMainServer.Use(prom.HandlerFunc)
	// Setup metrics endpoint at another server
	prom.SetMetricsPath(echoPrometheus)

	go func() { echoPrometheus.Logger.Fatal(echoPrometheus.Start(":9360")) }()

	echoMainServer.Logger.Fatal(echoMainServer.Start(":8080"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```
