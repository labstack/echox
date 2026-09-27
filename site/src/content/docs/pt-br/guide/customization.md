---
title: Customização
description: Customize logger, validator, binder, renderer, serializer e tratamento de erros do Echo.
sidebar:
  order: 12
---

Echo expõe um conjunto de campos na instância `Echo` que permite substituir o comportamento
embutido pelas suas próprias implementações.

## Logging

`Echo#Logger` escreve logs estruturados. O handler padrão emite JSON para `os.Stdout`.

### Logger customizado

O logger é um `*slog.Logger`, então você pode registrar qualquer handler `slog`:

```go
e.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

## Validator

`Echo#Validator` registra um validador para validação de payloads de request.

[Saiba mais](/pt-br/guide/request/#validate-data)

## Binder customizado

`Echo#Binder` registra um binder customizado para binding de payloads de request.

[Saiba mais](/pt-br/guide/binding/#custom-binder)

## Serializer JSON customizado

`Echo#JSONSerializer` registra um serializer JSON customizado. Veja `DefaultJSONSerializer`
em [json.go](https://github.com/labstack/echo/blob/master/json.go).

## Renderer

`Echo#Renderer` registra um renderer para renderização de templates.

[Saiba mais](/pt-br/guide/templates/)

<a id="http-error-handler"></a>

## Handler de erro HTTP

`Echo#HTTPErrorHandler` registra um handler de erro HTTP customizado.

[Saiba mais](/pt-br/guide/error-handling/)

## Callback de rota

`Echo#OnAddRoute` registra um callback chamado sempre que uma nova rota é adicionada ao
router.

## Extrator de IP

`Echo#IPExtractor` controla como o endereço IP real do cliente é determinado. Para
recuperá-lo de forma confiável e segura, sua aplicação precisa conhecer toda a sua
infraestrutura.

[Saiba mais](/pt-br/guide/ip-address/)
