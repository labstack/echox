// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auth_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/auth"
	"echo.labstack.com/pro/problem/echoadapter"
)

// ExampleNewBearer uses an HTTPS issuer fixture. The application supplies its
// real issuer, API audiences and trusted client; the token supplies no URL.
// docs:start ExampleNewBearer
func ExampleNewBearer() {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key:       public,
		KeyID:     "signing-1",
		Use:       "sig",
		Algorithm: "EdDSA",
	}}})
	if err != nil {
		panic(err)
	}
	issuer := httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(body); err != nil {
			panic(err)
		}
	}))
	defer issuer.Close()
	ctx := context.Background()
	bearer, err := auth.NewBearer(ctx, auth.BearerConfig{
		Namespace: "employees", Issuer: issuer.URL, Audiences: []string{"invoice-api"},
		JWKSURL: issuer.URL + "/jwks", Client: issuer.Client(),
		TenantClaim: "tenant_id", TrustTenantClaim: true, ClaimAllowlist: []string{
			"department",
		},
	})
	if err != nil {
		panic(err)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: private},
		(&jose.SignerOptions{}).WithType("at+jwt").WithHeader("kid", "signing-1"))
	if err != nil {
		panic(err)
	}
	claims, err := json.Marshal(map[string]any{
		"iss": issuer.URL, "sub": "alice", "aud": "invoice-api", "exp": time.Now().Add(
			time.Hour,
		).Unix(),
		"scope": "invoice:read", "tenant_id": "acme", "department": "billing",
	})
	if err != nil {
		panic(err)
	}
	signed, err := signer.Sign(claims)
	if err != nil {
		panic(err)
	}
	token, err := signed.CompactSerialize()
	if err != nil {
		panic(err)
	}
	middleware, err := (auth.Config{Schemes: map[string]auth.Scheme{"oauth": bearer},
		Fallback: auth.Require(map[string][]string{
			"oauth": {"invoice:read"},
		})}).ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.HTTPErrorHandler, err = (echoadapter.ErrorHandlerConfig{}).ToErrorHandler(
		e.HTTPErrorHandler,
	)
	if err != nil {
		panic(err)
	}
	e.GET("/invoices", func(c *echo.Context) error {
		principal, _ := auth.FromEcho(c)
		return c.String(
			http.StatusOK,
			principal.ID()+" tenant="+principal.TenantRestriction(),
		)
	}, middleware)
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/invoices", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	fmt.Println(w.Code, w.Body.String())
	// Output: 200 employees:alice tenant=acme
}

// docs:end ExampleNewBearer
