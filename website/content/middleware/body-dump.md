+++
title = "Body Dump Middleware"
description = "Body dump middleware for Echo"
[menu.main]
  name = "Body Dump"
  parent = "middleware"
+++

Body dump middleware captures the request and response payload and calls the registered handler. Generally used for debugging/logging purpose. Avoid using it if your request/response payload is huge e.g. file upload/download, but if you still need to, add an exception for your endpoints in the skipper function.

*Usage*

```go
e := echo.New()
e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
}))
```

## Custom Configuration

*Usage*

```go
e := echo.New()
e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{}))
```

## Configuration

```go
BodyDumpConfig struct {
  // Skipper defines a function to skip middleware.
  Skipper Skipper

  // Handler receives request and response payload.
  // Required.
  Handler BodyDumpHandler
}
```

*Default Configuration*

```go
DefaultBodyDumpConfig = BodyDumpConfig{
  Skipper: DefaultSkipper,
}
```
