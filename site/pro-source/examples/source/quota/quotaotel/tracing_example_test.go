// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotaotel_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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

// docs:start outbox-metrics
// ConfigureOutboxMetrics records delivery metrics after the exporter accepts an event.
func ConfigureOutboxMetrics(cfg quota.Config, metrics *quotaotel.Metrics, tracing, lifecycle quota.Hooks, deliver func(context.Context, quota.Event) error) quota.Config {
	cfg.Hooks = quota.JoinHooks(tracing, metrics.OutboxHooks(), lifecycle)
	cfg.Outbox = &quota.Outbox{Deliver: func(ctx context.Context, ev quota.Event) error {
		if err := deliver(ctx, ev); err != nil {
			return err
		}
		metrics.Count(ctx, ev)
		return nil
	}}
	return cfg
}

// docs:end outbox-metrics

func TestConfigureOutboxMetrics(t *testing.T) {
	deliveryErr := errors.New("delivery failed")
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{"success", nil, "quota.consumption {unit} {meter=\"tokens\"} 5\n"},
		{"failure", deliveryErr, "\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			t.Cleanup(func() {
				if err := provider.Shutdown(context.WithoutCancel(t.Context())); err != nil {
					t.Error(err)
				}
			})
			metrics, err := quotaotel.New(provider)
			if err != nil {
				t.Fatal(err)
			}
			event := quota.Event{Consumption: &quota.ConsumptionEvent{Meter: "tokens", Amount: 5}}
			deliverCalls := 0
			cfg := ConfigureOutboxMetrics(quota.Config{}, metrics, quota.Hooks{}, quota.Hooks{},
				func(ctx context.Context, ev quota.Event) error {
					deliverCalls++
					if ev.Consumption != event.Consumption {
						t.Fatal("delivery received a different event")
					}
					var observed metricdata.ResourceMetrics
					if err := reader.Collect(ctx, &observed); err != nil {
						t.Fatal(err)
					}
					if len(observed.ScopeMetrics) != 0 {
						t.Fatalf("counted before delivery succeeded: %+v", observed)
					}
					return tc.err
				})
			if err := cfg.Outbox.Deliver(t.Context(), event); !errors.Is(err, tc.err) {
				t.Fatalf("Deliver error = %v, want %v", err, tc.err)
			}
			if got := render(t, reader); got != tc.want || deliverCalls != 1 {
				t.Fatalf("metrics=%q deliveries=%d; want %q and one delivery", got, deliverCalls, tc.want)
			}
		})
	}
}
