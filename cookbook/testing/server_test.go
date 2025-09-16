package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, health(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	userHandler := NewUserHandler()

	// Assertions
	if assert.NoError(t, userHandler.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var user User
		err := json.Unmarshal(rec.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	userHandler := NewUserHandler()

	// Assertions
	err := userHandler.GetUser(c)
	assert.Error(t, err)

	// Check if it's an HTTP error with 404 status
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestCreateUser(t *testing.T) {
	// Setup
	userJSON := `{"name":"Jane Doe","email":"jane@example.com"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userHandler := NewUserHandler()

	// Assertions
	if assert.NoError(t, userHandler.CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var user User
		err := json.Unmarshal(rec.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID) // First user gets ID 1
		assert.Equal(t, "Jane Doe", user.Name)
		assert.Equal(t, "jane@example.com", user.Email)
	}
}

func TestCreateUserValidation(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		wantCode int
	}{
		{
			name:     "missing name",
			payload:  `{"email":"test@example.com"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing email",
			payload:  `{"name":"Test User"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid JSON",
			payload:  `{"name":"Test User"`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			userHandler := NewUserHandler()

			// Assertions
			err := userHandler.CreateUser(c)
			assert.Error(t, err)

			httpErr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, tt.wantCode, httpErr.Code)
		})
	}
}

func TestGetUsers(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userHandler := NewUserHandler()

	// Add a test user
	userHandler.users[1] = &User{ID: 1, Name: "Test User", Email: "test@example.com"}

	// Assertions
	if assert.NoError(t, userHandler.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var users []*User
		err := json.Unmarshal(rec.Body.Bytes(), &users)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "Test User", users[0].Name)
	}
}

// Integration test example
func TestUserAPIIntegration(t *testing.T) {
	// Setup Echo server
	e := echo.New()
	userHandler := NewUserHandler()

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUsers)
	e.GET("/users/:id", userHandler.GetUser)

	// Test creating a user
	userJSON := `{"name":"Integration Test","email":"integration@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	// Test getting all users
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var users []*User
	err := json.Unmarshal(rec.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Integration Test", users[0].Name)
}