+++
title = "JSONP Recipe"
description = "JSONP recipe for Echo"
[menu.main]
  name = "JSONP"
  parent = "cookbook"
+++

JSONP is a method that allows cross-domain server calls. You can read more about it at the JSON versus JSONP Tutorial.

## Server

`server.go`

{{< embed "jsonp/server.go" >}}

## Client

`index.html`

{{< embed "jsonp/public/index.html" >}}

## [Source Code]({{< source "jsonp" >}})

