// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"echo.labstack.com/pro/webhook"
	"echo.labstack.com/pro/webhook/webhookpg"
)

func seedWebhooks(ctx context.Context, pool *pgxpool.Pool, schema string) (result *webhookpg.Store, resultErr error) {
	var seeding atomic.Bool
	seeding.Store(true)
	cfg := webhookpg.StoreConfig{Schema: schema, Namespace: "commerce-events", Routing: webhookpg.RoutingConfig{
		MaxAttempts: 1,
		Sources: map[string]webhookpg.ProcessFunc{"billing": func(ctx context.Context, _ pgx.Tx, in webhook.Incoming) ([]webhook.PreparedOutgoing, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if seeding.Load() && in.Type == "invoice.invalid" {
				return nil, errors.New("simulated application rejection")
			}
			return nil, nil
		}},
		Destinations: map[string]webhookpg.DestinationConfig{"customer-events": {Version: "v1", ReplayWindow: 7 * 24 * time.Hour, RequestTimeout: time.Second,
			Send: func(ctx context.Context, d webhook.Delivery) (webhook.DeliveryResult, error) {
				if err := ctx.Err(); err != nil {
					return webhook.DeliveryResult{}, err
				}
				outcome := webhook.Accepted
				if seeding.Load() {
					switch d.Type {
					case "delivery.rejected":
						outcome = webhook.Rejected
					case "delivery.unknown":
						outcome = webhook.Unknown
					}
				}
				return webhook.DeliveryResult{Outcome: outcome, Code: "local_demo"}, nil
			},
		}},
	}}
	if err := webhookpg.Migrate(ctx, pool, cfg); err != nil {
		return nil, err
	}
	store, err := webhookpg.New(ctx, pool, cfg)
	if err != nil {
		return nil, err
	}
	for i := range 40 {
		direction, typ := "in", "invoice.paid"
		if i < 16 {
			if i%2 == 0 {
				typ = "invoice.invalid"
			}
		} else {
			direction = "out"
			typ = []string{"delivery.accepted", "delivery.rejected", "delivery.unknown"}[(i-16)%3]
		}
		if _, err = admitWebhook(ctx, pool, store, direction, demoTenants[i%4], fmt.Sprintf("seed-%03d", i), typ); err != nil {
			return nil, err
		}
	}
	worker, err := webhookpg.NewWorker(store, webhookpg.WorkerConfig{Concurrency: 4, ShutdownTimeout: time.Second})
	if err != nil {
		return nil, err
	}
	seedCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	done := make(chan error, 1)
	go func() { done <- worker.Run(seedCtx) }()
	joined := false
	defer func() {
		cancel()
		if !joined {
			if workerErr := <-done; workerErr != nil && !errors.Is(workerErr, seedCtx.Err()) {
				resultErr = errors.Join(resultErr, workerErr)
			}
		}
	}()
	tick := time.NewTicker(25 * time.Millisecond)
	defer tick.Stop()
	for {
		var remaining int
		if err = pool.QueryRow(seedCtx, "SELECT count(*) FROM "+pgx.Identifier{schema}.Sanitize()+".webhook_messages WHERE state IN ('ready','sending')").Scan(&remaining); err != nil {
			return nil, errors.New("webhook seed worker did not complete")
		}
		if remaining == 0 {
			break
		}
		select {
		case <-seedCtx.Done():
			return nil, errors.New("webhook seed worker timed out")
		case <-tick.C:
		}
	}
	// Finish the seed worker before adding review snapshots. Future replay
	// attempts succeed locally, including previously rejected/ambiguous payloads.
	cancel()
	workerErr := <-done
	joined = true
	if workerErr != nil && !errors.Is(workerErr, context.Canceled) {
		return nil, workerErr
	}
	seeding.Store(false)
	// No clock replacement is exported. Admission remains public; the two
	// intermediate-state snapshots below use SQL only in this fresh demo schema.
	// Creation times are untouched. Eligibility/claims are held for 24 hours so
	// the runtime worker can deliver replays without consuming the review rows.
	for _, direction := range []string{"in", "out"} {
		for _, state := range []string{"ready", "sending"} {
			for _, tenant := range demoTenants {
				receipt, admitErr := admitWebhook(ctx, pool, store, direction, tenant, direction+"-"+state+"-"+tenant, "invoice.paid")
				if admitErr != nil {
					return nil, admitErr
				}
				if err = webhookSnapshot(ctx, pool, schema, receipt.Ref, direction, state); err != nil {
					return nil, err
				}
			}
		}
	}
	for _, direction := range []string{"in", "out"} {
		states := []string{"ready", "sending", "done", "dead"}
		if direction == "out" {
			states = append(states, "uncertain")
		}
		for _, state := range states {
			page, inspectErr := store.Inspect(ctx, webhookpg.InspectQuery{Direction: direction, State: state})
			if inspectErr != nil {
				return nil, inspectErr
			}
			want := 8
			if state == "ready" || state == "sending" {
				want = 4
			}
			if len(page.Records) != want {
				return nil, fmt.Errorf("webhook seed %s/%s count = %d, want %d", direction, state, len(page.Records), want)
			}
		}
	}
	return store, nil
}

func admitWebhook(ctx context.Context, pool *pgxpool.Pool, store *webhookpg.Store, direction, tenant, key, typ string) (webhook.Receipt, error) {
	payload := []byte(`{"invoice":"inv_demo","amount":24900,"currency":"usd"}`)
	if direction == "in" {
		p, err := webhook.PrepareIncoming(webhook.Incoming{Tenant: tenant, Source: "billing", Account: "demo-account", EventID: key, Type: typ, Format: "json", Payload: payload})
		if err != nil {
			return webhook.Receipt{}, err
		}
		return store.Accept(ctx, p)
	}
	p, err := webhook.PrepareOutgoing(webhook.Outgoing{Tenant: tenant, Destination: "customer-events", Version: "v1", Key: key, Type: typ, Format: "json", Payload: payload})
	if err != nil {
		return webhook.Receipt{}, err
	}
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return webhook.Receipt{}, errors.New("outgoing demo transaction failed")
	}
	defer rollback(ctx, tx)
	if _, err = tx.Exec(ctx, "SET LOCAL synchronous_commit = on"); err != nil {
		return webhook.Receipt{}, errors.New("outgoing demo durability configuration failed")
	}
	receipt, err := store.Enqueue(ctx, tx, p)
	if err != nil {
		return webhook.Receipt{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return webhook.Receipt{}, errors.New("outgoing demo commit failed")
	}
	return receipt, nil
}

func webhookSnapshot(ctx context.Context, pool *pgxpool.Pool, schema, ref, direction, state string) error {
	var token [16]byte
	if _, err := rand.Read(token[:]); err != nil {
		return err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return errors.New("webhook snapshot transaction failed")
	}
	defer rollback(ctx, tx)
	s := pgx.Identifier{schema}.Sanitize()
	// Raise the fixture attempt budget to leave room after the retained receipt.
	_, err = tx.Exec(ctx, "UPDATE "+s+`.webhook_messages SET state=$2, max_attempts=8, cycle_attempts=1, total_attempts=1,
next_at=statement_timestamp()+interval '24 hours',
attempt_token=CASE WHEN $2='sending' THEN $3::uuid END,
claim_until=CASE WHEN $2='sending' THEN statement_timestamp()+interval '24 hours' END,
first_intent_at=CASE WHEN direction='out' THEN created_at END,
replay_deadline=CASE WHEN direction='out' THEN created_at+interval '7 days' END,
may_have_applied=(direction='out'), code=CASE WHEN $2='ready' THEN 'receiver_throttled' END
WHERE ref=$1::uuid`, ref, state, hex.EncodeToString(token[:]))
	if err != nil {
		return errors.New("webhook intermediate snapshot failed")
	}
	outcome := "retry"
	if direction == "in" {
		outcome = "failed"
	}
	if state == "sending" {
		outcome = "intent"
	}
	_, err = tx.Exec(ctx, "INSERT INTO "+s+`.webhook_attempts(ref,number,generation,token,started_at,finished_at,outcome,code)
SELECT ref,1,0,$2::uuid,created_at,CASE WHEN $3<>'intent' THEN clock_timestamp() END,$3,
CASE WHEN $3<>'intent' THEN 'receiver_throttled' END FROM `+s+".webhook_messages WHERE ref=$1::uuid", ref, hex.EncodeToString(token[:]), outcome)
	if err != nil {
		return errors.New("webhook snapshot receipt failed")
	}
	if err = tx.Commit(ctx); err != nil {
		return errors.New("webhook snapshot commit failed")
	}
	return nil
}
