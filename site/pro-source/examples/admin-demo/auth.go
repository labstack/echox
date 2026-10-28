// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/auth/authbind"
)

var demoTokenPattern = regexp.MustCompile(`^demo_[A-Za-z0-9]{22}_[A-Za-z0-9]{43}$`)

type demoAuth struct {
	credential   string
	middleware   []echo.MiddlewareFunc
	authenticate func(*echo.Context) (admin.Principal, error)
}

// An empty token enables browser mode with a freshly issued, process-local key.
// A supplied token requires explicit authentication on every console request.
func demoAuthentication(ctx context.Context, token, authority string) (demoAuth, error) {
	browserMode := token == ""
	if !browserMode && !demoTokenPattern.MatchString(token) {
		return demoAuth{}, errors.New("provide ADMIN_DEMO_API_KEY from auth.IssueKey with prefix demo; credential withheld")
	}
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		return demoAuth{}, err
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	opts := auth.KeyOptions{Kind: auth.Service, Scopes: []string{"admin:operator", "admin:global", "admin:reveal"}}
	var record auth.StoredKey
	if browserMode {
		var err error
		token, record, err = auth.IssueKey("demo", ring, "demo-operator", opts)
		if err != nil {
			return demoAuth{}, err
		}
	} else {
		// IssueKey only generates new credentials; auth has no import API. This
		// disposable demo must compute HMAC-SHA256 over the supplied credential
		// to import it under its fresh local pepper. Production instead persists
		// IssueKey's original StoredKey and manages that record's pepper ring.
		mac := hmac.New(sha256.New, pepper)
		_, _ = mac.Write([]byte(token))
		record = auth.StoredKey{ID: strings.Split(token, "_")[1], PepperVersion: 1, Subject: "demo-operator", Kind: opts.Kind, Scopes: opts.Scopes}
		copy(record.MAC[:], mac.Sum(nil))
	}
	store := &auth.MemoryKeyStore{Capacity: 1}
	if err := store.Put(ctx, record); err != nil {
		return demoAuth{}, err
	}
	//nolint:contextcheck // Construction only; returned middleware obtains its own request context.
	mw, err := (auth.Config{Schemes: map[string]auth.Scheme{"key": auth.APIKey{Namespace: "demo", Prefix: "demo", Store: store, Peppers: ring}}, Fallback: auth.Require(map[string][]string{"key": nil})}).ToMiddleware()
	if err != nil {
		return demoAuth{}, err
	}
	middleware := []echo.MiddlewareFunc{demoAuthority(authority), mw}
	if browserMode {
		middleware = []echo.MiddlewareFunc{demoAuthority(authority), demoBrowserKey(token), mw}
	}
	//nolint:contextcheck // Returns a callback; identity is read from that callback's request context.
	return demoAuth{credential: token, middleware: middleware, authenticate: authbind.Admin(authbind.AdminGrants{Operator: "admin:operator", Global: "admin:global", Reveal: "admin:reveal"})}, nil
}

func demoBrowserKey(token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if len(c.Request().Header.Values("Authorization")) == 0 {
				c.Request().Header.Set("Authorization", "ApiKey "+token)
			}
			return next(c)
		}
	}
}

// demoAuthority checks HTTP authority before any operator credential is injected.
// Only the actual listener address is trusted; forwarded headers have no effect.
func demoAuthority(authority string) echo.MiddlewareFunc {
	browserAuthority := demoBrowserAuthority(authority)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Host != authority && c.Request().Host != browserAuthority {
				return echo.NewHTTPError(http.StatusForbidden, "untrusted demo authority")
			}
			return next(c)
		}
	}
}

// demoBrowserAuthority omits HTTP's default port from an actual listener address,
// preserving the exact IP (including IPv6 brackets) and every nondefault port.
func demoBrowserAuthority(authority string) string {
	return strings.TrimSuffix(authority, ":80")
}
