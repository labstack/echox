---
title: Elige un middleware
description: Empieza por la tarea, sigue un ejemplo ejecutable y consulta la API exacta de Echo.
sidebar:
  order: 0
---

¿Nuevo en Echo? Crea primero un servidor con la [guía de inicio](../guide/quickstart/). Cada página explica cuándo usar el middleware, muestra cómo configurarlo y presenta campos y firmas tomados de una revisión concreta de Echo.

## ¿Qué necesitas hacer?

| Tarea | Empieza aquí | Pruébalo |
| --- | --- | --- |
| Ver peticiones y errores | [Request Logger](./logger/) y [Recover](./recover/) | Ejecuta el [ejemplo de logging](./logger/) |
| Permitir peticiones desde el navegador | [CORS](./cors/) | Usa orígenes de confianza explícitos |
| Proteger formularios | [CSRF](./csrf/) | Elige cómo guardar el token |
| Autenticar una API | [JWT](./jwt/) | Ejecuta el [recetario de JWT](../cookbook/jwt/) |
| Limitar peticiones abusivas | [Rate Limiter](./rate-limiter/) | Elige identificador y almacenamiento |
| Servir archivos | [Static](./static/) | Ejecuta el [ejemplo de archivos estáticos](./static/) |
| Reenviar peticiones | [Proxy](./proxy/) | Ejecuta el [recetario de proxy inverso](../cookbook/reverse-proxy/) |

## Del código a una petición real

1. Elige una tarea y copia el ejemplo de uso más pequeño.
2. Consulta la configuración para ver los campos y firmas de la revisión indicada de Echo. Lee los valores por defecto y las notas de seguridad antes de cambiar el comportamiento.
3. Ejecuta el ejemplo completo o la receta enlazada y envía una petición con `curl`.

Para una primera API, sigue con [rutas](../guide/routing/), [binding](../guide/binding/), [gestión de errores](../guide/error-handling/) y [pruebas](../guide/testing/). Para tráfico de producción, añade [Recover](./recover/) y [Request Logger](./logger/) antes de elegir los middleware de seguridad.

Integraciones como JWT se mantienen en módulos separados; sus páginas indican el paquete responsable. La API del middleware principal está en [`github.com/labstack/echo/v5/middleware`](https://pkg.go.dev/github.com/labstack/echo/v5/middleware).
