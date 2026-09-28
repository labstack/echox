---
title: Endereço IP
description: Recupere o endereço IP real do cliente com segurança atrás de proxies.
sidebar:
  order: 14
---

O endereço IP tem um papel fundamental em HTTP — ele é usado para controle de acesso,
auditoria, análise baseada em geografia e mais. Echo expõe `Context#RealIP()` para recuperá-lo.

Recuperar o IP _real_ do cliente não é trivial, especialmente quando proxies L7 ficam na
frente da sua aplicação. Nesse caso, o IP real precisa ser repassado por HTTP dos
proxies para sua app — mas você não deve confiar incondicionalmente em headers HTTP, ou
corre o risco de ser enganado. **Isso é um risco de segurança.**

Para recuperar o IP de forma confiável e segura, sua aplicação precisa conhecer toda a sua
infraestrutura. No Echo, você configura isso por meio de `Echo#IPExtractor`.

No v5, `Context#RealIP()` usa por padrão o endereço do par direto desde v5.1.0. O v4 ainda confia por padrão em `X-Forwarded-For` e `X-Real-IP` enviados por qualquer cliente; por isso, **sempre configure `Echo#IPExtractor`** conforme seus proxies. O esquema HTTP/HTTPS é configurado separadamente em [Esquema da requisição](/pt-br/guide/request-scheme/).

Comece com duas perguntas para encontrar a abordagem correta:

1. Você coloca algum proxy HTTP (L7) na frente da aplicação? Isso inclui load balancers
   de nuvem (como AWS ALB ou GCP HTTP LB) e proxies open-source (como Nginx, Envoy ou
   um gateway de ingresso Istio).
2. Se sim, qual header HTTP seus proxies usam para passar o IP do cliente para a
   aplicação?

## Caso 1: Sem proxy

Se não há proxy (a app encara a internet diretamente), o único endereço em que você pode
confiar é o da camada de rede. Todo header HTTP não é confiável porque
clientes têm controle total sobre eles.

Use `echo.ExtractIPDirect()`:

```go
e.IPExtractor = echo.ExtractIPDirect()
```

## Caso 2: Proxies usando o header X-Forwarded-For

[`X-Forwarded-For` (XFF)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For)
é o header mais comum para retransmitir IPs de clientes. A cada salto, o proxy acrescenta o
IP do request ao final do header.

```text
            ┌──────────┐            ┌──────────┐            ┌──────────┐
───────────>│ Proxy 1  │───────────>│ Proxy 2  │───────────>│ Your app │
            │ (IP: b)  │            │ (IP: c)  │            │          │
            └──────────┘            └──────────┘            └──────────┘

Case 1.
XFF:  ""                    "a"                     "a, b"
                                                    ~~~~~~
Case 2.
XFF:  "x"                   "x, a"                  "x, a, b"
                                                    ~~~~~~~~~
                                                    ↑ What your app will see
```

Nesse caso, pegue a **primeira leitura de IP não confiável a partir da direita**. Nunca pegue
a primeira da esquerda, pois o cliente a controla. Aqui "confiável" significa
que você tem certeza de que o IP pertence à sua infraestrutura. No exemplo acima, se `b` e
`c` são confiáveis, o IP do cliente é `a` nos dois casos — nunca `x`.

Use `ExtractIPFromXFFHeader(...TrustOption)`:

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader()
```

Por padrão, ele confia em endereços IP internos — loopback, link-local unicast,
private-use e unique local addresses de
[RFC 6890](https://datatracker.ietf.org/doc/html/rfc6890),
[RFC 4291](https://datatracker.ietf.org/doc/html/rfc4291) e
[RFC 4193](https://datatracker.ietf.org/doc/html/rfc4193). Controle isso com `TrustOption`s:

```go
e.IPExtractor = echo.ExtractIPFromXFFHeader(
	echo.TrustLoopback(false),   // e.g. IPv4 starting with 127.
	echo.TrustLinkLocal(false),  // e.g. IPv4 starting with 169.254.
	echo.TrustPrivateNet(false), // e.g. IPv4 starting with 10. or 192.168.
	echo.TrustIPRange(lbIPRange),
)
```

## Caso 3: Proxies usando o header X-Real-IP

`X-Real-IP` é outro header para retransmitir o IP do cliente, mas diferentemente de XFF ele carrega
apenas um único endereço.

Se seus proxies definem este header, use `ExtractIPFromRealIPHeader(...TrustOption)`:

```go
e.IPExtractor = echo.ExtractIPFromRealIPHeader()
```

Assim como com XFF, ele confia em endereços IP internos por padrão e aceita as mesmas
`TrustOption`s.

:::danger
**Nunca esqueça** de configurar o proxy mais externo (na borda da sua
infraestrutura) para **não repassar headers recebidos**. Caso contrário, um cliente pode
forjá-los, abrindo espaço para fraude.
:::

## Comportamento padrão

Sem `Echo#IPExtractor`, **o v5 usa o endereço do par direto** e ignora cabeçalhos de IP enviados pelo cliente. **O v4 mantém o comportamento legado**, que pode usar `X-Forwarded-For` ou `X-Real-IP` sem verificar o proxy. No v4, use `echo.ExtractIPDirect()` se não houver proxy, ou um extrator de cabeçalho com opções de confiança adequadas aos seus proxies.
