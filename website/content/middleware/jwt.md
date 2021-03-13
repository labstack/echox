+++
title = "JWT Middleware"
description = "JWT middleware for Echo"
[menu.main]
  name = "JWT"
  parent = "middleware"
+++

JWT provides a JSON Web Token (JWT) authentication middleware.

- For valid token, it sets the user in context and calls next handler.
- For invalid token, it sends "401 - Unauthorized" response.
- For missing or invalid `Authorization` header, it sends "400 - Bad Request".

*Usage*

`e.Use(middleware.JWT([]byte("secret")))`

## Custom Configuration

*Usage*

```go
e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
  SigningKey: []byte("secret"),
  TokenLookup: "query:token",
}))
```

## Configuration

```go
JWTConfig struct {
  // Skipper defines a function to skip middleware.
  Skipper Skipper

  // BeforeFunc defines a function which is executed just before the middleware.
  BeforeFunc BeforeFunc

  // SuccessHandler defines a function which is executed for a valid token.
  SuccessHandler JWTSuccessHandler

  // ErrorHandler defines a function which is executed for an invalid token.
  // It may be used to define a custom JWT error.
  ErrorHandler JWTErrorHandler

  // ErrorHandlerWithContext is almost identical to ErrorHandler, but it's passed the current context.
  ErrorHandlerWithContext JWTErrorHandlerWithContext

  // Signing key to validate token.
  // This is one of the three options to provide a token validation key.
  // The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
  // Required if neither user-defined KeyFunc nor SigningKeys is provided.
  SigningKey interface{}

  // Map of signing keys to validate token with kid field usage.
  // This is one of the three options to provide a token validation key.
  // The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
  // Required if neither user-defined KeyFunc nor SigningKey is provided.
  SigningKeys map[string]interface{}

  // Signing method used to check the token's signing algorithm.
  // Optional. Default value HS256.
  SigningMethod string

  // Context key to store user information from the token into context.
  // Optional. Default value "user".
  ContextKey string

  // Claims are extendable claims data defining token content.
  // Optional. Default value jwt.MapClaims
  Claims jwt.Claims

  // TokenLookup is a string in the form of "<source>:<name>" that is used
  // to extract token from the request.
  // Optional. Default value "header:Authorization".
  // Possible values:
  // - "header:<name>"
  // - "query:<name>"
  // - "param:<name>"
  // - "cookie:<name>"
  // - "form:<name>"
  TokenLookup string

  // AuthScheme to be used in the Authorization header.
  // Optional. Default value "Bearer".
  AuthScheme string

  // KeyFunc defines a user-defined function that supplies the public key for a token validation.
  // The function shall take care of verifying the signing algorithm and selecting the proper key.
  // A user-defined KeyFunc can be useful if tokens are issued by an external party.
  //
  // When a user-defined KeyFunc is provided, SigningKey, SigningKeys, and SigningMethod are ignored.
  // This is one of the three options to provide a token validation key.
  // The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
  // Required if neither SigningKeys nor SigningKey is provided.
  // Default to an internal implementation verifying the signing algorithm and selecting the proper key.
  KeyFunc jwt.Keyfunc
}
```

*Default Configuration*

```go
DefaultJWTConfig = JWTConfig{
  Skipper:       DefaultSkipper,
  SigningMethod: AlgorithmHS256,
  ContextKey:    "user",
  TokenLookup:   "header:" + echo.HeaderAuthorization,
  AuthScheme:    "Bearer",
  Claims:        jwt.MapClaims{},
}
```

## [Example](/cookbook/jwt)
