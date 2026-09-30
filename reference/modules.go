// SPDX-License-Identifier: MIT

// Package reference keeps the middleware packages used by the documentation's
// API extractor in the Go module graph. CI compiles them with the examples, and
// Dependabot updates their versions in reference/go.mod.
package reference

import (
	_ "github.com/labstack/echo-jwt/v5"
	_ "github.com/labstack/echo-otel/v5"
	_ "github.com/labstack/echo-prometheus"
)
