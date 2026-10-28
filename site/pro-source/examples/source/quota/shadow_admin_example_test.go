// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"fmt"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/quota"
)

// docs:start ExampleLimiter_ShadowReport
func ExampleLimiter_ShadowReport() {
	ctx := context.Background()
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Plans: []quota.Plan{
		{ID: "live", Limits: []quota.Limit{quota.PerDay("daily", 10)}},
		{ID: "candidate", Limits: []quota.Limit{quota.PerDay("daily", 1)}},
	}})
	defer func() {
		if err := l.Close(ctx); err != nil {
			panic(err)
		}
	}()
	r := quota.Request{Identity: quota.Identity{Tenant: "acme"}, Plan: "live", Shadow: "candidate"}
	for range 2 {
		if _, err := l.Allow(ctx, r); err != nil {
			panic(err)
		}
	}
	report := l.ShadowReport()
	fmt.Println(report[0].Evaluated, report[0].OnlyShadowDenied)
	// Output: 2 1
}

// docs:end ExampleLimiter_ShadowReport

// docs:start ExampleTopK
func ExampleTopK() {
	ctx := context.Background()
	top := quota.NewTopK(5)
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{quota.PerDay("requests", 100)}}}, Hooks: quota.JoinHooks(top.Hooks())})
	defer func() {
		if err := l.Close(ctx); err != nil {
			panic(err)
		}
	}()
	if _, err := l.Allow(ctx, quota.Request{Identity: quota.Identity{Tenant: "acme"}, Cost: quota.Cost{"requests": 7}}); err != nil {
		panic(err)
	}
	fmt.Println(top.Top("requests"))
	// Output: [{acme 7 0}]
}

// docs:end ExampleTopK

// docs:start ExampleAdminConfig_Register
func ExampleAdminConfig_Register() {
	ctx := context.Background()
	l := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{quota.PerDay("requests", 100)}}}})
	defer func() {
		if err := l.Close(ctx); err != nil {
			panic(err)
		}
	}()
	e := echo.New()
	err := (quota.AdminConfig{Limiter: l, Authorize: func(c *echo.Context) error {
		if c.Request().Header.Get("X-Admin") != "trusted" {
			return echo.ErrForbidden
		}
		return nil // Production: check authenticated admin claims.
	}}).Register(e.Group("/admin"))
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequestWithContext(ctx, "GET", "/admin/plans", nil)
	req.Header.Set("X-Admin", "trusted")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	fmt.Println(rec.Code, rec.Header().Get("Cache-Control"))
	// Output: 200 no-store
}

// docs:end ExampleAdminConfig_Register
