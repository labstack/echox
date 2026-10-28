// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package admin_test

import (
	"crypto/rand"
	"fmt"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaadmin"
)

// docs:start ExampleConfig_Mount
func ExampleConfig_Mount() {
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("calls", 100),
	}}}})
	if err != nil {
		fmt.Println(err)
		return
	}
	p, err := quotaadmin.New(quotaadmin.Config{Limiter: l})
	if err != nil {
		fmt.Println(err)
		return
	}
	// A fixture login binding; map the actual validated session in production.
	session := rand.Text()
	e := echo.New()
	cfg := admin.Config{Panels: []admin.Panel{p}, Instance: "node-1", Authenticate: func(
		*echo.Context,
	) (admin.Principal, error) {
		return admin.Principal{
			ID:      "alice",
			Session: session,
			Tenants: []string{"tenant-a"},
		}, nil
	}}
	if err = cfg.Mount(e, "/admin"); err != nil {
		fmt.Println(err)
		return
	}
	// Register remaining app routes before the final reservation check.
	fmt.Println(admin.VerifyRoutes(e, "/admin"))
	fmt.Println(admin.RouteKeys("/admin"))
	// Output:
	// <nil>
	// [echo_route_any /admin echo_route_any /admin/*]
}

// docs:end ExampleConfig_Mount
