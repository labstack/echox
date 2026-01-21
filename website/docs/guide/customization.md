---
description: Customization
slug: /customization
sidebar_position: 2
---

# Customization

## Logging

`Echo#Logger` - the default format for logging is JSON, which writes to `os.Stdout` by default.


### Custom Logger

Logging is implemented using `slog.Logger` interface which allows you to register
a custom logger using `Echo#Logger`.

```go
e.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```


## Validator

`Echo#Validator` can be used to register a validator for performing data validation
on request payload.

[Learn more](./request.md#validate-data)


## Custom Binder

`Echo#Binder` can be used to register a custom binder for binding request payload.

[Learn more](./binding#custom-binding)


## Custom JSON Serializer

`Echo#JSONSerializer` can be used to register a custom JSON serializer.

Have a look at `DefaultJSONSerializer` on [json.go](https://github.com/labstack/echo/blob/master/json.go).


## Renderer

`Echo#Renderer` can be used to register a renderer for template rendering.

[Learn more](./templates.md)


## HTTP Error Handler

`Echo#HTTPErrorHandler` can be used to register a custom http error handler.

[Learn more](./error-handling.md)


## HTTP Error Handler

`Echo#OnAddRoute` can be used to register a callback function that is invoked when a new route is added to the router.


## IP Extractor for finding real IP address

`Echo#IPExtractor` is used to retrieve IP address reliably/securely, you must let your application be aware of the entire 
architecture of your infrastructure. In Echo, this can be done by configuring `Echo#IPExtractor` appropriately.

[Learn more](./ip-address.md)

