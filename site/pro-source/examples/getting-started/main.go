// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

// Command getting-started runs a memory-only, loopback-only Echo Pro app.
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/auth/authbind"
	"echo.labstack.com/pro/problem/echoadapter"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaadmin"
	"echo.labstack.com/pro/tenant"
	"echo.labstack.com/pro/webhook"
	"echo.labstack.com/pro/webhook/webhookmem"
)

func main() {
	browser := flag.Bool(
		"browser",
		false,
		"grant local browser requests read-only admin access (not for production)",
	)
	keyFile := flag.String(
		"key-file",
		"",
		"new private file for the disposable API key (required)",
	)
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, runConfig{
		keyFile: *keyFile,
		browser: *browser,
		address: "127.0.0.1:8080",
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type runConfig struct {
	keyFile  string
	browser  bool
	address  string
	ready    func(net.Addr)
	verified func()
}

func run(ctx context.Context, cfg runConfig) (result error) {
	keyFile := cfg.keyFile
	if keyFile == "" {
		return errors.New(
			"provide -key-file pointing to a new private file for the disposable API key",
		)
	}
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		return err
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	key, record, err := auth.IssueKey("demo", ring, "invoice-worker", auth.KeyOptions{
		Kind: auth.Service, TenantRestriction: "acme", Scopes: []string{
			"invoice:read",
			"invoice:write",
			"admin:read",
		},
	})
	if err != nil {
		return err
	}
	keys := &auth.MemoryKeyStore{Capacity: 1}
	if err := keys.Put(ctx, record); err != nil {
		return err
	}
	// The demo delivers its disposable credential through a private local file,
	// never stdout. Refuse existing paths; remove only the file we created.
	//nolint:gosec // G304: CLI path, exclusively created with mode 0600.
	f, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer func() { result = errors.Join(result, os.Remove(keyFile)) }()
	_, writeErr := io.WriteString(f, key+"\n")
	if err := errors.Join(writeErr, f.Close()); err != nil {
		return err
	}
	// Bind first so the sender and console use the actual test port.
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", cfg.address)
	if err != nil {
		return err
	}
	defer func() {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			result = errors.Join(result, err)
		}
	}()
	authority := listener.Addr().String()
	// docs:start site-auth
	authCfg := auth.Config{
		Schemes: map[string]auth.Scheme{"key": auth.APIKey{
			Namespace: "app",
			Prefix:    "demo",
			Peppers:   ring,
			Store:     keys,
			Header:    "X-API-Key",
		}},
		Fallback: auth.Require(map[string][]string{"key": {"invoice:read"}}),
	}
	tenantPath, err := tenant.FromPathParam("tenant")
	if err != nil {
		return err
	}
	//nolint:contextcheck // Constructs a guard that uses each request context.
	guard, err := authbind.Middleware(authCfg, tenant.Config{
		Selectors: []tenant.Selector{tenant.FromCredential(), tenantPath}, Required: true,
	})
	if err != nil {
		return err
	}
	// docs:end site-auth
	// A daily counter makes the walkthrough repeatable without refill timing.
	//nolint:contextcheck // Constructor has no context API; no outbox is configured.
	// docs:start site-quota
	lim, err := quota.New(quota.Config{
		Store: quota.NewMemoryStore(quota.MemoryConfig{MaxKeys: 100}),
		Plans: []quota.Plan{{ID: "free", Limits: []quota.Limit{
			quota.PerDay("requests", 3),
		}}},
	})
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		result = errors.Join(result, lim.Close(closeCtx))
	}()
	//nolint:contextcheck // Builds a callback that uses the Echo request context.
	identity, err := authbind.Quota()
	if err != nil {
		return err
	}
	//nolint:contextcheck // Constructor only; callbacks use request cancellation.
	quotaMW, err := (quota.MiddlewareConfig{
		Limiter:  lim,
		Identity: identity,
	}).ToMiddleware()
	if err != nil {
		return err
	}
	// docs:end site-quota
	e := echo.New()
	//nolint:contextcheck // Constructor only; rendering uses the request context.
	e.HTTPErrorHandler, err = (echoadapter.ErrorHandlerConfig{}).ToErrorHandler(
		e.HTTPErrorHandler,
	)
	if err != nil {
		return err
	}
	//nolint:contextcheck // Registers verification that uses each request context.
	outbox, err := localWebhooks(e, authority, cfg.verified)
	if err != nil {
		return err
	}
	// docs:start site-group
	api := e.Group("/v1/tenants/:tenant", guard, quotaMW)
	api.GET("/invoices", func(c *echo.Context) error {
		id, _ := tenant.From(c.Request().Context())
		return c.JSON(http.StatusOK, map[string]any{
			"tenant":   id,
			"invoices": []string{"inv-1"},
		})
	})
	// docs:end site-group
	//nolint:contextcheck // Enqueues with each Echo request context.
	api.POST("/invoices", func(c *echo.Context) error {
		id, _ := tenant.From(c.Request().Context())
		invoice := "inv_" + rand.Text()
		// docs:start site-webhook
		prepared, err := webhook.PrepareOutgoing(webhook.Outgoing{
			Tenant: string(
				id,
			), Destination: "invoices", Version: "v1", Key: invoice + "/created",
			Type: "invoice.created", Format: "json", Payload: []byte(fmt.Sprintf(
				`{"invoice":%q}`,
				invoice,
			)),
		})
		if err != nil {
			return err
		}
		receipt, err := outbox.Enqueue(c.Request().Context(), prepared)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusAccepted, map[string]string{
			"invoice":     invoice,
			"webhook_ref": receipt.Ref,
		})
		// docs:end site-webhook
	})
	//nolint:contextcheck // Registers callbacks; authentication uses request contexts.
	if err := mountAdmin(e, lim, guard, key, cfg.browser, authority); err != nil {
		return err
	}
	worker, err := webhookmem.NewWorker(outbox, webhookmem.WorkerConfig{
		Concurrency:     1,
		ShutdownTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	workerCtx, cancelWorker := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- worker.Run(workerCtx) }()
	defer func() {
		cancelWorker()
		if err := <-done; !errors.Is(err, context.Canceled) {
			result = errors.Join(result, err)
		}
	}()
	fmt.Printf("API: http://%s/v1/tenants/acme/invoices\n", authority)
	fmt.Printf(
		"Admin: http://%s/admin (use -browser for local browser access)\n",
		authority,
	)
	sc := echo.StartConfig{
		Listener: listener, ListenerAddrFunc: cfg.ready, GracefulTimeout: 5 * time.Second,
		BeforeServeFunc: func(s *http.Server) error {
			s.ReadHeaderTimeout = 5 * time.Second
			s.ReadTimeout = 10 * time.Second
			s.WriteTimeout = 15 * time.Second
			s.MaxHeaderBytes = 16 << 10
			return nil
		},
	}
	if err := sc.Start(ctx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func localWebhooks(e *echo.Echo, authority string, verified func()) (
	*webhookmem.Store,
	error,
) {
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, err
	}
	secret := "whsec_" + base64.StdEncoding.EncodeToString(secretBytes)
	verify, err := (webhook.StandardVerifierConfig{
		Secrets: []string{
			secret,
		}, Source: "app", Tenant: "acme", Type: "invoice.created", Format: "json",
	}).ToVerifyFunc()
	if err != nil {
		return nil, err
	}
	// This receiver verifies bytes only; it performs no business effect. A real
	// receiver also needs inbox deduplication before applying an effect.
	e.POST("/hooks/invoices", func(c *echo.Context) error {
		r := c.Request()
		body, err := io.ReadAll(io.LimitReader(r.Body, webhook.MaxPayloadBytes+1))
		if err != nil {
			return echo.ErrBadRequest
		}
		if _, err := verify(r.Context(), r, body); err != nil {
			return echo.ErrBadRequest
		}
		fmt.Println("Webhook delivered: verified invoice.created")
		if verified != nil {
			verified()
		}
		return c.NoContent(http.StatusNoContent)
	})
	send, err := (webhook.HTTPSenderConfig{Endpoints: map[string]webhook.Endpoint{
		"invoices": {
			URL:     "http://" + authority + "/hooks/invoices",
			Secrets: []string{secret},
		},
	}}).ToSendFunc()
	if err != nil {
		return nil, err
	}
	return webhookmem.New(webhookmem.StoreConfig{MaxMessages: 100,
		Routing: webhookmem.RoutingConfig{
			Destinations: map[string]webhookmem.DestinationConfig{
				"invoices": {Version: "v1", Send: send, ReplayWindow: 24 * time.Hour},
			},
		},
	})
}

func mountAdmin(
	e *echo.Echo,
	lim *quota.Limiter,
	guard echo.MiddlewareFunc,
	key string,
	browser bool,
	authority string,
) error {
	panel, err := quotaadmin.New(quotaadmin.Config{Limiter: lim})
	if err != nil {
		return err
	}
	chain := []echo.MiddlewareFunc{guard}
	if browser {
		// Not for production: anyone reaching this loopback listener at its exact
		// authority gets the demo's read-only admin credential on console routes.
		chain = []echo.MiddlewareFunc{func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c *echo.Context) error {
				if c.Request().Host != authority {
					return echo.ErrForbidden
				}
				if c.Request().Header.Get("X-API-Key") == "" {
					c.Request().Header.Set("X-API-Key", key)
				}
				return next(c)
			}
		}, guard}
	}
	cfg := admin.Config{
		Panels: []admin.Panel{panel}, Instance: "local", Origins: []string{
			"http://" + authority,
		},
		Middleware: chain, Authenticate: authbind.Admin(authbind.AdminGrants{
			Viewer: "admin:read",
		}),
	}
	if err := cfg.Mount(e, "/admin"); err != nil {
		return err
	}
	return admin.VerifyRoutes(e, "/admin")
}
