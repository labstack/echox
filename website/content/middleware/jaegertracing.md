+++
title = "Jaeger Tracing Middleware"
description = "Jaeger tracing middleware for Labstack Echo"
[menu.main]
  name = "Jaeger"
  parent = "middleware"
+++

Trace requests on Echo framework with Jaeger Tracing Middleware.

*Usage*

```go
package main
import (
    "github.com/labstack/echo-contrib/jaegertracing"
    "github.com/labstack/echo/v4"
)
func main() {
    e := echo.New()
    // Enable tracing middleware
    c := jaegertracing.New(e, nil)
    defer c.Close()

    e.Logger.Fatal(e.Start(":1323"))
}
```

Enabling the tracing middleware creates a tracer and a root tracing span for every request.

## Custom Configuration

By default, traces are sent to `localhost` Jaeger instance. To configure an external Jaeger, start your application with environment variables.

*Usage*

```bash
$ JAEGER_AGENT_HOST=192.168.1.10 JAEGER_AGENT_PORT=6831 ./myserver
```

### Middleware Skipper

A middleware skipper can be passed to avoid tracing spans to certain URLs:

*Usage*

```go
package main
import (
	"strings"
    "github.com/labstack/echo-contrib/jaegertracing"
    "github.com/labstack/echo/v4"
)

// urlSkipper ignores metrics route on some middleware
func urlSkipper(c echo.Context) bool {
    if strings.HasPrefix(c.Path(), "/testurl") {
        return true
    }
    return false
}

func main() {
    e := echo.New()
    // Enable tracing middleware
    c := jaegertracing.New(e, urlSkipper)
    defer c.Close()

    e.Logger.Fatal(e.Start(":1323"))
}
```

### TraceFunction

This is a wrapper function that can be used to seamlessly add a span for
the duration of the invoked function. There is no need to change function arguments.

*Usage*

```go
package main
import (
    "github.com/labstack/echo-contrib/jaegertracing"
    "github.com/labstack/echo/v4"
    "net/http"
    "time"
)
func main() {
    e := echo.New()
    // Enable tracing middleware
    c := jaegertracing.New(e, nil)
    defer c.Close()
    e.GET("/", func(c echo.Context) error {
        // Wrap slowFunc on a new span to trace it's execution passing the function arguments
		jaegertracing.TraceFunction(c, slowFunc, "Test String")
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}

// A function to be wrapped. No need to change it's arguments due to tracing
func slowFunc(s string) {
	time.Sleep(200 * time.Millisecond)
	return
}
```

### CreateChildSpan

For more control over the Span, the function `CreateChildSpan` can be called
giving control on data to be appended to the span like log messages, baggages and tags.

*Usage*

```go
package main
import (
    "github.com/labstack/echo-contrib/jaegertracing"
    "github.com/labstack/echo/v4"
)
func main() {
    e := echo.New()
    // Enable tracing middleware
    c := jaegertracing.New(e, nil)
    defer c.Close()
    e.GET("/", func(c echo.Context) error {
        // Do something before creating the child span
        time.Sleep(40 * time.Millisecond)
        sp := jaegertracing.CreateChildSpan(c, "Child span for additional processing")
        defer sp.Finish()
        sp.LogEvent("Test log")
        sp.SetBaggageItem("Test baggage", "baggage")
        sp.SetTag("Test tag", "New Tag")
        time.Sleep(100 * time.Millisecond)
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}
```

*References*:

* Opentracing Library: <https://github.com/opentracing/opentracing-go>
* Jaeger configuration: <https://github.com/jaegertracing/jaeger-client-go#environment-variables>