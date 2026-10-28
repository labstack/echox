// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"fmt"
	"time"

	"echo.labstack.com/pro/quota"
)

// docs:start exampleLimiter
func exampleLimiter() *quota.Limiter {
	return quota.MustNew(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{Clock: time.Now}), Plans: []quota.Plan{{ID: "free", Limits: []quota.Limit{quota.PerMonth("requests", 1)}}}})
}

// docs:end exampleLimiter
// docs:start ExampleNew
func ExampleNew() {
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Plans: []quota.Plan{{ID: "free", Limits: []quota.Limit{quota.PerSecond("burst", 5), quota.PerMonth("requests", 10000)}}}})
	fmt.Println(err == nil)
	if l != nil {
		if err := l.Close(context.Background()); err != nil {
			panic(err)
		}
	}
	// Output: true
}

// docs:end ExampleNew
// docs:start ExampleLimiter_Allow
func ExampleLimiter_Allow() {
	l := exampleLimiter()
	d, err := l.Allow(context.Background(), quota.Request{Identity: quota.Identity{Tenant: "acme"}})
	fmt.Println(d.Allowed, err)
	if err := l.Close(context.Background()); err != nil {
		panic(err)
	}
	// Output: true <nil>
}

// docs:end ExampleLimiter_Allow
// docs:start ExampleLimiter_InstanceLocal
func ExampleLimiter_InstanceLocal() {
	l := exampleLimiter()
	fmt.Println(l.InstanceLocal())
	if err := l.Close(context.Background()); err != nil {
		panic(err)
	}
	// Output: true
}

// docs:end ExampleLimiter_InstanceLocal
// docs:start ExampleDecision_Explain
func ExampleDecision_Explain() {
	// An explicit decision keeps the date-bearing explanation reproducible.
	d := quota.Decision{Reason: quota.LimitExceeded, Plan: "free", Binding: 0,
		RetryAfter: 564 * time.Hour, Results: []quota.LimitResult{{Limit: "requests", Meter: "requests", Quota: 1,
			ResetAt: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}}}
	fmt.Println(d.Explain())
	// Output: denied: limit exceeded (plan "free"); limit "requests", meter requests, remaining 0 of 1; resets 2026-11-01T00:00:00Z; retry after 564h0m0s
}

// docs:end ExampleDecision_Explain
// docs:start ExampleLimiter_Usage
func ExampleLimiter_Usage() {
	l := exampleLimiter()
	u, err := l.Usage(context.Background(), quota.UsageRequest{Identity: quota.Identity{Tenant: "acme"}})
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Plan, u.Limits[0].Used, u.Limits[0].Remaining)
	if err := l.Close(context.Background()); err != nil {
		panic(err)
	}
	// Output: free 0 1
}

// docs:end ExampleLimiter_Usage
// docs:start ExampleLimiter_Reset
func ExampleLimiter_Reset() {
	l := exampleLimiter()
	r, err := l.Reset(context.Background(), quota.ResetRequest{Identity: quota.Identity{Tenant: "acme"}, Plan: "free"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Generations["requests"])
	if err := l.Close(context.Background()); err != nil {
		panic(err)
	}
	// Output: 1
}

// docs:end ExampleLimiter_Reset
