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

```sh
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
