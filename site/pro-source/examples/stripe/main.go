// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

// Command stripe demonstrates quota/outbox wiring. It is offline by default.
// -online makes owner-authorized Stripe requests using an injected environment.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotastripe"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	online := flag.Bool(
		"online",
		false,
		"make Stripe calls with owner-injected sandbox configuration",
	)
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := demo(ctx, *online); err != nil {
		return err
	}

	return nil
}

func demo(ctx context.Context, online bool) (resultErr error) {
	// Demonstration only: replace with a persistent RedisStore for paid usage.
	store, err := quotastripe.NewMemoryStore(quotastripe.StoreConfig{})
	if err != nil {
		return err
	}
	var caughtUp atomic.Bool
	cfg := quotastripe.Config{
		Namespace: "demo-single-process-history-1",
		Store:     store,
		Client:    &http.Client{},
		APIKey: func(context.Context) (string, error) {
			key := os.Getenv("STRIPE_SANDBOX_KEY")
			if key == "" {
				return "", errors.New("inject STRIPE_SANDBOX_KEY")
			}
			return key, nil
		},
		Customer: func(context.Context, quota.ConsumptionEvent) (string, error) {
			if !online {
				return "cus_demo", nil
			}
			customer := os.Getenv("STRIPE_SANDBOX_CUSTOMER_ID")
			if customer == "" {
				return "", errors.New("inject STRIPE_SANDBOX_CUSTOMER_ID")
			}
			return customer, nil
		},
		Drained: func(context.Context) (bool, error) { return caughtUp.Load(), nil },
		Meters: map[string]quotastripe.Meter{"tokens": {
			ID:          "mtr_demo",
			EventName:   "tokens",
			CustomerKey: "stripe_customer_id",
			ValueKey:    "value",
		}},
	}
	if online {
		m := cfg.Meters["tokens"]
		m.ID = os.Getenv("STRIPE_SANDBOX_METER_ID")
		m.EventName = os.Getenv("STRIPE_SANDBOX_EVENT_NAME")
		cfg.Meters["tokens"] = m
	}
	exporter, err := quotastripe.New(cfg)
	if err != nil {
		return err
	}
	accepted := make(chan struct{})
	var once sync.Once
	//nolint:contextcheck // Outbox owns its background delivery contexts.
	lim, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "demo", Limits: []quota.Limit{
		quota.PerDay("tokens", 100).Meter("tokens").Strict(),
	}}}, Outbox: &quota.Outbox{PollInterval: time.Millisecond, Deliver: func(
		ctx context.Context,
		ev quota.Event,
	) error {
		err := exporter.Deliver(ctx, ev)
		if err == nil {
			once.Do(func() { close(accepted) })
		}
		return err
	}}})
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if closed {
			return
		}
		call, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		resultErr = errors.Join(resultErr, lim.Close(call))
	}()
	periodStart := time.Now().UTC().Truncate(time.Minute)
	if _, err = lim.Allow(ctx, quota.Request{
		Identity: quota.Identity{Tenant: "demo"},
		Cost:     quota.Cost{"tokens": 12},
	}); err != nil {
		return err
	}
	select {
	case <-accepted:
	case <-ctx.Done():
		return ctx.Err()
	}
	// Stop all source admissions. This example has no reservations or future
	// settlements; Close joins outbox processing before affirming catch-up.
	if err := lim.Close(ctx); err != nil {
		return err
	}
	closed = true
	caughtUp.Store(true)
	status, err := exporter.Status(ctx)
	if err != nil {
		return err
	}
	fmt.Println("accepted open units:", status.States[quotastripe.Open].Amount)
	if !online {
		fmt.Println("offline demo: no Stripe API calls")
		return nil
	}
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	finished := make(chan error, 1)
	go func() { finished <- exporter.Run(runCtx) }()
	// A production service owns dispatcher startup/readiness. Flush reports
	// ErrNotRunning until preflight and the namespace lease have succeeded.
	periodEnd := periodStart.Add(time.Minute)
	timer := time.NewTimer(time.Until(periodEnd))
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
		cancel()
		runErr := <-finished
		if errors.Is(runErr, context.Canceled) {
			runErr = nil
		}
		return errors.Join(ctx.Err(), runErr)
	case err = <-finished:
		return err
	}
	result, flushErr := exporter.Flush(ctx, periodEnd)
	fmt.Println(
		"targeted/sent/unresolved:",
		result.Targeted,
		result.Sent,
		result.Unresolved,
	)
	var report quotastripe.Report
	if flushErr == nil {
		report, err = exporter.Reconcile(ctx, os.Getenv(
			"STRIPE_SANDBOX_CUSTOMER_ID",
		), cfg.Meters["tokens"].ID, periodStart, periodEnd)
		if err == nil {
			fmt.Println("aggregate verdict:", report.Verdict)
		}
	}
	cancel()
	runErr := <-finished
	if errors.Is(runErr, context.Canceled) {
		runErr = nil
	}
	return errors.Join(flushErr, err, runErr)
}
