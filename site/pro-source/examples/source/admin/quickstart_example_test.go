// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package admin_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaadmin"
)

func ExampleConfig_readOnly() {
	lim, err := quota.New(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
			quota.PerDay("calls", 100),
		}}},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := lim.Close(context.Background()); err != nil {
			panic(err)
		}
	}()
	panel, err := quotaadmin.New(quotaadmin.Config{Limiter: lim})
	if err != nil {
		panic(err)
	}
	operatorKey := rand.Text() // Demo only; load production credentials at runtime.
	e := echo.New()
	// docs:start site-anchor
	cfg := admin.Config{
		Panels: []admin.Panel{panel}, Instance: "local-demo",
		Authenticate: func(c *echo.Context) (admin.Principal, error) {
			if c.Request().Header.Get("Authorization") != "Bearer "+operatorKey {
				return admin.Principal{}, echo.ErrUnauthorized
			}
			// Viewer is the default role; no write actions are enabled.
			return admin.Principal{ID: "operator-1", Tenants: []string{"acme"},
				HeaderCredential: true}, nil
		},
	}
	if err := cfg.Mount(e, "/admin"); err != nil {
		panic(err)
	}
	if err := admin.VerifyRoutes(e, "/admin"); err != nil {
		panic(err)
	}
	// docs:end site-anchor
	for _, key := range []string{operatorKey, ""} {
		req := httptest.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"/admin",
			nil,
		)
		if key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)
		fmt.Println(res.Code)
	}
	// Output:
	// 200
	// 401
}
