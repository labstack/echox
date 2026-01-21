---
description: CSRF middleware
---

# CSRF

Cross-site request forgery, also known as one-click attack or session riding and
abbreviated as CSRF (sometimes pronounced sea-surf) or XSRF, is a type of malicious
exploit of a website where unauthorized commands are transmitted from a user that
the website trusts.

## Usage

```go
e.Use(middleware.CSRF())
```

## Custom Configuration

The CSRF middleware supports the [**Sec-Fetch-Site**](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site) header as a modern, defense-in-depth approach to [CSRF
protection](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#fetch-metadata-headers), implementing the OWASP-recommended Fetch Metadata API alongside the traditional token-based mechanism.

**How it works:**

Modern browsers automatically send the `Sec-Fetch-Site` header with all requests, indicating the relationship
between the request origin and the target. The middleware uses this to make security decisions:

- **`same-origin`** or **`none`**: Requests are allowed (exact origin match or direct user navigation)
- **`same-site`**: Falls back to token validation (e.g., subdomain to main domain)
- **`cross-site`**: Blocked by default with 403 error for unsafe methods (POST, PUT, DELETE, PATCH)

For browsers that don't send this header (older browsers), the middleware seamlessly falls back to
traditional token-based CSRF protection.

For `Sec-Fetch-Site` usage follow the configuration options:
- `TrustedOrigins []string`: Allowlist specific origins for cross-site requests (useful for OAuth callbacks, webhooks)
- `AllowSecFetchSiteFunc func(echo.Context) (bool, error)`: Custom logic for same-site/cross-site request validation

```go
e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
  // Allow OAuth callbacks from trusted provider
  TrustedOrigins: []string{"https://oauth-provider.com"},

  // Custom validation for same-site requests
  AllowSecFetchSiteFunc: func(c *echo.Context) (bool, error) {
      // Your custom authorization logic here
      return validateCustomAuth(c), nil
      // return true, err  // blocks request with error
      // return true, nil  // allows CSRF request through
      // return false, nil // falls back to legacy token logic
  },
}))
```

### Usage with token based CSRF protection

```go
e := echo.New()
e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
  TokenLookup: "header:X-XSRF-TOKEN",
}))
```

Example above uses `X-XSRF-TOKEN` request header to extract CSRF token.

*Example Configuration that reads token from Form*

```go
middleware.CSRFWithConfig(middleware.CSRFConfig{
	TokenLookup:    "form:_csrf",
	CookiePath:     "/",
	CookieDomain:   "example.com",
	CookieSecure:   true,
	CookieHTTPOnly: true,
	CookieSameSite: http.SameSiteStrictMode,
})
```

## Accessing CSRF Token

### Server-side

CSRF token can be accessed from `Echo#Context` using `ContextKey` and passed to
the client via template.

### Client-side

CSRF token can be accessed from CSRF cookie.

## Configuration

```go
type CSRFConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper
	// TrustedOrigin permits any request with `Sec-Fetch-Site` header whose `Origin` header
	// exactly matches the specified value.
	// Values should be formated as Origin header "scheme://host[:port]".
	//
	// See [Origin]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin
	// See [Sec-Fetch-Site]: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#fetch-metadata-headers
	TrustedOrigins []string

	// AllowSecFetchSameSite allows custom behaviour for `Sec-Fetch-Site` requests that are about to
	// fail with CRSF error, to be allowed or replaced with custom error.
	// This function applies to `Sec-Fetch-Site` values:
	// - `same-site` 		same registrable domain (subdomain and/or different port)
	// - `cross-site`		request originates from different site
	// See [Sec-Fetch-Site]: https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#fetch-metadata-headers
	AllowSecFetchSiteFunc func(c *echo.Context) (bool, error)

	// TokenLength is the length of the generated token.
	TokenLength uint8
	// Optional. Default value 32.

	// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:X-CSRF-Token".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// - "query:<name>"
	// - "form:<name>"
	// Multiple sources example:
	// - "header:X-CSRF-Token,query:csrf"
	TokenLookup string `yaml:"token_lookup"`

	// Generator defines a function to generate token.
	// Optional. Defaults tp randomString(TokenLength).
	Generator func() string

	// Context key to store generated CSRF token into context.
	// Optional. Default value "csrf".
	ContextKey string

	// Name of the CSRF cookie. This cookie will store CSRF token.
	// Optional. Default value "csrf".
	CookieName string

	// Domain of the CSRF cookie.
	// Optional. Default value none.
	CookieDomain string

	// Path of the CSRF cookie.
	// Optional. Default value none.
	CookiePath string

	// Max age (in seconds) of the CSRF cookie.
	// Optional. Default value 86400 (24hr).
	CookieMaxAge int

	// Indicates if CSRF cookie is secure.
	// Optional. Default value false.
	CookieSecure bool

	// Indicates if CSRF cookie is HTTP only.
	// Optional. Default value false.
	CookieHTTPOnly bool

	// Indicates SameSite mode of the CSRF cookie.
	// Optional. Default value SameSiteDefaultMode.
	CookieSameSite http.SameSite

	// ErrorHandler defines a function which is executed for returning custom errors.
	ErrorHandler func(c *echo.Context, err error) error
}
```

### Default Configuration

```go
var DefaultCSRFConfig = CSRFConfig{
	Skipper:        DefaultSkipper,
	TokenLength:    32,
	TokenLookup:    "header:" + echo.HeaderXCSRFToken,
	ContextKey:     "csrf",
	CookieName:     "_csrf",
	CookieMaxAge:   86400,
	CookieSameSite: http.SameSiteDefaultMode,
}
```
