// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/problem/echoadapter"
	"echo.labstack.com/pro/quota"
)

// docs:start ExampleMiddlewareConfig_authenticatedPlans
func ExampleMiddlewareConfig_authenticatedPlans() {
	if err := authenticatedPlans(); err != nil {
		panic(err)
	}
	// Output:
	// 200
	// 200
	// 429
	// 403
	// 401
}

// docs:end ExampleMiddlewareConfig_authenticatedPlans

// docs:start authenticatedPlans
func authenticatedPlans() error {
	ctx := context.Background()
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		return err
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	token, record, err := auth.IssueKey("demo", ring, "alice",
		auth.KeyOptions{Kind: auth.User, TenantRestriction: "acme"})
	if err != nil {
		return err
	}
	store := &auth.MemoryKeyStore{Capacity: 10}
	if err := store.Put(ctx, record); err != nil {
		return err
	}
	authMW, err := (auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "customers", Prefix: "demo", Peppers: ring, Store: store,
		}},
		Fallback: auth.Require(map[string][]string{"key": nil}),
	}).ToMiddleware()
	if err != nil {
		return err
	}
	plans := []quota.Plan{
		{ID: "free", Limits: []quota.Limit{quota.PerDay("requests", 1)}},
		{ID: "team", Limits: []quota.Limit{quota.PerDay("requests", 2)}},
	}
	lim, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: plans})
	if err != nil {
		return err
	}
	defer func() {
		if err := lim.Close(ctx); err != nil {
			panic(err)
		}
	}()
	e := echo.New()
	e.HTTPErrorHandler, err = (echoadapter.ErrorHandlerConfig{}).ToErrorHandler(
		e.HTTPErrorHandler,
	)
	if err != nil {
		return err
	}

	// docs:start site-plans
	// Application-owned subscription data; replace with your plan repository.
	// This demo map is immutable while requests run.
	planByTenant := map[string]string{"acme": "team"}
	cfg := quota.MiddlewareConfig{
		Limiter: lim,
		Identity: func(c *echo.Context) (quota.Identity, error) {
			p, ok := auth.FromEcho(c)
			if !ok {
				return quota.Identity{}, echo.ErrUnauthorized
			}
			tenant := p.TenantRestriction()
			if tenant == "" {
				return quota.Identity{}, echo.ErrForbidden
			}
			return quota.Identity{Tenant: tenant, User: p.ID(), APIKeyID: p.KeyID()}, nil
		},
		Plan: func(_ *echo.Context, id quota.Identity) (string, error) {
			plan, ok := planByTenant[id.Tenant]
			if !ok {
				return "", echo.ErrForbidden
			}
			return plan, nil
		},
	}
	mw, err := cfg.ToMiddleware()
	if err != nil {
		return err
	}
	api := e.Group("/v1", authMW, mw)
	api.GET("/hello", func(c *echo.Context) error { return c.NoContent(http.StatusOK) })
	// docs:end site-plans

	err = (quota.AdminConfig{
		Limiter: lim,
		Authorize: func(c *echo.Context) error {
			p, ok := auth.FromEcho(c)
			if !ok || !slices.Contains(p.Scopes(), "quota:admin") {
				return echo.ErrForbidden
			}
			return nil
		},
	}).Register(e.Group("/admin", authMW))
	if err != nil {
		return err
	}
	for _, path := range []string{"/v1/hello", "/v1/hello", "/v1/hello", "/admin/plans"} {
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set("Authorization", "ApiKey "+token)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)
		fmt.Println(res.Code)
	}
	res := httptest.NewRecorder()
	e.ServeHTTP(res, httptest.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/v1/hello",
		nil,
	))
	fmt.Println(res.Code)
	return nil
}

// docs:end authenticatedPlans
