package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Response represents a JSON response
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

// PerformanceMetrics holds basic performance data
type PerformanceMetrics struct {
	Goroutines   int           `json:"goroutines"`
	MemAllocMB   float64       `json:"mem_alloc_mb"`
	MemSysMB     float64       `json:"mem_sys_mb"`
	NumGC        uint32        `json:"num_gc"`
	ResponseTime time.Duration `json:"response_time_ms"`
}

// Simple in-memory cache for demonstration
type Cache struct {
	data map[string]any
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]any)}
}

func (c *Cache) Get(key string) (any, bool) {
	val, exists := c.data[key]
	return val, exists
}

func (c *Cache) Set(key string, value any) {
	c.data[key] = value
}

var cache = NewCache()

// Fast endpoint - minimal processing
func fast(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Message:   "Fast response",
		Timestamp: time.Now(),
	})
}

// Slow endpoint - simulates heavy processing
func slow(c echo.Context) error {
	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	// Simulate some CPU work
	result := 0
	for i := 0; i < 1000000; i++ {
		result += i
	}

	return c.JSON(http.StatusOK, Response{
		Message:   "Slow response with processing",
		Timestamp: time.Now(),
		Data:      map[string]int{"calculation_result": result},
	})
}

// Memory intensive endpoint
func memoryIntensive(c echo.Context) error {
	// Create a large slice to demonstrate memory usage
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i
	}

	// Sum to prevent optimization
	sum := 0
	for _, v := range data {
		sum += v
	}

	return c.JSON(http.StatusOK, Response{
		Message:   "Memory intensive operation completed",
		Timestamp: time.Now(),
		Data:      map[string]int{"sum": sum, "array_length": len(data)},
	})
}

// Cached endpoint - demonstrates caching for performance
func cached(c echo.Context) error {
	key := "expensive_data"

	// Check cache first
	if cached, exists := cache.Get(key); exists {
		return c.JSON(http.StatusOK, Response{
			Message:   "Data from cache",
			Timestamp: time.Now(),
			Data:      cached,
		})
	}

	// Simulate expensive operation
	time.Sleep(200 * time.Millisecond)
	expensiveData := map[string]any{
		"computed_at": time.Now(),
		"value":       "expensive computation result",
		"pi":          3.14159,
	}

	// Store in cache
	cache.Set(key, expensiveData)

	return c.JSON(http.StatusOK, Response{
		Message:   "Data computed and cached",
		Timestamp: time.Now(),
		Data:      expensiveData,
	})
}

// JSON processing endpoint
func jsonProcessing(c echo.Context) error {
	var input map[string]any
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	// Process the JSON (transform, validate, etc.)
	processed := make(map[string]any)
	for k, v := range input {
		processed[fmt.Sprintf("processed_%s", k)] = v
	}

	// Add metadata
	processed["processing_time"] = time.Now()
	processed["input_keys"] = len(input)

	return c.JSON(http.StatusOK, Response{
		Message:   "JSON processed successfully",
		Timestamp: time.Now(),
		Data:      processed,
	})
}

// Performance metrics endpoint
func metrics(c echo.Context) error {
	start := time.Now()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := PerformanceMetrics{
		Goroutines:   runtime.NumGoroutine(),
		MemAllocMB:   float64(m.Alloc) / 1024 / 1024,
		MemSysMB:     float64(m.Sys) / 1024 / 1024,
		NumGC:        m.NumGC,
		ResponseTime: time.Since(start),
	}

	return c.JSON(http.StatusOK, Response{
		Message:   "Performance metrics",
		Timestamp: time.Now(),
		Data:      metrics,
	})
}

// Custom middleware to measure response time
func responseTimeMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			responseTime := time.Since(start)
			c.Response().Header().Set("X-Response-Time", responseTime.String())

			return err
		}
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(responseTimeMiddleware())

	// Disable Echo's default Server header for performance
	e.HideBanner = true

	// Performance-related routes
	e.GET("/fast", fast)
	e.GET("/slow", slow)
	e.GET("/memory", memoryIntensive)
	e.GET("/cached", cached)
	e.POST("/json", jsonProcessing)
	e.GET("/metrics", metrics)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Start server
	fmt.Printf("Server starting on :1323\n")
	fmt.Printf("Try these endpoints:\n")
	fmt.Printf("  GET  /fast      - Fast endpoint\n")
	fmt.Printf("  GET  /slow      - Slow endpoint (100ms delay)\n")
	fmt.Printf("  GET  /memory    - Memory intensive\n")
	fmt.Printf("  GET  /cached    - Cached endpoint\n")
	fmt.Printf("  POST /json      - JSON processing\n")
	fmt.Printf("  GET  /metrics   - Performance metrics\n")
	fmt.Printf("  GET  /health    - Health check\n")

	e.Logger.Fatal(e.Start(":1323"))
}