---
title: Casbin Auth
description: Authorize requests with the Casbin access control library using a small custom middleware.
sidebar:
  order: 4
---

[Casbin](https://github.com/casbin/casbin) is a powerful, efficient open-source access
control library for Go. It supports enforcing authorization across many models:

- ACL (Access Control List)
- ACL with superuser
- ACL without users — useful for systems without authentication or user log-ins
- ACL without resources — target a type of resource (for example `write-article`, `read-log`) rather than an individual one
- RBAC (Role-Based Access Control)
- RBAC with resource roles — both users and resources can have roles
- RBAC with domains/tenants — users can have different role sets per domain/tenant
- ABAC (Attribute-Based Access Control)
- RESTful
- Deny-override — both allow and deny rules are supported, deny overrides allow

See the [API overview](https://casbin.org/docs/api-overview) and the
[Casbin documentation](https://casbin.org/docs/) for details.

## Dependencies

```bash
go get github.com/casbin/casbin/v3
```

```go
import (
	"github.com/casbin/casbin/v3"
)
```

## Implementation

Echo does not ship a Casbin middleware; the integration is a small wrapper around the
Casbin enforcer:

```go file=cookbook/casbin/server.go#middleware
```

## Example

Create a Casbin model file `auth_model.conf`:

```ini file=cookbook/casbin/auth_model.conf
```

Create a Casbin policy file `auth_policy.csv`:

```csv file=cookbook/casbin/auth_policy.csv
```

Authentication and authorization are separate concerns. Authenticate the user with
another middleware (such as JWT or Basic Auth), then supply a `userGetter` so Casbin
can authorize the request.

### With JWT

```go file=cookbook/casbin/server.go#jwt
```

Try it with:

```bash
curl -v "http://localhost:8080/dataset1/any" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
```

### With Basic Auth

```go
// BasicAuth middleware does authentication
e.Use(middleware.BasicAuth(func(c *echo.Context, user, password string) (bool, error) {
	return subtle.ConstantTimeCompare([]byte(user), []byte("alice")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1, nil
}))
basicAuthUser := func(c *echo.Context) (string, error) { // Basic auth user getter for Casbin authorization
	username, _, _ := c.Request().BasicAuth() // password is verified by the BasicAuth middleware above
	return username, nil
}
e.Use(NewCasbinMiddleware(ce, basicAuthUser)) // Casbin does authorization
```

Try it with:

```bash
# should pass
curl -v -u "alice:password" http://localhost:8080/dataset1/any
# should fail
curl -v -u "alice:password" http://localhost:8080/dataset2/resource2
```

### Full Casbin + JWT example

Run the [complete Casbin + JWT example](https://github.com/labstack/echox/tree/master/cookbook/casbin), using the model and policy files above.
