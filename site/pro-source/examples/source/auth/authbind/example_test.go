// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package authbind_test

import (
	"fmt"
	"time"

	mcpauth "github.com/modelcontextprotocol/go-sdk/auth"

	"echo.labstack.com/pro/auth/authbind"
)

// docs:start ExampleMCPMapping_Principal
func ExampleMCPMapping_Principal() {
	m := authbind.MCPMapping{Namespace: "tools", ExtraTenantKey: "verified_tenant"}
	info := &mcpauth.TokenInfo{UserID: "alice", Expiration: time.Now().Add(
		time.Hour,
	), Extra: map[string]any{
		"verified_tenant": "acme",
	}}
	p, err := m.Principal(info)
	if err != nil {
		panic(err)
	}
	fmt.Println(p.ID(), p.TenantRestriction())
	// Output: tools:alice acme
}

// docs:end ExampleMCPMapping_Principal
