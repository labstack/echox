package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	e := echo.New() // this Echo instance will serve route on port 8080

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		// labels of default metrics can be modified or added with `LabelFuncs` function
		LabelFuncs: map[string]echoprometheus.LabelValueFunc{
			"scheme": func(c *echo.Context, err error) string { // additional custom label
				return c.Scheme()
			},
			"host": func(c *echo.Context, err error) string { // overrides default 'host' label value
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

	e.GET("/hello", func(c *echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
