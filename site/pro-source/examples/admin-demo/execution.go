// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/execution/executionadmin"
)

var demoTenants = []string{"acme", "northstar", "paperkite", "orbit"}

type demoExecution struct {
	Namespace, Tenant, Subject, Scope, Key, Ref, State string
}

func seedExecutions(ctx context.Context, store execution.Store, namespace string) (*executionadmin.Panel, []demoExecution, error) {
	policies := map[string]execution.Policy{}
	for _, operation := range []string{"invoice.issue", "payment.capture", "shipment.book"} {
		policies[operation] = execution.Policy{SafetyVersion: "demo-v1", ReconcilerVersion: "demo-v1", Lease: time.Hour,
			Effects: []execution.Effect{{Name: "dispatch", Provider: "local-ledger/v1", RequestTimeout: time.Second, ReplayWindow: 7 * 24 * time.Hour}}}
	}
	journal, err := execution.New(execution.JournalConfig{Namespace: namespace, Store: store, Policies: policies})
	if err != nil {
		return nil, nil, err
	}
	// The local provider ledger records acceptance independently of the journal.
	// A probe checks that ledger; it never fabricates a verdict for an unknown send.
	accepted := make(map[[32]byte]bool)
	var rows []demoExecution
	for i := range 48 {
		operation := []string{"invoice.issue", "payment.capture", "shipment.book"}[i%3]
		tenant := demoTenants[i%len(demoTenants)]
		key := fmt.Sprintf("%s-%03d", namespace, i+1)
		a, answer, beginErr := journal.Begin(ctx, execution.Request{Operation: operation, Key: execution.Key{Tenant: tenant, Subject: "checkout-service", Value: key}, Fingerprint: sha256.Sum256([]byte(key))})
		if beginErr != nil {
			return nil, nil, beginErr
		}
		if a == nil {
			return nil, nil, errors.New("demo execution was not admitted")
		}
		record := answer.Record
		// 8 started, 24 completed, 8 failed before dispatch and 8 uncertain.
		switch i % 6 {
		case 0:
			// No dispatch: this remains a real started record for its one-hour lease.
		case 4:
			record, err = a.Fail(ctx, execution.Result{Status: 422, Codec: "json", Body: []byte(`{"code":"inventory_unavailable"}`)})
		default:
			var send execution.Token
			if _, err = rand.Read(send[:]); err != nil {
				return nil, nil, err
			}
			start := time.Now()
			permit, intentErr := a.Intent(ctx, "dispatch", "demo-payload/"+key, sha256.Sum256([]byte(key)), send)
			if intentErr != nil {
				return nil, nil, intentErr
			}
			if permit.ValidFor <= 0 || !time.Now().Before(start.Add(permit.ValidFor)) {
				return nil, nil, errors.New("local provider send permit expired")
			}
			accepted[sha256.Sum256([]byte(key))] = true
			if i%6 == 5 {
				record, err = a.Uncertain(ctx, execution.Evidence{Text: "Local provider accepted dispatch; simulated reply lost before journal completion"})
			} else {
				record, err = a.Complete(ctx, execution.Result{Status: 201, Codec: "json", Body: []byte(`{"status":"accepted"}`)})
			}
		}
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, demoExecution{namespace, tenant, "checkout-service", operation, key, record.Ref, string(record.State)})
	}
	panel, err := executionadmin.New(executionadmin.Config{Operator: journal.Operator(), EnableResolve: true,
		Reconciler: func(ctx context.Context, _ *execution.ReconcileClaim, effect execution.EffectRecord) (execution.Verdict, error) {
			if err := ctx.Err(); err != nil {
				return execution.Verdict{}, err
			}
			// The local ledger is indexed by the frozen request digest.
			if !accepted[effect.Digest] {
				return execution.Verdict{Effect: effect.Effect.Name, Outcome: execution.Unknown, Evidence: "No local acceptance receipt"}, nil
			}
			return execution.Verdict{Effect: effect.Effect.Name, Outcome: execution.EffectApplied,
				Proof: execution.Proof{Effect: effect.Effect.Name, Kind: execution.ProviderSuccess, Version: "demo-v1", Reference: "local-ledger/accepted", Authoritative: true, AllSends: true, Quiescent: true}, Evidence: "Local provider acceptance receipt found"}, nil
		},
		Completion: func(ctx context.Context, _ execution.Record, verdicts []execution.Verdict) (bool, *execution.Result, error) {
			if err := ctx.Err(); err != nil {
				return false, nil, err
			}
			for _, verdict := range verdicts {
				if verdict.Outcome != execution.EffectApplied {
					return false, nil, nil
				}
			}
			return true, &execution.Result{Status: 201, Codec: "json", Body: []byte(`{"status":"reconciled"}`)}, nil
		},
	})
	return panel, rows, err
}
