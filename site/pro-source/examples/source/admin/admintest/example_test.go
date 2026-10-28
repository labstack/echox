// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package admintest_test

import (
	"fmt"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/admin/admintest"
)

// docs:start ExampleTenantPrincipal
func ExampleTenantPrincipal() {
	p := admintest.TenantPrincipal("acme", admin.Viewer, false)
	fmt.Println(p.Tenants, p.Global, p.Reveal)
	// Output: [acme] false false
}

// docs:end ExampleTenantPrincipal
