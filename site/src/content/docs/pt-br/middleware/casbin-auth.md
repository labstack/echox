---
title: Casbin Auth
description: Autorize requests com a biblioteca de controle de acesso Casbin usando um pequeno middleware customizado.
sidebar:
  order: 4
---

[Casbin](https://github.com/casbin/casbin) é uma biblioteca de controle de acesso open-source
poderosa e eficiente para Go. Ela oferece suporte à aplicação de autorização em muitos modelos:

- ACL (Access Control List)
- ACL com superusuário
- ACL sem usuários — útil para sistemas sem autenticação ou logins de usuário
- ACL sem recursos — mire um tipo de recurso (por exemplo `write-article`, `read-log`) em vez de um individual
- RBAC (Role-Based Access Control)
- RBAC com roles de recurso — usuários e recursos podem ter roles
- RBAC com domínios/tenants — usuários podem ter conjuntos de roles diferentes por domínio/tenant
- ABAC (Attribute-Based Access Control)
- RESTful
- Deny-override — regras de allow e deny têm suporte; deny sobrescreve allow

Veja a [visão geral da API](https://casbin.org/docs/api-overview) e a
[documentação do Casbin](https://casbin.org/docs/) para detalhes.

## Dependências

```bash
go get github.com/casbin/casbin/v3
go get github.com/labstack/echo-jwt/v5
go get github.com/golang-jwt/jwt/v5
```

```go
import (
	"github.com/casbin/casbin/v3"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
)
```

## Implementação

Echo não inclui um middleware Casbin; a integração é um pequeno wrapper em torno do
enforcer do Casbin:

```go file=cookbook/casbin/server.go#middleware
```

## Exemplo

Crie um arquivo de modelo Casbin `auth_model.conf`:

```ini file=cookbook/casbin/auth_model.conf
```

Crie um arquivo de política Casbin `auth_policy.csv`:

```csv file=cookbook/casbin/auth_policy.csv
```

Carregue o modelo e a política em um enforcer do Casbin:

```go file=cookbook/casbin/server.go#enforcer
```

Autenticação e autorização são preocupações separadas. Autentique o usuário com
outro middleware (como JWT ou Basic Auth) e então forneça um `userGetter` para que Casbin
possa autorizar o request.

### Com JWT

```go file=cookbook/casbin/server.go#jwt
```

Teste com:

```bash
curl -v "http://localhost:8080/dataset1/any" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
```

### Com Basic Auth

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

Teste com:

```bash
# should pass
curl -v -u "alice:password" http://localhost:8080/dataset1/any
# should fail
curl -v -u "alice:password" http://localhost:8080/dataset2/resource2
```

### Exemplo completo de Casbin + JWT

Execute o [exemplo completo de Casbin + JWT](https://github.com/labstack/echox/tree/master/cookbook/casbin) com os arquivos de modelo e política acima.
