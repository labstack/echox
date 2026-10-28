// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package webhookpg_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"echo.labstack.com/pro/webhook"
	"echo.labstack.com/pro/webhook/webhookpg"
)

// docs:start ExampleNew
func ExampleNew() {
	ctx := context.Background()
	// A lazily connecting pool: New validates configuration before any I/O.
	pool, err := pgxpool.New(
		ctx,
		"postgres://webhook_runtime@127.0.0.1:1/app?sslmode=disable",
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()
	cfg := webhookpg.StoreConfig{Schema: "app", Routing: webhookpg.RoutingConfig{
		Sources: map[string]webhookpg.ProcessFunc{
			"billing": func(context.Context, pgx.Tx, webhook.Incoming) (
				[]webhook.PreparedOutgoing,
				error,
			) {
				return nil, nil
			},
		},
		Destinations: map[string]webhookpg.DestinationConfig{"credits": {
			Version: "v1",
			Send: func(context.Context, webhook.Delivery) (
				webhook.DeliveryResult,
				error,
			) {
				return webhook.DeliveryResult{Outcome: webhook.Accepted}, nil
			},
			// ReplayWindow is required: the receiver's documented dedupe horizon.
		}},
	}}
	_, err = webhookpg.New(ctx, pool, cfg)
	e, ok := errors.AsType[*webhook.ConfigError](err)
	fmt.Println(ok, e.Field)
	d := cfg.Routing.Destinations["credits"]
	d.ReplayWindow = 7 * 24 * time.Hour
	cfg.Routing.Destinations["credits"] = d
	_, err = webhookpg.New(ctx, nil, cfg)
	e, ok = errors.AsType[*webhook.ConfigError](err)
	fmt.Println(ok, e.Field)
	// Output:
	// true Routing.Destinations.ReplayWindow
	// true pool
}

// docs:end ExampleNew

// ExampleStore_Enqueue shows the required caller-transaction policy: READ
// COMMITTED, a synchronous COMMIT, business SQL and the enqueue in one
// transaction, and whole-transaction retries. It needs a database, so it is
// compiled but not run.
// docs:start ExampleStore_Enqueue
func ExampleStore_Enqueue() {
	ctx := context.Background()
	pool, store := applicationPool(), applicationStore()
	p, err := webhook.PrepareOutgoing(webhook.Outgoing{
		Tenant:      "t1",
		Destination: "credits",
		Version:     "v1",
		Key:         "invoice_42/grant",
		Type:        "credit.add",
		Format:      "json",
		Payload: []byte(
			`{"amount":100}`,
		),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Println(err)
		}
	}()
	// The outbox row is only as durable as this COMMIT.
	if _, err = tx.Exec(ctx, "SET LOCAL synchronous_commit = on"); err != nil {
		fmt.Println(err)
		return
	}
	const amount int64 = 100
	tag, err := tx.Exec(
		ctx,
		"INSERT INTO app.entitlements(tenant, invoice, action, amount) "+
			"VALUES ($1, $2, 'grant', $3) ON CONFLICT DO NOTHING",
		"t1",
		"invoice_42",
		amount,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	if tag.RowsAffected() == 0 {
		var recorded int64
		if err = tx.QueryRow(
			ctx,
			"SELECT amount FROM app.entitlements WHERE tenant = $1 AND "+
				"invoice = $2 AND action = 'grant'",
			"t1",
			"invoice_42",
		).Scan(&recorded); err != nil {
			fmt.Println(err)
			return
		}
		if recorded != amount {
			fmt.Println("conflicting grant")
			return
		}
		// The ledger and original command committed together. Outbox and
		// receiver keys can expire, so never enqueue an existing grant again.
	} else {
		// A transient error requires retrying the whole transaction with p.
		if _, err = store.Enqueue(ctx, tx, p); err != nil {
			fmt.Println(err)
			return
		}
	}
	// On an unknown COMMIT outcome, retry the whole transaction; the identical
	// ledger grant prevents another command, even after outbox pruning.
	if err = tx.Commit(ctx); err != nil {
		fmt.Println(err)
	}
}

// docs:end ExampleStore_Enqueue

// docs:start ExampleNewWorker
func ExampleNewWorker() {
	// NewWorker validates its bounds before using the Store and starts nothing.
	_, err := webhookpg.NewWorker(nil, webhookpg.WorkerConfig{Concurrency: 65})
	e, ok := errors.AsType[*webhook.ConfigError](err)
	fmt.Println(ok, e.Field)
	_, err = webhookpg.NewWorker(nil, webhookpg.WorkerConfig{
		Concurrency: 8,
		BatchSize:   4,
	})
	e, ok = errors.AsType[*webhook.ConfigError](err)
	fmt.Println(ok, e.Field)
	// Output:
	// true Concurrency
	// true store
}

// docs:end ExampleNewWorker

// ExampleWorker_Run runs one Worker until the process is interrupted. Run
// stops claiming at once, drains dispatched work for at most ShutdownTimeout
// and returns. It needs a database, so it is compiled but not run.
// docs:start ExampleWorker_Run
func ExampleWorker_Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	w, err := webhookpg.NewWorker(applicationStore(), webhookpg.WorkerConfig{
		Concurrency:     8,
		ShutdownTimeout: 20 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Schedule bounded maintenance separately; Run only claims and processes.
	if err = w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		// A stop error: an unknown claim COMMIT, a database clock regression,
		// missing permissions, a schema mismatch or an unregistered route.
		// Investigate before starting a new Worker.
		fmt.Println(err)
	}
}

// docs:end ExampleWorker_Run

// ExampleProcessFunc grants an invoice's credits once per business key, not
// per provider event ID: two provider events for the same invoice insert one
// ledger row, and the outgoing credit command's key derives from that business
// key. Errors roll back the whole transaction. It is compiled but not run.
// docs:start ExampleProcessFunc
func ExampleProcessFunc() {
	var grant webhookpg.ProcessFunc = grantInvoice
	_ = webhookpg.StoreConfig{Routing: webhookpg.RoutingConfig{
		Sources: map[string]webhookpg.ProcessFunc{"billing": grant},
	}}
}

// docs:end ExampleProcessFunc

// grantInvoice is the basic recipe's business-key grant processor.
// docs:start grantInvoice
func grantInvoice(ctx context.Context, tx pgx.Tx, in webhook.Incoming) (
	[]webhook.PreparedOutgoing,
	error,
) {
	invoice, amount := string(in.Payload), int64(100) // parse the verified payload
	tag, err := tx.Exec(
		ctx,
		`INSERT INTO app.entitlements(tenant, invoice, action, amount)
VALUES ($1, $2, 'grant', $3) ON CONFLICT DO NOTHING`,
		in.Tenant,
		invoice,
		amount,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		var recorded int64
		if err = tx.QueryRow(
			ctx,
			"SELECT amount FROM app.entitlements WHERE tenant = $1 AND "+
				"invoice = $2 AND action = 'grant'",
			in.Tenant,
			invoice,
		).Scan(&recorded); err != nil {
			return nil, err
		}
		if recorded != amount {
			return nil, errors.New("conflicting grant")
		}
		return nil, nil // the ledger and original command committed together
	}
	// Emit only with the first ledger insertion: outbox and receiver keys expire.
	p, err := webhook.PrepareOutgoing(webhook.Outgoing{
		Tenant:      in.Tenant,
		Destination: "credits",
		Version:     "v1",
		Key: strconv.Quote(in.Tenant) + "/" + strconv.Quote(
			invoice,
		) + "/grant",
		Type:    "credit.add",
		Format:  "json",
		Payload: []byte(`{"amount":100}`),
	})
	return []webhook.PreparedOutgoing{p}, err
}

// docs:end grantInvoice

// Application wiring for the documentation example.
type applicationDB interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

var (
	applicationPool  = func() applicationDB { return nil }
	applicationStore = func() *webhookpg.Store { return nil }
)
