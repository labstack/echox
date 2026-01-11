---
description: Request ID middleware
---

# Request ID

Request ID middleware generates a unique id for a request.

## Usage

```go
e.Use(middleware.RequestID())
```

*Example*

```go
    e := echo.New()

    e.Use(middleware.RequestID())

    e.GET("/", func(c *echo.Context) error {
        return c.String(http.StatusOK, c.Response().Header().Get(echo.HeaderXRequestID))
    })
    	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
```

## Custom Configuration

### Usage

```go
e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
  Generator: func() string {
    return customGenerator()
  },
}))
```

## Configuration

```go
RequestIDConfig struct {
    // Skipper defines a function to skip middleware.
    Skipper Skipper

    // Generator defines a function to generate an ID.
    // Optional. Default value random.String(32).
    Generator func() string

    // RequestIDHandler defines a function which is executed for a request id.
    RequestIDHandler func(echo.Context, string)

    // TargetHeader defines what header to look for to populate the id
    TargetHeader string
}
```

### Default Configuration

```go
DefaultRequestIDConfig = RequestIDConfig{
  Skipper:   DefaultSkipper,
  Generator: generator,
  TargetHeader: echo.HeaderXRequestID,
}
```

## Set ID

You can set the id from the requester with the `X-Request-ID`-Header

### Request

```sh
curl -H "X-Request-ID: 3" --compressed -v "http://localhost:1323/?my=param"
```

### Log

```js
{"time":"2017-11-13T20:26:28.6438003+01:00","id":"3","remote_ip":"::1","host":"localhost:1323","method":"GET","uri":"/?my=param","my":"param","status":200, "latency":0,"latency_human":"0s","bytes_in":0,"bytes_out":13}
```
