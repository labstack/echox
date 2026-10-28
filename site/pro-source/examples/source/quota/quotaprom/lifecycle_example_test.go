// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotaprom_test

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaprom"
)

// docs:start gathered
func gathered(reg *prometheus.Registry) {
	families, err := reg.Gather()
	if err != nil {
		panic(err)
	}
	for _, f := range families {
		fmt.Println(f.GetName())
	}
}

// docs:end gathered

// docs:start ExampleNewLifecycle
func ExampleNewLifecycle() {
	ctx := context.Background()
	reg := prometheus.NewRegistry()
	m, err := quotaprom.NewLifecycle(reg, quotaprom.LifecycleConfig{Names: quotaprom.NameAllowlist{Plans: []string{"team"}}})
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
	gathered(reg)
	// Output: quota_admissions_total
}

// docs:end ExampleNewLifecycle

// docs:start ExampleNewFailureCounter
func ExampleNewFailureCounter() {
	reg := prometheus.NewRegistry()
	f, err := quotaprom.NewFailureCounter(reg)
	if err != nil {
		panic(err)
	}
	f.Fail("application-specific")
	gathered(reg)
	// Output: quota_telemetry_failures_total
}

// docs:end ExampleNewFailureCounter
