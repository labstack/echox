// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotamcp_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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
