// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package webhookadmin_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/audit/auditadmin"
	"echo.labstack.com/pro/audit/auditpg"
	"echo.labstack.com/pro/webhook/webhookadmin"
	"echo.labstack.com/pro/webhook/webhookpg"
)

// New validates its configuration without database access and reports a
// *ConfigError naming the invalid field.
// docs:start ExampleNew
func ExampleNew() {
	_, err := webhookadmin.New(webhookadmin.Config{})
	fmt.Println(err, errors.Is(err, &webhookadmin.ConfigError{}))
	// Production passes the *webhookpg.Store returned by webhookpg.New.
	_, err = webhookadmin.New(webhookadmin.Config{ID: "Webhooks", Store: new(webhookpg.Store)})
	if e, ok := errors.AsType[*webhookadmin.ConfigError](err); ok {
		fmt.Println(e.Field)
	}
	// Output:
	// webhookadmin: invalid configuration: Store true
	// ID
}

// docs:end ExampleNew

// Example_mount wires an audited replay panel into an application. It needs
// PostgreSQL, so it is compiled but not run.
// docs:start Example_mount
func Example_mount() {
	ctx := context.Background()
	var (
		runtimePool *pgxpool.Pool  // connected as the webhook runtime role (grants.sql)
		auditStore  *auditpg.Store // a validated recorder-role audit store
		signingKey  []byte         // 32–64 caller-managed bytes
	)
	store, err := webhookpg.New(ctx, runtimePool, webhookpg.StoreConfig{Schema: "webhooks", Namespace: "billing"})
	if err != nil {
		fmt.Println(err)
		return
	}
	sink, err := auditadmin.New(auditStore)
	if err != nil {
		fmt.Println(err)
		return
	}
	panel, err := webhookadmin.New(webhookadmin.Config{Store: store, EnableReplay: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	e := echo.New()
	err = admin.Config{
		Authenticate: func(*echo.Context) (admin.Principal, error) {
			// Map the application's validated session; webhook views and
			// replay require a Global principal, and replay an Operator.
			return admin.Principal{}, errors.New("replace with the application's session mapping")
		},
		Panels:  []admin.Panel{panel},
		Keys:    [][]byte{signingKey},
		Audit:   sink,
		Origins: []string{"https://service.example"},
	}.Mount(e, "/admin")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(admin.VerifyRoutes(e, "/admin"))
}

// docs:end Example_mount
