// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

// Command admin demonstrates local-only, read-only customer authentication.
package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaadmin"
)

func run(ctx context.Context) error {
	//nolint:contextcheck // Construction has no context API; no Outbox is configured.
	l, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.PerDay("calls", 100),
		quota.Credits("credits").Meter("tokens"),
	}}}})
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
		defer cancel()
		if closeErr := l.Close(closeCtx); closeErr != nil {
			fmt.Fprintln(os.Stderr, "limiter close failed")
		}
	}()
	p, err := quotaadmin.New(quotaadmin.Config{Limiter: l})
	if err != nil {
		return err
	}
	// One process-local fixture login; use the validated login's session binding
	// for each request in a deployed application, never a shared fixed value.
	session := rand.Text()
	e := echo.New()
	prefix := "/admin"
	// The outer CORS middleware must let every console method enter its wrapper.
	cors, err := (middleware.CORSConfig{AllowOrigins: []string{
		"http://127.0.0.1:8080",
	}, Skipper: func(c *echo.Context) bool {
		path := c.Request().URL.Path
		return path == prefix || strings.HasPrefix(path, prefix+"/")
	}}).ToMiddleware()
	if err != nil {
		return err
	}
	e.Use(cors)
	cfg := admin.Config{Panels: []admin.Panel{
		p,
	}, Instance: "local-demo", Theme: admin.Theme{
		Environment: "Local demonstration",
	}, Authenticate: func(*echo.Context) (admin.Principal, error) {
		// Replace with validated customer session/bearer claims in a deployed app.
		// Fixed auth is confined to this loopback-only demonstration.
		return admin.Principal{
			ID:      "local-reader",
			Session: session,
			Tenants: []string{"demo"},
		}, nil
	}}
	//nolint:contextcheck // Mount only registers routes; Echo supplies request contexts.
	if err = cfg.Mount(e, prefix); err != nil {
		return err
	}
	if err = admin.VerifyRoutes(e, prefix); err != nil {
		return err
	}
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		err := server.Shutdown(shutdown)
		if serveErr := <-done; !errors.Is(serveErr, http.ErrServerClosed) {
			return errors.Join(err, serveErr)
		}
		return err
	}
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
