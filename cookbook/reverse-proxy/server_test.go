package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func TestGroupedProxy(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, r.URL.Path)
	}))
	defer upstream.Close()
	target, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	configureBlogProxy(e, []*middleware.ProxyTarget{{URL: target}})
	for _, path := range []string{"/blog", "/blog?draft=1"} {
		response := httptest.NewRecorder()
		e.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		want := "/blog/" + path[len("/blog"):]
		if response.Code != http.StatusPermanentRedirect || response.Header().Get("Location") != want {
			t.Fatalf("%s: got %d %q, want 308 %q", path, response.Code, response.Header().Get("Location"), want)
		}
		canonical, err := url.Parse(want)
		if err != nil {
			t.Fatal(err)
		}
		if socketPath := canonical.ResolveReference(&url.URL{Path: "ws"}).Path; socketPath != "/blog/ws" {
			t.Fatalf("relative WebSocket URL: got %q, want /blog/ws", socketPath)
		}
		followed := httptest.NewRecorder()
		e.ServeHTTP(followed, httptest.NewRequest(http.MethodGet, want, nil))
		if followed.Code != http.StatusOK || followed.Body.String() != "/" {
			t.Fatalf("%s: redirected request got %d %q, want 200 /", path, followed.Code, followed.Body.String())
		}
	}
	for _, path := range []string{"/blog/", "/blog/ws"} {
		response := httptest.NewRecorder()
		e.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		want := path[len("/blog"):]
		if response.Code != http.StatusOK || response.Body.String() != want {
			t.Fatalf("%s: got %d %q, want 200 %q", path, response.Code, response.Body.String(), want)
		}
	}
	response := httptest.NewRecorder()
	e.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/outside", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("outside group: got %d, want 404", response.Code)
	}
}

func TestGroupedProxyRedirectAllMethods(t *testing.T) {
	type upstreamRequest struct {
		method, path, query, body string
	}
	requests := make(chan upstreamRequest, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		requests <- upstreamRequest{r.Method, r.URL.Path, r.URL.RawQuery, string(body)}
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()
	target, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	configureBlogProxy(e, []*middleware.ProxyTarget{{URL: target}})
	for _, method := range []string{
		http.MethodGet, http.MethodHead, http.MethodPost, http.MethodOptions,
		http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodTrace,
	} {
		body := ""
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			body = "name=Echo"
		}
		response := httptest.NewRecorder()
		e.ServeHTTP(response, httptest.NewRequest(method, "/blog?draft=1", strings.NewReader(body)))
		location := response.Header().Get("Location")
		if response.Code != http.StatusPermanentRedirect || location != "/blog/?draft=1" {
			t.Fatalf("%s /blog: got %d %q, want 308 /blog/?draft=1", method, response.Code, location)
		}
		// Follow the 308 with the same method and body, including form POSTs.
		followed := httptest.NewRecorder()
		e.ServeHTTP(followed, httptest.NewRequest(method, location, strings.NewReader(body)))
		if followed.Code != http.StatusOK {
			t.Fatalf("%s /blog/: got %d, want 200", method, followed.Code)
		}
		if got := <-requests; got != (upstreamRequest{method, "/", "draft=1", body}) {
			t.Fatalf("%s redirected upstream request: got %+v", method, got)
		}
	}
}
