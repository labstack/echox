// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package webhook_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/webhook"
)

// inbox is an illustrative in-memory Acceptor; production uses webhookpg.Store.
type inbox struct {
	mu   sync.Mutex
	seen map[string][32]byte
}

func (
	b *inbox,
) Accept(_ context.Context, p webhook.PreparedIncoming) (webhook.Receipt, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	in := p.Incoming()
	id := in.Tenant + "\x00" + in.Account + "\x00" + in.EventID
	if old, ok := b.seen[id]; ok {
		if old != p.Digest() {
			return webhook.Receipt{}, webhook.ErrConflict
		}
		return webhook.Receipt{Digest: old, Duplicate: true}, nil
	}
	b.seen[id] = p.Digest()
	return webhook.Receipt{Digest: p.Digest(), RecordedAt: time.Now()}, nil
}

// docs:start ExampleReceiveConfig_ToMiddleware
func ExampleReceiveConfig_ToMiddleware() {
	secret := []byte("example-signing-key")
	sign := func(body string) string {
		m := hmac.New(sha256.New, secret)
		m.Write([]byte(body))
		return hex.EncodeToString(m.Sum(nil))
	}
	mw, err := webhook.ReceiveConfig{
		Inbox:  &inbox{seen: map[string][32]byte{}},
		Source: "billing",
		// Verify authenticates the original bytes, then maps a trusted tenant and
		// returns the minimized processing representation.
		Verify: func(_ context.Context, r *http.Request, body []byte) (
			webhook.Incoming,
			error,
		) {
			if !hmac.Equal([]byte(
				r.Header.Get("X-Signature"),
			), []byte(sign(string(body)))) {
				return webhook.Incoming{}, errors.New("invalid signature")
			}
			id, _, _ := strings.Cut(string(body), ":")
			return webhook.Incoming{
				Tenant:  "t1",
				Source:  "billing",
				Account: "acct",
				EventID: id,
				Type:    "invoice.paid",
				Format:  "text",
				Payload: body,
			}, nil
		},
	}.ToMiddleware()
	if err != nil {
		fmt.Println(err)
		return
	}
	e := echo.New()
	// The middleware is terminal: this handler is never called.
	// docs:start site-anchor
	e.POST("/webhooks/billing", func(c *echo.Context) error {
		return c.NoContent(http.StatusInternalServerError)
	}, mw)
	// docs:end site-anchor
	send := func(body, signature string) {
		req := httptest.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/webhooks/billing",
			strings.NewReader(body),
		)
		req.Header.Set("X-Signature", signature)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		fmt.Println(rec.Code)
	}
	send("evt_1:paid", sign("evt_1:paid")) // committed
	send("evt_1:paid", sign("evt_1:paid")) // identical duplicate
	send("evt_1:void", sign("evt_1:void")) // same event ID, changed content
	send("evt_2:paid", "forged")           // never persisted
	// Output:
	// 204
	// 204
	// 409
	// 400
}

// docs:end ExampleReceiveConfig_ToMiddleware

// docs:start ExamplePrepareOutgoing
func ExamplePrepareOutgoing() {
	out := webhook.Outgoing{Tenant: "t1", Destination: "credits", Version: "v1",
		Key: "invoice_42/grant", Type: "credit.add", Format: "json", Payload: []byte(
			`{"amount":100}`,
		)}
	a, err := webhook.PrepareOutgoing(out)
	if err != nil {
		fmt.Println(err)
		return
	}
	out.Payload = []byte(`{"amount":200}`) // the same key with changed content conflicts
	b, err := webhook.PrepareOutgoing(out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a.Valid(), a.Digest() == b.Digest())
	fmt.Println(a)
	// Output:
	// true false
	// webhook.PreparedOutgoing{Route:"credits" Type:"credit.add" Format:"json" Payload:[14 bytes]}
}

// docs:end ExamplePrepareOutgoing

// docs:start ExampleConfigError
func ExampleConfigError() {
	// The zero ReceiveConfig is invalid: Inbox, Verify and Source are required.
	_, err := webhook.ReceiveConfig{}.ToMiddleware()
	// errors.Is matches the category; Field is a diagnostic.
	fmt.Println(errors.Is(
		err,
		&webhook.ConfigError{},
	), errors.Is(err, &webhook.ValidationError{}))
	if e, ok := errors.AsType[*webhook.ConfigError](err); ok {
		fmt.Println(e.Field)
	}
	// Output:
	// true false
	// Inbox
}

// docs:end ExampleConfigError
