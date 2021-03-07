+++
title = "CORS Recipe"
description = "CORS recipe for Echo"
[menu.main]
  name = "CORS"
  identifier = "middleware-cors"
  parent = "cookbook"
+++

## Server using a list of allowed origins

`server.go`

{{< embed "cors/origin-list/server.go" >}}

## Server using a custom function to allow origins

`server.go`

{{< embed "cors/origin-func/server.go" >}}

## [Source Code]({{< source "cors" >}})
