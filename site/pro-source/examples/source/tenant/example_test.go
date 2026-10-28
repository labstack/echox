// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package tenant_test

import (
	"context"
	"fmt"

	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/tenant"
)

// docs:start ExampleStatic
func ExampleStatic() {
	p, err := auth.NewPrincipal(auth.PrincipalConfig{Namespace: "keys", Subject: "worker", Scheme: "key", Kind: auth.Service, Scopes: []string{"invoice:read"}})
	if err != nil {
		panic(err)
	}
	a, err := tenant.Static([]tenant.Grant{{Scope: "invoice:read", ResourceType: "invoice", Actions: []string{"read"}}})
	if err != nil {
		panic(err)
	}
	check := tenant.Check{Principal: p, Tenant: "acme", Action: "read", Resource: tenant.Resource{Type: "invoice", ID: "inv-1", Tenant: "acme"}}
	own, err := a.Authorize(context.Background(), check)
	if err != nil {
		panic(err)
	}
	check.Resource.Tenant = "other"
	cross, err := a.Authorize(context.Background(), check)
	if err != nil {
		panic(err)
	}
	fmt.Println(own == tenant.Allow, cross == tenant.Deny)
	// Output: true true
}

// docs:end ExampleStatic
