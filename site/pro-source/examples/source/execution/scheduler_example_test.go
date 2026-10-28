// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package execution_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"echo.labstack.com/pro/execution"
)

// docs:start ExampleNewReconcileScheduler
func ExampleNewReconcileScheduler() {
	now := time.Now()
	store := execution.NewMemoryStore(
		func() time.Time { return now },
	) // Tests inject a clock; production uses the default.
	p := execution.Policy{
		SafetyVersion:     "v1",
		Scope:             "write",
		ReconcilerVersion: "proof/v1",
		Effects: []execution.Effect{
			{
				Name:           "write",
				Provider:       "service/account/v1",
				RequestTimeout: time.Second,
			},
		},
	}
	j, err := execution.New(execution.JournalConfig{
		Namespace: "example",
		Store:     store,
		Policies:  map[string]execution.Policy{"write": p},
	})
	if err != nil {
		panic(err)
	}
	a, _, err := j.Begin(context.Background(), execution.Request{
		Operation: "write",
		Key: execution.Key{
			Tenant:  "tenant",
			Subject: "principal",
			Value:   "key",
		},
	})
	if err != nil {
		panic(err)
	}
	if _, err = a.Intent(
		context.Background(),
		"write",
		"durable/payload",
		[32]byte{1},
		execution.Token{1},
	); err != nil {
		panic(err)
	}
	r, err := a.Uncertain(
		context.Background(),
		execution.Evidence{},
	) // The provider reply was lost.
	if err != nil {
		panic(err)
	}
	now = now.Add(6 * time.Minute) // Past the send barrier: LastSendBy + 5 minutes.

	ctx, cancel := context.WithCancel(
		context.Background(),
	) // Application-owned lifecycle.
	defer cancel()
	// Retain one Operator: prepared handles, including a pending one returned by
	// Run, can only be retried through the Operator instance that prepared them.
	operator := j.Operator()
	passes := 0
	scheduler, err := execution.NewReconcileScheduler(execution.ReconcileSchedulerConfig{
		Operator: operator,
		Scopes: map[string]execution.ScheduledReconciler{"write": {
			Reconciler: func(
				_ context.Context,
				_ *execution.ReconcileClaim,
				_ execution.EffectRecord,
			) (
				execution.Verdict,
				error,
			) {
				// Probe using the stored request key. Never rerun business work.
				return execution.Verdict{
					Outcome: execution.EffectApplied,
					Proof: execution.Proof{
						Kind:          execution.ProviderSuccess,
						Version:       "proof/v1",
						Reference:     "redacted/success/1",
						Authoritative: true,
					},
				}, nil
			},
			Completion: func(context.Context, execution.Record, []execution.Verdict) (
				bool,
				*execution.Result,
				error,
			) {
				return true, nil, nil
			},
		}},
		PollInterval: time.Second,
		AuthorizeNamespace: func(context.Context, string) error {
			if passes++; passes > 1 {
				cancel() // Stop after one pass, as an application shutdown would.
			}
			return nil
		},
		AuthorizeTarget: func(_ context.Context, target execution.ReconcileTarget) (
			bool,
			error,
		) {
			return target.Tenant == "tenant", nil // Current per-target authorization.
		},
	})
	if err != nil {
		panic(err)
	}
	err = scheduler.Run(ctx) // Blocks; run it on an application-owned goroutine.
	if stopped, ok := errors.AsType[*execution.SchedulerError](err); ok {
		if handle, held := stopped.Pending(); held {
			// Reauthorize handle.Ref, then retry that handle.
			if _, retryErr := operator.Reconcile(
				context.Background(),
				handle.Prepared,
			); retryErr != nil {
				panic(retryErr)
			}
		}
	}
	fmt.Println(errors.Is(err, context.Canceled))
	got, err := operator.Get(context.Background(), r.Ref)
	if err != nil {
		panic(err)
	}
	fmt.Println(got.State, got.Proof, got.ReconcileAttempts, got.Receipts[0].ID[:4])
	// Output:
	// true
	// completed proven 1 sch_
}

// docs:end ExampleNewReconcileScheduler
