+++
title = "Load Balancing Recipe"
description = "Load balancing multiple Echo servers using a reverse proxy server like Nginx, HaProxy."
[menu.main]
  name = "Load Balancing"
  parent = "cookbook"
+++

This recipe demonstrates how you can use Nginx as a reverse proxy server and load balance between multiple Echo servers.

## How to setup Nginx proxy server with Echo?

### Step 1: Install Nginx

https://www.nginx.com/resources/wiki/start/topics/tutorials/install

### Step 2: Configure Nginx

Create a file `/etc/nginx/sites-enabled/localhost` with the following content:

{{< embed "load-balancing/nginx.conf" >}}

> Change listen, server_name, access_log per your need.

### Step 3: Restart Nginx

`service nginx restart`

### Step 4: Start upstream servers

- `cd upstream`
- `go run server.go server1 :8081`
- `go run server.go server2 :8082` 

### Step 5: Browse to https://localhost:8080

You should see a webpage being served from "server 1" or "server 2".

```sh
Hello from upstream server server1
```

## Source Code

[See]({{< source "load-balancing" >}})

`upstream/server.go`

{{< embed "load-balancing/upstream/server.go" >}}

