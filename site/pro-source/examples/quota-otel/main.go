// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

// Command quota-otel runs a bounded Echo request and a metered MCP tool call,
// emits spans and metrics to stdout, and shuts down all caller-owned resources.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotamcp"
	"echo.labstack.com/pro/quota/quotamcp/mcpotel"
	"echo.labstack.com/pro/quota/quotaotel"
)

func run(ctx context.Context, output io.Writer) (result error) {
	spanExporter, err := stdouttrace.New(stdouttrace.WithWriter(output))
	if err != nil {
		return err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(spanExporter))
	defer func() {
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		result = errors.Join(result, tp.Shutdown(shutdown))
	}()
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithWriter(output))
	if err != nil {
		return err
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(metricExporter),
	))
	defer func() {
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		result = errors.Join(result, mp.Shutdown(shutdown))
	}()
	names := quotaotel.NameAllowlist{Plans: []string{
		"team",
	}, Limits: []string{"daily", "tokens"}, Meters: []string{
		"requests",
		"tokens",
	}}
	failures, err := quotaotel.NewFailureCounter(mp)
	if err != nil {
		return err
	}
	tracing, err := quotaotel.Tracing(quotaotel.TraceConfig{
		Names:    names,
		Failures: failures,
	})
	if err != nil {
		return err
	}
	metrics, err := quotaotel.New(mp)
	if err != nil {
		return err
	}
	lifecycle, err := quotaotel.NewLifecycle(mp, quotaotel.LifecycleConfig{Names: names})
	if err != nil {
		return err
	}
	//nolint:contextcheck // Constructor has no context API; no outbox is configured.
	limiter, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("daily", 10),
		quota.PerDay("tokens", 1000).Meter("tokens"),
	}}}, Hooks: quota.JoinHooks(
		tracing,
		metrics.Hooks(),
		lifecycle.Hooks(),
	)})
	if err != nil {
		return err
	}
	defer func() {
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		result = errors.Join(result, limiter.Close(shutdown))
	}()
	tracer := tp.Tracer("quota-otel-example")
	e := echo.New()
	// Tracing runs first, so quota receives the span in the request context.
	//nolint:contextcheck // Uses each Echo request context.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			requestCtx, span := tracer.Start(
				c.Request().Context(),
				"HTTP "+c.Request().Method,
			)
			defer span.End()
			c.SetRequest(c.Request().WithContext(requestCtx))
			return next(c)
		}
	})
	//nolint:contextcheck // Registers middleware that uses each request context.
	mw, err := (quota.MiddlewareConfig{
		Limiter: limiter,
		Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "demo"}, nil
		},
	}).ToMiddleware()
	if err != nil {
		return err
	}
	e.GET("/", func(c *echo.Context) error {
		return c.String(
			http.StatusOK,
			"quota-otel demo\n",
		)
	}, mw)
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	server := &http.Server{
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	defer func() {
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		result = errors.Join(result, server.Shutdown(shutdown))
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			result = errors.Join(result, err)
		}
	}()
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://"+listener.Addr().String()+"/",
		nil,
	)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if _, err := io.Copy(io.Discard, response.Body); err != nil {
		return errors.Join(err, response.Body.Close())
	}
	if err := response.Body.Close(); err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP demonstration returned %d", response.StatusCode)
	}
	return toolCall(ctx, limiter, mp, tracer, failures)
}

func toolCall(
	ctx context.Context,
	limiter *quota.Limiter,
	mp *sdkmetric.MeterProvider,
	tracer trace.Tracer,
	failures quotaotel.FailureSink,
) (
	result error,
) {
	names := mcpotel.NameAllowlist{
		Plans:  []string{"team"},
		Tools:  []string{"search"},
		Meters: []string{"tokens"},
	}
	tracing, err := mcpotel.MCP(mcpotel.MCPTraceConfig{Names: names, Failures: failures})
	if err != nil {
		return err
	}
	metrics, err := mcpotel.NewMCP(mp, mcpotel.MCPMetricsConfig{
		Names:    names,
		Failures: failures,
	})
	if err != nil {
		return err
	}
	a := metrics.Adapters()
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "quota-otel-demo",
		Version: "v1",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Local demonstration",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (
		*mcp.CallToolResult,
		any,
		error,
	) {
		quotamcp.Charge(ctx, quota.Cost{"tokens": 12})
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{
			Text: "demo result",
		}}}, nil, nil
	})
	mw, err := (quotamcp.MiddlewareConfig{Limiter: limiter, Identity: func(
		context.Context,
		quotamcp.Caller,
	) (
		quota.Identity,
		error,
	) {
		return quota.Identity{Tenant: "demo"}, nil
	}, Policies: map[string]quotamcp.Policy{"search": {
		Cost:   quota.Cost{"tokens": 100},
		Settle: true,
	}}, OnCall: func(ctx context.Context, i quotamcp.CallInfo) {
		tracing.OnCall(ctx, i)
		a.OnCall(
			ctx,
			i,
		)
	}, OnMeteringError: func(ctx context.Context, f quotamcp.MeteringFailure) {
		tracing.OnMeteringError(ctx, f)
		a.OnMeteringError(ctx, f)
	}}).ToMiddleware()
	if err != nil {
		return err
	}
	server.AddReceivingMiddleware(mw)
	// In-memory transports do not propagate a client's span automatically. The
	// server session explicitly borrows this live application span for the demo.
	callCtx, span := tracer.Start(ctx, "MCP search")
	defer span.End()
	st, ct := mcp.NewInMemoryTransports()
	ss, err := server.Connect(callCtx, st, nil)
	if err != nil {
		return err
	}
	defer func() { result = errors.Join(result, ss.Close()) }()
	cs, err := mcp.NewClient(&mcp.Implementation{
		Name:    "demo-client",
		Version: "v1",
	}, nil).Connect(callCtx, ct, nil)
	if err != nil {
		return err
	}
	defer func() { result = errors.Join(result, cs.Close()) }()
	response, err := cs.CallTool(callCtx, &mcp.CallToolParams{Name: "search"})
	if err != nil {
		return err
	}
	if response.IsError {
		return errors.New("MCP demonstration was denied")
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
