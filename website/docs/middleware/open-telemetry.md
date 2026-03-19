---
description: Open-telemetry middleware
---

# Open-telemetry

Echo-opentelemetry is a middleware for Echo framework that provides OpenTelemetry instrumentation for HTTP requests.
https://github.com/labstack/echo-opentelemetry

Open-telemetry is a set of open-source tools that provide instrumentation for cloud native applications. 
* [OpenTelemetry Exporters](https://opentelemetry.io/docs/languages/go/exporters/)
* [OpenTelemetry HTTP spec](https://opentelemetry.io/docs/specs/semconv/http/)
* [HTTP metrics spec](https://opentelemetry.io/docs/specs/semconv/http/http-metrics/)



## Usage

Add OpenTelemetry middleware dependency with go modules

```bash
go get github.com/labstack/echo-opentelemetry
```

Use as an import statement

```go
e.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
    ServerName:          "my-server",
    TracerProvider:      tp,
    
    //Skipper:             nil,
    //OnNextError:         nil,
    //OnExtractionError:   nil,
    //MeterProvider:       nil,
    //Propagators:         nil,
    //SpanStartOptions:    nil,
    //SpanStartAttributes: nil,
    //SpanEndAttributes:   nil,
    //MetricAttributes:    nil,
    //Metrics:             nil,
}))
```
For configuration options documentation see [Config](https://github.com/labstack/echo-opentelemetry/blob/main/otel.go#L28) struct.

Add middleware in simplified form, by providing only the server name

```go
e.Use(echootel.NewMiddleware("app.example.com"))
```

Add middleware with configuration options

```go
e.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
  TracerProvider: tp,
}))
```

Retrieving the tracer from the Echo context
```go
tp, err := echo.ContextGet[trace.Tracer](c, echootel.TracerKey)
```

## Example

This example is exporting metrics and spans to stdout but you can use any exporter (OTLP etc).
Read about [OpenTelemetry exporters](https://opentelemetry.io/docs/languages/go/exporters)

```go reference
https://github.com/labstack/echo-opentelemetry/blob/main/example/main.go
```
