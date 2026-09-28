---
title: Esquema de la solicitud y proxies de confianza
description: Configura cómo Echo determina HTTP o HTTPS detrás de un proxy de confianza.
sidebar:
  order: 15
---

`Context#Scheme()` indica si la solicitud usa HTTP o HTTPS. Las redirecciones a
HTTPS, HSTS y el middleware Proxy dependen de este valor. Desde Echo v5.4.0 y
v4.16.0, Echo solo usa `X-Forwarded-Proto`, `X-Forwarded-Protocol`,
`X-Forwarded-Ssl` o `X-Url-Scheme` si el **par directo** tiene una dirección de
bucle local, de enlace local o privada, o si usa un socket Unix. En los demás
casos, la conexión determina el esquema. Así, un cliente público no puede enviar
`X-Forwarded-Proto: https` por HTTP para evitar la redirección.

## Elegir el extractor

`Echo#SchemeExtractor` controla esta decisión; en v5 también existe
`Config.SchemeExtractor`. El valor predeterminado es
`echo.ExtractSchemeFromHeaders()` con las direcciones de confianza anteriores.
`echo.ExtractSchemeDirect()` ignora las cabeceras de reenvío.
`echo.LegacySchemeExtractor()` restaura el comportamiento antiguo y no es seguro
si un cliente no confiable puede acceder a la aplicación o el proxy deja pasar
una cabecera enviada por el cliente.

Si existe `X-Forwarded-Proto`, Echo utiliza **solo su último valor** y devuelve
el esquema en minúsculas. Un valor inválido produce `http`; Echo no prueba otras
cabeceras de esquema. El proxy de confianza debe **sobrescribir** la cabecera,
no reenviar la que envió el cliente. En nginx, usa
`proxy_set_header X-Forwarded-Proto $scheme;`.

## Proxies con direcciones públicas

Confía explícitamente en los rangos de un proxy que se conecte desde una
dirección pública o desde `100.64.0.0/10`. De lo contrario, Echo ignorará su
`X-Forwarded-Proto`: las redirecciones HTTPS pueden entrar en bucle y Secure
puede dejar de enviar HSTS. Esto afecta, por ejemplo, a Cloudflare, CloudFront,
Azure Front Door y a los balanceadores HTTP(S) externos de GCP o GKE Ingress.
Confía solo en los rangos que utilice tu despliegue. Para los rangos de GCP:

```go
_, gclb1, _ := net.ParseCIDR("35.191.0.0/16")
_, gclb2, _ := net.ParseCIDR("130.211.0.0/22")
e.SchemeExtractor = echo.ExtractSchemeFromHeaders(
	echo.TrustIPRange(gclb1),
	echo.TrustIPRange(gclb2),
)
```

Importa `net` y `github.com/labstack/echo/v5` para este ejemplo; en v4 usa
`github.com/labstack/echo/v4`. Los proxies en el mismo equipo, en una red
privada o en un socket Unix funcionan con la configuración predeterminada.
Si un adaptador reemplaza `RemoteAddr` por la dirección del cliente, configura
el extractor según la topología real. La IP del cliente se configura por
separado en [Dirección IP](/es/guide/ip-address/).
