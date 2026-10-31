// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package webhookpg_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/audit/auditpg"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/redact"
	"echo.labstack.com/pro/webhook"
	"echo.labstack.com/pro/webhook/webhookpg"
)

// creditCommand is a frozen quota.AddCredits command. Everything the sender
// needs travels in the payload, including the plan resolved when the command
// was enqueued, so a later plan change cannot route a retry to another
// balance or make it ambiguous.
type creditCommand struct {
	ID         string            `json:"id"`
	Tenant     string            `json:"tenant"`
	Plan       string            `json:"plan"`
	Limit      string            `json:"limit"`
	Dimensions map[string]string `json:"dimensions,omitempty"`
	Amount     int64             `json:"amount"`
}

// prepareCredit freezes a credit grant for the outbox with its resolved plan.
// The outgoing key and the quota credit ID both derive from the business key
// (tenant, invoice, action) with unambiguous quoting, never from a provider
// event ID or attempt number. encoding/json sorts map keys, so the payload is
// deterministic.
func prepareCredit(tenant, plan, invoice string, amount int64) (
	webhook.PreparedOutgoing,
	error,
) {
	key := strconv.Quote(tenant) + "/" + strconv.Quote(invoice) + "/grant"
	payload, err := json.Marshal(creditCommand{
		ID:     key,
		Tenant: tenant,
		Plan:   plan,
		Limit:  "credits",
		Amount: amount,
	})
	if err != nil {
		return webhook.PreparedOutgoing{}, err
	}
	return webhook.PrepareOutgoing(webhook.Outgoing{
		Tenant:      tenant,
		Destination: "quota-credits",
		Version:     "v1",
		Key:         key,
		Type:        "credit.add",
		Format:      "json",
		Payload:     payload,
	})
}

// creditSender applies a frozen command with the same quota credit ID on every
// attempt, always against plans, the quota.Snapshot of the immutable plan
// configuration that the destination version names (capture it when the
// sender is built, from a Limiter loaded with that configuration, also after a
// restart). Quota deduplicates the ID for 30 days per balance, so a retry after
// a lost reply credits once, even across a Reload. A different amount under the
// ID and invalid input are permanent rejections; any other error is unknown and
// retried with the same command while the destination's ReplayWindow allows.
func creditSender(l *quota.Limiter, plans quota.Snapshot) webhook.SendFunc {
	return func(ctx context.Context, d webhook.Delivery) (webhook.DeliveryResult, error) {
		c, ok := decodeCredit(d.Payload)
		if !ok {
			return webhook.DeliveryResult{
				Outcome: webhook.Rejected,
				Code:    "malformed_command",
			}, nil
		}
		_, err := l.AddCredits(ctx, quota.CreditRequest{
			Snapshot: plans,
			Identity: quota.Identity{
				Tenant: c.Tenant,
			},
			Plan:       c.Plan,
			Limit:      c.Limit,
			Dimensions: c.Dimensions,
			Amount:     c.Amount,
			ID:         c.ID,
		})
		switch code := creditCode(err); code {
		case "":
			// Unknown result: retry the same ID inside the window.
			return webhook.DeliveryResult{}, err
		case "credited":
			return webhook.DeliveryResult{Outcome: webhook.Accepted, Code: code}, nil
		default:
			return webhook.DeliveryResult{Outcome: webhook.Rejected, Code: code}, nil
		}
	}
}

// decodeCredit reports whether payload is a well-formed frozen command.
func decodeCredit(payload []byte) (creditCommand, bool) {
	var c creditCommand
	err := json.Unmarshal(payload, &c)
	return c, err == nil && c.Plan != ""
}

// creditCode classifies an AddCredits result: applied (including an idempotent
// replay of the same ID and amount), a permanent rejection, or "" when the
// outcome is unknown.
func creditCode(err error) string {
	if err == nil {
		return "credited"
	}
	if errors.Is(err, quota.ErrCreditConflict) {
		return "credit_conflict"
	}
	if e, ok := errors.AsType[*quota.RequestError](err); ok && e != nil {
		return "credit_invalid"
	}
	return ""
}

// creditPlans is the immutable quota configuration that destination version
// "v1" names.
func creditPlans() []quota.Plan {
	return []quota.Plan{{ID: "pro", Limits: []quota.Limit{
		quota.Credits("credits").Meter("credits"),
	}}}
}

// Example_quotaCredits projects an entitlement into quota credits eventually:
// the inbox ProcessFunc commits its ledger row and a frozen credit command
// (with its resolved plan) together, and the worker delivers the command with
// creditSender, bound to the snapshot of the plan configuration that the
// destination version names. The destination's resend horizon (7 days) stays
// well inside quota's 30-day credit deduplication; after it, an unresolved
// command becomes uncertain and needs reconciliation. This example drives the
// sender directly: the first reply is lost, a Reload adds a plan, and the
// retry of the unchanged command still credits once.
func Example_quotaCredits() {
	ctx := context.Background()
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: creditPlans()})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := l.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()
	destination := webhookpg.DestinationConfig{
		Version:      "v1",
		Send:         creditSender(l, l.Snapshot()),
		ReplayWindow: 7 * 24 * time.Hour,
	}
	p, err := prepareCredit("t1", "pro", "inv_42", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	out := p.Outgoing()
	d := webhook.Delivery{
		Destination: out.Destination,
		Version:     out.Version,
		Key:         out.Key,
		Type:        out.Type,
		Format:      out.Format,
		Payload:     out.Payload,
		Attempt:     1,
	}
	r, err := destination.Send(ctx, d) // applied, but the reply is lost
	fmt.Println(r.Outcome, r.Code, err)
	// A valid Reload adds a second plan before the retry.
	team := quota.Plan{ID: "team", Limits: []quota.Limit{
		quota.Credits("credits").Meter("credits"),
	}}
	if err = l.Reload(append(creditPlans(), team)); err != nil {
		fmt.Println(err)
		return
	}
	d.Attempt = 2
	r, err = destination.Send(ctx, d) // the unchanged command, against the bound snapshot
	fmt.Println(r.Outcome, r.Code, err)
	// The same credit ID with another amount is a permanent conflict.
	changed, err := prepareCredit("t1", "pro", "inv_42", 150)
	if err != nil {
		fmt.Println(err)
		return
	}
	d.Payload = changed.Outgoing().Payload
	r, err = destination.Send(ctx, d)
	fmt.Println(r.Outcome, r.Code, err)
	res, err := l.AddCredits(ctx, quota.CreditRequest{
		Identity: quota.Identity{Tenant: "t1"},
		Plan:     "pro",
		Limit:    "credits",
		Amount:   100,
		ID:       out.Key,
	})
	fmt.Println("balance", res.Balance, "applied", res.Applied, err)
	// Output:
	// accepted credited <nil>
	// accepted credited <nil>
	// rejected credit_conflict <nil>
	// balance 100 applied false <nil>
}

// auditWriter is the part of *auditpg.Store the audit recipe uses.
type auditWriter interface {
	Write(context.Context, pgx.Tx, audit.PreparedEvent) (audit.Receipt, error)
}

// errGrantConflict reports a second provider event naming an existing
// business grant with a different amount.
var errGrantConflict = errors.New("conflicting grant for the business key")

// auditedGrant is the audit recipe's ProcessFunc. The business key (tenant,
// invoice, action) makes the grant idempotent across distinct provider
// events: a new grant inserts its ledger row, writes its audit event with
// auditpg.Write in the inbox transaction and emits the frozen credit command;
// a repeat with the identical amount succeeds without a second grant or audit
// event or credit command, even after outbox and receiver retention; a repeat with a
// different amount is an error, so the event retries and then dies instead of
// being marked done over the conflict. SQL errors are returned unchanged.
func auditedGrant(schema *audit.Schema, w auditWriter) webhookpg.ProcessFunc {
	return func(ctx context.Context, tx pgx.Tx, in webhook.Incoming) (
		[]webhook.PreparedOutgoing,
		error,
	) {
		var paid struct {
			Invoice string    `json:"invoice"`
			Plan    string    `json:"plan"`
			At      time.Time `json:"paid_at"`
			Amount  int64     `json:"amount"`
		}
		if err := json.Unmarshal(in.Payload, &paid); err != nil {
			return nil, err
		}
		credit, err := prepareCredit(in.Tenant, paid.Plan, paid.Invoice, paid.Amount)
		if err != nil {
			return nil, err
		}
		tag, err := tx.Exec(
			ctx,
			`INSERT INTO app.entitlements(tenant, invoice, action, amount)
VALUES ($1, $2, 'grant', $3) ON CONFLICT DO NOTHING`,
			in.Tenant,
			paid.Invoice,
			paid.Amount,
		)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			var recorded int64
			if err := tx.QueryRow(
				ctx,
				"SELECT amount FROM app.entitlements WHERE tenant = $1 AND "+
					"invoice = $2 AND action = 'grant'",
				in.Tenant,
				paid.Invoice,
			).Scan(&recorded); err != nil {
				return nil, err
			}
			if recorded != paid.Amount {
				return nil, errGrantConflict
			}
			return nil, nil // the ledger and original command committed together
		}
		id := sha256.Sum256([]byte(
			strconv.Quote(
				in.Source,
			) + "/" + strconv.Quote(in.Account) + "/" + strconv.Quote(in.EventID),
		))
		event, err := schema.Prepare(audit.Event{ID: "wh_" + hex.EncodeToString(
			id[:16],
		), OccurredAt: paid.At, Tenant: in.Tenant,
			Actor: audit.Actor{
				Kind: audit.System,
			}, Action: "entitlement.granted", Outcome: audit.Succeeded,
			Target: audit.Target{Kind: "invoice", ID: redact.Value{
				Class: redact.Identifier,
				Text:  paid.Invoice,
			}}})
		if err != nil {
			return nil, err
		}
		if _, err = w.Write(ctx, tx, event); err != nil {
			return nil, err
		}
		return []webhook.PreparedOutgoing{credit}, nil
	}
}

// ExampleProcessFunc_audit records an audit event inside the inbox
// transaction with auditpg.Write, so it commits with the business change, the
// outputs and the inbox completion, or rolls back with them. The audit ID and
// time derive from the verified event, so every attempt prepares the identical
// event; a conflicting repeat of the business grant is refused (see
// auditedGrant). Use the same database and commit policy as the webhook store.
// It needs a database, so it is compiled but not run.
func ExampleProcessFunc_audit() {
	schema, err := audit.NewSchema(audit.SchemaConfig{Policy: redact.StorePolicy{
		AllowMaskFallback: true,
	}})
	if err != nil {
		fmt.Println(err)
		return
	}
	var store auditWriter = applicationAuditStore()
	process := auditedGrant(schema, store)
	_ = webhookpg.StoreConfig{Routing: webhookpg.RoutingConfig{
		Sources: map[string]webhookpg.ProcessFunc{"billing": process},
	}}
}

// ExampleStore_Replay replays dead inbox records after the cause was fixed.
// The actor is projected once with a keyed HMAC and stored with the request,
// so a retry after an unknown outcome is identical even after key rotation. It
// needs a database, so it is compiled but not run.
func ExampleStore_Replay() {
	ctx := context.Background()
	store := applicationStore()
	keys, err := redact.NewKeyring([]redact.Key{{ID: "ops1", Material: operatorKey()}})
	if err != nil {
		fmt.Println(err)
		return
	}
	projector, err := redact.NewStoreProjector(redact.StorePolicy{}, keys)
	if err != nil {
		fmt.Println(err)
		return
	}
	actor, err := projector.Store(redact.Value{
		Class: redact.SensitiveDigest8,
		Text:  "operator@example.com",
	}, "", "webhook.replay")
	if err != nil {
		fmt.Println(err)
		return
	}
	page, err := store.Inspect(ctx, webhookpg.InspectQuery{
		Direction: "in",
		State:     "dead",
		Limit:     10,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range page.Records {
		req := webhookpg.ReplayRequest{
			Ref:        r.Ref,
			Generation: r.Generation,
			RequestID:  "ticket-981/" + r.Ref,
			Actor:      actor,
			Reason:     "callback_fixed",
		}
		res, err := store.Replay(ctx, req)
		if u, unknown := errors.AsType[*webhook.CommitUnknownError](
			err,
		); unknown && u != nil {
			res, err = store.Replay(
				ctx,
				req,
			) // identical retry returns the committed result
		}
		switch {
		case err == nil:
			fmt.Println("replayed generation", res.Generation)
		case errors.Is(err, webhook.ErrFenced):
			fmt.Println("changed since inspection; inspect again")
		case errors.Is(err, webhook.ErrReplayUnsafe):
			fmt.Println("not replayable; reconcile instead")
		default:
			fmt.Println(err)
		}
	}
}

// ExampleStore_Maintain runs bounded retention maintenance once a minute until
// ctx ends, draining backlog while MaintainResult.More is set. It needs a
// database, so it is compiled but not run.
func ExampleStore_Maintain() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	store := applicationStore()
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	for {
		for range 10 { // at most 1,000 records per minute
			res, err := store.Maintain(ctx, 100)
			if err != nil || !res.More {
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
	}
}

// ExampleStore_Diagnose reads the separately costed global diagnostics for
// monitoring. It needs a database, so it is compiled but not run.
func ExampleStore_Diagnose() {
	d, err := applicationStore().Diagnose(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	// The due and expired counts share one bounded scan: when it was capped they
	// are lower bounds and a zero means "not observed", never "absent".
	if d.EligibleCapped {
		fmt.Println(
			"more than 10,000 records (the diagnostics cap) are due or expired: backlog",
		)
	} else if d.OldestDue > 5*time.Minute && d.ExpiredClaims == 0 {
		fmt.Println("due work is ageing and nothing is claiming it: is a Worker running?")
	}
	if d.ExpiredClaims > 0 { // a positive observation holds even when capped
		fmt.Println("claims expired: workers crashed, paused or exceeded their claims")
	}
	fmt.Println("pressure", d.Pressure.Messages, "of", d.Pressure.MaxMessages)
}

var (
	applicationAuditStore = func() *auditpg.Store { return nil }
	operatorKey           = func() []byte { return make([]byte, 32) }
)

// docs:start http-outbox
// BEGIN production
// runHTTPOutbox uses caller-owned pools with deployment/runtime roles. Apply
// grants.sql after Migrate and before invoking this function with runtimePool.
// The application ledger is deployed separately, with a permanent UNIQUE key:
// CREATE TABLE app.entitlements (
//
//	tenant text NOT NULL, invoice text NOT NULL, action text NOT NULL,
//	amount bigint NOT NULL, PRIMARY KEY (tenant, invoice, action));
func runHTTPOutbox(
	ctx context.Context,
	ownerPool, runtimePool *pgxpool.Pool,
	endpoint webhook.Endpoint,
) error {
	send, err := (webhook.HTTPSenderConfig{Endpoints: map[string]webhook.Endpoint{
		"credits": endpoint,
	}}).ToSendFunc()
	if err != nil {
		return err
	}
	cfg := webhookpg.StoreConfig{Schema: "webhooks", Namespace: "billing",
		Routing: webhookpg.RoutingConfig{
			Destinations: map[string]webhookpg.DestinationConfig{
				"credits": {Version: "v1", Send: send, ReplayWindow: 7 * 24 * time.Hour},
			},
		},
	}
	if err := webhookpg.Migrate(ctx, ownerPool, cfg); err != nil {
		return err
	}
	store, err := webhookpg.New(ctx, runtimePool, cfg)
	if err != nil {
		return err
	}
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
		return err
	}
	if err := enqueueHTTPGrant(ctx, runtimePool, store, p); err != nil {
		return err
	}
	worker, err := webhookpg.NewWorker(store, webhookpg.WorkerConfig{Concurrency: 8})
	if err != nil {
		return err
	}
	// Run blocks; the application cancels ctx on shutdown and waits for Run.
	// Investigate non-cancellation stop errors before starting a new worker.
	err = worker.Run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func enqueueHTTPGrant(
	ctx context.Context,
	pool *pgxpool.Pool,
	store *webhookpg.Store,
	p webhook.PreparedOutgoing,
) (
	err error,
) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		if rollbackErr := tx.Rollback(cleanup); rollbackErr != nil && !errors.Is(
			rollbackErr,
			pgx.ErrTxClosed,
		) {
			err = errors.Join(err, rollbackErr)
		}
	}()
	if _, err = tx.Exec(ctx, "SET LOCAL synchronous_commit = on"); err != nil {
		return err
	}
	tag, err := tx.Exec(
		ctx,
		`INSERT INTO app.entitlements(tenant, invoice, action, amount)
VALUES ($1, $2, 'grant', $3) ON CONFLICT DO NOTHING`,
		"t1",
		"invoice_42",
		int64(100),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		var amount int64
		if err := tx.QueryRow(
			ctx,
			`SELECT amount FROM app.entitlements
WHERE tenant=$1 AND invoice=$2 AND action='grant'`,
			"t1",
			"invoice_42",
		).Scan(&amount); err != nil {
			return err
		}
		if amount != 100 {
			return errors.New("conflicting grant for the business key")
		}
		// Original ledger and command committed together. Never re-enqueue an
		// existing grant, even after outbox/receiver deduplication expires.
	} else {
		receipt, err := store.Enqueue(ctx, tx, p)
		if err != nil {
			return err
		}
		_ = receipt // Provisional until this transaction commits.
	}
	// On a transient failure or unknown COMMIT outcome, retry this ENTIRE
	// function with the same prepared p, never just Enqueue. The permanent
	// business-key ledger handles even a lost reply followed by outbox pruning.
	return tx.Commit(ctx)
}

// END production
// docs:end http-outbox

// Example_transactionalHTTP is compiled, but requires deployment/runtime pools
// and a receiver that deduplicates keys for at least the configured seven days.
func Example_transactionalHTTP() {
	_ = runHTTPOutbox
}
