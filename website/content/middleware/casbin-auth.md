+++
title = "Casbin Auth Middleware"
description = "Casbin Auth middleware for Echo. It supports access control models like ACL, RBAC, ABAC."
[menu.main]
  name = "Casbin Auth"
  parent = "middleware"
+++

[Casbin](https://github.com/casbin/casbin) is a powerful and efficient open-source access control library for Go. It provides support for enforcing authorization based on various models. So far, the access control models supported by Casbin are:

- ACL (Access Control List)
- ACL with superuser
- ACL without users: especially useful for systems that don't have authentication or user log-ins.
- ACL without resources: some scenarios may target for a type of resources instead of an individual resource by using permissions like write-article, read-log. It doesn't control the access to a specific article or log.
- RBAC (Role-Based Access Control)
- RBAC with resource roles: both users and resources can have roles (or groups) at the same time.
- RBAC with domains/tenants: users can have different role sets for different domains/tenants.
- ABAC (Attribute-Based Access Control)
- RESTful
- Deny-override: both allow and deny authorizations are supported, deny overrides the allow.

> Echo community contribution

## Dependencies

```go
import (
  "github.com/casbin/casbin"
  casbin_mw "github.com/labstack/echo-contrib/casbin"
)
```

*Usage*

```go
e := echo.New()
e.Use(casbin_mw.Middleware(casbin.NewEnforcer("casbin_auth_model.conf", "casbin_auth_policy.csv")))
```

For syntax, see: [Model.md](https://github.com/casbin/casbin/blob/master/Model.md).


## Custom Configuration

*Usage*

```go
e := echo.New()
ce := casbin.NewEnforcer("casbin_auth_model.conf", "")
ce.AddRoleForUser("alice", "admin")
ce.AddPolicy(...)
e.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
  Enforcer: ce,
}))
```

## Configuration

```go
// Config defines the config for CasbinAuth middleware.
Config struct {
  // Skipper defines a function to skip middleware.
  Skipper middleware.Skipper

  // Enforcer CasbinAuth main rule.
  // Required.
  Enforcer *casbin.Enforcer
}
```

*Default Configuration*

```go
// DefaultConfig is the default CasbinAuth middleware config.
DefaultConfig = Config{
  Skipper: middleware.DefaultSkipper,
}
```
