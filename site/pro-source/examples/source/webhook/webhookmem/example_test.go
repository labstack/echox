// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package webhookmem_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"echo.labstack.com/pro/webhook"
	"echo.labstack.com/pro/webhook/webhookmem"
)

func Example() {
	// BEGIN quickstart
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Public test fixture only. Load a crypto/rand-generated key from your
	// secret manager for real endpoints; never reuse keys between customers.
	//nolint:gosec // G101: public Standard Webhooks test vector.
	secret := "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	verify, err := (webhook.StandardVerifierConfig{
		Secrets: []string{secret}, Source: "billing", Tenant: "t1",
		Type: "credit.add", Format: "json",
	}).ToVerifyFunc()
	if err != nil {
		panic(err)
	}
	verified := make(chan webhook.Incoming, 1)
	receiver := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		body, err := io.ReadAll(io.LimitReader(r.Body, webhook.MaxPayloadBytes+1))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		in, err := verify(r.Context(), r, body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		select {
		case verified <- in:
		default:
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer receiver.Close()
	send, err := (webhook.HTTPSenderConfig{Endpoints: map[string]webhook.Endpoint{
		"credits": {URL: receiver.URL, Secrets: []string{secret}},
	}}).ToSendFunc()
	if err != nil {
		panic(err)
	}
	store, err := webhookmem.New(webhookmem.StoreConfig{MaxMessages: 100,
		Routing: webhookmem.RoutingConfig{
			Destinations: map[string]webhookmem.DestinationConfig{
				"credits": {Version: "v1", Send: send, ReplayWindow: 24 * time.Hour},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	p, err := webhook.PrepareOutgoing(webhook.Outgoing{
		Tenant: "t1", Destination: "credits", Version: "v1", Key: "invoice_42/grant",
		Type: "credit.add", Format: "json", Payload: []byte(`{"amount":100}`),
	})
	if err != nil {
		panic(err)
	}
	receipt, err := store.Enqueue(ctx, p)
	if err != nil {
		panic(err)
	}
	record, err := store.InspectRecord(ctx, receipt.Ref)
	if err != nil {
		panic(err)
	}
	fmt.Println(record.State) // ready: no worker has started
	worker, err := webhookmem.NewWorker(store, webhookmem.WorkerConfig{Concurrency: 1})
	if err != nil {
		panic(err)
	}
	done := make(chan error, 1)
	go func() { done <- worker.Run(ctx) }()
	select {
	case in := <-verified:
		fmt.Println("verified", in.Type, in.EventID)
	case <-ctx.Done():
		panic(ctx.Err())
	}
	// Receiver verification precedes the sender's recorded completion.
	poll := time.NewTicker(10 * time.Millisecond)
	defer poll.Stop()
	for {
		record, err = store.InspectRecord(ctx, receipt.Ref)
		if err != nil {
			panic(err)
		}
		if record.State == "done" {
			break
		}
		select {
		case <-poll.C:
		case <-ctx.Done():
			panic(ctx.Err())
		}
	}
	fmt.Println(record.State) // done: the sender reported acceptance
	cancel()
	// Join the caller-managed worker before releasing its dependencies.
	if err := <-done; !errors.Is(err, context.Canceled) {
		panic(err)
	}
	// END quickstart
	// Output:
	// ready
	// verified credit.add invoice_42/grant
	// done
}
