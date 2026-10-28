// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

// Command admin-demo runs the full console against disposable PostgreSQL schemas.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/audit"
	"echo.labstack.com/pro/audit/auditadmin"
	"echo.labstack.com/pro/audit/auditpg"
	"echo.labstack.com/pro/execution/pgstore"
	"echo.labstack.com/pro/openapi/openapiadmin"
	"echo.labstack.com/pro/webhook/webhookadmin"
)

func main() {
	// Do not put the URL in flag defaults: -help must never disclose credentials.
	var databaseURL string
	flag.StringVar(
		&databaseURL,
		"postgres-url",
		"",
		"disposable PostgreSQL URL (or ADMIN_DEMO_POSTGRES_URL)",
	)
	addr := flag.String("addr", "127.0.0.1:8080", "loopback listen address")
	flag.Parse()
	if databaseURL == "" {
		databaseURL = os.Getenv("ADMIN_DEMO_POSTGRES_URL")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, databaseURL, *addr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, databaseURL, addr string) error {
	if databaseURL == "" {
		return errors.New(
			"provide -postgres-url or ADMIN_DEMO_POSTGRES_URL; use the " +
				"assigned development database with schema creation privileges",
		)
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
		return errors.New("-addr must use a loopback IP and port")
	}
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
	if err != nil {
		return errors.New("cannot bind demo listen address")
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(
			closeErr,
			net.ErrClosed,
		) {
			fmt.Fprintln(os.Stderr, "demo listener close failed")
		}
	}()
	authority := listener.Addr().String()
	credential := os.Getenv("ADMIN_DEMO_API_KEY")
	authentication, err := demoAuthentication(ctx, credential, authority)
	if err != nil {
		return err
	}
	if credential == "" {
		fmt.Fprintln(
			os.Stderr,
			"browser demo mode: requests without Authorization use an "+
				"ephemeral local key",
		)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return errors.New("invalid PostgreSQL configuration (URL withheld)")
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		return errors.New(
			"cannot connect to PostgreSQL (URL and driver diagnostics withheld)",
		)
	}
	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		return err
	}
	prefix := "admin_demo_" + hex.EncodeToString(key[:6])
	var canCreateRoles bool
	if err = pool.QueryRow(
		ctx,
		"SELECT rolcreaterole OR rolsuper FROM pg_roles WHERE rolname=current_user",
	).Scan(&canCreateRoles); err != nil {
		return errors.New("cannot inspect demo role privileges")
	}
	var roles auditpg.Roles
	var as demoAuditStore
	var sink admin.AuditSink
	if canCreateRoles {
		roles = auditpg.Roles{
			Deployment:  prefix + "_owner",
			Recorder:    prefix + "_recorder",
			Relay:       prefix + "_relay",
			Maintenance: prefix + "_maintenance",
			Reader:      prefix + "_reader",
		}
		if err = createRoles(ctx, pool, roles); err != nil {
			return err
		}
	}
	defer cleanup(ctx, pool, prefix, roles)
	if canCreateRoles {
		ac := auditpg.SchemaConfig{Schema: prefix + "_audit", Roles: roles}
		if err = auditpg.Migrate(ctx, pool, ac); err != nil {
			return fmt.Errorf("audit migration: %w", err)
		}
		store, storeErr := auditpg.New(ctx, pool, ac)
		if storeErr != nil {
			return storeErr
		}
		as = store
		sink, err = auditadmin.New(store)
	} else {
		fmt.Fprintln(
			os.Stderr,
			"development owner mode: execution and webhook use one database "+
				"owner; role separation checks and durable PostgreSQL audit are "+
				"skipped; audit evidence is in memory and lost on exit",
		)
		store, storeErr := audit.NewMemoryStore(audit.MemoryStoreConfig{
			Capacity: audit.Capacity{
				Events: 1000,
				Bytes:  4 << 20,
			},
		})
		if storeErr != nil {
			return storeErr
		}
		as = store
		sink = &memoryAuditSink{store: store, semantics: make(map[string]string)}
	}
	if err != nil {
		return err
	}
	auditKey := make([]byte, 32)
	if _, err = rand.Read(auditKey); err != nil {
		return err
	}
	schema, err := newAuditSchema(auditKey)
	if err != nil {
		return err
	}
	if err = seedAudit(ctx, as, schema); err != nil {
		return fmt.Errorf("seed audit: %w", err)
	}
	reader, err := audit.NewReader(as, audit.ReaderConfig{Schema: schema, Access: as})
	if err != nil {
		return err
	}
	ap, err := auditadmin.NewPanel(auditadmin.PanelConfig{
		Reader:       reader,
		EnableReveal: true,
	})
	if err != nil {
		return err
	}
	ec := pgstore.StoreConfig{Schema: prefix + "_execution"}
	if err = pgstore.Migrate(ctx, pool, ec); err != nil {
		return fmt.Errorf("execution migration: %w", err)
	}
	es, err := pgstore.New(ctx, pool, ec)
	if err != nil {
		return err
	}
	ep, records, err := seedExecutions(ctx, es, "payments")
	if err != nil {
		return fmt.Errorf("seed payments: %w", err)
	}
	fp, fulfillment, err := seedExecutions(ctx, es, "fulfillment")
	if err != nil {
		return fmt.Errorf("seed fulfillment: %w", err)
	}
	ws, err := seedWebhooks(ctx, pool, prefix+"_webhook")
	if err != nil {
		return fmt.Errorf("seed webhooks: %w", err)
	}
	if err = verifySeededThroughput(ctx, es, ws); err != nil {
		return fmt.Errorf("seed throughput: %w", err)
	}

	wp, err := webhookadmin.New(webhookadmin.Config{Store: ws, EnableReplay: true})
	if err != nil {
		return err
	}
	lim, qp, err := seedQuota(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if closeErr := lim.Close(closeCtx); closeErr != nil {
			fmt.Fprintln(os.Stderr, "quota close failed")
		}
	}()
	sp, err := seedStripe(ctx)
	if err != nil {
		return err
	}
	op, err := openapiadmin.New(openapiadmin.Config{})
	if err != nil {
		return err
	}
	e := echo.New()
	cfg := admin.Config{Panels: []admin.Panel{
		qp,
		sp,
		op,
		ep,
		wp,
		ap,
	}, Instance: "admin-demo", Keys: [][]byte{key}, Audit: sink,
		Theme: admin.Theme{
			Environment: "Development demo · Operator · writes enabled",
			Tone:        "warning",
		},
		Authenticate:   authentication.authenticate,
		Middleware:     authentication.middleware,
		ActorBoundCSRF: true,
		Origins:        []string{"http://" + demoBrowserAuthority(authority)},
	}
	//nolint:contextcheck // Mount registers routes; request contexts originate in Echo.
	if err = cfg.Mount(e, "/admin"); err != nil {
		return err
	}
	// executionadmin binds one namespace. A second console exposes the other.
	cfg.Panels = []admin.Panel{qp, sp, op, fp, wp, ap}
	//nolint:contextcheck // Mount registers routes; request contexts originate in Echo.
	if err = cfg.Mount(e, "/admin-fulfillment"); err != nil {
		return err
	}
	if err = mountSamples(ctx, e, lim, op, append(records, fulfillment...)); err != nil {
		return err
	}
	for _, mount := range []string{"/admin", "/admin-fulfillment"} {
		if err = admin.VerifyRoutes(e, mount); err != nil {
			return err
		}
	}
	// Exercise the actual console protocol at startup: CSRF, preview, confirmation,
	// commit and durable intent/outcome evidence, rather than fabricated admin logs.
	if err = seedAdminActions(ctx, e, authentication.credential, authority); err != nil {
		return fmt.Errorf("demo action smoke check: %w", err)
	}
	return serve(ctx, listener, e, ws)
}

func roleNames(r auditpg.Roles) []string {
	return []string{r.Deployment, r.Recorder, r.Relay, r.Maintenance, r.Reader}
}

func createRoles(ctx context.Context, pool *pgxpool.Pool, roles auditpg.Roles) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return errors.New("demo role bootstrap failed")
	}
	defer rollback(ctx, tx)
	for _, name := range roleNames(roles) {
		if _, err = tx.Exec(
			ctx,
			"CREATE ROLE "+pgx.Identifier{name}.Sanitize()+" NOLOGIN",
		); err != nil {
			return errors.New(
				"demo needs a disposable database administrator with CREATE ROLE " +
					"privileges",
			)
		}
	}
	var database, user string
	if err = tx.QueryRow(
		ctx,
		"SELECT current_database(), current_user",
	).Scan(&database, &user); err != nil {
		return errors.New("demo role bootstrap failed")
	}
	for _, sql := range []string{
		"GRANT " + pgx.Identifier{roles.Deployment}.Sanitize() + " TO " + pgx.Identifier{
			user,
		}.Sanitize(),
		"GRANT CREATE ON DATABASE " + pgx.Identifier{
			database,
		}.Sanitize() + " TO " + pgx.Identifier{roles.Deployment}.Sanitize(),
		"GRANT EXECUTE ON FUNCTION pg_catalog.pg_control_system() TO " + pgx.Identifier{
			roles.Deployment,
		}.Sanitize(),
	} {
		if _, err = tx.Exec(ctx, sql); err != nil {
			return errors.New("demo role grants failed")
		}
	}
	return tx.Commit(ctx)
}

func rollback(ctx context.Context, tx pgx.Tx) {
	closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	if err := tx.Rollback(closeCtx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		fmt.Fprintln(os.Stderr, "demo transaction rollback failed")
	}
}

func cleanup(
	ctx context.Context,
	pool *pgxpool.Pool,
	prefix string,
	roles auditpg.Roles,
) {
	closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 20*time.Second)
	defer cancel()
	for _, suffix := range []string{"_execution", "_webhook", "_audit"} {
		if _, err := pool.Exec(
			closeCtx,
			"DROP SCHEMA IF EXISTS "+pgx.Identifier{
				prefix + suffix,
			}.Sanitize()+" CASCADE",
		); err != nil {
			fmt.Fprintln(
				os.Stderr,
				"demo schema cleanup failed; remove only this demo run's schemas",
			)
			return
		}
	}
	for _, name := range roleNames(roles) {
		if name == "" {
			continue
		}
		for _, sql := range []string{
			"DROP OWNED BY " + pgx.Identifier{name}.Sanitize(),
			"DROP ROLE " + pgx.Identifier{name}.Sanitize(),
		} {
			if _, err := pool.Exec(closeCtx, sql); err != nil {
				fmt.Fprintln(
					os.Stderr,
					"demo role cleanup failed; stop the disposable fixture to remove it",
				)
				return
			}
		}
	}
}
