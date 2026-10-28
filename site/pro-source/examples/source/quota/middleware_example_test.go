// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quota_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/quota"
)

// docs:start closeExampleLimiter
func closeExampleLimiter(l *quota.Limiter) {
	if err := l.Close(context.Background()); err != nil {
		panic(err)
	}
}

// docs:end closeExampleLimiter

// docs:start ExampleMiddlewareConfig_ToMiddleware
func ExampleMiddlewareConfig_ToMiddleware() {
	// Freeze token refill for reproducible output, starting at the current calendar.
	now := time.Now()
	lim, err := quota.New(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{
			Clock: func() time.Time { return now },
		}),
		Plans: []quota.Plan{{ID: "free", Limits: []quota.Limit{
			quota.PerSecond("burst", 5).Burst(10),
			quota.PerMonth("requests", 10000),
		}}},
	})
	if err != nil {
		panic(err)
	}
	defer closeExampleLimiter(lim)
	cfg := quota.MiddlewareConfig{Limiter: lim, Identity: func(
		c *echo.Context,
	) (quota.Identity, error) {
		return quota.Identity{Tenant: c.Request().Header.Get("X-Tenant")}, nil
	}}
	mw, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}
	usage, err := cfg.UsageHandler()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.Group("/v1", mw).GET(
		"/hello",
		func(c *echo.Context) error { return c.NoContent(204) },
	)
	e.GET("/usage", usage)
	request := func(path string) *httptest.ResponseRecorder {
		req := httptest.NewRequestWithContext(context.Background(), "GET", path, nil)
		req.Header.Set(
			"X-Tenant",
			"acme",
		) // Replace with authenticated claims in production.
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec
	}
	fmt.Println(request("/v1/hello").Code)
	for range 9 {
		request("/v1/hello")
	}
	rec := request("/v1/hello")
	var p struct {
		Policies []string `json:"violated-policies"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		panic(err)
	}
	fmt.Println(rec.Code, p.Policies, rec.Header().Get("Retry-After"))
	fmt.Println(request("/usage").Code)
	// Output:
	// 204
	// 429 [burst] 1
	// 200
}

// docs:end ExampleMiddlewareConfig_ToMiddleware

// docs:start ExampleMiddlewareConfig_ToMiddleware_plans
func ExampleMiddlewareConfig_ToMiddleware_plans() {
	lim := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{
		{ID: "free", Limits: []quota.Limit{
			quota.PerSecond("burst", 5).Burst(10),
			quota.PerMonth("requests", 10000),
		}},
		{ID: "team", Limits: []quota.Limit{quota.PerSecond(
			"burst",
			50,
		).Burst(100), quota.PerMonth(
			"requests",
			1000000,
		)}},
	}})
	defer closeExampleLimiter(lim)
	type claims struct{ org, user, plan string }
	// Your authentication middleware sets verified claims before calling next.
	auth := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("claims", claims{
				"acme",
				"alice",
				"team",
			})
			return next(c)
		}
	}
	getClaims := func(c *echo.Context) (claims, error) {
		cl, ok := c.Get("claims").(claims)
		if !ok {
			return cl, errors.New("missing authenticated claims")
		}
		return cl, nil
	}
	mw := quota.MustMiddleware(quota.MiddlewareConfig{
		Limiter: lim,
		Identity: func(c *echo.Context) (quota.Identity, error) {
			cl, err := getClaims(c)
			return quota.Identity{Tenant: cl.org, User: cl.user}, err
		},
		Plan: func(c *echo.Context, _ quota.Identity) (string, error) {
			cl, err := getClaims(
				c,
			)
			return cl.plan, err
		},
	})
	e := echo.New()
	e.Group("/v1", auth, mw).GET("/hello", func(c *echo.Context) error {
		d, ok := quota.DecisionFrom(c)
		if !ok {
			panic("missing quota decision")
		}
		return c.String(200, d.Plan)
	})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequestWithContext(
		context.Background(),
		"GET",
		"/v1/hello",
		nil,
	))
	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 team
}

// docs:end ExampleMiddlewareConfig_ToMiddleware_plans

// docs:start ExampleByIP
func ExampleByIP() {
	guard := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{MaxKeys: 100000},
	), Plans: []quota.Plan{{ID: "anon", Limits: []quota.Limit{
		quota.PerSecond("ip", 20).Burst(40).By(quota.IP),
	}}}})
	defer closeExampleLimiter(guard)
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()
	e.Pre(quota.MustMiddleware(quota.MiddlewareConfig{
		Limiter:  guard,
		Identity: quota.ByIP,
		Headers:  quota.NoHeaders,
	}))
	e.GET("/", func(c *echo.Context) error { return c.NoContent(204) })
	req := httptest.NewRequestWithContext(context.Background(), "GET", "/", nil)
	req.RemoteAddr = "[2001:db8:1:2::1234]:443"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	fmt.Println(rec.Code, rec.Header().Get("RateLimit") == "")
	// Output: 204 true
}

// docs:end ExampleByIP

// docs:start ExampleRouteCosts
func ExampleRouteCosts() {
	lim := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "free", Limits: []quota.Limit{
		quota.PerDay("requests", 100),
	}}}})
	defer closeExampleLimiter(lim)
	costs := map[string]int64{"GET /v1/search": 5, "POST /v1/export": 100}
	e := echo.New()
	e.Use(quota.MustMiddleware(quota.MiddlewareConfig{
		Limiter: lim,
		Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "acme"}, nil
		},
		Cost: quota.RouteCosts(costs),
	}))
	handler := func(c *echo.Context) error {
		d, ok := quota.DecisionFrom(c)
		if !ok {
			panic("missing quota decision")
		}
		return c.String(200, fmt.Sprint(d.Results[0].Remaining))
	}
	e.GET("/v1/search", handler)
	e.POST("/v1/export", handler)
	if err := quota.CheckRouteCosts(e, costs); err != nil {
		panic(err)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequestWithContext(
		context.Background(),
		"GET",
		"/v1/search",
		nil,
	))
	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 95
}

// docs:end ExampleRouteCosts

// docs:start ExampleMiddlewareConfig_UsageHandler
func ExampleMiddlewareConfig_UsageHandler() {
	lim := exampleLimiter()
	defer closeExampleLimiter(lim)
	cfg := quota.MiddlewareConfig{
		Limiter: lim,
		Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "acme"}, nil
		},
	}
	h, err := cfg.UsageHandler()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.GET("/usage", h)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequestWithContext(
		context.Background(),
		"GET",
		"/usage?tenant=other",
		nil,
	))
	var usage struct {
		Plan   string
		Limits []struct{ Used, Remaining int64 }
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &usage); err != nil {
		panic(err)
	}
	fmt.Println(rec.Code, usage.Plan, usage.Limits[0].Used, usage.Limits[0].Remaining)
	// Output: 200 free 0 1
}

// docs:end ExampleMiddlewareConfig_UsageHandler

// docs:start ExampleIPKey
func ExampleIPKey() {
	key, err := quota.IPKey("2001:db8:abcd:1::1234", 64)
	fmt.Println(key, err)
	// Output: 2001:db8:abcd:1::/64 <nil>
}

// docs:end ExampleIPKey

// docs:start ExampleMiddlewareConfig_ToMiddleware_exclusive
func ExampleMiddlewareConfig_ToMiddleware_exclusive() {
	lim := quota.MustNew(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "free", Limits: []quota.Limit{
		quota.PerDay("requests", 100),
	}}}})
	defer closeExampleLimiter(lim)
	cfg := quota.MiddlewareConfig{
		Limiter: lim,
		Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "acme"}, nil
		},
		Exclusive: true,
	}
	boundary := quota.MustMiddleware(cfg)
	cfg.Exclusive = false
	legacy := quota.MustMiddleware(cfg)
	c := echo.New().NewContext(httptest.NewRequestWithContext(
		context.Background(),
		"GET",
		"/",
		nil,
	), httptest.NewRecorder())
	// Remove this legacy stage when adopting the chain; it is refused here.
	err := boundary(legacy(func(*echo.Context) error { return nil }))(c)
	fmt.Println(errors.Is(err, quota.ErrDuplicateAdmission))
	d, ok := quota.DecisionFrom(c)
	fmt.Println(ok, d.Allowed, d.Results[0].Remaining)
	// Output:
	// true
	// true true 99
}

// docs:end ExampleMiddlewareConfig_ToMiddleware_exclusive
