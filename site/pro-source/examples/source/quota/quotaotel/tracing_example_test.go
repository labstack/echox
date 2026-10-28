// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotaotel_test

import (
	"context"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaotel"
)

// docs:start ExampleTracing
func ExampleTracing() {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx, span := provider.Tracer("example").Start(context.Background(), "request")
	hooks, err := quotaotel.Tracing(quotaotel.TraceConfig{
		Names: quotaotel.NameAllowlist{
			Plans:  []string{"standard"},
			Limits: []string{"daily"},
		},
	})
	if err != nil {
		panic(err)
	}
	limiter, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "standard", Limits: []quota.Limit{
		quota.PerDay("daily", 10),
	}}}, Hooks: hooks})
	if err != nil {
		panic(err)
	}
	d, err := limiter.Allow(ctx, quota.Request{
		Identity: quota.Identity{Tenant: "example"},
	})
	if err != nil {
		panic(err)
	}
	span.End()
	fmt.Println(d.Allowed, exporter.GetSpans()[0].Events[0].Name)
	if err := limiter.Close(ctx); err != nil {
		panic(err)
	}
	if err := provider.Shutdown(ctx); err != nil {
		panic(err)
	}
	// Output: true quota.decision
}

// docs:end ExampleTracing
