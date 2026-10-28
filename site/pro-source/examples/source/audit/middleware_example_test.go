// SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
// SPDX-FileCopyrightText: © 2026 LabStack LLC

package audit_test

import (
	"context"
	"fmt"
	"net/http/httptest"

	"github.com/labstack/echo/v5"

	"echo.labstack.com/pro/audit"
)

// docs:start ExampleMiddlewareConfig_ToMiddleware
func ExampleMiddlewareConfig_ToMiddleware() {
	e := echo.New()
	// docs:start site-anchor
	mw, err := (audit.MiddlewareConfig{
		Actor: func(*echo.Context) (audit.Actor, error) {
			return audit.Actor{Kind: audit.System}, nil
		},
	}).ToMiddleware()
	if err != nil {
		fmt.Println(err)
		return
	}
	e.Use(mw)
	e.GET("/jobs/:id", func(c *echo.Context) error {
		scope, err := audit.ScopeFromContext(c.Request().Context())
		if err != nil {
			return err
		}
		fmt.Println(scope.Actor().Kind, scope.Source().Route)
		return nil
	})
	// docs:end site-anchor
	request := httptest.NewRequestWithContext(
		context.Background(),
		"GET",
		"/jobs/42",
		nil,
	)
	request.RemoteAddr = "192.0.2.1:1234"
	e.ServeHTTP(httptest.NewRecorder(), request)
	// Output: system /jobs/:id
}

// docs:end ExampleMiddlewareConfig_ToMiddleware
