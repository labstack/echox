// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/quota"
)

// docs:start ExamplePerBillingPeriod
func ExamplePerBillingPeriod() {
	ctx := context.Background()
	lim, err := quota.New(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "subscription", Limits: []quota.Limit{
			quota.PerBillingPeriod("requests", 2),
		}}},
	})
	if err != nil {
		panic(err)
	}
	defer closeExampleLimiter(lim)
	// Trusted subscription data, immutable in this fixture. A repository lookup
	// uses c.Request().Context(); the request body never selects billing bounds.
	period := quota.Period{Start: time.Date(
		2026,
		10,
		15,
		0,
		0,
		0,
		0,
		time.UTC,
	), End: time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)}
	middleware, err := (quota.MiddlewareConfig{
		Limiter: lim, Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "acme"}, nil
		},
		Period: func(*echo.Context, quota.Identity) (
			quota.Period,
			error,
		) {
			return period, nil
		},
	}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.GET("/invoices", func(c *echo.Context) error {
		return c.NoContent(
			http.StatusNoContent,
		)
	}, middleware)
	for range 3 {
		w := httptest.NewRecorder()
		e.ServeHTTP(w, httptest.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"/invoices",
			nil,
		))
		fmt.Println(w.Code)
	}
	// Output:
	// 204
	// 204
	// 429
}

// docs:end ExamplePerBillingPeriod

// docs:start ExampleConcurrent
func ExampleConcurrent() {
	ctx := context.Background()
	lim, err := quota.New(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "jobs", Limits: []quota.Limit{
			quota.Concurrent("running", 1),
		}}},
	})
	if err != nil {
		panic(err)
	}
	defer closeExampleLimiter(lim)
	request := quota.Request{Identity: quota.Identity{Tenant: "acme"}}
	first, err := lim.Allow(ctx, request)
	if err != nil {
		panic(err)
	}
	second, err := lim.Allow(ctx, request)
	if err != nil {
		panic(err)
	}
	if err := lim.Release(ctx, first); err != nil {
		panic(err)
	}
	third, err := lim.Allow(ctx, request)
	if err != nil {
		panic(err)
	}
	fmt.Println(first.Allowed, second.Allowed, third.Allowed)
	if err := lim.Release(ctx, third); err != nil {
		panic(err)
	}
	// Output: true false true
}

// docs:end ExampleConcurrent

// docs:start ExampleCredits
func ExampleCredits() {
	ctx := context.Background()
	lim, err := quota.New(quota.Config{Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "paid", Limits: []quota.Limit{
			quota.Credits("balance").Meter("tokens"),
		}}},
	})
	if err != nil {
		panic(err)
	}
	defer closeExampleLimiter(lim)
	credit := quota.CreditRequest{
		Identity: quota.Identity{Tenant: "acme"},
		Plan:     "paid",
		Limit:    "balance",
		ID:       "invoice-42",
		Amount:   100,
	}
	for range 2 {
		result, err := lim.AddCredits(ctx, credit)
		if err != nil {
			panic(err)
		}
		fmt.Println(result.Balance, result.Applied)
	}
	credit.Amount = 200
	_, err = lim.AddCredits(ctx, credit)
	fmt.Println(errors.Is(err, quota.ErrCreditConflict))
	// Output:
	// 100 true
	// 100 false
	// true
}

// docs:end ExampleCredits
