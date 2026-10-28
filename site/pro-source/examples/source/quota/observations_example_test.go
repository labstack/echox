// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"fmt"

	"echo.labstack.com/pro/quota"
)

// docs:start ExampleHooks
func ExampleHooks() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	hooks := quota.JoinHooks(quota.Hooks{
		OnDecision: func(_ context.Context, _ quota.Request, d quota.Decision) {
			fmt.Println(
				"decision",
				d.Op,
				d.Allowed,
			)
		},
		OnConsumption: func(_ context.Context, e quota.ConsumptionEvent) {
			fmt.Println(
				"consumption",
				e.Meter,
				e.Amount,
			)
		},
		OnThreshold: func(_ context.Context, e quota.ThresholdEvent) {
			fmt.Println(
				"threshold",
				e.Limit,
				e.Percent,
			)
		},
	}, quota.Hooks{
		OnSettlement: func(_ context.Context, o quota.SettlementObservation) {
			fmt.Println("settlement", o.Source, o.Outcome, o.Refunded)
		},
		OnCredit: func(_ context.Context, o quota.CreditObservation) {
			fmt.Println(
				"credit",
				o.Outcome,
			)
		},
		OnRelease: func(_ context.Context, o quota.ReleaseObservation) {
			fmt.Println(
				"release",
				o.Outcome,
			)
		},
		OnReset: func(_ context.Context, o quota.ResetObservation) {
			fmt.Println(
				"reset",
				o.Outcome,
				len(o.Limits),
			)
		},
		OnOperation: func(_ context.Context, o quota.OperationObservation) {
			fmt.Println("operation", o.Op, o.Outcome, o.CrossingCount)
		},
	})
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "p", Limits: []quota.Limit{
		quota.PerDay("day", 100).Thresholds(50),
		quota.PerDay("tokens", 100).Meter("tokens"),
		quota.Credits("credits").Meter("credits"),
	}}}, Hooks: hooks})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := l.Close(ctx); err != nil {
			panic(err)
		}
	}()
	id := quota.Identity{Tenant: "example"}
	if _, err := l.Allow(ctx, quota.Request{
		Identity: id,
		Cost:     quota.Cost{},
	}); err != nil {
		panic(err)
	}
	if _, err := l.AddCredits(ctx, quota.CreditRequest{
		Identity: id,
		Limit:    "credits",
		Amount:   10,
		ID:       "invoice",
	}); err != nil {
		panic(err)
	}
	res, d, err := l.Reserve(ctx, quota.Request{Identity: id, Cost: quota.Cost{
		"requests": 60,
		"tokens":   10,
	}})
	if err != nil || !d.Allowed {
		panic("reserve failed")
	}
	for range 2 {
		if _, err := l.Settle(ctx, res, quota.Cost{
			"requests": 50,
			"tokens":   5,
		}); err != nil {
			panic(err)
		}
	}
	if err := l.Release(ctx, d); err != nil {
		panic(err)
	}
	if _, err := l.Reset(ctx, quota.ResetRequest{Identity: id, Plan: "p"}); err != nil {
		panic(err)
	}
	// Output:
	// decision allow true
	// operation allow local_noop 0
	// credit applied
	// operation credit applied 0
	// decision reserve true
	// operation reserve applied 0
	// consumption requests 50
	// consumption tokens 5
	// threshold day 50
	// settlement call applied true
	// operation settle applied 1
	// settlement call replayed true
	// operation settle replayed 0
	// release local_noop
	// operation release local_noop 0
	// reset applied 2
	// operation reset applied 0
}

// docs:end ExampleHooks
