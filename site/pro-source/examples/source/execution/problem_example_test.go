// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package execution_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/execution"
	"echo.labstack.com/pro/problem"
	"echo.labstack.com/pro/problem/echoadapter"
)

// orderRejected is an application error; any type with Problem() qualifies.
type orderRejected struct{}

func (orderRejected) Error() string { return "order rejected" }
func (orderRejected) Problem() problem.Details {
	return problem.Details{
		Type:    "https://shop.example/problems/order-rejected",
		Title:   "Order rejected",
		TitleJA: "注文は受け付けられません",
		Status:  422,
		Code:    "order_rejected",
	}
}

// Execution renders a handler's error by calling Echo's HTTPErrorHandler
// inside its private buffer, before classification and freezing. Install the
// problem adapter there, with the previous handler as the explicit fallback;
// problem middleware inside execution cannot render. Replay sends only
// allowlisted stored headers, so list Content-Type, Content-Language and
// Cache-Control. The frozen bytes replay unchanged after a language change.
// docs:start ExampleMiddlewareConfig_ToMiddleware_problemAdapter
func ExampleMiddlewareConfig_ToMiddleware_problemAdapter() {
	cfg := exampleHTTPConfig()
	cfg.ReplayHeaders = []string{"Content-Type", "Content-Language", "Cache-Control"}
	mw, err := cfg.ToMiddleware()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	fallback := e.HTTPErrorHandler
	install := func(lang problem.Language) {
		h, err := echoadapter.ErrorHandlerConfig{
			Render: problem.RenderConfig{Language: lang},
		}.ToErrorHandler(fallback)
		if err != nil {
			panic(err)
		}
		e.HTTPErrorHandler = h
	}
	install(problem.English)
	e.POST("/write", func(c *echo.Context) error {
		// A business refusal is classified, then frozen with its bytes.
		if err := execution.Classify(c, execution.Failed); err != nil {
			return err
		}
		return orderRejected{}
	}, mw)
	show := func(w *httptest.ResponseRecorder) {
		h := w.Header()
		fmt.Println(w.Code, h.Get("Content-Type"), h.Get("Content-Language"),
			h.Get("Cache-Control"), h.Get("Idempotent-Replayed") == "true")
		var body bytes.Buffer
		if err := json.Indent(&body, w.Body.Bytes(), "", "  "); err != nil {
			panic(err)
		}
		fmt.Print(body.String())
	}
	show(exampleHTTPRequest(context.Background(), e))
	install(problem.Japanese)
	show(exampleHTTPRequest(context.Background(), e))
	// Output:
	// 422 application/problem+json en no-store false
	// {
	//   "type": "https://shop.example/problems/order-rejected",
	//   "title": "Order rejected",
	//   "status": 422,
	//   "code": "order_rejected"
	// }
	// 422 application/problem+json en no-store true
	// {
	//   "type": "https://shop.example/problems/order-rejected",
	//   "title": "Order rejected",
	//   "status": 422,
	//   "code": "order_rejected"
	// }
}

// docs:end ExampleMiddlewareConfig_ToMiddleware_problemAdapter
