// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package webhook_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"echo.labstack.com/pro/webhook"
)

// docs:start ExampleHTTPSenderConfig_ToSendFunc
func ExampleHTTPSenderConfig_ToSendFunc() {
	// This is a public test key, not a production secret.
	secret := "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw" //nolint:gosec // G101: public Standard Webhooks test vector.
	verify, err := (webhook.StandardVerifierConfig{Secrets: []string{secret},
		Source: "billing", Type: "invoice.paid", Format: "json"}).ToVerifyFunc()
	if err != nil {
		panic(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, webhook.MaxPayloadBytes+1))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err = verify(r.Context(), r, body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	send, err := (webhook.HTTPSenderConfig{Endpoints: map[string]webhook.Endpoint{
		"billing": {URL: server.URL, Secrets: []string{secret}},
	}}).ToSendFunc()
	if err != nil {
		panic(err)
	}
	result, err := send(context.Background(), webhook.Delivery{Destination: "billing", Key: "invoice_42/paid", Format: "json", Payload: []byte(`{"amount":100}`)})
	fmt.Println(result.Outcome, result.Code, err)
	// Output: accepted http_204 <nil>
}

// docs:end ExampleHTTPSenderConfig_ToSendFunc

// docs:start ExampleStandardID
func ExampleStandardID() {
	fmt.Println(webhook.StandardID("invoice_42/grant"))
	fmt.Println(webhook.StandardID("invoice.42"))
	// Output:
	// true
	// false
}

// docs:end ExampleStandardID
