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
	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/redact"
)

// docs:start ExampleSchema_businessEvent
func ExampleSchema_businessEvent() {
	ctx := context.Background()
	material := make([]byte, 32)
	if _, err := rand.Read(material); err != nil {
		panic(err)
	}
	keys, err := redact.NewKeyring([]redact.Key{{ID: "v1", Material: material}})
	if err != nil {
		panic(err)
	}
	store, err := audit.NewMemoryStore(audit.MemoryStoreConfig{
		Capacity: audit.Capacity{Events: 100, Bytes: 1 << 20},
	})
	if err != nil {
		panic(err)
	}
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		panic(err)
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	token, record, err := auth.IssueKey("demo", ring, "alice", auth.KeyOptions{
		Kind: auth.User, TenantRestriction: "acme", Scopes: []string{"credit:record"},
	})
	if err != nil {
		panic(err)
	}
	keyStore := &auth.MemoryKeyStore{Capacity: 10}
	if err := keyStore.Put(ctx, record); err != nil {
		panic(err)
	}
	authMW, err := (auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "operators", Prefix: "demo", Peppers: ring, Store: keyStore,
		}},
		Fallback: auth.Require(map[string][]string{"key": {"credit:record"}}),
	}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	// BEGIN business event
	schema, err := audit.NewSchema(audit.SchemaConfig{
		Keys: keys, Fields: map[string]redact.Class{"amount": redact.Public},
	})
	if err != nil {
		panic(err)
	}
	e.Group("/v1", authMW).POST("/credits/record", func(c *echo.Context) error {
		p, ok := auth.FromEcho(c)
		if !ok {
			return echo.ErrUnauthorized
		}
		if p.Kind() != auth.User || p.TenantRestriction() != "acme" {
			return echo.ErrForbidden
		}
		// Server-confirmed grant to this tenant's customer; no client actor input.
		prepared, err := schema.Prepare(audit.Event{
			ID: "grant_1234567890a", Tenant: p.TenantRestriction(),
			Actor: audit.Actor{Kind: audit.User, ID: redact.Value{
				Class: redact.Identifier,
				Text:  p.ID(),
			}},
			Target: audit.Target{Kind: "customer", ID: redact.Value{
				Class: redact.Identifier,
				Text:  "customer_42",
			}},
			Action: "credit.granted", Outcome: audit.Succeeded,
			Metadata: map[string]redact.Value{"amount": {
				Class: redact.Public,
				Text:  "100",
			}},
		})
		if err != nil {
			return err
		}
		if _, err := store.Record(c.Request().Context(), prepared); err != nil {
			return err
		}
		return c.NoContent(http.StatusCreated)
	})
	// END business event
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/v1/credits/record", nil)
	req.Header.Set("Authorization", "ApiKey "+token)
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)
	fmt.Println(res.Code)
	reader, err := store.Reader(audit.ReaderConfig{Schema: schema})
	if err != nil {
		panic(err)
	}
	entry, err := reader.Get(ctx, "acme", "grant_1234567890a")
	if err != nil {
		panic(err)
	}
	fmt.Println(
		entry.Actor.Kind,
		entry.Actor.ID,
		entry.Target.Kind,
		entry.Target.ID,
		entry.Metadata["amount"],
	)
	// Output:
	// 201
	// user operators:alice customer customer_42 100
}

// docs:end ExampleSchema_businessEvent
