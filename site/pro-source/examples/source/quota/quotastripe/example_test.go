// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotastripe_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotastripe"
)

// docs:start ExampleExporter_Deliver
func ExampleExporter_Deliver() {
	ctx := context.Background()
	store, err := quotastripe.NewMemoryStore(quotastripe.StoreConfig{})
	if err != nil {
		fmt.Println(err)
		return
	}
	exporter, err := quotastripe.New(quotastripe.Config{Namespace: "example-state-1", Store: store, Client: &http.Client{}, APIKey: func(context.Context) (string, error) { return "owner-injected-key", nil }, Customer: func(context.Context, quota.ConsumptionEvent) (string, error) { return "cus_example", nil }, Meters: map[string]quotastripe.Meter{"tokens": {ID: "mtr_example", EventName: "tokens", CustomerKey: "stripe_customer_id", ValueKey: "value"}}})
	if err != nil {
		fmt.Println(err)
		return
	}
	event := quota.ConsumptionEvent{ID: "unique-source-event", Meter: "tokens", Amount: 12, At: time.Now().UTC().Truncate(time.Millisecond)}
	if err = exporter.Deliver(ctx, quota.Event{ID: event.ID, Consumption: &event}); err != nil {
		fmt.Println(err)
		return
	}
	if err = exporter.Deliver(ctx, quota.Event{ID: event.ID, Consumption: &event}); err != nil {
		fmt.Println(err)
		return
	}
	status, err := exporter.Status(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(status.States[quotastripe.Open].Amount)
	page, err := exporter.Inspect(ctx, quotastripe.Filter{Limit: 10}, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(page.Decisions))
	// Output:
	// 12
	// 1
}

// docs:end ExampleExporter_Deliver

// docs:start ExampleIdentifier
func ExampleIdentifier() {
	id, err := quotastripe.Identifier(quotastripe.Group{Namespace: "history-1", Customer: "cus_a", MeterID: "m_a"}, 0, 1)
	fmt.Println(len(id), err)
	// Output: 36 <nil>
}

// docs:end ExampleIdentifier

// docs:start ExampleExporter_Flush
func ExampleExporter_Flush() {
	ctx := context.Background()
	store, err := quotastripe.NewMemoryStore(quotastripe.StoreConfig{})
	if err != nil {
		fmt.Println(err)
		return
	}
	e, err := quotastripe.New(quotastripe.Config{Namespace: "history-1", Store: store, Client: &http.Client{}, APIKey: func(context.Context) (string, error) { return "injected", nil }, Customer: func(context.Context, quota.ConsumptionEvent) (string, error) { return "cus_a", nil }, Meters: map[string]quotastripe.Meter{"tokens": {ID: "m_a", EventName: "tokens", CustomerKey: "customer", ValueKey: "value"}}})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = e.Flush(ctx, time.Now())
	fmt.Println(err)
	// Output: quotastripe: dispatcher not running
}

// docs:end ExampleExporter_Flush

// docs:start exampleExporter
func exampleExporter() (*quotastripe.Exporter, error) {
	s, err := quotastripe.NewMemoryStore(quotastripe.StoreConfig{})
	if err != nil {
		return nil, err
	}
	return quotastripe.New(quotastripe.Config{Namespace: "example-history", Store: s, Client: &http.Client{Transport: exampleTransport{}}, APIKey: func(context.Context) (string, error) { return "injected", nil }, Customer: func(context.Context, quota.ConsumptionEvent) (string, error) { return "cus_example", nil }, Meters: map[string]quotastripe.Meter{"tokens": {ID: "m_a", EventName: "tokens", CustomerKey: "customer", ValueKey: "value"}}})
}

// docs:end exampleExporter

type exampleTransport struct{}

func (exampleTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"object":"list","has_more":false,"data":[]}`)), Request: r}, nil
}

// docs:start ExampleExporter_Status
func ExampleExporter_Status() {
	e, err := exampleExporter()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	d := quota.ConsumptionEvent{ID: "event", Meter: "tokens", At: time.Now().UTC().Truncate(time.Millisecond), Amount: 12}
	if err := e.Deliver(ctx, quota.Event{ID: d.ID, Consumption: &d}); err != nil {
		fmt.Println(err)
		return
	}
	status, err := e.Status(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("open units:", status.States[quotastripe.Open].Amount)
	fmt.Println("fenced groups:", status.FencedGroups)
	// Output:
	// open units: 12
	// fenced groups: 0
}

// docs:end ExampleExporter_Status

// docs:start ExampleExporter_Run
func ExampleExporter_Run() {
	e, err := exampleExporter()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	// Cancellation before startup performs no network requests. Services normally
	// keep this context live until shutdown and join the Run goroutine.
	cancel()
	done := make(chan error, 1)
	go func() { done <- e.Run(ctx) }()
	fmt.Println(<-done)
	// Output: context canceled
}

// docs:end ExampleExporter_Run

// docs:start ExampleExporter_Reconcile
func ExampleExporter_Reconcile() {
	e, err := exampleExporter()
	if err != nil {
		fmt.Println(err)
		return
	}
	end := time.Now().UTC().Truncate(time.Minute)
	r, err := e.Reconcile(context.Background(), "cus_example", "m_a", end.Add(-time.Hour), end)
	if err != nil {
		fmt.Println(err)
		return
	}
	// A Stripe total of zero does not establish coverage or source catch-up.
	fmt.Println(r.Verdict, r.InputComplete)
	// Output: incomplete unknown
}

// docs:end ExampleExporter_Reconcile

// docs:start ExampleExporter_Inspect
func ExampleExporter_Inspect() {
	e, err := exampleExporter()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	event := quota.ConsumptionEvent{ID: "bounded-example", Meter: "tokens", Amount: 4, At: time.Now().UTC().Truncate(time.Millisecond)}
	if err = e.Deliver(ctx, quota.Event{ID: event.ID, Consumption: &event}); err != nil {
		fmt.Println(err)
		return
	}
	page, err := e.Inspect(ctx, quotastripe.Filter{Kind: "batches", State: quotastripe.Open, Budget: quotastripe.InspectionBudget{Groups: 10, Snapshots: 1, DecodedBytes: 1 << 20}}, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Keep calling with page.Next, even for empty pages, until it is empty.
	fmt.Println("batches:", len(page.Batches), "snapshots:", page.Work.Snapshots)
	// Output: batches: 1 snapshots: 1
}

// docs:end ExampleExporter_Inspect

// docs:start ExampleExporter_ReconcileBounded
func ExampleExporter_ReconcileBounded() {
	e, err := exampleExporter()
	if err != nil {
		fmt.Println(err)
		return
	}
	end := time.Now().UTC().Truncate(time.Minute)
	report, err := e.ReconcileBounded(context.Background(), "cus_example", "m_a", end.Add(-time.Hour), end, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	// A bounded read with missing local coverage still cannot establish a match.
	fmt.Println(report.Verdict, report.InputComplete)
	// Output: incomplete unknown
}

// docs:end ExampleExporter_ReconcileBounded

// ExampleExporter_FlushOnce shows that durable barriers need a local Run.
// docs:start ExampleExporter_FlushOnce
func ExampleExporter_FlushOnce() {
	store, err := quotastripe.NewMemoryStore(quotastripe.StoreConfig{})
	if err != nil {
		panic(err)
	}
	exporter, err := quotastripe.New(quotastripe.Config{Namespace: "example", Store: store, Client: &http.Client{}, APIKey: func(context.Context) (string, error) { return "example-key", nil }, Customer: func(context.Context, quota.ConsumptionEvent) (string, error) { return "customer", nil }, Meters: map[string]quotastripe.Meter{"requests": {ID: "meter", EventName: "requests", CustomerKey: "customer", ValueKey: "value"}}})
	if err != nil {
		panic(err)
	}
	_, err = exporter.FlushOnce(context.Background(), quotastripe.FlushRequest{ID: "period-close-1", Before: time.Now().UTC().Truncate(time.Minute)})
	fmt.Println(errors.Is(err, quotastripe.ErrNotRunning))
	// Output: true
}

// docs:end ExampleExporter_FlushOnce
