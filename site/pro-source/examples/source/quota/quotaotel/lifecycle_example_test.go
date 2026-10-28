// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotaotel_test

import (
	"context"
	"fmt"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaotel"
)

// docs:start collected
func collected(reader *sdkmetric.ManualReader) {
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		panic(err)
	}
	for _, scope := range rm.ScopeMetrics {
		for _, m := range scope.Metrics {
			fmt.Println(m.Name, m.Unit)
		}
	}
}

// docs:end collected

// docs:start ExampleNewLifecycle
func ExampleNewLifecycle() {
	ctx := context.Background()
	r := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(r))
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	m, err := quotaotel.NewLifecycle(mp, quotaotel.LifecycleConfig{Names: quotaotel.NameAllowlist{Plans: []string{"team"}}})
	if err != nil {
		panic(err)
	}
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{quota.PerDay("calls", 10)}}}, Hooks: m.Hooks()})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := l.Close(ctx); err != nil {
			panic(err)
		}
	}()
	if _, err := l.Allow(ctx, quota.Request{Identity: quota.Identity{Tenant: "acme"}}); err != nil {
		panic(err)
	}
	collected(r)
	// Output: quota.admissions {decision}
}

// docs:end ExampleNewLifecycle

// docs:start ExampleNewFailureCounter
func ExampleNewFailureCounter() {
	ctx := context.Background()
	r := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(r))
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	f, err := quotaotel.NewFailureCounter(mp)
	if err != nil {
		panic(err)
	}
	// Share f with Tracing and MCP; unknown components normalize to unknown.
	f.Fail("application-specific")
	collected(r)
	// Output: quota.telemetry.failures {failure}
}

// docs:end ExampleNewFailureCounter
