// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotaadmin_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/admin"
	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaadmin"
)

// docs:start ExamplePanel_creditWrites
func ExamplePanel_creditWrites() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lim, err := quota.New(quota.Config{Store: quota.NewMemoryStore(
		quota.MemoryConfig{},
	), Plans: []quota.Plan{{ID: "team", Limits: []quota.Limit{
		quota.Credits("credits"),
	}}}})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := lim.Close(context.WithoutCancel(ctx)); err != nil {
			fmt.Println(err)
		}
	}()
	panel, err := quotaadmin.New(quotaadmin.Config{Limiter: lim, EnableCredits: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		fmt.Println(err)
		return
	}
	e := echo.New()
	// Only this disposable example opts out; production supplies auditadmin.New(store).
	err = (admin.Config{Panels: []admin.Panel{panel}, Instance: "example", Keys: [][]byte{
		key,
	}, AllowUnauditedWrites: true, Authenticate: func(
		*echo.Context,
	) (admin.Principal, error) {
		return admin.Principal{
			ID:               "operator",
			Role:             admin.Operator,
			Tenants:          []string{"acme"},
			HeaderCredential: true,
		}, nil
	}}).Mount(e, "/admin")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = admin.VerifyRoutes(e, "/admin"); err != nil {
		fmt.Println(err)
		return
	}
	fields := map[string]string{
		"mode":       "preview",
		"action_id":  "00000000000000000000000000000001",
		"tenant":     "acme",
		"plan":       "team",
		"limit":      "credits",
		"amount":     "10",
		"invoice_id": "ticket-57",
	}
	post := func() *httptest.ResponseRecorder {
		body, err := json.Marshal(fields)
		if err != nil {
			panic(err)
		}
		request := httptest.NewRequestWithContext(
			ctx,
			"POST",
			"http://console.example/admin/api/v1/quota/credit-topup",
			bytes.NewReader(body),
		)
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e.ServeHTTP(w, request)
		return w
	}
	preview := post()
	var confirmation struct {
		Token string `json:"confirmation"`
	}
	if err = json.Unmarshal(preview.Body.Bytes(), &confirmation); err != nil {
		fmt.Println(err)
		return
	}
	fields["mode"] = "commit"
	fields["confirmation"] = confirmation.Token
	fields["reason"] = "Invoice payment"
	committed := post()
	var answer struct {
		Result quota.CreditResult `json:"result"`
		Audit  struct {
			Outcome string `json:"outcome"`
		} `json:"audit"`
	}
	if err = json.Unmarshal(committed.Body.Bytes(), &answer); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(
		committed.Code,
		answer.Result.Balance,
		answer.Result.Applied,
		answer.Audit.Outcome,
	)
	// Output: 200 10 true applied
}

// docs:end ExamplePanel_creditWrites
