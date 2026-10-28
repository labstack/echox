// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auditadmin_test

import (
	"fmt"
	"strings"
	"time"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/audit/auditadmin"
	"echo.labstack.com/pro/redact"
)

// docs:start ExamplePrepare
func ExamplePrepare() {
	p, err := auditadmin.Prepare(admin.AuditRecord{
		OperationID: strings.Repeat(
			"a",
			64,
		),
		EmissionID: strings.Repeat(
			"b",
			64,
		),
		SemanticDigest: strings.Repeat(
			"c",
			64,
		),
		Phase:  "intent",
		Actor:  "operator",
		Panel:  "quota",
		Action: "adjust",
		Effect: admin.Read,
		At: time.Date(
			2026,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		),
	})
	fmt.Println(err, p.Valid())
	// Output: <nil> true
}

// docs:end ExamplePrepare

// docs:start ExampleNewPanel
func ExampleNewPanel() {
	// Synthetic example key; obtain production keys from a secret manager.
	keys, err := redact.NewKeyring([]redact.Key{{ID: "k1", Material: []byte(
		"01234567890123456789012345678901",
	)}})
	if err != nil {
		panic(err)
	}
	schema, err := audit.NewSchema(audit.SchemaConfig{Keys: keys})
	if err != nil {
		panic(err)
	}
	store, err := audit.NewMemoryStore(audit.MemoryStoreConfig{Capacity: audit.Capacity{
		Events: 100,
		Bytes:  1 << 20,
	}})
	if err != nil {
		panic(err)
	}
	reader, err := store.Reader(audit.ReaderConfig{Schema: schema, Access: store})
	if err != nil {
		panic(err)
	}
	panel, err := auditadmin.NewPanel(auditadmin.PanelConfig{
		Reader:       reader,
		EnableReveal: true,
	})
	if err != nil {
		panic(err)
	}
	// Mount panel in admin.Config.Panels with Keys and authenticated principals.
	// Reveal requires the principal's separate Reveal grant on every request.
	descriptor := panel.AdminPanel().Describe()
	fmt.Println(descriptor.ID, len(descriptor.Views), len(descriptor.Actions))
	// Output: audit 3 0
}

// docs:end ExampleNewPanel
