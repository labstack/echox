---
title: CRUD
description: Crea, lee, actualiza y elimina recursos con Echo y binding JSON.
sidebar:
  order: 2
---

Una API CRUD (create, read, update, delete) completa respaldada por un store en memoria.
Los handlers vinculan el body JSON del request a un struct, acceden al store bajo un lock y
devuelven JSON. Un usuario inexistente devuelve `404`, y un id inválido o un nombre vacío devuelve `400`.

## Servidor

```go file=cookbook/crud/server.go title="cookbook/crud/server.go"
```

## Cliente

### Crear usuario

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

### Obtener usuario

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

### Listar usuarios

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

### Actualizar usuario

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

### Eliminar usuario

Request:

```sh
curl -X DELETE localhost:1323/users/1
```

Response: `204 No Content`.

### Usuario inexistente

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
