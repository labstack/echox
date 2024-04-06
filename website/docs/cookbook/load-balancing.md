---
description: Load balancing recipe
---

# Load Balancing

This recipe demonstrates how you can use Nginx as a reverse proxy server and load balance between multiple Echo servers.

## Echo

```go reference
https://github.com/labstack/echox/blob/master/cookbook/load-balancing/upstream/server.go
```

### Start servers

- `cd upstream`
- `go run server.go server1 :8081`
- `go run server.go server2 :8082` 

## Nginx

### 1) Install Nginx

https://www.nginx.com/resources/wiki/start/topics/tutorials/install

### 2) Configure Nginx

Create a file `/etc/nginx/sites-enabled/localhost` with the following content:

``` reference
https://github.com/labstack/echox/blob/master/cookbook/load-balancing/nginx.conf
```

:::info

Change listen, server_name, access_log per your need.

:::

### 3) Restart Nginx

```sh
service nginx restart
```

Browse to https://localhost:8080, and you should see a webpage being served from either "server 1" or "server 2".

```sh
Hello from upstream server server1
```
