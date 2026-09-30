---
title: Casbin Auth
description: Autoriza requests con la biblioteca de control de acceso Casbin usando un middleware personalizado pequeño.
sidebar:
  order: 4
---

[Casbin](https://github.com/casbin/casbin) es una biblioteca de control de acceso open-source
potente y eficiente para Go. Soporta aplicar autorización en muchos modelos:

- ACL (Access Control List)
- ACL con superuser
- ACL sin usuarios, útil para sistemas sin autenticación o log-ins de usuario
- ACL sin recursos: apunta a un tipo de recurso (por ejemplo `write-article`, `read-log`) en lugar de a uno individual
- RBAC (Role-Based Access Control)
- RBAC con roles de recursos: tanto usuarios como recursos pueden tener roles
- RBAC con dominios/tenants: los usuarios pueden tener conjuntos de roles distintos por dominio/tenant
- ABAC (Attribute-Based Access Control)
- RESTful
- Deny-override: se soportan reglas allow y deny; deny sobrescribe allow

Consulta el [resumen de API](https://casbin.org/docs/api-overview) y la
[documentación de Casbin](https://casbin.org/docs/) para más detalles.

## Dependencias

```bash
go get github.com/casbin/casbin/v3
```

```go
import (
	"github.com/casbin/casbin/v3"
)
```

## Implementación

Echo no incluye un middleware Casbin; la integración es un wrapper pequeño alrededor del
enforcer de Casbin:

```go file=cookbook/casbin/server.go#middleware
```

## Ejemplo

Crea un archivo de modelo Casbin `auth_model.conf`:

```ini file=cookbook/casbin/auth_model.conf
```

Crea un archivo de policy Casbin `auth_policy.csv`:

```csv file=cookbook/casbin/auth_policy.csv
```

La autenticación y la autorización son responsabilidades separadas. Autentica al usuario con
otro middleware (como JWT o Basic Auth), y luego proporciona un `userGetter` para que Casbin
pueda autorizar el request.

### Con JWT

```go file=cookbook/casbin/server.go#jwt
```

Pruébalo con:

```bash
curl -v "http://localhost:8080/dataset1/any" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
```

### Con Basic Auth

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

Pruébalo con:

```bash
# should pass
curl -v -u "alice:password" http://localhost:8080/dataset1/any
# should fail
curl -v -u "alice:password" http://localhost:8080/dataset2/resource2
```

### Ejemplo completo de Casbin + JWT

Ejecuta el [ejemplo completo de Casbin + JWT](https://github.com/labstack/echox/tree/master/cookbook/casbin) con los archivos de modelo y política anteriores.
