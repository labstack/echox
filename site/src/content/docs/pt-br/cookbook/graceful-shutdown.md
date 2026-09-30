---
title: Encerramento graceful
description: Drene requests em andamento antes de parar o servidor em um sinal de interrupção.
sidebar:
  order: 8
---

Um graceful shutdown permite que requests em andamento terminem antes de o processo sair. A
abordagem mais simples é passar um contexto cancelável para `StartConfig.Start` e definir
um `GracefulTimeout`. Quando o contexto é cancelado por um sinal de interrupção, Echo
para de aceitar novas conexões e aguarda até o timeout para que requests ativos
terminem.

## Servidor

```go file=cookbook/graceful-shutdown/server.go
```

## Usando um servidor HTTP customizado

Se você gerencia o `http.Server` por conta própria, inicie-o em uma goroutine, aguarde no
contexto de sinal e então chame `Shutdown` com um timeout:

```go file=cookbook/graceful-shutdown/custom-server/server.go
```

```sh
go run ./custom-server
```
