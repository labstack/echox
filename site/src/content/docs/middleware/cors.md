---
title: CORS
description: Cross-Origin Resource Sharing middleware for secure cross-domain access control.
sidebar:
  order: 6
---

CORS middleware implements the [CORS](https://fetch.spec.whatwg.org/#http-cors-protocol) specification. CORS gives
web servers cross-domain access controls, which enable secure cross-domain data transfers.

All core middleware lives in the `middleware` package:

```go
import "github.com/labstack/echo/v5/middleware"
```

## Usage

```go
e.Use(middleware.CORS("https://example.com", "https://subdomain.example.com"))
```

## Custom configuration

```go
e := echo.New()
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
}))
```

## Configuration

```go
type CORSConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	// AllowOrigins determines the value of the Access-Control-Allow-Origin response
	// header, defining the list of origins that may access the resource.
	//
	// An origin consists of: scheme + "://" + host + optional ":" + port.
	// A wildcard may be used, but it must be set explicitly as []string{"*"}.
	// Example: `https://example.com`, `http://example.com:8080`, `*`.
	//
	// Security: use extreme caution when handling the origin and carefully validate any
	// logic. Attackers may register hostile domain names. See
	// https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
	//
	// Mandatory.
	AllowOrigins []string

	// UnsafeAllowOriginFunc is an optional custom function to validate the origin. It
	// takes the origin and returns the allowed origin, whether it is allowed, and an
	// error (returned immediately by the handler). If set, AllowOrigins is ignored.
	//
	// Security: use extreme caution when handling the origin. Attackers may register
	// hostile (sub)domain names.
	//
	// Sub-domain check example:
	//	UnsafeAllowOriginFunc: func(c *echo.Context, origin string) (string, bool, error) {
	//		if strings.HasSuffix(origin, ".example.com") {
	//			return origin, true, nil
	//		}
	//		return "", false, nil
	//	}
	//
	// Optional.
	UnsafeAllowOriginFunc func(c *echo.Context, origin string) (allowedOrigin string, allowed bool, err error)

	// AllowMethods determines the value of the Access-Control-Allow-Methods response
	// header, used in response to a preflight request.
	//
	// Optional. Defaults to GET, HEAD, PUT, PATCH, POST, DELETE. If left empty, the
	// middleware fills the preflight Access-Control-Allow-Methods header from the
	// `Allow` header that the router set into the context.
	AllowMethods []string

	// AllowHeaders determines the value of the Access-Control-Allow-Headers response
	// header, indicating which HTTP headers can be used in the actual request.
	//
	// Optional. Defaults to an empty list.
	AllowHeaders []string

	// AllowCredentials determines the value of the Access-Control-Allow-Credentials
	// response header, indicating whether the response can be exposed when the
	// credentials mode is true.
	//
	// Optional. Default value false, in which case the header is not set.
	//
	// Security: avoid using AllowCredentials = true together with AllowOrigins = *.
	AllowCredentials bool

	// ExposeHeaders determines the value of Access-Control-Expose-Headers, the list of
	// headers clients are allowed to access.
	//
	// Optional. Default value []string{}, in which case the header is not set.
	ExposeHeaders []string

	// MaxAge determines the value of the Access-Control-Max-Age response header, how long
	// (in seconds) the results of a preflight request can be cached. The header is set
	// only if MaxAge != 0; a negative value sends "0", instructing browsers not to cache.
	//
	// Optional. Default value 0 — the header is not sent.
	MaxAge int
}
```

### Default configuration

```go
// Effective defaults applied when fields are left unset.
CORSConfig{
	Skipper:      DefaultSkipper,
	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}
```

## Security

A wildcard origin (`AllowOrigins: []string{"*"}`) combined with `AllowCredentials: true`
is dangerous: it would reflect **any** request's `Origin` back in
`Access-Control-Allow-Origin`, letting a page on any site make credentialed cross-origin
requests to your API (see [Exploiting CORS misconfigurations](https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html)).

Echo refuses this combination rather than building an insecure middleware: `CORS` and
`CORSWithConfig` **panic**, and `CORSConfig.ToMiddleware()` returns an error. To allow
credentialed requests, list the trusted origins explicitly:

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins:     []string{"https://example.com"},
	AllowCredentials: true,
}))
```

For dynamic origin validation, use `UnsafeAllowOriginFunc` and validate each origin
carefully — attackers may register look-alike or hostile (sub)domain names.

### Behind a reverse proxy

When you enable this middleware on every layer of a chained-proxy setup — e.g. an Echo
gateway whose handler proxies to an upstream Echo service that *also* runs `CORS` — the
response ends up with duplicated headers, because the reverse proxy copies the upstream's
CORS headers on top of the ones the gateway already set:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Origin: *
Vary: Origin
Vary: Origin
```

A duplicated `Access-Control-Allow-Origin` is rejected by browsers (only one is allowed).
This is a configuration issue, not something the CORS middleware deduplicates for you.

Fix it at the boundary, not in the middleware:

- **Own CORS at the edge only.** Internal upstream services that sit behind a gateway
  generally should not run CORS at all — let the outermost service own it.
- **Strip upstream CORS headers when proxying.** If an upstream must keep CORS for
  standalone use, remove the CORS headers it set before they are copied, using
  [`httputil.ReverseProxy.ModifyResponse`](https://pkg.go.dev/net/http/httputil#ReverseProxy).
  Delete every `Access-Control-*` header, and remove only the `Origin` token from `Vary`
  so unrelated caching directives (e.g. `Accept-Encoding`) are preserved:

  ```go
  proxy := httputil.NewSingleHostReverseProxy(upstream)
  proxy.ModifyResponse = func(res *http.Response) error {
  	for k := range res.Header {
  		if strings.HasPrefix(k, "Access-Control-") {
  			res.Header.Del(k)
  		}
  	}

  	// Drop only the "Origin" token from Vary, keeping other directives.
  	var kept []string
  	for _, v := range res.Header.Values(echo.HeaderVary) {
  		for _, tok := range strings.Split(v, ",") {
  			if tok = strings.TrimSpace(tok); tok != "" && !strings.EqualFold(tok, echo.HeaderOrigin) {
  				kept = append(kept, tok)
  			}
  		}
  	}
  	res.Header.Del(echo.HeaderVary)
  	for _, v := range kept {
  		res.Header.Add(echo.HeaderVary, v)
  	}
  	return nil
  }
  ```
