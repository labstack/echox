package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name         string
		username     string
		password     string
		expectedCode int
		expectToken  bool
	}{
		{
			name:         "valid credentials",
			username:     "jon",
			password:     "shhh!",
			expectedCode: http.StatusOK,
			expectToken:  true,
		},
		{
			name:         "invalid username",
			username:     "invalid",
			password:     "shhh!",
			expectedCode: http.StatusUnauthorized,
			expectToken:  false,
		},
		{
			name:         "invalid password",
			username:     "jon",
			password:     "invalid",
			expectedCode: http.StatusUnauthorized,
			expectToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create form data
			formData := "username=" + tt.username + "&password=" + tt.password
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := login(c)

			if tt.expectedCode == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)

				if tt.expectToken {
					var response map[string]interface{}
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					require.NoError(t, err)
					assert.Contains(t, response, "token")
					assert.NotEmpty(t, response["token"])
				}
			} else {
				require.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				require.True(t, ok)
				assert.Equal(t, tt.expectedCode, httpErr.Code)
			}
		})
	}
}

func TestJWTTokenGeneration(t *testing.T) {
	// Set environment variable for test
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	e := echo.New()
	formData := "username=jon&password=shhh!"
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := login(c)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	tokenString := response["token"].(string)

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(*jwtCustomClaims)
	require.True(t, ok)
	assert.Equal(t, "Jon Snow", claims.Name)
	assert.True(t, claims.Admin)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestAccessible(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := accessible(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Accessible", rec.Body.String())
}

func TestRestricted(t *testing.T) {
	// Create a valid token
	claims := &jwtCustomClaims{
		"Test User",
		false,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set the token in context (simulating JWT middleware)
	c.Set("user", token)

	err := restricted(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Welcome Test User!", rec.Body.String())
}

func TestJWTIntegration(t *testing.T) {
	// Set consistent secret for both login and verification
	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Setup Echo server with all routes
	e := echo.New()
	e.HideBanner = true

	e.POST("/login", login)
	e.GET("/", accessible)

	// Restricted group with JWT middleware
	r := e.Group("/restricted")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("", restricted)

	// Test login and get token
	formData := "username=jon&password=shhh!"
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var loginResponse map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &loginResponse)
	require.NoError(t, err)

	token := loginResponse["token"].(string)

	// Test accessing restricted endpoint with token
	req = httptest.NewRequest(http.MethodGet, "/restricted", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Welcome Jon Snow!", rec.Body.String())

	// Test accessing restricted endpoint without token
	req = httptest.NewRequest(http.MethodGet, "/restricted", nil)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestEnvironmentVariableSecret(t *testing.T) {
	// Test with custom secret
	os.Setenv("JWT_SECRET", "custom-secret")
	defer os.Unsetenv("JWT_SECRET")

	e := echo.New()
	formData := "username=jon&password=shhh!"
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := login(c)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	tokenString := response["token"].(string)

	// Verify token was signed with custom secret
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("custom-secret"), nil
	})

	require.NoError(t, err)
	require.True(t, token.Valid)

	// Test with default secret should fail
	_, err = jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	assert.Error(t, err)
}