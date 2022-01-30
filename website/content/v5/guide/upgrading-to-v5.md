+++ title = "Upgrading to V5"
description = "Upgrading to V5"
toc = true
[menu.v5_main]
name = "Upgrading to V5"
parent = "guide"
+++

## Error and handling of errors

### Error handler

* Changed default error handler signature:
  ```go
  type HTTPErrorHandler func(err error, c echo.Context) // <-- previously
  type HTTPErrorHandler func(c echo.Context, err error) // <-- now
  ```
* Changed `echo.NotFoundHandler` to private. Replaced places where it was called with actual error. 404 (
  echo.ErrNotFound)
  should now be handled in global error handler and not be replaced.
* Changed `echo.MethodNotAllowedHandler` to private. Replaced places where it was called with actual error. 405 (
  echo.ErrMethodNotAllowed)
  should now be handled in global error handler and not be replaced.
  
  If you need to get route name even for 404/405 you can access it with:
  ```go
  e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        log.Printf("executing middleware for route named: %s", c.RouteInfo().Name())
        return next(c)
    }
  })
  ```
  will be values of constants`echo.NotFoundRouteName` for 404 and `echo.MethodNotAllowedRouteName` for 405

### `echo.HTTPError` struct

* Added: function `echo.NewHTTPErrorWithInternal()` to help creating errors with internal error.
* Added: `echo.HTTPError.WithInternal(err error)` method to help creating errors with internal error. Example:
  ```go
  return echo.ErrUnauthorized.WithInternal(err)
  ```
* Removed `echo.HTTPError.SetInternal()` method as it allowed to mutate global errors. Use `HTTPError.WithInternal(err)`
  instead. 

## `echo.Echo` struct

* Removed field `Echo.ListenerNetwork`. Now part of `echo.StartConfig` as field `ListenerNetwork`
* Removed field `Echo.HidePort`. Now part of `echo.StartConfig` as field `HidePort`
* Removed field `Echo.HideBanner`. Now part of `echo.StartConfig` as field `HideBanner`
* Removed field `Echo.DisableHTTP2`. Now part of `echo.StartConfig` as field `DisableHTTP2`
* Removed field `Echo.AutoTLSManager`. AutoTLS is now started with `echo.StartConfig.StartAutoTLS()`. Can not be
  directly passed anymore. Create custom `http.Server` to have same functionality.
* Removed field `Echo.TLSListener`. Now part of `echo.StartConfig`. TLS can be configured
  with `StartConfig.TLSConfigFunc` callback
* Removed field `Echo.Listener`. Replaced with `echo.StartConfig` implementation. Can not be directly passed anymore.
  Create custom `http.Server` to have same functionality.
* Removed field `Echo.TLSServer`. Replaced with `echo.StartConfig` implementation. Can not be directly passed anymore.
  Create custom `http.Server` to have same functionality.
* Removed field `Echo.Server`. Replaced with `echo.StartConfig` implementation. Can not be directly passed anymore.
* Removed field `Echo.notFoundHandler`.
* Removed field `Echo.colorer`. No colored logger messages anymore. Removed as it is not useful in production env.
* Removed field `Echo.StdLogger` . Now part of `echo.StartConfig`
* Removed field `Echo.startupMutex`. Removed as server functionality is in `echo.StartConfig` struct now.


* Moved Echo method `Echo.DefaultHTTPErrorHandler(err error, c Context)` to be separate function that is created
  with `DefaultHTTPErrorHandler(exposeError bool) HTTPErrorHandler` function. Note functionality that `e.Debug` provided
  is now available with `exposeError bool` argument.

* Changed `Echo.Router()` returns now `Router` interface instead of `*Router` struct.
* Changed `Echo.Routers()` returns now `map[string]Router` instead of `map[string]*Router`.

* Added field `Echo.NewContextFunc func(e *Echo, pathParamAllocSize int) EditableContext` function called to create new context
  instance. 

  Echo now supports custom `echo.Context` implementations to serve requests (`echo.Echo.ServeHTTP()`). Setup context creation callback
  with `Echo.NewContextFunc` and when request is being served and context being borrowed from the context pool, then
  context will be created by that callback. Note: Using custom context is slightly slower than default context as cast
  to interface is slightly slower than to struct.

  Example:
  ```go
  type myCustomContext struct {
    echo.DefaultContext
  }

  func (c *myCustomContext) QueryParam(name string) string {
    return "prefix_" + c.DefaultContext.QueryParam(name)
  }

  func main() {
    e := echo.New()
    e.NewContextFunc = func(ec *Echo, pathParamAllocSize int) echo.ServableContext {
      p := make(echo.PathParams, pathParamAllocSize)
      return &myCustomContext{
        DefaultContext: *NewDefaultContext(ec, pathParamAllocSize),
      }
    }
  }
  ```
* Added `Echo.RouterFor(host string) Router` to get vhost Router instance by its hostname. This is useful for example cases
  when you want to list all registered routes with that vhost.
* Added `Echo.ResetRouterCreator(creator func(e *Echo) Router)` to allow users set custom callback to use/create custom
  Router implementations. This allows you to change router implementation that Echo is using.
  ```go
  e.ResetRouterCreator(func(e *Echo) Router {
    return &mySpecialRouter{}
  })
  ```
* Added `Echo.AddRoute(route Routable) (RouteInfo, error)` which allows adding route using `echo.Route` structure
  or `Routable` interface and to handle invalid routes with error instead of panic as other method do.
  ```go
    route := echo.Route{
      Method:      http.MethodGet,
      Path:        "/users",
      Handler:     handler,
      Middlewares: []echo.MiddlewareFunc{middleware.Logger()},
      Name:        "my_users", // <-- now only way to add name to the route
    }
    e.AddRoute(route) // ignore returned values
    ri, err := e.AddRoute(route) // or check for error
  ```
* Changed Echo route adding method signatures to return immutable `RouteInfo` instead of reference to added route. Now
  only way to add route with predefined name is: see above for example.

    * `Echo.CONNECT` returns now `RouteInfo` instead of `*Route`.

      Signature is now: `CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) RouteInfo`

    * `Echo.DELETE` returns now `RouteInfo` instead of `*Route`.
    * `Echo.GET` returns now `RouteInfo` instead of `*Route`.
    * `Echo.HEAD` returns now `RouteInfo` instead of `*Route`.
    * `Echo.OPTIONS` returns now `RouteInfo` instead of `*Route`.
    * `Echo.PATCH` returns now `RouteInfo` instead of `*Route`.
    * `Echo.POST` returns now `RouteInfo` instead of `*Route`.
    * `Echo.PUT` returns now `RouteInfo` instead of `*Route`.
    * `Echo.TRACE` returns now `RouteInfo` instead of `*Route`.
    * `Echo.Any` returns now `[]RouteInfo` instead of `[]*Route`.
    * `Echo.Match` returns now `[]RouteInfo` instead of `[]*Route`.
    * `Echo.File` returns now `RouteInfo` instead of `*Route`.
    * `Echo.Add` returns now `RouteInfo` instead of `*Route`.
* Changed `Echo.Static` returns `RouteInfo` and change to accept middlewares as arguments.
  NOW: `Static(prefix, root string, middleware ...MiddlewareFunc) RouteInfo`

* Removed `Echo.Routes() []*Route`. Equivalent is `e.Router().Routes()`
* Removed `Echo.Reverse(name string, params ...interface{}) string`. Equivalent
  is `e.Router().Routes().Reverse("routeName", 123)`
* Removed `Echo.URI(handler HandlerFunc, params ...interface{}) string`. Equivalent same as `Reverse`
* Removed `Echo.URL(h HandlerFunc, params ...interface{}) string`. Equivalent same as `Reverse`

### Methods related to running HTTP(S) server

All method related to starting `http.Server` of `echo.Echo` have been moved to the new `echo.StartConfig` struct.

* Only `Echo.Start(address string) error` was left in place for creating short demos/poc/examples.
  Details:
  * Internally calls `StartConfig.Start()`
  * Supports now `ctrl-c` to shutdown the server (note: `Echo.Shutdown` is removed)
  
  Example:
  ```go
  // only start* method that was left to `e` instance. 
  // Shorthand for demos/examples/proof-of-concepts
  if err := e.Start(":8080"); err != http.ErrServerClosed {
    log.Fatal(err)
  }
  
  sc := echo.StartConfig{ // `StartConfig` contains now fields related to running server
    Address: ":8080",
  }
  // functionally identical to `e.Start(":8080")`
  if err := sc.Start(e); err != http.ErrServerClosed {
    log.Fatal(err)
  }
  ```
* Removed `Echo.Shutdown(ctx stdContext.Context) error`. When using `StartConfig` started server it can be shutdown
  with `StartConfig.GracefulShutdown` struct that can have `context.Context` which closing will close server.

  Replacement example:
  ```go
  gracefulCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
  defer cancel()
	
  sc := echo.StartConfig{
      Address: ":8080", 
      GracefulContext: gracefulCtx,
      //GracefulTimeout:  0, // defaults to 10 seconds
  }
  if err := sc.Start(e); err != http.ErrServerClosed {
    log.Fatal(err)
  }
  ```
* Removed `Echo.Close() error`. When using `StartConfig` started server it can be shutdown
  with `StartConfig.GracefulShutdown` struct that can have `context.Context` which closing will close server.
* Removed `Echo.TLSListenerAddr() net.Addr`. To get listener address when `StartConfig` is started use
  callback `StartConfig.ListenerAddrFunc func(addr net.Addr)` which allows to access listener before server is started.

  Replacement example:
  ```go
  sc := echo.StartConfig{
    Address: ":0",
	ListenerAddrFunc: func(addr net.Addr) {
      log.Printf("server listening at: %v", addr.String())
	},
  } 
  ```
* Removed `Echo.ListenerAddr() net.Addr`. To get listener address when `StartConfig` is started use
  callback `StartConfig.ListenerAddrFunc func(addr net.Addr)` which allows to access listener before server is started.
* Removed `Echo.StartServer(s *http.Server) (err error)`. Equivalent would be using
  callback `StartConfig.BeforeServeFunc func(s *http.Server) error` and starting server with `StartConfig.Start(e)` or
  using `http.Server` directly.
* Removed `Echo.StartH2CServer(address string, h2s *http2.Server) (err error)`

  Replacement example:
  ```go
  // `e.StartH2CServer(address string, h2s *http2.Server)`  replacement would be:
  sc = echo.StartConfig{
      Address: ":8080",
      BeforeServeFunc: func (s *http.Server) error {
          h2s := &http2.Server{
              MaxConcurrentStreams: 250,
              MaxReadFrameSize:     1048576,
              IdleTimeout:          10 * time.Second,
          }
          s.Handler = h2c.NewHandler(e, h2s)
          return nil
      },
  }
  if err := sc.Start(e); err != http.ErrServerClosed {
    log.Fatal(err)
  }
  ```
* Removed `Echo.StartAutoTLS(address string) error`, replaced with `StartConfig.StartAutoTLS(e *Echo) error`

  Example:
  ```go
  // `e.StartAutoTLS(address string)` replacement would be:
  autoTLSManager := autocert.Manager{
    Prompt: autocert.AcceptTOS,
    Cache:  autocert.DirCache("/var/www/.cache"),
  }
  sc = echo.StartConfig{
    Address: ":443",
    TLSConfigFunc: func (tlsConfig *tls.Config) {
      tlsConfig.GetCertificate = autoTLSManager.GetCertificate
      tlsConfig.NextProtos = []string{acme.ALPNProto}
    },
  }
  if err := sc.Start(e); err != http.ErrServerClosed {
    log.Fatal(err)
  }
  ```
* Removed `Echo.StartTLS(address string, certFile, keyFile interface{}) (err error)`, replace
  with `StartConfig.StartTLS(e *Echo, certFile, keyFile interface{}) error`

  Replacement example:
  ```go
  // `e.StartTLS(address string, certFile, keyFile interface{})`  replacement would be:
  if err := sc.StartTLS(e, certFile, keyFile); err != http.ErrServerClosed {
    log.Fatal(err)
  }
  ```

* Removed internal support for custom `tcpKeepAliveListener` for setting TCP keep-alive timeouts. Since Go 1.13 this can be
  configured as:
  ```go
  func main() {
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
      return c.String(http.StatusOK, "OK")
    })

    server := &http.Server{Handler: e}

    lc := net.ListenConfig{KeepAlive: 15 * time.Second}
    ln, err := lc.Listen(context.Background(), "tcp", ":8080")
    if err != nil {
      panic(err)
    }
    defer ln.Close()

    if err := server.Serve(ln); err != http.ErrServerClosed {
      log.Fatalln(err)
    }
  }
  ```


### `echo.StartConfig` provides advanced ways to configure server

StartConfig provides advanced ways to configure server. New config struct fields are following:
```go
// StartConfig is for creating configured http.Server instance to start serve http(s) requests 
// with given Echo instance
type StartConfig struct {
	// Address for the server to listen on (if not using custom listener)
	Address string

	// ListenerNetwork allows setting listener network (see net.Listen for allowed values)
	// Optional: defaults to "tcp"
	ListenerNetwork string

	// CertFilesystem is file system used to load certificates and keys (if certs/keys are given 
	// as paths)
	CertFilesystem fs.FS

	// DisableHTTP2 disables supports for HTTP2 in TLS server
	DisableHTTP2 bool

	// HideBanner does not log Echo banner on server startup
	HideBanner bool

	// HidePort does not log port on server startup
	HidePort bool

	// GracefulContext is context that completion signals graceful shutdown start
	GracefulContext stdContext.Context

	// GracefulTimeout is period which server allows listeners to finish serving ongoing requests.
	// If this time is exceeded process is exited
	// Defaults to 10 seconds
	GracefulTimeout time.Duration

	// OnShutdownError allows customization of what happens when (graceful) server Shutdown 
	// method returns an error.
	// Defaults to calling e.logger.Error(err)
	OnShutdownError func(err error)

	// TLSConfigFunc allows modifying TLS configuration before listener is created with it.
	TLSConfigFunc func(tlsConfig *tls.Config)

	// ListenerAddrFunc allows getting listener address before server starts serving 
	// requests on listener. Useful when address is set as random (`:0`) port.
	ListenerAddrFunc func(addr net.Addr)

	// BeforeServeFunc allows customizing/accessing server before server starts serving 
	// requests on listener.
	BeforeServeFunc func(s *http.Server) error
}
```

See example:
```go
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	gracefulCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	
	sc := echo.StartConfig{
		Address: ":0", // start on random free port
		//ListenerNetwork:  "",
		//CertFilesystem:   nil,
		//DisableHTTP2:     false,
		//HideBanner:       false,
		//HidePort:         false,
		GracefulContext: gracefulCtx,
		//GracefulTimeout:  0, // defaults to 10 seconds
		//OnShutdownError:  nil,
		//TLSConfigFunc:    nil,
		ListenerAddrFunc: func(addr net.Addr) {
			log.Printf("server listening at: %v", addr.String())
		},
		//BeforeServeFunc:  nil,
	}
	err := sc.Start(e)
	//err := sc.StartTLS(e, "cert", "key")
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
```


## `echo.Group` has similar changes as `echo.Echo`

* Changed `Group.Use` does not add anymore 2 additional (`""` and `"/*"`) ANY routes. This was
  meant to handle situations when you create group without adding any routes to it and wondering why any of the middleware
  was not executed. Usually when you are dealing with `Static` or `Logger` middleware. To remediate this - these 2 
  hidden routes were added so middleware would be executed.
  This hidden behaviour is removed as nothing stops you to create two or some group with identical paths with different
  middlewares and in that case last group that calls `g.Use` will override in Router previous route additional routes
  creating hard to understand bug or "feature".
  
  If you want group middlewares to match requests - you need to create routes explicitly - possibly `ANY /*` route. 

  Example:
  ```go
  g := e.Group("/api", middleware.Logger()) // `Group` also calls internally `Use`
  g.Use(middleware.Recover())
  g.Any("/*", func(c echo.Context) error {
    return echo.ErrNotFound
  })
  ```
  This `ANY /*` route will match all requests to `/api/*` and therefore execute Logger and Recover middlewares
* Changed route adding method signatures to return error instead of added route reference
    * `Group.CONNECT` returns now `RouteInfo` instead of `*Route`.
    * `Group.DELETE` returns now `RouteInfo` instead of `*Route`.
    * `Group.GET` returns now `RouteInfo` instead of `*Route`.
    * `Group.HEAD` returns now `RouteInfo` instead of `*Route`.
    * `Group.OPTIONS` returns now `RouteInfo` instead of `*Route`.
    * `Group.PATCH` returns now `RouteInfo` instead of `*Route`.
    * `Group.POST` returns now `RouteInfo` instead of `*Route`.
    * `Group.PUT` returns now `RouteInfo` instead of `*Route`.
    * `Group.TRACE` returns now `RouteInfo` instead of `*Route`.
    * `Group.Any` returns now `[]RouteInfo` instead of `[]*Route`.
    * `Group.Match` returns now `[]RouteInfo` instead of `[]*Route`.
    * `Group.Add` returns now `RouteInfo` instead of `*Route`.
* Added `Group.AddRoute(route Routable) (RouteInfo, error)` which allows adding route using `echo.Route` structure
  or `Routable` interface.
* Changed `Group.Static` returns error and change to accept middlewares as arguments.
  NOW: `Static(prefix, root string, middleware ...MiddlewareFunc) RouteInfo`
* Changed `Group.File` returns error and change to accept middlewares as arguments.
  NOW: `File(path, file string, middleware ...MiddlewareFunc) RouteInfo `

## `echo.Context` interface has multiple breaking changes.

* Renamed `echo.context` to `echo.DefaultContext`. You can now create your own Context implementations that embeds and 
  possibly extends `DefaultContext` behaviour. See `Echo.NewContextFunc` above for example.

* Added `FileFS(file string, filesystem fs.FS) erro` to serve file from given filesystem.
* Added `QueryParamDefault(name, defaultValue string) string` allows getting default value if query param does not exist.
* Added `FormValueDefault(name, defaultValue string) string` allows getting default value if form value does not exist.
* Added `RouteInfo() RouteInfo` to return route info to which this request was matched to. As `RouteInfo` is interface
  then users can create custom implementation that can be added to router through `echo.Echo.AddRoute(Routable)`
  method. `Routable` is interface that is responsible for creating `RouteInfo` inside router when route is being added.

* Renamed `FormParams() (url.Values, error)` to `FormValues() (url.Values, error)`.
* Renamed `Param(name string) string` to `PathParam(name string) string`.
* Renamed `ParamNames() []string` to `PathParams() PathParams`
* Renamed `SetParamNames(names ...string)` to `SetPathParams(params PathParams)`

  Reasoning for that is: Method name `Param` is too ambiguous. Context interface already has methods like:

  > QueryParam(name string) string
  > FormValue(name string) string

  which explicitly describe where that `*Param` or `*Value` data comes/originates from. `PathParam` is explicit as those
  methods.

* Removed `ParamValues() []string` as `PathParams` result contains both name and value in `PathParam` struct
* Removed `Logger() Logger`. Usage of Echo own logger is discouraged and use whatever logger library you choose and if
  needed put it into context with `e.Set()` and fetch with `e.Get()`
* Removed `SetLogger(l Logger)`
* Removed `Handler()` method.
* Removed `SetHandler()` method.
* Removed `Error(err error)` method. This method was confusing. Return errors instead of writing error response immediately to client.
  Echo assumes that you return errors and allow upstream middleware to possibly handle them. If you absolutely need to handle error
  then global error handler can be accessed as:
  ```go
  c.Echo().HTTPErrorHandler(err)
  ```

Move to `echo.RoutableContext`

* `SetParamValues(values ...string)`

Note: `RoutableContext` interface is only for Echo internal usage to setup Context after route match has been found,
just before handler chain will be called. Contains setters that should not be touched in middleware or handler.

Note: internally Context now 2 distinct fields to hold path parameters. `context.pathParams` lifecycle (allocation) is
handled by Echo to reduce allocation. This field is not directly accessible to user. `context.currentParams` is meant
for cases when middlewares are mutating path params and get copy of `pathParams` values on first usage. This field is
settable for user.

## Router

### Router interface

Introduced interface for router. This will allow using custom router implementations with Echo.

```go
// Router is interface for routing request contexts to registered routes.
type Router interface {
	// Add registers Routable with the Router and returns registered RouteInfo
	Add(routable Routable) (RouteInfo, error)
	// Remove removes route from the Router
	Remove(method string, path string) error
	// Routes returns information about all registered routes
	Routes() Routes

	// Route searches Router for matching route and applies it to the given context. In case when 
	// no matching method was not found (405) or no matching route exists for path (404), router 
	// will return its implementation of 405/404 handler function.
	// When implementing custom Router make sure to populate Context during Router.Route call with:
	// * RoutableContext.SetPath
	// * RoutableContext.SetRawPathParams (IMPORTANT! with same slice pointer that 
	//   c.RawPathParams() returns)
	// * RoutableContext.SetRouteInfo
	// And optionally can set additional information to Context with RoutableContext.Set
	Route(c RoutableContext) HandlerFunc
}
```

### Routable interface

Routable is interface for registering Route with Router. During route registration process the Router will convert
Routable to RouteInfo with ToRouteInfo method. By creating custom implementation of Routable additional information
about registered route can be stored in Routes (i.e. privileges used with route etc.)

```go
// Routable is interface for registering Route with Router. During route registration process the 
// Router will convert Routable to RouteInfo with ToRouteInfo method. By creating custom 
// implementation of Routable additional information about registered route can be stored in 
// Routes (i.e. privileges used with route etc.)
type Routable interface {
	// ToRouteInfo converts Routable to RouteInfo
	//
	// This method is meant to be used by Router after it parses url for path parameters, to 
	// store information about route just added.
	ToRouteInfo(params []string) RouteInfo
	// ToRoute converts Routable to Route which Router uses to register the method handler 
	// for path.
	//
	// This method is meant to be used by Router to get fields (including handler and 
	// middleware functions) needed to add Route to Router.
	ToRoute() Route
	// ForGroup recreates routable with added group prefix and group middlewares it is
	// grouped to.
	//
	// Is necessary for Echo.Group to be able to add/register Routable with Router and
	// having group prefix and group middlewares included in actually registered Route.
	ForGroup(pathPrefix string, middlewares []MiddlewareFunc) Routable
}
```

### RouteInfo interface:

When you add route with `Echo.AddRoute(route Routable) error` that route is added to router
with `Add(routable Routable)` and `RouteInfo` of that route (which is created by `Routable.ToRouteInfo()` method) is
accessible from `router.Routes()`. This all means that you can pass in your own custom types
(than can have extra information about route included in that type) and you are able to have your own routeinfo
implementations when you access ask for route list with `e.Router().Routes()`.

```go
// RouteInfo describes registered route base fields.
// Method+Path pair uniquely identifies the Route. Name can have duplicates.
type RouteInfo interface {
	Method() string
	Path() string
	Name() string

	Params() []string
	Reverse(params ...interface{}) string

	// NOTE: handler and middlewares are not exposed because handler could be already wrapping
	// middlewares and therefore it is not always 100% known if handler function already wraps
	// middlewares or not. In Echo handler could be one function or several functions wrapping 
	// each other.
}
```

### RoutableConfig interface

RoutableConfig is meant to limit methods that Router implementation should be able to change. Although router could
cast provided context to any type it wants, this interface communicates which methods should be used.

```go
// RoutableContext is additional interface that structures implementing Context must implement.
// Methods inside this interface are meant for request routing purposes and should not be used 
// in middlewares.
type RoutableContext interface {
	// Request returns `*http.Request`.
	Request() *http.Request

	// RawPathParams returns raw path pathParams value. Allocation of PathParams is handled by
	// Context.
	RawPathParams() *PathParams

	// SetRawPathParams replaces any existing param values with new values for this context 
	// lifetime (request).
	// Do not set any other value than what you got from RawPathParams as allocation of PathParams
	// is handled by Context.
	SetRawPathParams(params *PathParams)

	// SetPath sets the registered path for the handler.
	SetPath(p string)

	// SetRouteInfo sets the route info of this request to the context.
	SetRouteInfo(ri RouteInfo)

	// Set saves data in the context. Allows router to store arbitrary (that only router has 
	// access to) data in context for later use in middlewares/handler.
	Set(key string, val interface{})
}
```

### Other router changes

* Changed: It is now possible to add same path with different params names i.e. `GET /users/:id`
  and `POST /users/:userID`
* Changed: any arbitrary http method type is now supported.
  ```go
    e.Add("PROPPATCH", "/users", func(c echo.Context) error { return c.String(200, "OK") })
    // or
    e.AddRoute(echo.Route{
      Method:  "PROPPATCH",
      Path:    "/users",
      Handler: func(c echo.Context) error { return c.String(200, "OK") },
    })
  ```
  
Default router behaviour can be configured:

* Added: `RouterConfig.AllowOverwritingRoute` will instruct `DefaultRouter` (not to) to generate error when duplicate
  route (method+path) is added
* Added: `RouterConfig.UnescapePathParamValues` will instruct `DefaultRouter` to unescape path parameter values for
  matched route. People did this with middlewares and possible messed up maxParams size with that.
* Added: `RouterConfig.UseEscapedPathForMatching` will instruct `DefaultRouter` to use escaped path for matching routes.
* Added: `RouterConfig.NotFoundHandler` will instruct `DefaultRouter` to use custom not found handler (404) cases.
* Added: `RouterConfig.MethodNotAllowedHandler` will instruct `DefaultRouter` to use custom handler for method not found handler (405) cases.
* Added: `RouterConfig.OptionsMethodHandler` will instruct `DefaultRouter` to use custom handler for OPTIONS requests.
* Example for configuring default router with custom options:
  ```go
  func main() {
    e := echo.New()

	rc := echo.RouterConfig{
		AllowOverwritingRoute:     false,
		UnescapePathParamValues:   false,
		UseEscapedPathForMatching: false,
		NotFoundHandler:           nil,
		MethodNotAllowedHandler:   nil,
		OptionsMethodHandler:      nil,
	}
	e.ResetRouterCreator(func(_ *echo.Echo) echo.Router {
		return echo.NewRouter(rc)
	})
	
 	if err := e.Start(":8080"); err != http.ErrServerClosed {
  		log.Fatal(err)
  	}
  }
  ```

## Logger

Logger is trimmed down and used only for logging Echo internal stuff and/or used as `http.Server.Errorlog` field for
capturing http server logged errors. For logging in handlers/middlewares use logger library of your own choosing.

```go
// Logger defines the logging interface that Echo uses internally in few places.
// For logging in handlers use your own logger instance (dependency injected or package/public
// variable) from logging framework of your choice.
type Logger interface {
	// Write provides writer interface for http.Server `ErrorLog` and for logging startup messages.
	// `http.Server.ErrorLog` logs errors from accepting connections, unexpected behavior from 
	// handlers, and underlying FileSystem errors.
	// `logger` middleware will use this method to write its JSON payload.
	Write(p []byte) (n int, err error)
	// Error logs the error
	Error(err error)
}
```

## Binder

* `e.Bind()` and `DefaultBinder` to be precise do not clear values when binding for base types (ints, bool, float).
  ```go
    type SearchOpts struct {
      ID       int    `query:"id"`
    }

    opts := SearchOpts{
      Length:   100,
    }
    err := c.Bind(&opts)
  ```
  Previously `opts.Length` would be `0` after bind if that field did not exists in query. Now it stays `100` (value is
  not touched).
* Moved following `echo.DefaultBinder` methods to be functions instead:
  * `echo.BindPathParams(c Context, i interface{}) error` 
  * `echo.BindQueryParams(c Context, i interface{}) error` 
  * `echo.BindHeaders(c Context, i interface{}) error` 

  Example:
  ```go
  if err := (&echo.DefaultBinder{}).BindBody(c, &body); err != nil {  // <-- previously
  if err := echo.BindBody(c, &body); err != nil {  // <-- now
  ```

## Middlewares

### General

* Introduced interface for creating middleware functions without panics on configuration errors
  ```go
  type MiddlewareConfigurator interface {
    ToMiddleware() (MiddlewareFunc, error)
  }
  ```
  All middleware configuration structs now implement that interface and allow you to create middleware function and
  check if there are errors. Existing creator functions stay as they are - same signature and panicing on invalid
  configuration.

### BasicAuth

* Removed `DefaultBasicAuthConfig` struct. Middleware creator func will set now default values.
* Added `BasicAuthConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create basic auth middleware without panics.
* Changed `middleware.BasicAuthConfig.Validator` signature.
  Example:
  ```go
  BasicAuthValidator func(string, string, echo.Context) (bool, error) // <-- previously
  BasicAuthValidator func(c echo.Context, user string, password string) (bool, error) // <-- now
  ```

### BodyDump

* Removed `DefaultBodyDumpConfig` struct. Middleware creator func will set now default values.
* Added `BodyDumpConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create body dump middleware without panics.
* Changed `middleware.BodyDumpConfig.Handler` signature.
  Example:
  ```go
  BodyDumpHandler func(echo.Context, []byte, []byte) // <-- previously
  BodyDumpHandler func(c echo.Context, reqBody []byte, resBody []byte) // <-- now
  ```
  
### BodyLimit

* Removed `DefaultBodyLimitConfig` struct. Middleware creator func will set now default values.
* Added `BodyLimitConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create body limit middleware without panics.
* Changed `BodyLimitConfig.LimitBytes` from `string` to `int64` (allowed size in bytes).

### Gzip

* Removed `DefaultGzipConfig` struct. Middleware creator func will set now default values.
* Added `GzipConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create Gzip middleware without panics.

### CORS

* Added `CORSConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create CORS middleware without panics.

### CSRF

* Added `CSRFConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create CSRF middleware without panics.
* Added `CSRFConfig.Generator func() string` a function to generate token.
* Changed how errors are in case of token extraction failure (400). Errors are now of type `middleware.ValueExtractorError`.
  If you rely on error messages this could be important change.

### Decompress

* Removed `DefaultDecompressConfig` struct. Middleware creator func will set now default values.
* Added `DecompressConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create decompress middleware without panics.

### JWT

JWT middleware was refactored not to depend any JWT library. To use JWT user must provide own token extraction method.

* Changed how errors are in case of extraction failure (400). Errors are now of type `middleware.ValueExtractorError`
* Changed cases when no token was found/extracted (`ErrJWTMissing`) to return `http.StatusUnauthorized 401` status code
  instead of `http.StatusBadRequest 400`. If you rely on error messages this could be important change.
* Renamed `ErrorHandlerWithContext` to `ErrorHandler`. Signatures
  is now `JWTErrorHandlerWithContext func(c echo.Context, err error) error`
* Removed configuration fields related to specific JWT token parsing implementations
    * Removed `JWTConfig.SigningKey` field.
    * Removed `JWTConfig.SigningKeys` field.
    * Removed `JWTConfig.SigningMethod` field.
    * Removed `JWTConfig.Claims` field.

  Added example in `jwt_external_test.go` how to use JWT library with `golang-jwt/jwt` library. Middleware is to be
  created as:
  ```go
  // see `jwt_external_test.go` for `CreateJWTGoParseTokenFunc` implementation
  e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
      ParseTokenFunc: CreateJWTGoParseTokenFunc([]byte("secret"), nil), 
  }))
  e.Use(middleware.JWT(CreateJWTGoParseTokenFunc([]byte("secret"), nil))) // or shorter way
  ```
* Added `JWTConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create JWT middleware without panics.

### KeyAuth

* Removed `KeyAuthConfig.AuthScheme` field. Use `KeyAuthConfig.KeyLookup = "header:Authorization:Bearer "` instead (NB:
  space at the end is important).
* Changed how errors are in case of extraction failure (400). Errors are now of type `middleware.ValueExtractorError`. 
  If you rely on error messages this could be important change.
* Changed `middleware.KeyAuthConfig.Validator` signature.
  Example:
  ```go
  KeyAuthValidator func(auth string, c echo.Context) (bool, error) // <-- previously
  KeyAuthValidator func(c echo.Context, key string) (bool, error) // <-- now
  ```
* Changed `middleware.KeyAuthConfig.ErrorHandler` signature.
  Example:
  ```go
  KeyAuthErrorHandler func(err error, c echo.Context) error // <-- previously
  KeyAuthErrorHandler func(c echo.Context, err error) error // <-- now
  ```
* Added `KeyAuthConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create key auth middleware without panics.

### Logger

* Removed coloring support.
* Changed middleware not to call `c.Error()` and return error to be handled up in caller chain (or global error handler)
* Added `LoggerConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create logger middleware without panics.

### MethodOverride

* Added `MethodOverrideConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create method override middleware without
  panics.

### Proxy

* Added `ProxyConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create proxy middleware without panics.

### RateLimiter

* Added `RateLimiterConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create rate limiter middleware without
  panics.

### Recover

* Removed `RecoverConfig.LogLevel` field.
* Changed middleware not to call `c.Error()` and `c.Logger()`. Recovered error/value will bubble up and be handled
  eventually in global error handler.
* Added `RecoverConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create recover middleware without panics.

### RequestID

* Removed `DefaultRequestIDConfig` struct
* Added `RequestIDConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create request ID middleware without panics.

### Redirect

* Added `RedirectConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create redirect middleware without panics.
* Added default configs for:
  * `middleware.RedirectHTTPSConfig` is the HTTPS Redirect middleware config
  * `middleware.RedirectHTTPSWWWConfig` is the HTTPS WWW Redirect middleware config.
  * `middleware.RedirectNonHTTPSWWWConfig` is the non HTTPS WWW Redirect middleware config.
  * `middleware.RedirectWWWConfig` is the WWW Redirect middleware config.
  * `middleware.RedirectNonWWWConfig` is the non WWW Redirect middleware config.

### Rewrite

* Added `RewriteConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create redirect middleware without panics.

### Secure

* Added `SecureConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create secure middleware without panics.

### (Add/Remove)trailing slash

* Changed `TrailingSlashConfig` to `AddTrailingSlashConfig` and `RemoveTrailingSlashConfig`
* Added `AddTrailingSlashConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create add trailing slash middleware
  without panics.
* Added `RemoveTrailingSlashConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create remove trailing slash
  middleware without panics.

### Static

* Added `StaticConfig.Filesystem` field which will be `fs.FS` used to serve file. Defaults to `echo.Echo.Filesystem`
  value.
* Added `StaticConfig.DisablePathUnescaping` field disable path escaping when serving the file.
* Added `StaticConfig.DirectoryListTemplate` field to allow developer to specify custom html template for directory
  listings.
* Added `StaticConfig.ToMiddleware() (echo.MiddlewareFunc, error)` to create static middleware without panics.

### Timeout

* Moved to https://github.com/labstack/echo-contrib repository

## General changes

* Bumped minimal Go version to 1.16 as `fs.FS` filesystem was introduced in that version as we use it for field
  `echo.Filesystem` that is used in methods `e.File()` and `e.Static()`.
* Changed logged greeting messages to be single line to work better with structured logging
  ```json
  {"time":"2022-01-30T23:44:29.952718371+02:00","level":"INFO","prefix":"echo",
  "message":"Echo (v5.0.X). High performance, minimalist Go web framework https://echo.labstack.com"}
  
  {"time":"2022-01-30T23:44:29.952771385+02:00","level":"INFO","prefix":"echo",
  "message":"http(s) server started on [::]:8080"}
  ```
* Changed function fields that have `echo.Context` as an argument have it as a first argument. Most notable is default
  error handler signature change:
  ```go
  type HTTPErrorHandler func(err error, c Context) // <-- previously
  type HTTPErrorHandler func(c Context, err error) // <-- now
  ```
  This is to normalize how callback function are defined. Previously `Context` was sometimes first sometimes at some
  other place.
* Changed how type definitions are written from
  ```go
  type (
    // Binder is the interface that wraps the Bind method.
    Binder interface {
      Bind(i interface{}, c Context) error
    }
  )
  ```
  to

  ```go
  // Binder is the interface that wraps the Bind method.
  type Binder interface {
    Bind(c Context, i interface{}) error
  }
  ```

  This is because in that way it is easier to copy/paste Echo code when there are multiple definitions in one `()`
  block.

