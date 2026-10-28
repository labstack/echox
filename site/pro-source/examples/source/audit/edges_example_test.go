// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package audit_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/redact"
)

// docs:start ExampleDenialConfig_ToMiddleware
func ExampleDenialConfig_ToMiddleware() {
	schema, err := audit.NewSchema(audit.SchemaConfig{Policy: redact.StorePolicy{
		AllowMaskFallback: true,
	}})
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
	denials, err := (audit.DenialConfig{Schema: schema, Recorder: store}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	// Denial recording is outside this group's authorization middleware.
	group := e.Group("/private", denials, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(*echo.Context) error { return echo.NewHTTPError(403, "Forbidden") }
	})
	group.GET("/jobs", func(c *echo.Context) error { return c.NoContent(200) })
	ctx := context.TODO()
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequestWithContext(ctx, "GET", "/private/jobs", nil))
	events, err := store.TestEvents(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(w.Code, len(events))
	// Output: 403 1
}

// docs:end ExampleDenialConfig_ToMiddleware

// docs:start ExampleSlogExporter
func ExampleSlogExporter() {
	schema, err := audit.NewSchema(audit.SchemaConfig{Policy: redact.StorePolicy{
		AllowMaskFallback: true,
	}})
	if err != nil {
		panic(err)
	}
	p, err := schema.Prepare(audit.Event{
		Actor:   audit.Actor{Kind: audit.System},
		Action:  "job.completed",
		Outcome: audit.Succeeded,
	})
	if err != nil {
		panic(err)
	}
	exporter := &audit.SlogExporter{Logger: slog.New(
		slog.NewJSONHandler(io.Discard, nil),
	)}
	err = exporter.Export(context.TODO(), []audit.Line{{Delivery: audit.Delivery{
		Lineage:  "example",
		Position: audit.Position{XID: 1, Seq: 1},
		Event:    p,
	}}})
	fmt.Println(err, exporter.Representation())
	// Output: <nil> masked
}

// docs:end ExampleSlogExporter
