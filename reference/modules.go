//go:build docs_external

// SPDX-License-Identifier: MIT

// Package reference keeps the middleware packages used by the documentation's
// API extractor in the Go module graph. Source preparation enables docs_external
// to compile them against stable Echo. Next and proposed Echo workspaces compile
// only the Echo examples and extractor; Dependabot tracks all module versions.
package reference

import (
	_ "github.com/labstack/echo-jwt/v5"
	_ "github.com/labstack/echo-otel/v5"
	_ "github.com/labstack/echo-prometheus"
)
