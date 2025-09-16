package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Benchmark tests
func BenchmarkFast(b *testing.B) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/fast", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = fast(c)
	}
}

func BenchmarkSlow(b *testing.B) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = slow(c)
	}
}

func BenchmarkMemoryIntensive(b *testing.B) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/memory", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = memoryIntensive(c)
	}
}

func BenchmarkJSONProcessing(b *testing.B) {
	e := echo.New()
	jsonPayload := `{"name":"test","value":123,"data":{"nested":"value"}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/json", strings.NewReader(jsonPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = jsonProcessing(c)
	}
}

func BenchmarkCachedEndpoint(b *testing.B) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cached", nil)

	// Warm up cache
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = cached(c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = cached(c)
	}
}

// Performance tests
func TestResponseTime(t *testing.T) {
	tests := []struct {
		name           string
		handler        echo.HandlerFunc
		needsWarmup    bool
		maxDuration    time.Duration
		description    string
	}{
		{
			name:        "fast endpoint",
			handler:     fast,
			needsWarmup: false,
			maxDuration: 10 * time.Millisecond,
			description: "should respond within 10ms",
		},
		{
			name:        "cached endpoint (after warmup)",
			handler:     cached,
			needsWarmup: true,
			maxDuration: 50 * time.Millisecond,
			description: "should respond quickly from cache",
		},
	}

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Warm up cache if needed
			if tt.needsWarmup {
				req := httptest.NewRequest(http.MethodGet, "/cached", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				_ = cached(c)
			}

			// Measure response time
			start := time.Now()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := tt.handler(c)
			duration := time.Since(start)

			assert.NoError(t, err)
			assert.LessOrEqual(t, duration, tt.maxDuration,
				"Handler %s took %v, expected max %v", tt.name, duration, tt.maxDuration)
		})
	}
}

func TestCacheEffectiveness(t *testing.T) {
	// Clear cache to ensure fresh test state
	cache.Clear()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cached", nil)

	// First call - should be slow (cache miss)
	start1 := time.Now()
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req, rec1)
	err1 := cached(c1)
	duration1 := time.Since(start1)

	assert.NoError(t, err1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	var resp1 Response
	err := json.Unmarshal(rec1.Body.Bytes(), &resp1)
	assert.NoError(t, err)
	assert.Contains(t, resp1.Message, "computed")

	// Second call - should be fast (cache hit)
	start2 := time.Now()
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req, rec2)
	err2 := cached(c2)
	duration2 := time.Since(start2)

	assert.NoError(t, err2)
	assert.Equal(t, http.StatusOK, rec2.Code)

	var resp2 Response
	err = json.Unmarshal(rec2.Body.Bytes(), &resp2)
	assert.NoError(t, err)
	assert.Contains(t, resp2.Message, "cache")

	// Cache hit should be significantly faster
	assert.Less(t, duration2, duration1,
		"Cached response (%v) should be faster than computed response (%v)",
		duration2, duration1)
}

func TestMemoryUsage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/memory", nil)

	// Run memory intensive operation
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := memoryIntensive(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "Memory intensive")

	// Verify response contains expected data
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "sum")
	assert.Contains(t, data, "array_length")
}

func TestMetricsEndpoint(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := metrics(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify metrics structure
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "goroutines")
	assert.Contains(t, data, "mem_alloc_mb")
	assert.Contains(t, data, "num_gc")

	// Verify reasonable values
	goroutines := data["goroutines"].(float64)
	assert.Greater(t, goroutines, 0.0)

	memAlloc := data["mem_alloc_mb"].(float64)
	assert.GreaterOrEqual(t, memAlloc, 0.0)
}

// Example of load testing function
func TestConcurrentRequests(t *testing.T) {
	e := echo.New()
	e.Use(responseTimeMiddleware())
	e.GET("/fast", fast)

	server := httptest.NewServer(e)
	defer server.Close()

	// Test concurrent requests
	concurrency := 10
	requests := 100

	client := &http.Client{Timeout: 5 * time.Second}

	// Channel to collect results
	results := make(chan bool, requests)

	// Launch concurrent requests
	for i := 0; i < requests; i++ {
		go func() {
			resp, err := client.Get(server.URL + "/fast")
			if err != nil {
				results <- false
				return
			}
			defer resp.Body.Close()
			results <- resp.StatusCode == http.StatusOK
		}()

		// Add slight delay every batch to control concurrency
		if (i+1)%concurrency == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Collect results
	successCount := 0
	for i := 0; i < requests; i++ {
		if <-results {
			successCount++
		}
	}

	// All requests should succeed
	assert.Equal(t, requests, successCount)
}