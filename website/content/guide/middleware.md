---
title: Middleware
menu:
  side:
    parent: guide
    weight: 5
---

**Middleware** is a function chained in the HTTP request-response cycle with access
to `Echo#Context` which it uses to perform a specific action, for example, logging
every request or limiting the number of requests.

Handler is processed in the end after all middleware are finished executing.

### Middleware Levels

#### Root Level (Before router)

`Echo#Pre()` can be used to register a middleware which is executed before router
processes the request. It is helpful to make any changes to the request properties,
for example, adding or removing a trailing slash from the path so it matches the
route.

*Note*: As router has not processed the request, middleware at this level won't
have access to any path related API from `echo.Context`.

#### Root Level (After router)

Most of the time you will register a middleware at this level using `Echo#Use()`.
This middleware is executed after router processes the request and has full access
to `echo.Context` API.

The following built-in middleware should be registered at this level:

- Logger
- Gzip
- Recover
- BasicAuth
- Static

#### Group Level

When creating a new group, you can register middleware just for that group. For
example, you can have an admin group which is secured by registering a BasicAuth
middleware for it.

*Example*:

```go
e := echo.New()
admin := e.Group("/admin", middleware.BasicAuth())
```

You can also add a middleware after creating a group via `admin.Use()`.

#### Route Level

When defining a new route, you can optionally register middleware just for it.

*Example*:

```go
e := echo.New()
e.Get("/", <Handler>, <Middleware...>)
```

### Logger Middleware

Logger middleware logs the information about each HTTP request.

#### Configuration

```go
LoggerConfig struct {
  // Format is the log format which can be constructed using the following tags:
  //
  // - time_rfc3339
  // - remote_ip
  // - uri
  // - method
  // - path
  // - status
  // - response_time
  // - response_size
  //
  // Example "${remote_id} ${status}"
  //
  // Optional with default value as `DefaultLoggerConfig.Format`.
  Format string

  // Output is the writer where logs are written.
  // Optional with default value as `DefaultLoggerConfig.Output`.
  Output io.Writer
}
```

#### Default Configuration

```go
DefaultLoggerConfig = LoggerConfig{
  Format: "time=${time_rfc3339}, remote_ip=${remote_ip}, method=${method}, " +
    "uri=${uri}, status=${status}, took=${response_time}, sent=${response_size} bytes\n",
  color:  color.New(),
  Output: os.Stdout,
}
```

##### Usage

`e.Use(middleware.Logger())`

##### Sample Output

`time=2016-03-22T10:33:59-07:00, remote_ip=::1, method=GET, uri=/hello, status=200, took=54.957µs, sent=20 bytes`

#### Custom Configuration

##### Usage

```go
e := echo.New()
e.Use(middleware.LoggerFromConfig(middleware.LoggerConfig{
  Format: "method=${method}, uri=${uri}, status=${status}\n",
}))
```

Example above uses a `Format` which logs request method and request URI.

##### Sample Output

`method=GET, uri=/hello, status=200`

### Recover Middleware

Recover middleware recovers from panics anywhere in the chain, prints stack trace
and handles the control to the centralized
[HTTPErrorHandler]({{< relref "guide/customization.md#http-error-handler">}}).

#### Configuration

```go
RecoverConfig struct {
  // StackSize is the stack size to be printed.
  // Optional with default value as `DefaultRecoverConfig.StackSize`.
  StackSize int

  // StackAll is a flag to format stack traces of all other goroutines into
  // buffer after the trace for the current goroutine, or not.
  // Required.
  StackAll bool

  // PrintStack is a flag to print stack or not.
  // Required.
  PrintStack bool
}
```

#### Default Configuration

```go
DefaultRecoverConfig = RecoverConfig{
	StackSize:  4 << 10, // 4 KB
	StackAll:   true,
	PrintStack: true,
}
```

##### Usage

`e.Use(middleware.Recover())`

#### Custom Configuration

##### Usage

```go
e := echo.New()
e.Use(middleware.RecoverFromConfig(middleware.RecoverConfig{
  StackSize:  1 << 10, // 1 KB
  StackAll:   middleware.DefaultRecoverConfig.StackAll,
  PrintStack: middleware.DefaultRecoverConfig.PrintStack,
}))
```

Example above uses a `StackSize` of 1 KB and sets StackAll to false. For `PrintStack`
it uses the default value.

### Gzip Middleware

Gzip middleware compresses HTTP response using gzip compression scheme.

#### Configuration

```go
GzipConfig struct {
  // Level is the gzip level.
  // Optional with default value as `DefaultGzipConfig.Level`.
  Level int
}
```

#### Default Configuration

```go
DefaultGzipConfig = GzipConfig{
  Level: -1,
}
```

##### Usage

`e.Use(middleware.Gzip())`

#### Custom Configuration

##### Usage

```go
e := echo.New()
e.Use(middleware.GzipFromConfig(middleware.GzipConfig{
  Level: 5
}))
```

Example above uses a `Level` 5 for gzip compression.

### BasicAuth Middleware

BasicAuth middleware provides an HTTP basic authentication.
For valid credentials it calls the next handler.
For invalid credentials, it sends "401 - Unauthorized" response.

#### Configuration

```go
BasicAuthConfig struct {
  // AuthFunc is the function to validate basic auth credentials.
  AuthFunc BasicAuthFunc
}
```

##### Usage

```go
e := echo.New()
e.Use(middleware.BasicAuth(func(username, password string) bool {
	if username == "joe" && password == "secret" {
		return true
	}
	return false
}))
```

### [Static Middleware]({{< relref "guide/static-files.md#using-static-middleware">}})

### AddTrailingSlash Middleware

AddTrailingSlash add a trailing slash to the request URL path.

##### Usage

```go
e := echo.New()
e.Pre(middleware.AddTrailingSlash())
```

### RemoveTrailingSlash Middleware

RemoveTrailingSlash removes a trailing slash from the request URL path.

##### Usage

```go
e := echo.New()
e.Pre(middleware.RequestTrailingSlash())
```

### Writing a custom middleware

*TBD*...
