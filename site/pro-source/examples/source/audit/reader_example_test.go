// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package audit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/redact"
)

// docs:start ExampleReader_Query
func ExampleReader_Query() {
	// Synthetic example key. Load production keys through your secret manager.
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
		Events: 10,
		Bytes:  1 << 20,
	}})
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	p, err := schema.Prepare(audit.Event{ID: "example_1234567890", OccurredAt: time.Date(
		2026,
		1,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	), Tenant: "tenant1", Actor: audit.Actor{
		Kind: audit.System,
	}, Action: "job.completed", Outcome: audit.Succeeded})
	if err != nil {
		panic(err)
	}
	if _, err := store.Record(ctx, p); err != nil {
		panic(err)
	}
	r, err := store.Reader(audit.ReaderConfig{Schema: schema})
	if err != nil {
		panic(err)
	}
	// Authorize tenant1 before every call, including calls with Page.Next.
	page, err := r.Query(ctx, audit.QueryRequest{Filter: audit.Filter{Tenant: "tenant1"}})
	if err != nil {
		panic(err)
	}
	entry, err := r.Get(ctx, "tenant1", page.Entries[0].ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(entry.Action, entry.Outcome, page.Next == "")
	// Output: job.completed succeeded true
}

// docs:end ExampleReader_Query

// docs:start ExampleReader_Export
func ExampleReader_Export() {
	// Synthetic example key; production keys come from a secret manager.
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
		Events: 10,
		Bytes:  1 << 20,
	}})
	if err != nil {
		panic(err)
	}
	// PostgreSQL deployments use separate reader-role and recorder-role stores.
	r, err := store.Reader(audit.ReaderConfig{Schema: schema, Access: store})
	if err != nil {
		panic(err)
	}
	var output bytes.Buffer
	// Authorize this tenant and the representation before calling Export.
	result, err := r.Export(context.TODO(), audit.ExportRequest{Filter: audit.Filter{
		Tenant: "tenant1",
		Action: "job.",
	}, Actor: audit.Actor{Kind: audit.System}, Reason: "investigation"}, &output)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	var trailer struct {
		Schema    string `json:"schema"`
		Integrity string `json:"integrity"`
	}
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &trailer); err != nil {
		panic(err)
	}
	fmt.Printf(
		"%s: %d events, integrity=%s\n",
		trailer.Schema,
		result.Events,
		trailer.Integrity,
	)
	// Output:
	// pro.audit/v1/export-trailer: 0 events, integrity=unauthenticated
}

// docs:end ExampleReader_Export
