// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package executionadmin_test

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/execution/executionadmin"
)

// docs:start ExampleNew
func ExampleNew() {
	journal, err := execution.New(execution.JournalConfig{Namespace: "payments", Store: execution.NewMemoryStore(nil), Policies: map[string]execution.Policy{
		"transfer": {SafetyVersion: "v1", ReconcilerVersion: "v1", Effects: []execution.Effect{{Name: "transfer", Provider: "bank/v1", RequestTimeout: time.Second}}},
	}})
	if err != nil {
		fmt.Println(err)
		return
	}
	panel, err := executionadmin.New(executionadmin.Config{Operator: journal.Operator()})
	if err != nil {
		fmt.Println(err)
		return
	}
	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		fmt.Println(err)
		return
	}
	e := echo.New()
	cfg := admin.Config{Panels: []admin.Panel{panel}, Keys: [][]byte{key}, Authenticate: func(*echo.Context) (admin.Principal, error) {
		// Replace this fixture with the application's authenticated session mapping.
		return admin.Principal{ID: "operator", Role: admin.Operator, Global: true, HeaderCredential: true}, nil
	}}
	if err = cfg.Mount(e, "/admin"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(admin.VerifyRoutes(e, "/admin"))
	// Output: <nil>
}

// docs:end ExampleNew
