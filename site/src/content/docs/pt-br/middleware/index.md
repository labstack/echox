---
title: Escolha um middleware
description: Comece pela tarefa, siga um exemplo executável e consulte a API exata do Echo.
sidebar:
  order: 0
---

Está começando com Echo? Crie primeiro um servidor com o [guia de início](../guide/quickstart/). Cada página explica quando usar o middleware, mostra a configuração e apresenta campos e assinaturas extraídos de uma revisão específica do Echo.

## O que você quer fazer?

| Tarefa | Comece aqui | Experimente |
| --- | --- | --- |
| Ver requisições e falhas | [Request Logger](./logger/) e [Recover](./recover/) | Execute o [exemplo de log](./logger/) |
| Permitir chamadas do navegador | [CORS](./cors/) | Informe explicitamente as origens confiáveis |
| Proteger formulários | [CSRF](./csrf/) | Escolha como armazenar o token |
| Autenticar uma API | [JWT](./jwt/) | Execute o [exemplo de JWT](../cookbook/jwt/) |
| Limitar clientes abusivos | [Rate Limiter](./rate-limiter/) | Escolha identificador e armazenamento |
| Servir arquivos | [Static](./static/) | Execute o [exemplo de arquivos estáticos](./static/) |
| Encaminhar requisições | [Proxy](./proxy/) | Execute o [exemplo de proxy reverso](../cookbook/reverse-proxy/) |

## Do código a uma requisição real

1. Escolha a tarefa acima e copie o menor exemplo de uso.
2. Consulte a configuração para ver os campos e assinaturas da revisão indicada do Echo. Leia os valores padrão e as notas de segurança antes de alterar o comportamento.
3. Execute o exemplo completo ou a receita vinculada e envie uma requisição com `curl`.

Para criar sua primeira API, siga por [roteamento](../guide/routing/), [binding](../guide/binding/), [tratamento de erros](../guide/error-handling/) e [testes](../guide/testing/). Para tráfego de produção, adicione [Recover](./recover/) e [Request Logger](./logger/) antes de escolher os middlewares de segurança.

Integrações como JWT são mantidas em módulos separados; suas páginas identificam o pacote responsável. A API principal de middleware fica em [`github.com/labstack/echo/v5/middleware`](https://pkg.go.dev/github.com/labstack/echo/v5/middleware).
