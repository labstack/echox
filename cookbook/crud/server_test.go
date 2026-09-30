package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

func newTestServer(t *testing.T) *echo.Echo {
	t.Helper()
	lock.Lock()
	users = map[int]*user{}
	seq = 1
	lock.Unlock()

	e := newServer()
	e.Logger = slog.New(slog.DiscardHandler)
	return e
}

func request(t *testing.T, e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func decodeUser(t *testing.T, rec *httptest.ResponseRecorder) user {
	t.Helper()
	var u user
	if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil {
		t.Fatalf("decode user: %v (body %q)", err, rec.Body.String())
	}
	return u
}

func expectStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("expected status %d, got %d: %s", want, rec.Code, rec.Body.String())
	}
}

func TestCreateUserIgnoresBodyID(t *testing.T) {
	e := newTestServer(t)
	expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":"first"}`), http.StatusCreated)

	rec := request(t, e, http.MethodPost, "/users", `{"id":1,"name":"second"}`)
	expectStatus(t, rec, http.StatusCreated)
	if got := decodeUser(t, rec); got.ID != 2 {
		t.Fatalf("expected new user to get id 2, got %d", got.ID)
	}

	rec = request(t, e, http.MethodGet, "/users/1", "")
	expectStatus(t, rec, http.StatusOK)
	if got := decodeUser(t, rec); got.Name != "first" {
		t.Fatalf("expected user 1 to keep name %q, got %q", "first", got.Name)
	}
}

func TestUpdateUser(t *testing.T) {
	e := newTestServer(t)
	expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":"before"}`), http.StatusCreated)

	rec := request(t, e, http.MethodPut, "/users/1", `{"name":"after"}`)
	expectStatus(t, rec, http.StatusOK)
	if got := decodeUser(t, rec); got.Name != "after" {
		t.Fatalf("expected updated name %q, got %q", "after", got.Name)
	}

	rec = request(t, e, http.MethodGet, "/users/1", "")
	expectStatus(t, rec, http.StatusOK)
	if got := decodeUser(t, rec); got.Name != "after" {
		t.Fatalf("expected stored name %q, got %q", "after", got.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	e := newTestServer(t)
	expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":"gone"}`), http.StatusCreated)

	expectStatus(t, request(t, e, http.MethodDelete, "/users/1", ""), http.StatusNoContent)
	expectStatus(t, request(t, e, http.MethodGet, "/users/1", ""), http.StatusNotFound)
}

func TestUserNotFoundAndInvalidID(t *testing.T) {
	e := newTestServer(t)
	tests := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodGet, "/users/9999", "", http.StatusNotFound},
		{http.MethodPut, "/users/9999", `{"name":"x"}`, http.StatusNotFound},
		{http.MethodDelete, "/users/9999", "", http.StatusNotFound},
		{http.MethodGet, "/users/abc", "", http.StatusBadRequest},
		{http.MethodPut, "/users/abc", `{"name":"x"}`, http.StatusBadRequest},
		{http.MethodDelete, "/users/abc", "", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			expectStatus(t, request(t, e, tt.method, tt.path, tt.body), tt.want)
		})
	}
}
