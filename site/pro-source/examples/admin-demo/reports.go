// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/openapi"
	"echo.labstack.com/pro/openapi/openapiadmin"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotastripe"
	"echo.labstack.com/pro/quota/quotastripe/stripeadmin"
	"echo.labstack.com/pro/redact"
)

func newAuditSchema(material []byte) (*audit.Schema, error) {
	keys, err := redact.NewKeyring([]redact.Key{{ID: "demo", Material: material}})
	if err != nil {
		return nil, err
	}
	return audit.NewSchema(audit.SchemaConfig{Keys: keys, Fields: map[string]redact.Class{"summary": redact.Public},
		Policy: redact.StorePolicy{Modes: map[redact.Class]redact.Mode{redact.SensitiveIP: redact.Raw}}})
}

func seedAudit(ctx context.Context, store audit.Recorder, schema *audit.Schema) error {
	now := time.Now().UTC()
	for day := range 7 {
		for n := range 12 + day*5 {
			tenant := demoTenants[n%4]
			event, err := schema.Prepare(audit.Event{
				OccurredAt: now.Add(-time.Duration(6-day)*24*time.Hour - time.Duration(n*17+1)*time.Minute), Tenant: tenant,
				Actor:    audit.Actor{Kind: []audit.ActorKind{audit.User, audit.Service, audit.System}[n%3], ID: redact.Value{Class: redact.Identifier, Text: []string{"morgan", "billing-worker", "scheduler"}[n%3]}},
				Action:   []string{"invoice.created", "payment.captured", "shipment.booked", "access.refused"}[n%4],
				Outcome:  []audit.Outcome{audit.Succeeded, audit.Succeeded, audit.Failed, audit.Denied, audit.Unknown}[n%5],
				Target:   audit.Target{Kind: "invoice", ID: redact.Value{Class: redact.Identifier, Text: fmt.Sprintf("inv_%s_%d_%03d", tenant, day, n)}},
				Source:   audit.Source{IP: redact.Value{Class: redact.SensitiveIP, Text: fmt.Sprintf("192.0.2.%d", 10+n%50)}, Route: "/api/invoices", Method: "POST", RequestID: fmt.Sprintf("req-%d-%03d", day, n)},
				Metadata: map[string]redact.Value{"summary": {Class: redact.Public, Text: "Application event from the local commerce demo"}},
			})
			if err != nil {
				return err
			}
			if _, err = store.Record(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}

// offlineTransport is a hard network boundary, even if a future change calls
// an exporter API that normally contacts Stripe. No key or dispatcher is used.
type offlineTransport struct{}

func (offlineTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("stripe network access is disabled in admin-demo")
}

func seedStripe(ctx context.Context) (*stripeadmin.Panel, error) {
	store, err := quotastripe.NewMemoryStore(quotastripe.StoreConfig{})
	if err != nil {
		return nil, err
	}
	exporter, err := quotastripe.New(quotastripe.Config{Namespace: "offline-demo", Store: store, Client: &http.Client{Transport: offlineTransport{}},
		APIKey: func(context.Context) (string, error) { return "", errors.New("no Stripe key in offline demo") },
		Customer: func(ctx context.Context, event quota.ConsumptionEvent) (string, error) {
			return "cus_demo_" + event.Identity.Tenant, ctx.Err()
		},
		Meters: map[string]quotastripe.Meter{"tokens": {ID: "mtr_demo_tokens", EventName: "tokens", CustomerKey: "customer", ValueKey: "value"}},
	})
	if err != nil {
		return nil, err
	}
	for i := range 24 {
		event := quota.ConsumptionEvent{ID: fmt.Sprintf("offline-usage-%03d", i), Identity: quota.Identity{Tenant: demoTenants[i%4]}, Plan: "team", Meter: "tokens", Amount: int64(80 + i*19), At: time.Now().UTC().Add(-time.Duration(i*11+1) * time.Minute).Truncate(time.Millisecond)}
		if err = exporter.Deliver(ctx, quota.Event{ID: event.ID, Consumption: &event}); err != nil {
			return nil, err
		}
	}
	// Status and inspection use the real MemoryStore/Exporter without Run.
	// Reconcile/flush require a Stripe API or live dispatcher and are omitted.
	return stripeadmin.New(stripeadmin.Config{Status: exporter, Inspection: exporter})
}

func mountSamples(ctx context.Context, e *echo.Echo, lim *quota.Limiter, panel *openapiadmin.Panel, records []demoExecution) error {
	for _, route := range []struct{ method, path string }{{"GET", "/api/orders"}, {"POST", "/api/orders"}, {"GET", "/api/orders/:id"}, {"GET", "/api/v1/invoices"}, {"POST", "/api/refunds"}, {"GET", "/healthz"}} {
		if _, err := e.AddRoute(echo.Route{Method: route.method, Path: route.path, Handler: func(c *echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "local demo"})
		}}); err != nil {
			return err
		}
	}
	type recordLink struct {
		demoExecution
		URL string
	}
	links := make([]recordLink, 0, len(records))
	for _, record := range records {
		prefix := "/admin"
		if record.Namespace == "fulfillment" {
			prefix = "/admin-fulfillment"
		}
		links = append(links, recordLink{record, prefix + "/execution/record?" + url.Values{"tenant": {record.Tenant}, "ref": {record.Ref}}.Encode()})
	}
	index, err := template.New("demo").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><title>Echo Pro console demo</title>
<h1>Echo Pro console demo</h1><p><a href="/admin">Open admin console · payments</a> · <a href="/admin-fulfillment">Fulfillment namespace</a></p>
<p>Signed in as Demo Operator, all tenants, reveal and writes enabled. Throwaway data; no external services are called.</p>
<p>Quota: use tenant acme, northstar, paperkite or orbit; plan team or enterprise; credit limit credits; meters requests and tokens.</p>
<p><a href="/admin/quota/usage?tenant=acme&amp;plan=team">Acme usage</a> · <a href="/admin/quota/top?meter=requests">Top consumers</a> · <a href="/admin/quota/denials">Denials</a> · <a href="/admin/audit/events">Audit events</a> · <a href="/admin/webhook/outbox">Webhook outbox</a></p>
<p>Audit events span seven days. Execution and webhook creation times are current: their PostgreSQL stores export no clock hooks. Scheduled-ready and sending webhook review snapshots use demo-only SQL, held for 24 hours. Replays run through a real worker with local callbacks. Stripe shows local backlog and decisions; network actions are disabled.</p>
<p>OpenAPI is a checked startup snapshot with deliberate uncovered routes and a sunset advisory; it is not an installed enforcement chain.</p>
<h2>Execution seed directory</h2><p>Use these links for started, completed and failed records; uncertain records are also linked in the console. Lookup subject: checkout-service; scope and key appear below. A started record has a one-hour lease.</p>
<table><thead><tr><th>Namespace</th><th>Tenant</th><th>Scope</th><th>Key</th><th>Seed state</th></tr></thead><tbody>
{{range .}}<tr><td>{{.Namespace}}</td><td>{{.Tenant}}</td><td>{{.Scope}}</td><td><a href="{{.URL}}">{{.Key}}</a></td><td>{{.State}}</td></tr>{{end}}</tbody></table></html>`)
	if err != nil {
		return err
	}
	if _, err = e.AddRoute(echo.Route{Method: "GET", Path: "/", Handler: func(c *echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Response().Header().Set("Cache-Control", "no-store")
		return index.Execute(c.Response(), links)
	}}); err != nil {
		return err
	}
	routes := map[string]openapi.Policy{"GET /healthz": {Exempt: "local health probe", Public: true}}
	for _, prefix := range []string{"/admin", "/admin-fulfillment"} {
		for _, route := range admin.RouteKeys(prefix) {
			routes[route] = openapi.Policy{Exempt: "console owns authentication and CSRF", Public: true}
		}
	}
	// Check deliberately retains missing policies/routes and response findings.
	// BindChecked labels the report honestly without installing broken coverage.
	report, err := (openapi.MiddlewareConfig{Spec: []byte(sampleSpec), Quota: quota.MiddlewareConfig{Limiter: lim,
		Identity: func(*echo.Context) (quota.Identity, error) { return quota.Identity{Tenant: "acme"}, nil },
		Plan:     func(*echo.Context, quota.Identity) (string, error) { return "team", nil }},
		AuthMode: openapi.AuthBySpec, Routes: routes, Auth: func(next echo.HandlerFunc) echo.HandlerFunc { return next },
	}).Check(ctx, e)
	if err != nil {
		return err
	}
	if report.Incomplete || len(report.Operations) == 0 || len(report.Findings) == 0 || len(report.Advisories) == 0 {
		return errors.New("sample OpenAPI report is missing expected review content")
	}
	return panel.BindChecked(report)
}

const sampleSpec = `{
 "openapi":"3.1.2", "info":{"title":"Commerce demo","version":"1.0"},
 "paths":{
  "/api/orders":{"get":{"operationId":"listOrders","security":[],"x-labstack-quota":{"cost":{"requests":1}},"responses":{"200":{"description":"Orders"}}},
                  "post":{"operationId":"createOrder","security":[],"x-labstack-quota":{"cost":{"requests":5,"tokens":20}},"responses":{"201":{"description":"Created"}}}},
  "/api/orders/{id}":{"get":{"operationId":"getOrder","security":[],"responses":{"200":{"description":"Order"}}}},
  "/api/v1/invoices":{"get":{"operationId":"legacyInvoices","security":[],"deprecated":true,
    "x-labstack-deprecation":{"at":"2025-01-01T00:00:00Z","sunset":"2025-06-01T00:00:00Z","link":"https://example.com/migrate/invoices"},
    "x-labstack-quota":{},"responses":{"default":{"content":{"application/problem+json":{}}}}}},
  "/api/reports":{"get":{"operationId":"exportReport","security":[],"x-labstack-quota":{"cost":{"requests":10}},"responses":{"200":{"description":"Not mounted yet"}}}}
 }
}`
