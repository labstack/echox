// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package executiontest_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/execution/executiontest"
)

// A store implementation's conformance test runs CheckStore for the base
// Store contract, and CheckScheduler only for a store advertising the optional
// scheduled-reconciliation capability. The factory receives each case's
// *testing.T, so per-case fixtures can register t.Cleanup.
// docs:start ExampleCheckStore
func ExampleCheckStore() {
	newStore := func(t *testing.T, now func() time.Time) execution.Store {
		t.Helper()

		store := executiontest.NewModel(now)
		t.Cleanup(func() {}) // Release per-case fixtures here.
		return store
	}
	// Call from func TestStoreConformance(t *testing.T).
	conformance := func(t *testing.T) {
		t.Helper()

		executiontest.CheckStore(t, newStore)
		executiontest.CheckScheduler(t, newStore)
	}
	_ = conformance

	var store execution.Store = executiontest.NewModel(nil)
	scheduled, ok := store.(execution.ScheduledReconcileStore)
	fmt.Println("run CheckScheduler:", ok && scheduled.SupportsScheduledReconcile())
	// Output: run CheckScheduler: true
}

// docs:end ExampleCheckStore

// The reference harness runs transactional scripts against the Model. Pass a
// factory returning a TxHarness to CheckTransactional for a real backend.
// docs:start ExampleNewModelHarness
func ExampleNewModelHarness() {
	ctx := context.Background()
	now := time.Unix(1700000000, 0).UTC()
	h := executiontest.NewModelHarness(func() time.Time { return now }, execution.Limits{})
	script := executiontest.TxScript{
		Key:         execution.Key{Tenant: "tenant", Subject: "subject", Scope: executiontest.TxScope, Value: "order-1"},
		Fingerprint: [32]byte{1},
		Writes:      1,
		Decision:    execution.Applied,
		Result:      &execution.Result{Status: 201, Body: []byte(`{"id":1}`)},
	}
	done, invoked, err := h.RunTx(ctx, script)
	fmt.Println(invoked, err, done.Record.State, done.Answer.Replayed)
	replay, invoked, err := h.RunTx(ctx, script)
	fmt.Println(invoked, err, replay.Answer.Replayed, string(replay.Result.Body))
	writes, err := h.Ledger(ctx, "order-1")
	fmt.Println("committed writes:", writes, err)
	// Output:
	// true <nil> completed false
	// false <nil> true {"id":1}
	// committed writes: 1 <nil>
}

// docs:end ExampleNewModelHarness

// The fake provider deduplicates committed permits and refuses a send without one.
// docs:start ExampleProvider
func ExampleProvider() {
	ctx := context.Background()
	p := executiontest.NewProvider(nil)
	permit := execution.SendPermit{Key: execution.EffectKey("ref", "charge"), ValidFor: time.Second}
	applied, err := p.Send(ctx, permit, [32]byte{1})
	fmt.Println(applied, err)
	_, err = p.Send(ctx, execution.SendPermit{}, [32]byte{1})
	fmt.Println(errors.Is(err, &execution.ConfigError{}))
	// Output:
	// true <nil>
	// true
}

// docs:end ExampleProvider
