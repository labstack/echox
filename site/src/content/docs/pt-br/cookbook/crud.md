---
title: CRUD
description: Crie, leia, atualize e exclua recursos com Echo e binding JSON.
sidebar:
  order: 2
---

Uma API CRUD completa (create, read, update, delete) baseada em um store em memória.
Os handlers fazem binding do body JSON do request em uma struct, acessam o store sob um
lock e retornam JSON. Um usuário inexistente retorna `404`, e um id inválido ou um nome vazio retorna `400`.

## Servidor

```go file=cookbook/crud/server.go title="cookbook/crud/server.go"
```

## Cliente

### Criar usuário

Request:

```sh
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe Smith"}' \
  localhost:1323/users
```

Response:

```json
{
  "id": 1,
  "name": "Joe Smith"
}
```

### Obter usuário

Request:

```sh
curl localhost:1323/users/1
```

Response:

```json
{
  "id": 1,
  "name": "Joe Smith"
}
```

### Listar usuários

Request:

```sh
curl localhost:1323/users
```

Response:

```json
[
  {
    "id": 1,
    "name": "Joe Smith"
  }
]
```

### Atualizar usuário

Request:

```sh
curl -X PUT \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe"}' \
  localhost:1323/users/1
```

Response:

```json
{
  "id": 1,
  "name": "Joe"
}
```

### Excluir usuário

Request:

```sh
curl -X DELETE localhost:1323/users/1
```

Response: `204 No Content`.

### Usuário inexistente

Request:

```sh
curl localhost:1323/users/9999
```

Response (`404 Not Found`):

```json
{
  "message": "user not found"
}
```
