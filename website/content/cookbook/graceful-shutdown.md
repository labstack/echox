+++
title = "Graceful Shutdown Recipe"
description = "Graceful shutdown recipe for Echo"
[menu.main]
  name = "Graceful Shutdown"
  parent = "cookbook"
+++

## Using [http.Server#Shutdown()](https://golang.org/pkg/net/http/#Server.Shutdown)

`server.go`

{{< embed "graceful-shutdown/server.go" >}}

> Requires go1.8+

## Source Code

[See]({{< source "graceful-shutdown" >}})