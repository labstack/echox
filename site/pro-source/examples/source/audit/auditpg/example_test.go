// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auditpg_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"echo.labstack.com/pro/redact"

	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/audit/auditpg"
)

// docs:start ExampleNew
func ExampleNew() {
	// Construction refuses a missing caller-owned pool before any I/O.
	_, err := openAuditStore(context.Background(), nil, nil)
	e, ok := errors.AsType[*audit.ConfigError](err)
	fmt.Println(ok, e.Field)
	// Output: true pool
}

// docs:end ExampleNew

// docs:start ExampleFork
func ExampleFork() {
	// Recovery needs no constructed Store; invalid evidence is rejected before I/O.
	err := auditpg.Fork(
		context.Background(),
		nil,
		auditpg.SchemaConfig{},
		auditpg.ForkRequest{
			Reason: "physical_restore",
		},
	)
	e, ok := errors.AsType[*audit.ValidationError](err)
	fmt.Println(ok, e.Field)
	// Output: true deployment evidence
}

// docs:end ExampleFork

// docs:start ExampleStore_Record
func ExampleStore_Record() {
	// The same intent/outcome recipe accepts auditpg.Store as its Recorder.
	// A bounded MemoryStore makes the failure path executable without a database.
	ctx := context.Background()
	schema, err := audit.NewSchema(audit.SchemaConfig{Policy: redact.StorePolicy{
		AllowMaskFallback: true,
	}})
	if err != nil {
		fmt.Println(err)
		return
	}
	event := audit.Event{
		ID:          "intent_123456789",
		Actor:       audit.Actor{Kind: audit.System},
		Action:      "job.completed",
		Outcome:     audit.Unknown,
		Phase:       audit.Intent,
		OperationID: "operation_123456",
	}
	intent, err := schema.Prepare(event)
	if err != nil {
		fmt.Println(err)
		return
	}
	event.ID = "outcome_123456789"
	event.Phase = audit.OutcomePhase
	event.Outcome = audit.Succeeded
	outcome, err := schema.Prepare(event)
	if err != nil {
		fmt.Println(err)
		return
	}
	store, err := audit.NewMemoryStore(audit.MemoryStoreConfig{Capacity: audit.Capacity{
		Events: 1,
		Bytes:  16384,
	}})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = recordEffect(
		ctx,
		store,
		intent,
		outcome,
		func() error { fmt.Println("effect performed"); return nil },
	)
	e, ok := errors.AsType[*audit.OutcomeUnrecordedError](err)
	fmt.Println(ok, errors.Is(e, audit.ErrCapacity))
	// Output:
	// effect performed
	// true true
}

// docs:end ExampleStore_Record
// docs:start recordEffect
func recordEffect(
	ctx context.Context,
	recorder audit.Recorder,
	intent, outcome audit.PreparedEvent,
	effect func() error,
) error {
	if _, err := recorder.Record(ctx, intent); err != nil {
		return err
	}
	if err := effect(); err != nil {
		return err
	}
	completion, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	if _, err := recorder.Record(completion, outcome); err != nil {
		return &audit.OutcomeUnrecordedError{OperationID: "operation_123456", Cause: err}
	}
	return nil
}

// docs:end recordEffect

// docs:start postgres-pools
// openAuditStore uses pre-provisioned roles and caller-owned pools connected
// to the same database. The deployment pool must be authorized to SET ROLE to
// the trusted NOLOGIN owner. The recorder pool uses the recorder runtime role.
func openAuditStore(ctx context.Context, deploymentPool, recorderPool *pgxpool.Pool) (
	*auditpg.Store,
	error,
) {
	cfg := auditpg.SchemaConfig{
		Schema: "business_audit",
		Roles: auditpg.Roles{
			Deployment:  "audit_deployment",
			Recorder:    "audit_recorder",
			Relay:       "audit_relay",
			Maintenance: "audit_maintenance",
			Reader:      "audit_reader",
		},
		Timeout: 5 * time.Second, MigrateTimeout: time.Minute,
	}
	if err := auditpg.Migrate(ctx, deploymentPool, cfg); err != nil {
		return nil, err
	}
	return auditpg.New(ctx, recorderPool, cfg)
}

// docs:end postgres-pools
