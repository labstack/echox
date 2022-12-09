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

### Serving custom Prometheus Metrics
*Usage*

When creating a new Prometheus middleware, you can pass down a list of extra []*Metric array you can use.

```go
package main

import (
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	prom "github.com/prometheus/client_golang/prometheus"
)

const cKeyMetrics = "custom_metrics"

// See the NewMetrics func for proper descriptions and prometheus names!
// In case you add a metric here later, make sure to include it in the
// MetricsList method or you'll going to have a bad time.
type Metrics struct {
	customCnt *prometheus.Metric
	customDur *prometheus.Metric
}

// Needed by echo-contrib so echo can register and collect these metrics
func (m *Metrics) MetricList() []*prometheus.Metric {
	return []*prometheus.Metric{
		// ADD EVERY METRIC HERE!
		m.customCnt,
		m.customDur,
	}
}

// Creates and populates a new Metrics struct
// This is where all the prometheus metrics, names and labels are specified
func NewMetrics() *Metrics {
	return &Metrics{
		customCnt: &prometheus.Metric{
			Name:        "custom_total",
			Description: "Custom counter events.",
			Type:        "counter_vec",
			Args:        []string{"label_one", "label_two"},
		},
		customDur: &prometheus.Metric{
			Name:        "custom_duration_seconds",
			Description: "Custom duration observations.",
			Type:        "histogram_vec",
			Args:        []string{"label_one", "label_two"},
			Buckets:     prom.DefBuckets, // or your Buckets
		},
	}
}

// This will push your metrics object into every request context for later use
func (m *Metrics) AddCustomMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(cKeyMetrics, m)
		return next(c)
	}
}

func (m *Metrics) IncCustomCnt(labelOne, labelTwo string) {
	labels := prom.Labels{"label_one": labelOne, "label_two": labelTwo}
	m.customCnt.MetricCollector.(*prom.CounterVec).With(labels).Inc()
}

func (m *Metrics) ObserveCustomDur(labelOne, labelTwo string, d time.Duration) {
	labels := prom.Labels{"label_one": labelOne, "label_two": labelTwo}
	m.customCnt.MetricCollector.(*prom.HistogramVec).With(labels).Observe(d.Seconds())
}

func main() {
	e := echo.New()
	m := NewMetrics()

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil, m.MetricList())
	p.Use(e)

	e.Use(m.AddCustomMetricsMiddleware)

	e.GET("/custom", func(c echo.Context) error {
		metrics := c.Get(cKeyMetrics).(*Metrics)
		metrics.IncCustomCnt("any", "value")
		return nil
	})

	e.Logger.Fatal(e.Start(":1323"))
}
```

### Skipping certain URLs
*Usage*

A middleware skipper can be passed to avoid generating metrics to certain URLs:

```go
package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/prometheus"
)

// urlSkipper middleware ignores metrics on some route
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
