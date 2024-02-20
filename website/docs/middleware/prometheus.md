---
description: Prometheus metrics middleware
---

# Prometheus

:::note

Echo community contribution

:::

Prometheus middleware generates metrics for HTTP requests. 

There are 2 versions of Prometheus middleware:

- latest (recommended) https://github.com/labstack/echo-contrib/blob/master/echoprometheus/prometheus.go
- old (deprecated) https://github.com/labstack/echo-contrib/blob/master/prometheus/prometheus.go)

Migration guide from old to newer middleware can found [here](https://github.com/labstack/echo-contrib/blob/master/echoprometheus/README.md).

## Usage

- Add needed module `go get -u github.com/labstack/echo-contrib`
- Add Prometheus middleware and metrics serving route
    ```go
    e := echo.New()
    e.Use(echoprometheus.NewMiddleware("myapp")) // adds middleware to gather metrics
    e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
    ```

## Examples

Serve metric from the same server as where metrics is gathered

```go
package main

import (
	"errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(echoprometheus.NewMiddleware("myapp")) // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
	
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```

Serve metrics on a separate port

```go
func main() {
	app := echo.New() // this Echo instance will serve route on port 8080
	app.Use(echoprometheus.NewMiddleware("myapp")) // adds middleware to gather metrics

	go func() {
		metrics := echo.New() // this Echo will run on separate port 8081
		metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
		if err := metrics.Start(":8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	app.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	if err := app.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```

*Sample output (for first example)*

```bash
curl http://localhost:8080/metrics

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

Using custom metrics with Prometheus default registry:

```go
package main

import (
	"errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

func main() {
	e := echo.New()

	customCounter := prometheus.NewCounter( // create new counter metric. This is replacement for `prometheus.Metric` struct
		prometheus.CounterOpts{
			Name: "custom_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
	)
	if err := prometheus.Register(customCounter); err != nil { // register your new counter metric with default metrics registry
		log.Fatal(err)
	}

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			customCounter.Inc() // use our custom metric in middleware. after every request increment the counter
		},
	}))
	e.GET("/metrics", echoprometheus.NewHandler()) // register route for getting gathered metrics

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```

or create your own registry and register custom metrics with that:

```go
package main

import (
	"errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

func main() {
	e := echo.New()

	customRegistry := prometheus.NewRegistry() // create custom registry for your custom metrics
	customCounter := prometheus.NewCounter(    // create new counter metric. This is replacement for `prometheus.Metric` struct
		prometheus.CounterOpts{
			Name: "custom_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
	)
	if err := customRegistry.Register(customCounter); err != nil { // register your new counter metric with metrics registry
		log.Fatal(err)
	}

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			customCounter.Inc() // use our custom metric in middleware. after every request increment the counter
		},
		Registerer: customRegistry, // use our custom registry instead of default Prometheus registry
	}))
	e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{Gatherer: customRegistry})) // register route for getting gathered metrics data from our custom Registry

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```

### Skipping URL(s)

*Usage*

A middleware skipper can be passed to avoid generating metrics to certain URL(s)

```go
package main

import (
	"errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

func main() {
	e := echo.New()

	mwConfig := echoprometheus.MiddlewareConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/testurl")
		}, // does not gather metrics metrics on routes starting with `/testurl`
	}
	e.Use(echoprometheus.NewMiddlewareWithConfig(mwConfig)) // adds middleware to gather metrics

	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```

## Complex Scenarios

Example: modify default `echoprometheus` metrics definitions

```go
package main

import (
	"errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

func main() {
	e := echo.New()

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		// labels of default metrics can be modified or added with `LabelFuncs` function
		LabelFuncs: map[string]echoprometheus.LabelValueFunc{
			"scheme": func(c echo.Context, err error) string { // additional custom label
				return c.Scheme()
			},
			"host": func(c echo.Context, err error) string { // overrides default 'host' label value
				return "y_" + c.Request().Host
			},
		},
		// The `echoprometheus` middleware registers the following metrics by default:
		// - Histogram: request_duration_seconds
		// - Histogram: response_size_bytes
		// - Histogram: request_size_bytes
		// - Counter: requests_total
		// which can be modified with `HistogramOptsFunc` and `CounterOptsFunc` functions
		HistogramOptsFunc: func(opts prometheus.HistogramOpts) prometheus.HistogramOpts {
			if opts.Name == "request_duration_seconds" {
				opts.Buckets = []float64{1000.0, 10_000.0, 100_000.0, 1_000_000.0} // 1KB ,10KB, 100KB, 1MB
			}
			return opts
		},
		CounterOptsFunc: func(opts prometheus.CounterOpts) prometheus.CounterOpts {
			if opts.Name == "requests_total" {
				opts.ConstLabels = prometheus.Labels{"my_const": "123"}
			}
			return opts
		},
	})) // adds middleware to gather metrics

	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
```
