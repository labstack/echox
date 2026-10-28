// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package pgstore_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/execution/pgstore"
)

// ExampleTxRunner_RunTx commits an order row with its execution binding and
// stored result in one transaction. A replay returns the stored result
// without running the callback; NotApplied rolls the order back but still
// records a binding whose SafeNewKey is true.
//
// This example is compile-only: it needs a reachable PostgreSQL database
// whose "execution" schema was already created with Migrate, plus the
// application's own items(id, stock) and orders(item) tables.
// docs:start ExampleTxRunner_RunTx
func ExampleTxRunner_RunTx() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres:///orders?default_query_exec_mode=exec")
	if err != nil {
		return
	}
	defer pool.Close()
	store, err := pgstore.New(ctx, pool, pgstore.StoreConfig{
		Schema:    "execution",
		Namespace: "orders",
	})
	if err != nil {
		return
	}
	runner, err := pgstore.NewTxRunner(store, pgstore.TxRunnerConfig{
		Namespace: "orders",
		Policies: map[string]execution.Policy{
			"create-order": {Mode: execution.TransactionalMode, SafetyVersion: "v1"},
		},
	})
	if err != nil {
		return
	}
	key := execution.Key{Tenant: "tenant-1", Subject: "user-1", Value: "client-key"}
	done, err := runner.RunTx(ctx, pgstore.TxRequest{
		Operation:   "create-order",
		Key:         key,
		Fingerprint: [32]byte{1},
	},
		func(ctx context.Context, tx *pgstore.Tx) (pgstore.TxOutcome, error) {
			var stock int
			if err := tx.QueryRow(
				ctx,
				"SELECT stock FROM items WHERE id=$1 FOR UPDATE",
				"item-1",
			).Scan(&stock); err != nil {
				return pgstore.TxOutcome{}, err
			}
			if stock < 1 {
				// No protected work remains. The record says no effect applied.
				return pgstore.TxOutcome{
					Decision: execution.NotApplied,
					Result: &execution.Result{
						Status: 409,
						Body: []byte(
							`{"error":"out_of_stock"}`,
						),
					},
				}, nil
			}
			if _, err := tx.Exec(
				ctx,
				"UPDATE items SET stock=stock-1 WHERE id=$1",
				"item-1",
			); err != nil {
				return pgstore.TxOutcome{}, err
			}
			if _, err := tx.Exec(
				ctx,
				"INSERT INTO orders(item) VALUES ($1)",
				"item-1",
			); err != nil {
				return pgstore.TxOutcome{}, err
			}
			return pgstore.TxOutcome{
				Decision: execution.Applied,
				Result: &execution.Result{
					Status: 201,
					Body:   []byte(`{"created":true}`),
				},
			}, nil
		})
	if unknown, ok := errors.AsType[*pgstore.TxCommitUnknownError](
		err,
	); ok && unknown != nil {
		return // Retry with the identical request; never with a new key.
	}
	if err != nil {
		return
	}
	fmt.Println(done.Answer.Status, done.Answer.Replayed)
}

// docs:end ExampleTxRunner_RunTx

// Neither a rollback nor an unknown COMMIT authorizes a new key. An unknown
// COMMIT is unavailable: retry the identical request to replay or execute.
// docs:start ExampleTxCommitUnknownError
func ExampleTxCommitUnknownError() {
	var unknown error = &pgstore.TxCommitUnknownError{}
	var rolledBack error = &pgstore.TxError{
		Phase:             pgstore.TxPhaseCallback,
		Cause:             context.Canceled,
		RollbackConfirmed: true,
	}
	fmt.Println(errors.Is(unknown, execution.ErrUnavailable), unknown)
	fmt.Println(errors.Is(rolledBack, context.Canceled), rolledBack)
	// Output:
	// true execution: transactional commit outcome unknown
	// true execution: transactional callback failed
}

// docs:end ExampleTxCommitUnknownError

// Configuration is validated before any I/O, and errors are typed for
// inspection with errors.Is and errors.AsType. New validates MigrateTimeout
// even though only Migrate uses it.
// docs:start ExampleNewTxRunner
func ExampleNewTxRunner() {
	_, err := pgstore.NewTxRunner(nil, pgstore.TxRunnerConfig{Namespace: "orders"})
	if e, ok := errors.AsType[*execution.ConfigError](err); ok {
		fmt.Println("runner:", e.Field)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(
		ctx,
		"postgres://localhost:1/unused",
	) // Connects lazily; never dialed here.
	if err != nil {
		return
	}
	defer pool.Close()
	_, err = pgstore.New(ctx, pool, pgstore.StoreConfig{MigrateTimeout: 2 * time.Hour})
	if e, ok := errors.AsType[*execution.ConfigError](err); ok {
		fmt.Println("store:", e.Field)
	}
	var rolledBack error = &pgstore.TxError{
		Phase: pgstore.TxPhaseFinalize,
		Cause: execution.ErrUnavailable,
	}
	if e, ok := errors.AsType[*pgstore.TxError](rolledBack); ok {
		fmt.Println(errors.Is(
			rolledBack,
			&pgstore.TxError{},
		), e.Phase, errors.Is(rolledBack, execution.ErrUnavailable))
	}
	// Output:
	// runner: store
	// store: MigrateTimeout
	// true finalize true
}

// docs:end ExampleNewTxRunner
