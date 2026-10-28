// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/webhook/webhookpg"
)

func serve(
	ctx context.Context,
	listener net.Listener,
	e *echo.Echo,
	store *webhookpg.Store,
) error {
	worker, err := webhookpg.NewWorker(store, webhookpg.WorkerConfig{
		Concurrency:     4,
		ShutdownTimeout: 3 * time.Second,
	})
	if err != nil {
		return err
	}
	work, cancel := context.WithCancel(ctx)
	defer cancel()
	workerDone := make(chan error, 1)
	serverDone := make(chan error, 1)
	server := &http.Server{
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    16 << 10,
	}
	go func() { workerDone <- worker.Run(work) }()
	go func() { serverDone <- server.Serve(listener) }()
	fmt.Printf(
		"Admin demo: http://%s/admin\n"+
			"Seed directory: http://%s/\n"+
			"Fulfillment console: http://%s/admin-fulfillment\n",
		listener.Addr(),
		listener.Addr(),
		listener.Addr(),
	)
	fmt.Println(
		"Ctrl-C shuts down the server/worker and removes this run's " +
			"schemas and roles. Stop any PostgreSQL service you started for " +
			"this demo.",
	)
	var result error
	workerJoined, serverJoined := false, false
	select {
	case <-ctx.Done():
	case result = <-workerDone:
		workerJoined = true
	case result = <-serverDone:
		serverJoined = true
	}
	cancel()
	shutdown, stop := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer stop()
	if err = server.Shutdown(shutdown); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			fmt.Fprintln(os.Stderr, "demo server close failed")
		}
		result = err
	}
	if !workerJoined {
		if workerErr := <-workerDone; workerErr != nil && !errors.Is(
			workerErr,
			context.Canceled,
		) {
			result = workerErr
		}
	}
	if !serverJoined {
		if serverErr := <-serverDone; serverErr != nil && !errors.Is(
			serverErr,
			http.ErrServerClosed,
		) {
			result = serverErr
		}
	}
	if errors.Is(result, context.Canceled) || errors.Is(result, http.ErrServerClosed) {
		return nil
	}
	return result
}

func seedAdminActions(
	ctx context.Context,
	e *echo.Echo,
	credential, authority string,
) error {
	if err := consoleAction(
		ctx,
		e,
		credential,
		authority,
		"/admin/api/v1/quota/credit-topup",
		map[string]string{
			"tenant":     "acme",
			"plan":       "team",
			"limit":      "credits",
			"amount":     "250",
			"invoice_id": "inv_demo_paid",
		},
	); err != nil {
		return err
	}
	// Use a separate partition so the reset smoke check preserves review usage.
	return consoleAction(
		ctx,
		e,
		credential,
		authority,
		"/admin/api/v1/quota/reset",
		map[string]string{
			"tenant": "demo-smoke",
			"plan":   "team",
			"limits": "daily",
		},
	)
}

// consoleAction follows the browser protocol against the mounted Echo handler.
// Tokens stay in process memory, and responses never reach logs or disk.
func consoleAction(
	ctx context.Context,
	e *echo.Echo,
	credential, authority, path string,
	input map[string]string,
) error {
	request := func(method string, fields map[string]string, into any) error {
		var body bytes.Buffer
		if fields != nil {
			if err := json.NewEncoder(&body).Encode(fields); err != nil {
				return err
			}
		}
		r := httptest.NewRequestWithContext(ctx, method, "http://"+authority+path, &body)
		r.Header.Set("Sec-Fetch-Site", "same-origin")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "ApiKey "+credential)
		w := httptest.NewRecorder()
		e.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			var problem struct{ Code string }
			if err := json.Unmarshal(w.Body.Bytes(), &problem); err != nil {
				return fmt.Errorf("console %s %s returned HTTP %d", method, path, w.Code)
			}
			return fmt.Errorf(
				"console %s %s returned HTTP %d (%s)",
				method,
				path,
				w.Code,
				problem.Code,
			)
		}
		if into != nil {
			return json.Unmarshal(w.Body.Bytes(), into)
		}
		return nil
	}
	var form struct{ Controls map[string]string }
	if err := request(http.MethodGet, nil, &form); err != nil {
		return err
	}
	fields := form.Controls
	if fields == nil || fields["action_id"] == "" {
		return errors.New("console did not provide header-authenticated action controls")
	}
	for k, v := range input {
		fields[k] = v
	}
	var preview struct {
		Confirmation     string
		Typed            string
		Acknowledgements map[string]string
		CommitInput      map[string]string `json:"commit_input"`
	}
	if err := request(http.MethodPost, fields, &preview); err != nil {
		return err
	}
	if preview.Confirmation == "" {
		return errors.New("console preview did not provide a confirmation")
	}
	for k, v := range preview.CommitInput {
		fields[k] = v
	}
	fields["mode"] = "commit"
	fields["confirmation"] = preview.Confirmation
	fields["typed"] = preview.Typed
	fields["reason"] = "demo_startup_check"
	for k := range preview.Acknowledgements {
		fields["ack."+k] = "yes"
	}
	return request(http.MethodPost, fields, nil)
}
