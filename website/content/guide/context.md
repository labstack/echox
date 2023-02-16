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

`Context` should not be accessed concurrently. It has functions that are dangerous to execute from multiple goroutines. For example have a look at [Request](https://github.com/labstack/echo/blob/0ce73028d0815e0ecec80964cc2da42d98fafa33/context.go#L231) and [SetRequest](https://github.com/labstack/echo/blob/0ce73028d0815e0ecec80964cc2da42d98fafa33/context.go#L235).

If you would like to see for yourself, add the route below to your code and compile with race detector enabled: `go build -race`. Then send an HTTP request to it, and see what happens.

```go
func(c echo.Context) error {
	go func() {
		for i := 0; i < 100; i++ {
			r := c.Request()
			r, _ = http.NewRequest("GET", "/", nil)
			c.SetRequest(r)
		}
	}()

	for i := 0; i < 100; i++ {
		r := c.Request()
		r, _ = http.NewRequest("GET", "/", nil)
		c.SetRequest(r)
	}
	return nil
}
```

### Solutions

#### Use a channel

```go
func(c echo.Context) error {
	ca := make(chan string)
	r := c.Request()
	method := r.Method

	go func() {
		// This function must not touch the Context.

		fmt.Printf("Method: %s\n", method)

		// Do some long running operations...

		ca <- "Hey!"
	}()

	result := <-ca
	return c.String(http.StatusOK, "Result: "+result)
}
```
