// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auth_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/problem/echoadapter"
)

func ExampleConfig_quickstart() {
	ctx := context.Background()
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		panic(err)
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	token, record, err := auth.IssueKey("demo", ring, "billing-worker",
		auth.KeyOptions{
			Kind:              auth.Service,
			TenantRestriction: "acme",
			Scopes:            []string{"invoice:read"},
		})
	if err != nil {
		panic(err)
	}
	store := &auth.MemoryKeyStore{Capacity: 10}
	if err := store.Put(ctx, record); err != nil {
		panic(err)
	}

	// docs:start site-anchor
	cfg := auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "billing", Prefix: "demo", Peppers: ring, Store: store,
			Header: "X-API-Key",
		}},
		Fallback: auth.Require(map[string][]string{"key": {"invoice:read"}}),
	}
	requireKey, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	render := echoadapter.ErrorHandlerConfig{}
	e.HTTPErrorHandler, err = render.ToErrorHandler(e.HTTPErrorHandler)
	if err != nil {
		panic(err)
	}
	e.Group("/v1", requireKey).GET("/invoices", func(c *echo.Context) error {
		p, _ := auth.FromEcho(c)
		return c.String(http.StatusOK, p.ID()+" tenant="+p.TenantRestriction())
	})
	// docs:end site-anchor
	for _, key := range []string{token, ""} {
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/invoices", nil)
		if key != "" {
			req.Header.Set("X-API-Key", key)
		}
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)
		fmt.Println(res.Code, res.Header().Get("Content-Type"))
	}
	// Output:
	// 200 text/plain; charset=UTF-8
	// 401 application/problem+json
}
