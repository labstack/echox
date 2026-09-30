package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
