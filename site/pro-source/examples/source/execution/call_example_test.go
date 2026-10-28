// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package execution_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"echo.labstack.com/pro/execution"
)

// docs:start ExampleAttempt_Call
func ExampleAttempt_Call() {
	ctx := context.Background()
	journal, err := execution.New(execution.JournalConfig{
		Namespace: "example",
		Store: execution.NewMemoryStore(
			nil,
		),
		Policies: map[string]execution.Policy{
			"dispatch": {
				Scope:             "dispatch",
				SafetyVersion:     "v1",
				ReconcilerVersion: "v1",
				Effects: []execution.Effect{
					{
						Name:           "send",
						Provider:       "example/v1",
						RequestTimeout: time.Second,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	attempt, _, err := journal.Begin(ctx, execution.Request{
		Operation: "dispatch",
		Key: execution.Key{
			Tenant:  "example",
			Subject: "worker",
			Scope:   "dispatch",
			Value:   "job-1",
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	var id execution.Token
	if _, err = rand.Read(id[:]); err != nil {
		fmt.Println(err)
		return
	}
	admittedAt := time.Now() // capture before the store call, so its latency counts
	permit, err := attempt.Intent(ctx, "send", "owned-payload-reference", [32]byte{1}, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	timing, err := attempt.Call(ctx, "send", permit, admittedAt, func(
		call context.Context,
	) (string, string, error) {
		// Replace with one provider call using permit.Key and the identical request.
		// This example callback performs no external work and observes cancellation.
		if err := call.Err(); err != nil {
			return "uncertain", "", err
		}
		return "succeeded", "", nil
	})
	if err != nil {
		fmt.Println(err)
		return
	} // retry only timing persistence, never the callback
	fmt.Println(timing.Outcome)
	// Business finalization is separate from diagnostic timing.
	if _, err = attempt.Complete(ctx, execution.Result{Status: 204}); err != nil {
		fmt.Println(err)
	}
	// Output: succeeded
}

// docs:end ExampleAttempt_Call
