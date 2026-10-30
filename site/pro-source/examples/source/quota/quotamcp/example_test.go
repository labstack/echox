// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotamcp_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotamcp"
)

type exampleTenant struct{}

// Example meters a settled tool: each call reserves 4000 tokens and settles
// the amount the handler charges, and the plan allows two calls a day.
// docs:start Example
func Example() {
	ctx := context.Background()
	lim := quota.MustNew(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
			quota.PerDay("tool_calls", 2),
			quota.PerMonth("tokens", 100_000).Meter("tokens"),
		}}},
	})
	defer func() {
		if err := lim.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	srv := mcp.NewServer(&mcp.Implementation{Name: "example", Version: "v1"}, nil)
	mcp.AddTool(srv, &mcp.Tool{Name: "summarize"}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		quotamcp.Charge(ctx, quota.Cost{"tokens": 1200}) // from the model provider's usage report
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "summary"}}}, nil, nil
	})
	srv.AddReceivingMiddleware(quotamcp.MustMiddleware(quotamcp.MiddlewareConfig{
		Limiter: lim,
		Identity: func(ctx context.Context, _ quotamcp.Caller) (quota.Identity, error) {
			tenant, ok := ctx.Value(exampleTenant{}).(string)
			if !ok {
				return quota.Identity{}, errors.New("unauthenticated")
			}
			return quota.Identity{Tenant: tenant}, nil
		},
		Policies: map[string]quotamcp.Policy{
			"summarize": {Cost: quota.Cost{"tokens": 4000}, Settle: true},
		},
	}))

	st, ct := mcp.NewInMemoryTransports()
	ss, err := srv.Connect(context.WithValue(ctx, exampleTenant{}, "acme"), st, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v1"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := errors.Join(cs.Close(), ss.Close()); err != nil {
			fmt.Println(err)
		}
	}()

	for range 3 {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "summarize"})
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.IsError {
			denial, _ := res.Meta[quotamcp.MetaKey].(map[string]any) //nolint:errcheck // a missing entry prints <nil>
			fmt.Println("denied:", denial["reason"])
			continue
		}
		for _, c := range res.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				fmt.Println(tc.Text)
			}
		}
	}
	u, err := lim.Usage(ctx, quota.UsageRequest{Identity: quota.Identity{Tenant: "acme"}})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, l := range u.Limits {
		fmt.Println(l.Limit, l.Used)
	}
	// Output:
	// summary
	// summary
	// denied: limit_exceeded
	// tool_calls 2
	// tokens 2400
}

// docs:end Example

// Example_echo mounts a stateless MCP endpoint in Echo. The Echo middleware
// authenticates and stores the tenant in the request context, which reaches
// quotamcp in stateless mode.
// docs:start Example_echo
func Example_echo() {
	lim := quota.MustNew(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{quota.PerMinute("tool_calls", 60)}}},
	})
	srv := mcp.NewServer(&mcp.Implementation{Name: "example", Version: "v1"}, nil)
	srv.AddReceivingMiddleware(quotamcp.MustMiddleware(quotamcp.MiddlewareConfig{
		Limiter: lim,
		Identity: func(ctx context.Context, _ quotamcp.Caller) (quota.Identity, error) {
			tenant, ok := ctx.Value(exampleTenant{}).(string)
			if !ok {
				return quota.Identity{}, errors.New("unauthenticated")
			}
			return quota.Identity{Tenant: tenant}, nil
		},
	}))
	authenticate := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			tenant := c.Request().Header.Get("X-Tenant") // replace with real authentication
			if tenant == "" {
				return echo.ErrUnauthorized
			}
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), exampleTenant{}, tenant)))
			return next(c)
		}
	}
	h := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, &mcp.StreamableHTTPOptions{
		Stateless:                    true, // required for 2026-07-28 clients
		PropagateRequestCancellation: true,
	})
	e := echo.New()
	e.Any("/mcp", echo.WrapHandler(h), authenticate)
}

// docs:end Example_echo

// ExampleMiddlewareConfig_Validate rejects a plan set before Reload when it
// would let an attempt-counted tool spend credits.
// docs:start ExampleMiddlewareConfig_Validate
func ExampleMiddlewareConfig_Validate() {
	cfg := quotamcp.MiddlewareConfig{Policies: map[string]quotamcp.Policy{"search": {}}}
	next := []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("tool_calls", 1000),
		quota.Credits("prepaid"), // counts the requests meter
	}}}
	fmt.Println(cfg.Validate(next))
	// Output:
	// quota: Default.Attempt.requests: credit meter charged on admission; put credit meters only in the Cost of a settled tool
}

// docs:end ExampleMiddlewareConfig_Validate

// docs:start charge-first
// RegisterChargeFirstSearch registers a tool and returns its metering policy.
func RegisterChargeFirstSearch(srv *mcp.Server, search func(context.Context, string) (string, int64, error)) quotamcp.Policy {
	mcp.AddTool(srv, &mcp.Tool{Name: "search"}, func(ctx context.Context, _ *mcp.CallToolRequest, args struct {
		Query string `json:"query"`
	}) (*mcp.CallToolResult, any, error) {
		quotamcp.Charge(ctx, nil) // Mark started before calling the provider.
		text, tokens, err := search(ctx, args.Query)
		if err != nil {
			return nil, nil, err
		} // Keep the reservation if usage is unknown.
		quotamcp.Charge(ctx, quota.Cost{"tokens": tokens})
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})
	// Assign this policy to MiddlewareConfig.Policies["search"].
	return quotamcp.Policy{
		Cost: quota.Cost{"tokens": 100}, Settle: true, ChargeFirst: true,
	}
}

// docs:end charge-first

func TestRegisterChargeFirstSearch(t *testing.T) {
	rec := &recorder{}
	l := newLimiter(t, rec, quota.Plan{ID: "team", Limits: []quota.Limit{
		quota.PerDay("tool_calls", 10),
		quota.PerMonth("tokens", 1000).Meter("tokens"),
	}})
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v1"}, nil)
	var calls atomic.Int64
	policy := RegisterChargeFirstSearch(srv, func(_ context.Context, query string) (string, int64, error) {
		calls.Add(1)
		return "results for " + query, 37, nil
	})
	cfg := config(l, rec)
	cfg.Policies = map[string]quotamcp.Policy{"search": policy}
	mw, err := cfg.ToMiddleware()
	if err != nil {
		t.Fatal(err)
	}
	srv.AddReceivingMiddleware(mw)
	cs := connect(t, srv, nil, "")
	res := call(t, cs, "search", map[string]any{"query": "invoices"})
	if res.IsError || text(res) != "results for invoices" || calls.Load() != 1 {
		t.Fatalf("successful search: result=%+v calls=%d", res, calls.Load())
	}
	if got := usage(t, l, "acme", "tokens"); got != 37 {
		t.Fatalf("tokens = %d, want reported usage 37 instead of reservation 100", got)
	}
	if res := call(t, cs, "search", map[string]any{"query": 123}); !res.IsError {
		t.Fatal("invalid arguments accepted")
	}
	if calls.Load() != 1 {
		t.Fatalf("search calls = %d, want invalid arguments rejected before handler", calls.Load())
	}
	if got := usage(t, l, "acme", "tokens"); got != 37 {
		t.Fatalf("tokens = %d, want invalid call's reservation refunded, leaving 37", got)
	}
	if got := usage(t, l, "acme", "tool_calls"); got != 2 {
		t.Fatalf("attempts = %d, want both calls admitted", got)
	}
}
