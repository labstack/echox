---
title: Choose middleware
description: Start with the problem you need to solve, then follow a runnable example and the exact Echo API.
sidebar:
  order: 0
---

New to Echo? Build a server with the [quickstart](../guide/quickstart/) first. Each middleware page explains when to use it, shows usage, and ends with fields and function signatures drawn from a recorded Echo revision.

## What are you trying to do?

| Task | Start here | Try it |
| --- | --- | --- |
| See requests and failures | [Request Logger](./logger/) and [Recover](./recover/) | Run the [logging example](./logger/) |
| Allow a browser app to call your API | [CORS](./cors/) | Start with explicit trusted origins |
| Protect browser forms | [CSRF](./csrf/) | Match token storage to your client |
| Authenticate an API | [JWT](./jwt/) | Run the [JWT cookbook](../cookbook/jwt/) |
| Slow abusive clients | [Rate Limiter](./rate-limiter/) | Choose an identifier and store |
| Serve a website or assets | [Static](./static/) | Run the [static example](./static/) |
| Forward requests to another service | [Proxy](./proxy/) | Run the [reverse proxy cookbook](../cookbook/reverse-proxy/) |

## From code to a working request

1. Pick the task above and copy the smallest usage example.
2. Use its configuration section for the exact fields and signatures in the recorded Echo revision. Read the page's defaults and security notes before changing behavior.
3. Run its linked complete example or cookbook recipe, then send a request with `curl`.

For a first API, continue through [routing](../guide/routing/), [binding](../guide/binding/), [error handling](../guide/error-handling/), and [testing](../guide/testing/). For production traffic, add [Recover](./recover/) and [Request Logger](./logger/) before choosing security middleware for your clients.

External integrations such as JWT are maintained in separate modules; their pages identify the owning package. Echo's core middleware API is in [`github.com/labstack/echo/v5/middleware`](https://pkg.go.dev/github.com/labstack/echo/v5/middleware).
