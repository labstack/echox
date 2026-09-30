package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/labstack/echo/v5"
)

func newTestServer(t *testing.T) *echo.Echo {
	t.Helper()
	lock.Lock()
	users = map[int]user{}
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

func TestListUsers(t *testing.T) {
	e := newTestServer(t)
	rec := request(t, e, http.MethodGet, "/users", "")
	expectStatus(t, rec, http.StatusOK)
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Fatalf("expected empty list [], got %s", got)
	}

	for _, name := range []string{"a", "b", "c"} {
		expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":"`+name+`"}`), http.StatusCreated)
	}

	rec = request(t, e, http.MethodGet, "/users", "")
	expectStatus(t, rec, http.StatusOK)
	var list []user
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode users: %v (body %q)", err, rec.Body.String())
	}
	if len(list) != 3 || list[0].ID != 1 || list[1].ID != 2 || list[2].ID != 3 {
		t.Fatalf("expected users 1, 2, 3 in order, got %+v", list)
	}
}

func TestErrors(t *testing.T) {
	e := newTestServer(t)
	expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":"keep"}`), http.StatusCreated)

	tests := []struct {
		method, path, body string
		status             int
		message            string
	}{
		{http.MethodGet, "/users/9999", "", http.StatusNotFound, "user not found"},
		{http.MethodPut, "/users/9999", `{"name":"x"}`, http.StatusNotFound, "user not found"},
		{http.MethodDelete, "/users/9999", "", http.StatusNotFound, "user not found"},
		{http.MethodGet, "/users/abc", "", http.StatusBadRequest, "invalid user id"},
		{http.MethodPut, "/users/abc", `{"name":"x"}`, http.StatusBadRequest, "invalid user id"},
		{http.MethodDelete, "/users/abc", "", http.StatusBadRequest, "invalid user id"},
		{http.MethodPost, "/users", `{}`, http.StatusBadRequest, "name is required"},
		{http.MethodPut, "/users/1", "", http.StatusBadRequest, "name is required"},
		{http.MethodPut, "/users/1", `{}`, http.StatusBadRequest, "name is required"},
		{http.MethodPut, "/users/1", `{"name":"  "}`, http.StatusBadRequest, "name is required"},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path+" "+tt.body, func(t *testing.T) {
			rec := request(t, e, tt.method, tt.path, tt.body)
			expectStatus(t, rec, tt.status)
			var got struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || got.Message != tt.message {
				t.Fatalf("expected message %q, got body %q", tt.message, rec.Body.String())
			}
		})
	}

	// malformed JSON is rejected by Bind
	expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":`), http.StatusBadRequest)

	// the failed requests left user 1 unchanged
	rec := request(t, e, http.MethodGet, "/users/1", "")
	expectStatus(t, rec, http.StatusOK)
	if got := decodeUser(t, rec); got.Name != "keep" {
		t.Fatalf("expected stored name %q, got %q", "keep", got.Name)
	}
}

func TestConcurrentRequests(t *testing.T) {
	// Run with -race: handlers read and write the store from many goroutines at once.
	e := newTestServer(t)
	expectStatus(t, request(t, e, http.MethodPost, "/users", `{"name":"base"}`), http.StatusCreated)

	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			request(t, e, http.MethodPost, "/users", `{"name":"u"}`)
			request(t, e, http.MethodPut, "/users/1", `{"name":"n`+strconv.Itoa(i)+`"}`)
			request(t, e, http.MethodGet, "/users/1", "")
			request(t, e, http.MethodGet, "/users", "")
			request(t, e, http.MethodDelete, "/users/"+strconv.Itoa(i+2), "")
		})
	}
	wg.Wait()

	expectStatus(t, request(t, e, http.MethodGet, "/users/1", ""), http.StatusOK)
}
