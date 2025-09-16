---
description: Performance testing and optimization
---

# Performance & Benchmarking

This example demonstrates performance testing, benchmarking, and optimization techniques for Echo applications, including response time measurement, memory profiling, and load testing strategies.

## Key Features

- **Benchmark Testing**: Go's built-in benchmarking for performance measurement
- **Response Time Tracking**: Custom middleware to measure request duration
- **Memory Profiling**: Runtime memory usage monitoring
- **Caching Strategies**: In-memory caching for improved performance
- **Load Testing**: Concurrent request handling verification
- **Performance Metrics**: Real-time application performance data

## Server Implementation

```go reference
https://github.com/labstack/echox/blob/master/cookbook/performance/server.go
```

## Benchmark Tests

```go reference
https://github.com/labstack/echox/blob/master/cookbook/performance/server_test.go
```

## Performance Testing Patterns

### 1. Benchmark Testing

Use Go's built-in benchmarking to measure handler performance:

```go
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
```

### 2. Response Time Middleware

Track response times with custom middleware:

```go
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
```

### 3. Cache Performance Testing

Verify cache effectiveness by comparing response times:

```go
// First call - cache miss (slower)
duration1 := measureResponseTime(handler)

// Second call - cache hit (faster)
duration2 := measureResponseTime(handler)

assert.Less(t, duration2, duration1)
```

### 4. Memory Usage Monitoring

Monitor runtime memory statistics:

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)

metrics := PerformanceMetrics{
    MemAllocMB: float64(m.Alloc) / 1024 / 1024,
    MemSysMB:   float64(m.Sys) / 1024 / 1024,
    NumGC:      m.NumGC,
}
```

### 5. Concurrent Load Testing

Test application under concurrent load:

```go
func TestConcurrentRequests(t *testing.T) {
    server := httptest.NewServer(e)
    defer server.Close()

    concurrency := 10
    requests := 100

    for i := 0; i < requests; i++ {
        go func() {
            resp, err := client.Get(server.URL + "/endpoint")
            // Handle response
        }()
    }
}
```

## Running Performance Tests

### Benchmark Tests
```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkFast

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem

# Run benchmarks multiple times for better accuracy
go test -bench=. -count=5

# CPU profiling during benchmarks
go test -bench=. -cpuprofile=cpu.prof

# Memory profiling during benchmarks
go test -bench=. -memprofile=mem.prof
```

### Performance Tests
```bash
# Run performance tests
go test -run TestResponseTime

# Run with verbose output
go test -v -run TestPerformance

# Run load tests
go test -run TestConcurrentRequests
```

### Profiling
```bash
# CPU profiling
go tool pprof cpu.prof

# Memory profiling
go tool pprof mem.prof

# Heap profiling of running server
go tool pprof http://localhost:1323/debug/pprof/heap
```

## Performance Optimization Tips

### 1. Reduce Allocations
- Reuse objects and slices when possible
- Use sync.Pool for expensive-to-create objects
- Minimize string concatenation in hot paths

### 2. Efficient JSON Handling
- Use json.Decoder for streaming large payloads
- Consider alternative JSON libraries like jsoniter
- Avoid unnecessary marshaling/unmarshaling

### 3. Middleware Optimization
- Order middleware by frequency of use
- Use conditional middleware when appropriate
- Avoid expensive operations in frequently called middleware

### 4. Caching Strategies
- Implement response caching for expensive operations
- Use appropriate cache TTL values
- Consider distributed caching for scaled applications

### 5. Database Optimization
- Use connection pooling
- Implement query result caching
- Optimize database queries and indexes

## Load Testing Tools

### External Tools
```bash
# Apache Bench
ab -n 1000 -c 10 http://localhost:1323/fast

# wrk
wrk -t4 -c100 -d30s http://localhost:1323/fast

# Vegeta
vegeta attack -duration=30s -rate=100 | vegeta report
```

### Monitoring Endpoints

The example includes these monitoring endpoints:
- `GET /metrics` - Runtime performance metrics
- `GET /health` - Basic health check
- All endpoints include `X-Response-Time` header

## Best Practices

1. **Baseline First**: Establish performance baselines before optimization
2. **Profile Before Optimizing**: Use profiling to identify actual bottlenecks
3. **Test Realistic Scenarios**: Use realistic data sizes and request patterns
4. **Monitor in Production**: Implement monitoring and alerting
5. **Gradual Optimization**: Make incremental improvements and measure impact
6. **Documentation**: Document performance requirements and test results