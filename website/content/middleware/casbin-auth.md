+++
title = "Casbin Auth Middleware"
description = "An authorization library that supports access control models like ACL, RBAC, ABAC "
[menu.main]
  name = "CasbinAuth"
  parent = "middleware"
+++

[Casbin](https://github.com/casbin/casbin) is a powerful and efficient open-source access control library for Golang projects. It provides support for enforcing authorization based on various models. By far, the access control models supported by Casbin are:

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

*Usage*

```go
e := echo.New()
e.Use(middleware.CasbinAuth(casbin.NewEnforcer("casbin_auth_model.conf", "casbin_auth_policy.csv")))
```

- Syntax for models See: [Model.md](https://github.com/casbin/casbin/blob/master/Model.md)


## Custom Configuration

*Usage*

```go
e := echo.New()
ce := casbin.NewEnforcer("casbin_auth_model.conf", "")
ce.AddRoleForUser("alice", "admin")
ce.AddPolicy(...)

e.Use(middleware.CasbinAuth(ce))
// or
e.Use(middleware.CasbinAuthWithConfig(casbin.CasbinAuthConfig(
    Enforcer: ce,
)))
```

## Configuration

```go
type (
  // CasbinAuthConfig defines the config for CasbinAuth middleware.
  CasbinAuthConfig struct {
    // Skipper defines a function to skip middleware.
    Skipper  Skipper
    // Enforcer CasbinAuth main rule.
    // Required.
    Enforcer *casbin.Enforcer
  }
)
```

*Default Configuration*

```go
// DefaultCasbinAuthConfig is the default CasbinAuth middleware config.
DefaultCasbinAuthConfig = CasbinAuthConfig{
    Skipper: DefaultSkipper,
}
```

## [Example]({{< ref "cookbook/casbin-auth.md">}})
