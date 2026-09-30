package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	customServer := flag.Bool("custom-server", false, "use a standard http.Server")
	flag.Parse()
	e := echo.New()
	// docs:start handler
	e.GET("/request", func(c *echo.Context) error {
		req := c.Request()
		format := `
			<code>
				Protocol: %s<br>
				Host: %s<br>
				Remote Address: %s<br>
				Method: %s<br>
				Path: %s<br>
			</code>
		`
		return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
	})
	// docs:end handler
	if *customServer {
		customHTTPServer(e)
		return
	}
	// docs:start start-tls
	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.StartTLS(context.Background(), e, "cert.pem", "key.pem"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
	// docs:end start-tls
}

// customHTTPServer demonstrates configuring TLS with a standard HTTP server.
func customHTTPServer(e *echo.Echo) {
	// docs:start custom-server
	s := http.Server{
		Addr:    ":1323",
		Handler: e,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		// ReadTimeout: 30 * time.Second, // use custom timeouts
	}
	if err := s.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Error("failed to start server", "error", err)
	}
	// docs:end custom-server
}
