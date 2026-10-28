// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/quota"
)

// docs:start ExampleLimiter_Settle
func ExampleLimiter_Settle() {
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "pro", Limits: []quota.Limit{
		quota.PerMonth("tokens", 1000).Meter("tokens"),
	}}}})
	ctx := context.Background()
	r, d, err := l.Reserve(ctx, quota.Request{
		Identity: quota.Identity{Tenant: "acme"},
		Cost:     quota.Cost{"tokens": 100},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed:", d.Allowed)
	settled, err := l.Settle(ctx, r, quota.Cost{"tokens": 60})
	if err != nil {
		panic(err)
	}
	fmt.Println("final tokens:", settled.Actual["tokens"])
	if err = l.Close(ctx); err != nil {
		panic(err)
	}
	// Output:
	// allowed: true
	// final tokens: 60
}

// docs:end ExampleLimiter_Settle
// docs:start ExampleReservation_MarshalText
func ExampleReservation_MarshalText() {
	ctx := context.Background()
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "pro", Limits: []quota.Limit{
		quota.PerMonth("tokens", 1000).Meter("tokens"),
	}}}})
	r, _, err := l.Reserve(ctx, quota.Request{
		Identity: quota.Identity{Tenant: "acme"},
		Cost:     quota.Cost{"tokens": 100},
	})
	if err != nil {
		panic(err)
	}
	text, err := r.MarshalText()
	if err != nil {
		panic(err)
	}
	var job quota.Reservation
	if err = job.UnmarshalText(text); err != nil {
		panic(err)
	}
	settled, err := l.Settle(ctx, job, quota.Cost{"tokens": 60})
	if err != nil {
		panic(err)
	}
	fmt.Println("job tokens:", settled.Actual["tokens"])
	if err = l.Close(ctx); err != nil {
		panic(err)
	}
	// Output: job tokens: 60
}

// docs:end ExampleReservation_MarshalText
// docs:start ExampleCharge
func ExampleCharge() {
	ctx := context.Background()
	events := make(chan quota.ConsumptionEvent, 2)
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "pro", Limits: []quota.Limit{quota.PerMinute(
		"rpm",
		600,
	), quota.PerMinute("tpm", 200000).Meter("tokens"), quota.PerMonth(
		"tokens_month",
		50000000,
	).Meter("tokens").Strict().Thresholds(80, 100)}}}, Hooks: quota.Hooks{
		OnConsumption: func(_ context.Context, e quota.ConsumptionEvent) { events <- e },
	}})
	e := echo.New()
	e.Use(quota.MustMiddleware(quota.MiddlewareConfig{
		Limiter: l,
		Settle:  true,
		Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "acme"}, nil
		},
		Cost: func(*echo.Context) (quota.Cost, error) {
			return quota.Cost{
				"requests": 1,
				"tokens":   100,
			}, nil
		},
	}))
	e.GET("/chat", func(c *echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/event-stream")
		// Accumulate the provider's total tokens while streaming; the last Charge wins.
		if _, err := fmt.Fprint(c.Response(), "data: hello\n\n"); err != nil {
			return err
		}
		if err := http.NewResponseController(c.Response()).Flush(); err != nil {
			return err
		}
		quota.Charge(
			c,
			quota.Cost{"tokens": 60},
		) // requests is omitted, so its reserved 1 is kept.
		return nil
	})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequestWithContext(ctx, http.MethodGet, "/chat", nil))
	if err := l.Close(ctx); err != nil {
		panic(err)
	} // Drain before exiting the process.
	fmt.Print(rec.Body.String())
	for range 2 {
		event := <-events
		if event.Meter == "tokens" {
			fmt.Println("billed tokens:", event.Amount)
		}
	}
	// Output:
	// data: hello
	//
	// billed tokens: 60
}

// docs:end ExampleCharge

// docs:start ExampleLimit_Thresholds
func ExampleLimit_Thresholds() {
	events := make(chan quota.ThresholdEvent, 1)
	ctx := context.Background()
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "pro", Limits: []quota.Limit{
		quota.PerMonth("tokens", 100).Meter("tokens").Thresholds(80, 100),
	}}}, Hooks: quota.Hooks{
		OnThreshold: func(_ context.Context, e quota.ThresholdEvent) { events <- e },
	}})
	d, err := l.Allow(ctx, quota.Request{
		Identity: quota.Identity{Tenant: "acme"},
		Cost:     quota.Cost{"tokens": 80},
	})
	if err != nil || !d.Allowed {
		panic("admission failed")
	}
	fmt.Println("threshold:", (<-events).Percent)
	if err = l.Close(ctx); err != nil {
		panic(err)
	}
	// Output: threshold: 80
}

// docs:end ExampleLimit_Thresholds

// docs:start ExampleLimiter_Reserve
func ExampleLimiter_Reserve() {
	ctx := context.Background()
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "pro", Limits: []quota.Limit{
		quota.PerMonth("tokens", 1000).Meter("tokens"),
	}}}})
	r, d, err := l.Reserve(ctx, quota.Request{
		Identity: quota.Identity{Tenant: "acme"},
		Cost:     quota.Cost{"tokens": 100},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed:", d.Allowed, "reserved:", r.Reserved["tokens"])
	if _, err = l.Settle(ctx, r, nil); err != nil {
		panic(err)
	}
	if err = l.Close(ctx); err != nil {
		panic(err)
	}
	// Output: allowed: true reserved: 100
}

// docs:end ExampleLimiter_Reserve
