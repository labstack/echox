// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auth_test

import (
	"crypto/rand"
	"fmt"

	"echo.labstack.com/pro/auth"
)

// docs:start ExampleIssueKey
func ExampleIssueKey() {
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		panic(err)
	}
	token, record, err := auth.IssueKey("prod", auth.PepperRing{{
		Version: 7,
		Key:     pepper,
	}}, "billing-worker", auth.KeyOptions{Kind: auth.Service, TenantRestriction: "acme"})
	if err != nil {
		panic(err)
	} // Deliver token once through a secure channel; persist only record.
	fmt.Println(
		len(token),
		len(record.ID),
		record.PepperVersion,
		record.TenantRestriction,
	)
	// Output: 71 22 7 acme
}

// docs:end ExampleIssueKey
// docs:start ExampleNewPrincipal
func ExampleNewPrincipal() {
	p, err := auth.NewPrincipal(auth.PrincipalConfig{
		Namespace: "employees",
		Subject:   "alice",
		Scheme:    "oauth",
		Kind:      auth.User,
		Scopes:    []string{"invoice:read"},
	})
	if err != nil {
		panic(err)
	}
	copy := p.Scopes()
	copy[0] = "invoice:write"
	fmt.Println(p.ID(), p.Scopes()[0])
	// Output: employees:alice invoice:read
}

// docs:end ExampleNewPrincipal
