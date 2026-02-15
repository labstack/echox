package main

import (
	"log/slog"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

// NewCasbinMiddleware returns middleware for [Casbin](https://casbin.org/).
func NewCasbinMiddleware(enforcer *casbin.Enforcer, userGetter func(*echo.Context) (string, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			username, err := userGetter(c)
			if err != nil {
				return echo.ErrUnauthorized.Wrap(err)
			}
			if pass, err := enforcer.Enforce(username, c.Request().URL.Path, c.Request().Method); err != nil {
				return echo.ErrInternalServerError.Wrap(err)
			} else if !pass {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}
			return next(c)
		}
	}
}

/*
Test with:
curl -v "http://localhost:8080/dataset1/any" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
*/
func main() {
	e := echo.New()

	ce, err := casbin.NewEnforcer("auth_model.conf", "auth_policy.csv")
	if err != nil {
		slog.Error("failed to initialize Casbin enforcer", "error", err)
	}

	// BasicAuth middleware does authentication
	// - should pass `curl -v -u "alice:password" http://localhost:8080/dataset1/any`
	// - should fail `curl -v -u "alice:password" http://localhost:8080/dataset2/resource2`
	//e.Use(middleware.BasicAuth(func(c *echo.Context, user string, password string) (bool, error) {
	//	return subtle.ConstantTimeCompare([]byte(user), []byte("alice")) == 1 &&
	//		subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1, nil
	//}))
	//basicAuthUser := func(c *echo.Context) (string, error) { // basic auth user getter for Casbin authorization
	//	username, _, _ := c.Request().BasicAuth() // NB: authorization (PASSWORD check) must be done somewhere!!!
	//	return username, nil
	//}
	//e.Use(NewCasbinMiddleware(ce, basicAuthUser)) // Casbin does authorization

	e.Use(echojwt.JWT([]byte("secret")))               // JWT middleware does authentication
	jwtUser := func(c *echo.Context) (string, error) { // JWT user getter for Casbin authorization
		token, err := echo.ContextGet[*jwt.Token](c, "user")
		if err != nil {
			return "", err
		}
		return token.Claims.GetSubject()
	}
	e.Use(NewCasbinMiddleware(ce, jwtUser)) // Casbin does authorization

	e.GET("/*", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
