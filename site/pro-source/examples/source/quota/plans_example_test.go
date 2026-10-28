// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"echo.labstack.com/pro/quota"
)

// docs:start ExampleParsePlans
func ExampleParsePlans() {
	plans, err := quota.ParsePlans(strings.NewReader(
		`[{"id":"team","limits":[{"id":"daily","per":"day","n":10}]}]`,
	))
	if err != nil {
		panic(err)
	}
	fmt.Println(plans[0].ID)
	// Output: team
}

// docs:end ExampleParsePlans
// docs:start ExampleLimiter_Reload
func ExampleLimiter_Reload() {
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("daily", 10),
	}}}})
	if err := l.Reload([]quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("daily", 20),
	}}}); err != nil {
		panic(err)
	}
	fmt.Println(l.Plans()[0].ID)
	// Output: team
}

// docs:end ExampleLimiter_Reload
// docs:start ExampleLimit_ByDimension
func ExampleLimit_ByDimension() {
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("model", 10).ByDimension("model"),
	}}}})
	d, err := l.Allow(context.TODO(), quota.Request{Identity: quota.Identity{
		Tenant: "acme",
	}, Dimensions: map[string]string{"model": "small"}, Overrides: []quota.Override{{
		Limit: "model",
		N:     20,
	}}})
	if err != nil {
		panic(err)
	}
	fmt.Println(d.Allowed, d.Results[0].Quota)
	// Output: true 20
}

// docs:end ExampleLimit_ByDimension
// docs:start ExampleLimit_In
func ExampleLimit_In() {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	_, err = quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("daily", 100).In(loc),
		quota.PerMonth("monthly", 1000).In(loc).Anchor(31),
	}}}})
	fmt.Println(err)
	// Output: <nil>
}

// docs:end ExampleLimit_In
// docs:start ExampleJoinHooks
func ExampleJoinHooks() {
	h := quota.JoinHooks(quota.Hooks{
		OnConsumption: func(context.Context, quota.ConsumptionEvent) {
			fmt.Println(
				"billing",
			)
		},
	}, quota.Hooks{
		OnConsumption: func(context.Context, quota.ConsumptionEvent) {
			fmt.Println(
				"metrics",
			)
		},
	})
	h.OnConsumption(context.TODO(), quota.ConsumptionEvent{})
	// Output:
	// billing
	// metrics
}

// docs:end ExampleJoinHooks

// docs:start ExamplePlan_MarshalJSON
func ExamplePlan_MarshalJSON() {
	p := quota.Plan{ID: "team", Limits: []quota.Limit{quota.PerDay("daily", 10)}}
	b, err := p.MarshalJSON()
	if err != nil {
		panic(err)
	}
	var next quota.Plan
	if err = json.Unmarshal(b, &next); err != nil {
		panic(err)
	}
	fmt.Println(next.ID)
	// Output: team
}

// docs:end ExamplePlan_MarshalJSON
// docs:start ExampleLimit_Anchor
func ExampleLimit_Anchor() {
	_, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerMonth("monthly", 100).Anchor(31),
	}}}})
	fmt.Println(err)
	// Output: <nil>
}

// docs:end ExampleLimit_Anchor
