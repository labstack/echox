// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package execution_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"echo.labstack.com/pro/execution"
)

// docs:start ExampleJournal_Do_transfer
func ExampleJournal_Do_transfer() {
	journal, err := execution.New(execution.JournalConfig{Namespace: "payments", Store: execution.NewMemoryStore(nil), Policies: map[string]execution.Policy{
		"transfer": {SafetyVersion: "v1", ReconcilerVersion: "proof/v1", Effects: []execution.Effect{
			{Name: "debit", Provider: "bank/account/v1", ReplayWindow: time.Hour, RequestTimeout: time.Second},
			{Name: "credit", Provider: "bank/account/v1", ReplayWindow: time.Hour, RequestTimeout: time.Second},
		}},
	}})
	if err != nil {
		fmt.Println(err)
		return
	}
	request := execution.Request{Operation: "transfer", Key: execution.Key{Tenant: "tenant", Subject: "subject", Value: "transfer-42"}, Fingerprint: sha256.Sum256([]byte("transfer 42"))}
	done, err := journal.Do(context.Background(), request, func(ctx context.Context, a *execution.Attempt) (execution.Result, error) {
		for _, name := range []string{"debit", "credit"} {
			// Allocate a distinct token for each send; retain it for lost-reply retries.
			var sendID execution.Token
			if _, err := rand.Read(sendID[:]); err != nil {
				return execution.Result{}, err
			}
			digest := sha256.Sum256([]byte(name + " 42"))
			// Anchor the send deadline before the call, so latency only shortens it.
			start := time.Now()
			permit, err := a.Intent(ctx, name, "payload/transfer-42/"+name, digest, sendID)
			if err != nil {
				return execution.Result{}, err
			}
			// Use permit.Key and the local monotonic deadline for one identical provider
			// request, with SDK retries/redirects disabled. The example omits network I/O.
			deadline := start.Add(permit.ValidFor)
			if !time.Now().Before(deadline) {
				return execution.Result{}, &execution.AmbiguityError{}
			}
			fmt.Println(name, permit.SendID == sendID)
		}
		return execution.Result{Status: 200, Codec: "text", Body: []byte("transferred")}, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(done.Record.State, done.Record.SafeNewKey)
	// Output:
	// debit true
	// credit true
	// completed false
}

// docs:end ExampleJournal_Do_transfer

// docs:start ExampleJournal_Do_conditional
func ExampleJournal_Do_conditional() {
	journal, err := execution.New(execution.JournalConfig{Namespace: "optional", Store: execution.NewMemoryStore(nil), Policies: map[string]execution.Policy{
		"update": {SafetyVersion: "v1", ReconcilerVersion: "proof/v1", Effects: []execution.Effect{{Name: "write", Provider: "service/v1", RequestTimeout: time.Second}}},
	}})
	if err != nil {
		fmt.Println(err)
		return
	}
	done, err := journal.Do(context.Background(), execution.Request{Operation: "update", Key: execution.Key{Tenant: "t", Subject: "s", Value: "k"}}, func(context.Context, *execution.Attempt) (execution.Result, error) {
		// The condition is already satisfied, so no provider send or Intent is needed.
		return execution.Result{Status: 200, Body: []byte("already current")}, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(done.Record.State, done.Record.DuplicateRisk)
	// Output: completed false
}

// docs:end ExampleJournal_Do_conditional

// docs:start ExampleParseKey
func ExampleParseKey() {
	key, err := execution.ParseKey([]string{`"transfer-42"`})
	fmt.Println(key, err)
	// Output: transfer-42 <nil>
}

// docs:end ExampleParseKey

// docs:start exampleOperator
func exampleOperator() (*execution.Operator, string, error) {
	journal, err := execution.New(execution.JournalConfig{Namespace: "example", Store: execution.NewMemoryStore(nil), Policies: map[string]execution.Policy{
		"operation": {SafetyVersion: "v1", ReconcilerVersion: "v1", Effects: []execution.Effect{{Name: "write", Provider: "service/v1", RequestTimeout: time.Second}}},
	}})
	if err != nil {
		return nil, "", err
	}
	done, err := journal.Do(context.Background(), execution.Request{Operation: "operation", Key: execution.Key{Tenant: "tenant", Subject: "subject", Scope: "operation", Value: "example-key"}}, func(context.Context, *execution.Attempt) (execution.Result, error) {
		return execution.Result{Status: 200}, nil
	})
	return journal.Operator(), done.Record.Ref, err
}

// docs:end exampleOperator
// docs:start ExampleOperator_Target
func ExampleOperator_Target() {
	operator, ref, err := exampleOperator()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Authenticate and authorize the returned tenant before reading the record.
	target, err := operator.Target(context.Background(), ref)
	fmt.Println(target.Ref == ref, target.Namespace, target.Tenant, target.Result == nil, err)
	// Output: true example tenant true <nil>
}

// docs:end ExampleOperator_Target
// docs:start ExampleOperator_FindTargetByKey
func ExampleOperator_FindTargetByKey() {
	operator, ref, err := exampleOperator()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Accept the key in a bounded POST body, under trusted tenant authorization.
	target, err := operator.FindTargetByKey(context.Background(), execution.Key{Tenant: "tenant", Subject: "subject", Scope: "operation", Value: "example-key"})
	fmt.Println(target.Ref == ref, len(target.Effects), target.Result == nil, err)
	// Output: true 0 true <nil>
}

// docs:end ExampleOperator_FindTargetByKey
// docs:start ExampleDuplicateRiskAcceptance
func ExampleDuplicateRiskAcceptance() {
	acceptance := execution.DuplicateRiskAcceptance()
	fmt.Println(acceptance.Version, acceptance.Text)
	// Output: v1 accepted or paused sends may still apply
}

// docs:end ExampleDuplicateRiskAcceptance
