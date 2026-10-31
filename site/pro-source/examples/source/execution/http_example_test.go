// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package execution_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/quota"
)

// docs:start exampleHTTPConfig
func exampleHTTPConfig() execution.MiddlewareConfig {
	j, err := execution.New(execution.JournalConfig{
		Namespace: "demo",
		Store: execution.NewMemoryStore(
			nil,
		),
		Policies: map[string]execution.Policy{
			"write": {
				SafetyVersion:     "v1",
				ReconcilerVersion: "proof/v1",
				Effects: []execution.Effect{
					{
						Name:           "write",
						Provider:       "service/account/v1",
						RequestTimeout: time.Second,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return execution.MiddlewareConfig{
		Journal: j, Operation: "write",
		// In production these callbacks use the authenticated principal and the
		// operation's current authorization policy; never a client identity header.
		Identity: func(
			*echo.Context,
		) (string, string, error) {
			return "tenant", "subject", nil
		},
		ReplayCheck: func(*echo.Context) error { return nil },
	}
}

// docs:end exampleHTTPConfig

// docs:start exampleHTTPRequest
func exampleHTTPRequest(ctx context.Context, e *echo.Echo) *httptest.ResponseRecorder {
	r := httptest.NewRequestWithContext(ctx, http.MethodPost, "/write", nil)
	r.Header.Set("Idempotency-Key", `"write-42"`)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	return w
}

// docs:end exampleHTTPRequest

// docs:start ExampleMiddlewareConfig_ToMiddleware
func ExampleMiddlewareConfig_ToMiddleware() {
	cfg := exampleHTTPConfig()
	mw, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	calls := 0
	// docs:start site-anchor
	e.POST("/write", func(c *echo.Context) error {
		calls++
		// A conditional success needs no external send or Intent.
		if err := execution.Classify(c, execution.Completed); err != nil {
			return err
		}
		return c.String(200, "already current")
	}, mw)
	// docs:end site-anchor
	first := exampleHTTPRequest(context.Background(), e)
	second := exampleHTTPRequest(context.Background(), e)
	fmt.Println(first.Code, first.Body.String())
	fmt.Println(second.Code, second.Header().Get("Idempotent-Replayed"), calls)
	// Output:
	// 200 already current
	// 200 true 1
}

// docs:end ExampleMiddlewareConfig_ToMiddleware

// docs:start ExampleMiddlewareConfig_Prepare
func ExampleMiddlewareConfig_Prepare() {
	ctx := context.Background()
	events := make(chan quota.ConsumptionEvent, 8)
	limiter, err := quota.New(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{}),
		Plans: []quota.Plan{{ID: "pro", Limits: []quota.Limit{
			quota.PerMonth("requests", 100),
			quota.PerMonth("tokens", 10000).Meter("tokens"),
		}}},
		Hooks: quota.Hooks{
			OnConsumption: func(
				_ context.Context,
				event quota.ConsumptionEvent,
			) {
				events <- event
			},
		},
	})
	if err != nil {
		panic(err)
	}
	cfg := exampleHTTPConfig()
	cfg.NotExecuted = func(c *echo.Context) {
		// Zero refunds work. Omitted requests keeps this HTTP attempt charged.
		quota.Charge(c, quota.Cost{"tokens": 0})
	}
	prepare, err := cfg.Prepare()
	if err != nil {
		panic(err)
	}
	execute, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}
	admit, err := (quota.MiddlewareConfig{Limiter: limiter, Settle: true,
		Identity: func(
			*echo.Context,
		) (quota.Identity, error) {
			return quota.Identity{Tenant: "tenant"}, nil
		},
		Cost: func(*echo.Context) (quota.Cost, error) {
			return quota.Cost{
				"requests": 1,
				"tokens":   100,
			}, nil
		},
	}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	// Body/auth precede this chain; a validator goes between admit and execute.
	e.POST("/write", func(c *echo.Context) error {
		quota.Charge(c, quota.Cost{"tokens": 60})
		if err := execution.Classify(c, execution.Completed); err != nil {
			return err
		}
		return c.String(200, "done")
	}, prepare, admit, execute)
	first := exampleHTTPRequest(ctx, e)
	second := exampleHTTPRequest(ctx, e)
	// Settlement is asynchronous. Shutdown drains it without a guessed sleep.
	if err := limiter.Close(ctx); err != nil {
		panic(err)
	}
	close(events)
	amounts := quota.Cost{}
	for event := range events {
		amounts[event.Meter] += event.Amount
	}
	fmt.Println(first.Code, second.Code, second.Header().Get("Idempotent-Replayed"))
	fmt.Println("requests:", amounts["requests"], "tokens:", amounts["tokens"])
	// Output:
	// 200 200 true
	// requests: 2 tokens: 60
}

// docs:end ExampleMiddlewareConfig_Prepare

// docs:start ExampleAttemptFrom
func ExampleAttemptFrom() {
	ctx := context.Background()
	cfg := exampleHTTPConfig()
	mw, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.POST("/write", func(c *echo.Context) error {
		attempt := execution.AttemptFrom(c)
		digest := sha256.Sum256([]byte("durable request 42"))
		start := time.Now() // Anchor the send deadline before the call.
		permit, err := attempt.Intent(
			c.Request().Context(),
			"write",
			"payload/write-42",
			digest,
			execution.Token{1},
		)
		if err != nil {
			return err
		} // A refusal never authorizes a send.
		deadline := start.Add(permit.ValidFor)
		if !time.Now().Before(deadline) {
			return &execution.AmbiguityError{}
		}
		// One identical provider request uses permit.Key and this deadline, with
		// SDK retries/redirects disabled. Here the provider succeeds, then the
		// handler is interrupted before the journal records completion.
		panic("simulated interruption after provider success")
	}, mw)
	func() {
		defer func() {
			if recover() == nil {
				panic("expected interruption")
			}
		}()
		exampleHTTPRequest(ctx, e)
	}()
	w := exampleHTTPRequest(ctx, e)
	var problem execution.Problem
	if err := json.Unmarshal(w.Body.Bytes(), &problem); err != nil {
		panic(err)
	}
	fmt.Println(w.Code, problem.Code, problem.SafeNewKey)
	// A trusted operator authorizes this target tenant before accessing the ref,
	// verifies provider success, and reconstructs its HTTP result contract.
	operator := cfg.Journal.Operator()
	record, err := operator.Get(ctx, problem.OperationRef)
	if err != nil {
		panic(err)
	}
	prepared, err := operator.PrepareResolve(ctx, execution.ResolveRequest{
		Ref:              record.Ref,
		ExpectedRevision: record.Revision,
		ID:               "decision-42",
		Decision:         execution.Applied,
		Actor:            "trusted-operator",
		Reason:           "provider success verified",
		Result: &execution.Result{
			Status: 200,
			Codec:  "http/v1",
			Headers: map[string][]string{
				"Content-Type": {"text/plain"},
			},
			Body: []byte("recovered"),
		},
	})
	if err != nil {
		panic(err)
	}
	if _, err := operator.Resolve(ctx, prepared); err != nil {
		panic(err)
	}
	w = exampleHTTPRequest(ctx, e)
	fmt.Println(w.Code, w.Body.String(), w.Header().Get("Idempotent-Replayed"))
	// Output:
	// 409 operation_uncertain false
	// 200 recovered true
}

// docs:end ExampleAttemptFrom

// docs:start ExampleClassify
func ExampleClassify() {
	ctx := context.Background()
	cfg := exampleHTTPConfig()
	mw, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.POST("/write", func(c *echo.Context) error {
		if err := execution.Classify(c, execution.Failed); err != nil {
			return err
		}
		return echo.NewHTTPError(422, "business rule rejected")
	}, mw)
	first := exampleHTTPRequest(ctx, e)
	second := exampleHTTPRequest(ctx, e)
	fmt.Println(first.Code, second.Code, second.Header().Get("Idempotent-Replayed"))
	// Output: 422 422 true
}

// docs:end ExampleClassify
