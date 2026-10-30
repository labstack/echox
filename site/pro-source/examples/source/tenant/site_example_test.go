// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package tenant_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/auth/authbind"
	"echo.labstack.com/pro/problem/echoadapter"
	"echo.labstack.com/pro/tenant"
	"echo.labstack.com/pro/tenant/tenanttest"
)

// docs:start ExampleAuthorizeConfig_ToMiddleware
func ExampleAuthorizeConfig_ToMiddleware() {
	ctx := context.Background()
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		panic(err)
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	token, record, err := auth.IssueKey("demo", ring, "invoice-reader", auth.KeyOptions{
		Kind: auth.Service, TenantRestriction: "acme", Scopes: []string{"invoice:read"},
	})
	if err != nil {
		panic(err)
	}
	store := &auth.MemoryKeyStore{Capacity: 1}
	if err := store.Put(ctx, record); err != nil {
		panic(err)
	}
	path, err := tenant.FromPathParam("tenant")
	if err != nil {
		panic(err)
	}
	guard, err := authbind.Middleware(auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "invoices",
			Prefix:    "demo",
			Peppers:   ring,
			Store:     store,
			Header:    "X-API-Key",
		}},
		Fallback: auth.Require(map[string][]string{"key": {"invoice:read"}}),
	}, tenant.Config{Selectors: []tenant.Selector{
		tenant.FromCredential(),
		path,
	}, Required: true})
	if err != nil {
		panic(err)
	}
	policy, err := tenant.Static([]tenant.Grant{{
		Scope:        "invoice:read",
		ResourceType: "invoice",
		Actions:      []string{"read"},
	}})
	if err != nil {
		panic(err)
	}
	// This immutable map stands in for a repository lookup made with the
	// resolver's c.Request().Context(). Ownership comes from stored data.
	owners := map[string]tenant.ID{"inv-1": "acme", "inv-2": "globex"}
	check, err := (tenant.AuthorizeConfig{Authorizer: policy, Action: "read",
		Resource: func(c *echo.Context) (tenant.Resource, error) {
			id := c.Param("invoice")
			return tenant.Resource{Type: "invoice", ID: id, Tenant: owners[id]}, nil
		},
	}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.HTTPErrorHandler, err = (echoadapter.ErrorHandlerConfig{}).ToErrorHandler(
		e.HTTPErrorHandler,
	)
	if err != nil {
		panic(err)
	}
	e.Group(
		"/tenants/:tenant",
		guard,
	).GET("/invoices/:invoice", func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	}, check)
	for _, path := range []string{
		"/tenants/acme/invoices/inv-1",
		"/tenants/acme/invoices/inv-2",
		"/tenants/globex/invoices/inv-2",
	} {
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		r.Header.Set("X-API-Key", token)
		w := httptest.NewRecorder()
		e.ServeHTTP(w, r)
		fmt.Println(w.Code)
	}
	// Output:
	// 204
	// 403
	// 403
}

// docs:end ExampleAuthorizeConfig_ToMiddleware

// ExampleCerbos uses a local protocol fixture, not a running policy engine.
// docs:start ExampleCerbos
func ExampleCerbos() {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		if err := json.NewEncoder(
			w,
		).Encode(map[string]any{"results": []any{map[string]any{
			"resource": map[string]string{"id": "inv-1", "kind": "invoice"},
			"actions":  map[string]string{"read": "EFFECT_ALLOW"},
		}}}); err != nil {
			panic(err)
		}
	}))
	defer server.Close()
	policy, err := tenant.Cerbos(tenant.CerbosConfig{
		HTTP: tenant.HTTPConfig{
			URL:           server.URL,
			Client:        server.Client(),
			MaxConcurrent: 8,
		},
		Principal: func(tenant.Check) (tenant.CerbosPrincipal, error) {
			// Map roles from trusted application data, never arbitrary claims.
			return tenant.CerbosPrincipal{Roles: []string{"invoice-reader"}}, nil
		},
	})
	if err != nil {
		panic(err)
	}
	principal, err := auth.NewPrincipal(auth.PrincipalConfig{
		Namespace: "employees",
		Subject:   "alice",
		Scheme:    "oauth",
		Kind:      auth.User,
	})
	if err != nil {
		panic(err)
	}
	decision, err := policy.Authorize(context.Background(), tenant.Check{
		Principal: principal,
		Tenant:    "acme",
		Action:    "read",
		Resource:  tenant.Resource{Type: "invoice", ID: "inv-1", Tenant: "acme"},
	})
	fmt.Println(decision == tenant.Allow, err)
	// Output: true <nil>
}

// docs:end ExampleCerbos

// ExampleOpenFGA uses a local protocol fixture with a pinned model identifier.
// docs:start ExampleOpenFGA
func ExampleOpenFGA() {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		if err := json.NewEncoder(w).Encode(
			map[string]bool{"allowed": true},
		); err != nil {
			panic(err)
		}
	}))
	defer server.Close()
	policy, err := tenant.OpenFGA(tenant.OpenFGAConfig{
		HTTP:    tenant.HTTPConfig{URL: server.URL, Client: server.Client()},
		StoreID: "invoice-store", AuthorizationModelID: "reviewed-model-v1",
		ActionRelations: map[string]string{"read": "reader"},
	})
	if err != nil {
		panic(err)
	}
	principal, err := auth.NewPrincipal(auth.PrincipalConfig{
		Namespace: "employees",
		Subject:   "alice",
		Scheme:    "oauth",
		Kind:      auth.User,
	})
	if err != nil {
		panic(err)
	}
	decision, err := policy.Authorize(context.Background(), tenant.Check{
		Principal: principal,
		Tenant:    "acme",
		Action:    "read",
		Resource:  tenant.Resource{Type: "invoice", ID: "inv-1", Tenant: "acme"},
	})
	fmt.Println(decision == tenant.Allow, err)
	// Output: true <nil>
}

// docs:end ExampleOpenFGA

// docs:start ExampleConfig_selectors
func ExampleConfig_selectors() {
	ctx := context.Background()
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		panic(err)
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	token, record, err := auth.IssueKey(
		"demo",
		ring,
		"worker",
		auth.KeyOptions{Kind: auth.Service},
	)
	if err != nil {
		panic(err)
	}
	keys := &auth.MemoryKeyStore{Capacity: 1}
	if err := keys.Put(ctx, record); err != nil {
		panic(err)
	}
	path, err := tenant.FromPathParam("tenant")
	if err != nil {
		panic(err)
	}
	header, err := tenant.FromHeader("X-Tenant")
	if err != nil {
		panic(err)
	}
	host, err := tenant.FromSubdomain("api.example")
	if err != nil {
		panic(err)
	}
	guard, err := authbind.Middleware(auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "services",
			Prefix:    "demo",
			Peppers:   ring,
			Store:     keys,
		}},
		Fallback: auth.Require(map[string][]string{"key": nil}),
	}, tenant.Config{
		Selectors: []tenant.Selector{path, header, host}, Required: true,
		Member: func(ctx context.Context, p auth.Principal, id tenant.ID) (bool, error) {
			// Replace this immutable fixture with a context-aware membership read.
			if err := ctx.Err(); err != nil {
				return false, err
			}
			return p.ID() == "services:worker" && id == "acme", nil
		},
	})
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.GET("/tenants/:tenant", func(c *echo.Context) error {
		return c.NoContent(
			http.StatusNoContent,
		)
	}, guard)
	for _, tc := range []struct{ path, header, host string }{
		{"acme", "acme", "acme"}, {"acme", "other", "acme"},
		{"acme", "bad space", "acme"}, {"other", "other", "other"},
	} {
		r := httptest.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://"+tc.host+".api.example/tenants/"+tc.path,
			nil,
		)
		r.Header.Set("Authorization", "ApiKey "+token)
		r.Header.Set("X-Tenant", tc.header)
		w := httptest.NewRecorder()
		e.ServeHTTP(w, r)
		fmt.Println(w.Code)
	}
	// Output:
	// 204
	// 403
	// 400
	// 403
}

// docs:end ExampleConfig_selectors

// docs:start ExampleConfig_anonymous
func ExampleConfig_anonymous() {
	// docs:start site-anchor
	selector, err := tenant.FromPathParam("tenant")
	if err != nil {
		panic(err)
	}
	middleware, err := (tenant.Config{
		Required:          true,
		AnonymousSelector: selector,
	}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.GET("/public/:tenant", func(c *echo.Context) error {
		id, _ := tenant.From(c.Request().Context())
		return c.String(http.StatusOK, string(id))
	}, middleware)
	// docs:end site-anchor
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/public/acme",
		nil,
	))
	fmt.Println(w.Code, w.Body.String())
	// Output: 200 acme
}

// docs:end ExampleConfig_anonymous

// docs:start ExampleCasbin
// exampleEnforcer supplies a local model fixture with Casbin's Enforce shape.
// A deployed application supplies its own concurrency-safe Casbin enforcer.
type exampleEnforcer func(...any) (bool, error)

func (f exampleEnforcer) Enforce(args ...any) (bool, error) { return f(args...) }

func ExampleCasbin() {
	principal, err := auth.NewPrincipal(auth.PrincipalConfig{
		Namespace: "employees", Subject: "alice", Scheme: "oauth", Kind: auth.User,
	})
	if err != nil {
		panic(err)
	}
	check := tenant.Check{Principal: principal, Tenant: "acme", Action: "read",
		Resource: tenant.Resource{Type: "invoice", ID: "inv-1", Tenant: "acme"}}
	// This model expects subject, domain, object and action in that order.
	model := exampleEnforcer(func(args ...any) (bool, error) {
		return len(args) == 4 && args[0] == "employees:alice" &&
			args[1] == "acme" && args[2] == "invoice:inv-1" && args[3] == "read", nil
	})
	policy, err := tenant.Casbin(model, nil)
	if err != nil {
		panic(err)
	}
	decision, err := policy.Authorize(context.Background(), check)
	fmt.Println(decision == tenant.Allow, err)
	// A different model accepts subject, a tenant-qualified object and action.
	model = exampleEnforcer(func(args ...any) (bool, error) {
		return len(args) == 3 && args[0] == "employees:alice" &&
			args[1] == "invoice:acme/inv-1" && args[2] == "read", nil
	})
	policy, err = tenant.Casbin(model, func(c tenant.Check) ([]any, error) {
		return []any{
			c.Principal.ID(),
			c.Resource.Type + ":" + string(c.Tenant) + "/" + c.Resource.ID,
			c.Action,
		}, nil
	})
	if err != nil {
		panic(err)
	}
	decision, err = policy.Authorize(context.Background(), check)
	fmt.Println(decision == tenant.Allow, err)
	// Output:
	// true <nil>
	// true <nil>
}

// docs:end ExampleCasbin

// docs:start isolation-test
func TestSiteTenantIsolation(t *testing.T) {
	principal, err := auth.NewPrincipal(auth.PrincipalConfig{
		Namespace: "employees", Subject: "alice", Scheme: "oauth", Kind: auth.User,
		Scopes: []string{"invoice:read"},
	})
	if err != nil {
		t.Fatal(err)
	}
	policy, err := tenant.Static([]tenant.Grant{{
		Scope:        "invoice:read",
		ResourceType: "invoice",
		Actions:      []string{"read"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	own := tenant.Check{Principal: principal, Tenant: "acme", Action: "read",
		Resource: tenant.Resource{Type: "invoice", ID: "inv-1", Tenant: "acme"}}
	foreign := own
	foreign.Resource.Tenant = "globex"
	tenanttest.Authorizer(t, policy, []tenanttest.Case{
		{Name: "own invoice", Check: own, Want: tenant.Allow},
		{Name: "foreign invoice", Check: foreign, Want: tenant.Deny},
	})
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		t.Fatal(err)
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	credential, record, err := auth.IssueKey("test", ring, "alice", auth.KeyOptions{
		Kind: auth.User, TenantRestriction: "acme", Scopes: []string{"invoice:read"},
	})
	if err != nil {
		t.Fatal(err)
	}
	keys := &auth.MemoryKeyStore{Capacity: 1}
	if err := keys.Put(t.Context(), record); err != nil {
		t.Fatal(err)
	}
	guard, err := authbind.Middleware(auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "employees", Prefix: "test", Peppers: ring, Store: keys,
		}},
		Fallback: auth.Require(map[string][]string{"key": {"invoice:read"}}),
	}, tenant.Config{Selectors: []tenant.Selector{
		tenant.FromCredential(),
	}, Required: true})
	if err != nil {
		t.Fatal(err)
	}
	// The middleware must refuse ownership before policy or business work.
	resourceCalls, policyCalls, handlerCalls := 0, 0, 0
	check, err := (tenant.AuthorizeConfig{
		Authorizer: tenant.AuthorizerFunc(func(context.Context, tenant.Check) (
			tenant.Decision,
			error,
		) {
			policyCalls++
			return tenant.Allow, nil
		}),
		Action: "read",
		Resource: func(*echo.Context) (tenant.Resource, error) {
			resourceCalls++
			return foreign.Resource, nil
		},
	}).ToMiddleware()
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.GET("/invoices/inv-1", func(c *echo.Context) error {
		handlerCalls++
		return c.NoContent(http.StatusNoContent)
	}, guard, check)
	r := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"/invoices/inv-1",
		nil,
	)
	r.Header.Set("Authorization", "ApiKey "+credential)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden || resourceCalls != 1 {
		t.Fatalf(
			"status=%d resource lookups=%d; want 403 and one lookup",
			w.Code,
			resourceCalls,
		)
	}
	tenanttest.Refused(t, w.Code, policyCalls, handlerCalls)
}

// docs:end isolation-test

// docs:start missing-resource
// MountInvoice checks existence within the selected tenant before authorization.
func MountInvoice(e *echo.Echo, exists func(context.Context, tenant.ID, string) (bool, error), authorize echo.MiddlewareFunc, handler echo.HandlerFunc) {
	checkExists := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := c.Request().Context()
			id, ok := tenant.From(ctx)
			if !ok {
				return echo.ErrForbidden
			}
			found, err := exists(ctx, id, c.Param("id"))
			if err != nil {
				return err
			}
			if !found {
				return echo.ErrNotFound
			}
			return next(c)
		}
	}
	// Install authentication and tenant selection before this route.
	e.GET("/invoices/:id", handler, checkExists, authorize)
}

// docs:end missing-resource

func TestMountInvoice(t *testing.T) {
	for _, tc := range []struct {
		name  string
		found bool
		want  int
	}{
		{"missing", false, http.StatusNotFound},
		{"existing", true, http.StatusNoContent},
	} {
		t.Run(tc.name, func(t *testing.T) {
			authenticate, token := identity(t, "acme")
			selectTenant, err := (tenant.Config{
				Selectors: []tenant.Selector{tenant.FromCredential()}, Required: true,
			}).ToMiddleware()
			if err != nil {
				t.Fatal(err)
			}
			e := echo.New()
			e.Use(authenticate, selectTenant)
			existsCalls, authorizeCalls, handlerCalls := 0, 0, 0
			exists := func(ctx context.Context, id tenant.ID, invoice string) (bool, error) {
				existsCalls++
				if _, ok := auth.From(ctx); !ok || id != "acme" || invoice != "inv-1" {
					t.Fatalf("lookup without authenticated tenant/invoice: tenant=%q invoice=%q", id, invoice)
				}
				return tc.found, nil
			}
			authorize := func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c *echo.Context) error {
					authorizeCalls++
					return next(c)
				}
			}
			MountInvoice(e, exists, authorize, func(c *echo.Context) error {
				handlerCalls++
				return c.NoContent(http.StatusNoContent)
			})
			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/invoices/inv-1", nil)
			r.Header.Set("Authorization", "ApiKey "+token)
			w := httptest.NewRecorder()
			e.ServeHTTP(w, r)
			if w.Code != tc.want || existsCalls != 1 {
				t.Fatalf("status=%d lookups=%d; want %d and one lookup", w.Code, existsCalls, tc.want)
			}
			if !tc.found {
				tenanttest.Refused(t, w.Code, authorizeCalls, handlerCalls)
			} else if authorizeCalls != 1 || handlerCalls != 1 {
				t.Fatalf("authorize=%d handler=%d; want one each", authorizeCalls, handlerCalls)
			}
		})
	}
}
