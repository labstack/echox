---
title: Esquema da requisição e proxies confiáveis
description: Configure como o Echo identifica HTTP ou HTTPS atrás de um proxy confiável.
sidebar:
  order: 15
---

`Context#Scheme()` informa se a requisição usa HTTP ou HTTPS. Redirecionamentos
HTTPS, HSTS e o middleware Proxy dependem desse valor. Desde o Echo v5.4.0 e
v4.16.0, o Echo só usa `X-Forwarded-Proto`, `X-Forwarded-Protocol`,
`X-Forwarded-Ssl` ou `X-Url-Scheme` quando o **par conectado diretamente** usa
um endereço de loopback, link-local ou privado, ou um socket Unix. Caso contrário,
a própria conexão determina o esquema. Assim, um cliente público não pode enviar
`X-Forwarded-Proto: https` por HTTP para evitar o redirecionamento.

## Escolha do extrator

`Echo#SchemeExtractor` controla essa decisão; no v5 também há
`Config.SchemeExtractor`. O padrão `echo.ExtractSchemeFromHeaders()` confia apenas
nos pares diretos citados acima. `echo.ExtractSchemeDirect()` ignora os cabeçalhos
de encaminhamento. `echo.LegacySchemeExtractor()` restaura o comportamento antigo,
que é inseguro se um cliente não confiável puder acessar a aplicação ou se o proxy
repassar um cabeçalho enviado pelo cliente.

Quando `X-Forwarded-Proto` está presente, o Echo usa **somente o último valor**
e retorna o esquema em letras minúsculas. Um valor inválido resulta em `http`;
o Echo não tenta outro cabeçalho de esquema. O proxy confiável deve
**sobrescrever** `X-Forwarded-Proto` com o esquema observado, nunca repassar o
valor do cliente. No nginx, configure
`proxy_set_header X-Forwarded-Proto $scheme;`.

## Proxies com endereços públicos

Confie explicitamente nas faixas de um proxy que se conecta de um endereço
público ou de `100.64.0.0/10`. Sem isso, o Echo ignora seu `X-Forwarded-Proto`:
redirecionamentos HTTPS podem entrar em loop e o middleware Secure pode deixar
de enviar HSTS. Isso inclui Cloudflare, CloudFront, Azure Front Door,
balanceadores HTTP(S) externos do GCP e GKE Ingress. Confie apenas nas faixas
usadas pela sua implantação. Para as faixas do GCP:

```go
_, gclb1, _ := net.ParseCIDR("35.191.0.0/16")
_, gclb2, _ := net.ParseCIDR("130.211.0.0/22")
e.SchemeExtractor = echo.ExtractSchemeFromHeaders(
	echo.TrustIPRange(gclb1),
	echo.TrustIPRange(gclb2),
)
```

Importe `net` e `github.com/labstack/echo/v5` para este exemplo; no v4, use
`github.com/labstack/echo/v4`. Proxies no mesmo host, em rede privada ou em
socket Unix continuam funcionando com a configuração padrão. Se um adaptador
substitui `RemoteAddr` pelo endereço do cliente, configure o extrator para a
topologia real. O IP do cliente é configurado separadamente em
[Endereço IP](/pt-br/guide/ip-address/).
