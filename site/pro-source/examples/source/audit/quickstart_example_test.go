// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package audit_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/redact"
)

func ExampleMemoryStore_quickstart() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	keys, err := redact.NewKeyring([]redact.Key{{ID: "v1", Material: key}})
	if err != nil {
		panic(err)
	}
	schema, err := audit.NewSchema(audit.SchemaConfig{Keys: keys})
	if err != nil {
		panic(err)
	}
	store, err := audit.NewMemoryStore(audit.MemoryStoreConfig{
		Capacity: audit.Capacity{Events: 100, Bytes: 1 << 20},
	})
	if err != nil {
		panic(err)
	}
	reader, err := store.Reader(audit.ReaderConfig{Schema: schema})
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.POST("/jobs/complete", func(c *echo.Context) error {
		// This demo represents a trusted system job, not a client-supplied actor.
		prepared, err := schema.Prepare(audit.Event{
			ID: "event_1234567890a", Tenant: "acme",
			Actor:  audit.Actor{Kind: audit.System},
			Action: "job.completed", Outcome: audit.Succeeded,
		})
		if err != nil {
			return err
		}
		receipt, err := store.Record(c.Request().Context(), prepared)
		if err != nil {
			return err
		}
		return c.String(http.StatusCreated, receipt.ID)
	})
	ctx := context.Background()
	res := httptest.NewRecorder()
	e.ServeHTTP(res, httptest.NewRequestWithContext(ctx, http.MethodPost, "/jobs/complete", nil))
	fmt.Println(res.Code, res.Body.String())
	// The application must authorize this exact tenant before reading.
	entry, err := reader.Get(ctx, "acme", "event_1234567890a")
	if err != nil {
		panic(err)
	}
	fmt.Println(entry.Action, entry.Outcome)
	// Output:
	// 201 event_1234567890a
	// job.completed succeeded
}
