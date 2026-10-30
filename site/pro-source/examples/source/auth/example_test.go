// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package auth_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/auth"
)

// docs:start ExampleIssueKey
func ExampleIssueKey() {
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		panic(err)
	}
	token, record, err := auth.IssueKey("prod", auth.PepperRing{{
		Version: 7,
		Key:     pepper,
	}}, "billing-worker", auth.KeyOptions{Kind: auth.Service, TenantRestriction: "acme"})
	if err != nil {
		panic(err)
	} // Deliver token once through a secure channel; persist only record.
	fmt.Println(
		len(token),
		len(record.ID),
		record.PepperVersion,
		record.TenantRestriction,
	)
	// Output: 71 22 7 acme
}

// docs:end ExampleIssueKey
// docs:start ExampleNewPrincipal
func ExampleNewPrincipal() {
	p, err := auth.NewPrincipal(auth.PrincipalConfig{
		Namespace: "employees",
		Subject:   "alice",
		Scheme:    "oauth",
		Kind:      auth.User,
		Scopes:    []string{"invoice:read"},
	})
	if err != nil {
		panic(err)
	}
	copy := p.Scopes()
	copy[0] = "invoice:write"
	fmt.Println(p.ID(), p.Scopes()[0])
	// Output: employees:alice invoice:read
}

// docs:end ExampleNewPrincipal

// docs:start either-credential
// MountEitherCredential adds one route accepting either configured scheme.
func MountEitherCredential(e *echo.Echo, key auth.APIKey, jwt *auth.Bearer) error {
	guard, err := (auth.Config{
		Schemes: map[string]auth.Scheme{"key": key, "jwt": jwt},
		Fallback: auth.Require(
			map[string][]string{"key": {"invoice:read"}},
			map[string][]string{"jwt": {"invoice:read"}},
		),
	}).ToMiddleware()
	if err != nil {
		return err
	}
	e.GET("/invoices", func(c *echo.Context) error { return c.NoContent(204) }, guard)
	return nil
}

// docs:end either-credential

func TestMountEitherCredential(t *testing.T) {
	// Use the HTTPS issuer pattern from ExampleNewBearer; newJWKS belongs to
	// the internal auth test package and is unavailable in auth_test.
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key: public, KeyID: "signing-1", Use: "sig", Algorithm: "EdDSA",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	issuer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(body); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(issuer.Close)
	now := time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC)
	bearer, err := auth.NewBearer(t.Context(), auth.BearerConfig{
		Namespace: "employees", Issuer: issuer.URL, Audiences: []string{"invoice-api"},
		JWKSURL: issuer.URL + "/jwks", Client: issuer.Client(), Now: func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: private},
		(&jose.SignerOptions{}).WithType("at+jwt").WithHeader("kid", "signing-1"))
	if err != nil {
		t.Fatal(err)
	}
	claims, err := json.Marshal(map[string]any{
		"iss": issuer.URL, "sub": "alice", "aud": "invoice-api",
		"exp": now.Add(time.Hour).Unix(), "scope": "invoice:read",
	})
	if err != nil {
		t.Fatal(err)
	}
	signed, err := signer.Sign(claims)
	if err != nil {
		t.Fatal(err)
	}
	jwt, err := signed.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	pepper := make([]byte, 32)
	if _, err := rand.Read(pepper); err != nil {
		t.Fatal(err)
	}
	ring := auth.PepperRing{{Version: 1, Key: pepper}}
	token, record, err := auth.IssueKey("test", ring, "worker", auth.KeyOptions{
		Kind: auth.Service, Scopes: []string{"invoice:read"},
	})
	if err != nil {
		t.Fatal(err)
	}
	keys := &auth.MemoryKeyStore{Capacity: 1}
	if err := keys.Put(t.Context(), record); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	if err := MountEitherCredential(e, auth.APIKey{
		Namespace: "services", Prefix: "test", Peppers: ring, Store: keys,
	}, bearer); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, authorization string
		want                int
	}{
		{"api key", "ApiKey " + token, http.StatusNoContent},
		{"jwt", "Bearer " + jwt, http.StatusNoContent},
		{"no credential", "", http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/invoices", nil)
			if tc.authorization != "" {
				r.Header.Set("Authorization", tc.authorization)
			}
			w := httptest.NewRecorder()
			e.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("status = %d, want %d: %s", w.Code, tc.want, w.Body)
			}
		})
	}
}
