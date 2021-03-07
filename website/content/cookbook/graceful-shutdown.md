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

## Using [grace](https://github.com/facebookgo/grace)

`server.go`

{{< embed "graceful-shutdown/grace/server.go" >}}

## Using [graceful](https://github.com/tylerb/graceful)

`server.go`

{{< embed "graceful-shutdown/graceful/server.go" >}}

## [Source Code]({{< source "graceful-shutdown" >}})
