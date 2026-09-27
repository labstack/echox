---
title: Personalización
description: Personaliza el logger, validator, binder, renderer, serializer y manejo de errores de Echo.
sidebar:
  order: 12
---

Echo expone un conjunto de campos en la instancia `Echo` que te permiten reemplazar el
comportamiento integrado con tus propias implementaciones.

## Logging

`Echo#Logger` escribe logs estructurados. El handler por defecto emite JSON a `os.Stdout`.

### Logger personalizado

El logger es un `*slog.Logger`, por lo que puedes registrar cualquier handler de `slog`:

```go
e.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

## Validator

`Echo#Validator` registra un validator para validar payloads de request.

[Aprende más](/es/guide/request/#validate-data)

## Binder personalizado

`Echo#Binder` registra un binder personalizado para binding de payloads de request.

[Aprende más](/es/guide/binding/#custom-binder)

## Serializer JSON personalizado

`Echo#JSONSerializer` registra un serializer JSON personalizado. Consulta `DefaultJSONSerializer`
en [json.go](https://github.com/labstack/echo/blob/master/json.go).

## Renderer

`Echo#Renderer` registra un renderer para renderizado de templates.

[Aprende más](/es/guide/templates/)

<a id="http-error-handler"></a>

## Handler de errores HTTP

`Echo#HTTPErrorHandler` registra un handler de errores HTTP personalizado.

[Aprende más](/es/guide/error-handling/)

## Callback de ruta

`Echo#OnAddRoute` registra un callback que se invoca cada vez que se agrega una nueva ruta al
router.

## Extractor de IP

`Echo#IPExtractor` controla cómo se determina la dirección IP real del cliente. Para
obtenerla de forma fiable y segura, tu aplicación debe conocer toda su infraestructura.

[Aprende más](/es/guide/ip-address/)
