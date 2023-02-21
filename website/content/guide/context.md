+++
title = "Context"
description = "Context in Echo"
[menu.main]
  name = "Context"
  parent = "guide"
+++

`echo.Context` represents the context of the current HTTP request. It holds request and
response reference, path, path parameters, data, registered handler and APIs to read
request and write response. As Context is an interface, it is easy to extend it with
custom APIs.

## Extending Context

**Define a custom context**

```go
type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}
```

**Create a middleware to extend default context**

```go
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{c}
		return next(cc)
	}
})
```

> This middleware should be registered before any other middleware.

> The custom context cannot be defined in a middleware before the router ran (Pre)

**Use in handler**

```go
e.GET("/", func(c echo.Context) error {
	cc := c.(*CustomContext)
	cc.Foo()
	cc.Bar()
	return cc.String(200, "OK")
})
```

## A Note on Concurrency

`Context` must not be accessed out of the goroutine handling the request. There are two reasons:

1. `Context` has functions that are dangerous to execute from multiple goroutines. Therefore, only one goroutine should access it.
2. Echo uses a pool to create `Context`'s. When the request handling finishes, Echo returns the `Context` to the pool.

See issue [1908](https://github.com/labstack/echo/issues/1908) for a "cautionary tale" caused by this reason. Concurrency is complicated. Beware of this pitfall when working with goroutines.

### Solutions

#### Use a channel

A better design might be to use a channel:

```go
func(c echo.Context) error {
	ca := make(chan string, 1) // To prevent this channel from blocking, size is set to 1.
	r := c.Request()
	method := r.Method

	go func() {
		// This function must not touch the Context.

		fmt.Printf("Method: %s\n", method)

		// Do some long running operations...

		ca <- "Hey!"
	}()

	select {
	case result := <-ca:
		return c.String(http.StatusOK, "Result: "+result)
	case <-c.Request().Context().Done(): // Check context.
		// If it reaches here, this means that context was canceled (a timeout was reached, etc.).
		return nil
	}
}
```
