---
description: Testing Echo applications
---

# Testing

This example demonstrates comprehensive testing strategies for Echo applications, including unit tests, integration tests, and best practices for testing HTTP handlers.

## Key Features

- **Unit Testing**: Test individual handlers in isolation
- **Integration Testing**: Test complete request/response cycles
- **Validation Testing**: Test input validation and error handling
- **HTTP Testing**: Use `httptest` package for testing HTTP functionality
- **Table-Driven Tests**: Efficient testing of multiple scenarios

## Server Implementation

```go reference
https://github.com/labstack/echox/blob/master/cookbook/testing/server.go
```

## Test Suite

```go reference
https://github.com/labstack/echox/blob/master/cookbook/testing/server_test.go
```

## Testing Patterns

### 1. Unit Testing Handlers

Test individual handlers by creating Echo context manually:

```go
func TestHealth(t *testing.T) {
    // Setup
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Execute
    if assert.NoError(t, health(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        // Additional assertions...
    }
}
```

### 2. Testing with Path Parameters

Set path parameters using `SetParamNames` and `SetParamValues`:

```go
c.SetParamNames("id")
c.SetParamValues("1")
```

### 3. Testing JSON Requests

Create requests with JSON payloads and proper content type:

```go
userJSON := `{"name":"Jane Doe","email":"jane@example.com"}`
req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
```

### 4. Table-Driven Tests

Test multiple scenarios efficiently:

```go
tests := []struct {
    name     string
    payload  string
    wantCode int
}{
    {"missing name", `{"email":"test@example.com"}`, http.StatusBadRequest},
    {"missing email", `{"name":"Test User"}`, http.StatusBadRequest},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### 5. Integration Testing

Test the complete server using `ServeHTTP`:

```go
e.ServeHTTP(rec, req)
```

## Running Tests

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run specific test
go test -run TestHealth

# Run tests with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Testing Best Practices

1. **Isolation**: Each test should be independent and not rely on other tests
2. **Setup/Teardown**: Use proper setup and cleanup for each test
3. **Clear Names**: Use descriptive test names that explain what is being tested
4. **Edge Cases**: Test both success and failure scenarios
5. **Coverage**: Aim for high test coverage of critical paths
6. **Fast Tests**: Keep tests fast by avoiding external dependencies when possible

## Dependencies

Add testing dependencies to your `go.mod`:

```go
require (
    github.com/stretchr/testify v1.8.4
)
```

The `testify` package provides helpful assertion functions that make tests more readable and maintainable.