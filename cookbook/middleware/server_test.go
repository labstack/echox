package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServerHeaderMiddleware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Handler that uses the middleware
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Apply middleware
	middlewareHandler := ServerHeader(handler)

	// Execute
	if assert.NoError(t, middlewareHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Echo/4.11", rec.Header().Get(echo.HeaderServer))
	}
}

func TestStatsMiddleware(t *testing.T) {
	// Setup
	stats := NewStats()
	e := echo.New()

	// Create a simple handler
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Apply stats middleware
	middlewareHandler := stats.Process(handler)

	// Test multiple requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		assert.NoError(t, middlewareHandler(c))
	}

	// Check stats
	assert.Equal(t, uint64(5), stats.RequestCount)
	assert.Equal(t, 5, stats.Statuses["200"])
}

func TestStatsHandler(t *testing.T) {
	// Setup
	stats := NewStats()
	e := echo.New()

	// Add some test data
	stats.RequestCount = 10
	stats.Statuses["200"] = 8
	stats.Statuses["404"] = 2

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	if assert.NoError(t, stats.Handle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response Stats
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), response.RequestCount)
		assert.Equal(t, 8, response.Statuses["200"])
		assert.Equal(t, 2, response.Statuses["404"])
	}
}

func TestMainRoutes(t *testing.T) {
	// Setup Echo server with all routes
	e := echo.New()

	// Stats
	s := NewStats()
	e.Use(s.Process)
	e.GET("/stats", s.Handle)

	// Server header
	e.Use(ServerHeader)

	// Handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedCode   int
		expectedHeader string
	}{
		{
			name:           "root endpoint",
			method:         http.MethodGet,
			path:           "/",
			expectedCode:   http.StatusOK,
			expectedHeader: "Echo/4.11",
		},
		{
			name:         "stats endpoint",
			method:       http.MethodGet,
			path:         "/stats",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, rec.Header().Get(echo.HeaderServer))
			}
		})
	}
}

func BenchmarkStatsMiddleware(b *testing.B) {
	stats := NewStats()
	e := echo.New()

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	middlewareHandler := stats.Process(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = middlewareHandler(c)
	}
}