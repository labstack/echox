// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package quotaclient_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"echo.labstack.com/pro/quota"
	"echo.labstack.com/pro/quota/quotaclient"
)

// Branch on the stable code rather than the title or detail text.
// docs:start ExampleParseProblem
func ExampleParseProblem() {
	resp := &http.Response{
		StatusCode: http.StatusPaymentRequired,
		Header:     http.Header{"Content-Type": {"application/problem+json"}},
		Body: io.NopCloser(strings.NewReader(`{"type":"https://echo.labstack.com/pro/problems/insufficient-credits",` +
			`"title":"Insufficient credits","status":402,"code":"insufficient_credits","violated-policies":["credits"]}`)),
	}
	p, err := quotaclient.ParseProblem(resp)
	if err != nil || p == nil {
		fmt.Println("not a quota problem")
		return
	}
	fmt.Println(p.Status, p.Code, p.Code == quota.CodeInsufficientCredits, p.ViolatedPolicies)
	// Output: 402 insufficient_credits true [credits]
}

// docs:end ExampleParseProblem
