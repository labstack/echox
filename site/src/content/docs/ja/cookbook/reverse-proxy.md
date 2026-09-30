---
title: リバースプロキシ
description: 上流アプリケーションの前段で Echo をリバースプロキシ兼ロードバランサーとして使います。
sidebar:
  order: 19
---

このレシピでは、WordPress、Node.js、Java、Python、Ruby、Go などのアプリケーションの前段で、
Echo をリバースプロキシ兼ロードバランサーとして使う方法を示します。簡単のため、ここでの上流は
WebSocket も処理する Go サーバーです。

## 1) 上流ターゲット URL を特定する

```go file=cookbook/reverse-proxy/server.go#targets
```

## 2) 上流ターゲットでプロキシミドルウェアを設定する

下のスニペットは round-robin ロードバランシングを使います。
`middleware.NewRandomBalancer()` も使えます。

```go file=cookbook/reverse-proxy/server.go#middleware
```

サブルートにプロキシを設定するには、`Echo#Group()` を使います。

```go file=cookbook/reverse-proxy/server.go#grouped-proxy
```

```sh
go run . -grouped
```

## 3) 上流サーバーを起動する

```sh
cd upstream
go run server.go server1 :8081
go run server.go server2 :8082
```

## 4) プロキシサーバーを起動する

```sh
go run server.go
```

`http://localhost:1323` にアクセスすると、HTTP リクエストは "server 1" から、
WebSocket リクエストは "server 2" から配信される Web ページが表示されます。

```sh
HTTP

Hello from upstream server server1

WebSocket

Hello from upstream server server2!
Hello from upstream server server2!
Hello from upstream server server2!
```

## ソースコード

### 上流サーバー

```go file=cookbook/reverse-proxy/upstream/server.go
```

### プロキシサーバー

```go file=cookbook/reverse-proxy/server.go
```

## Echo のセキュリティ更新 (v5.4.0 / v4.16.0)

Proxy は `Context#Scheme()` から `X-Forwarded-Proto` を設定し、古いスキームヘッダーを削除します。v5 の `X-Real-IP` は `Context#RealIP()` から設定されます。Echo の前段に別のプロキシがある場合は[スキーム](/ja/guide/request-scheme/)と [IP](/ja/guide/ip-address/) の抽出器を設定し、上流へ正しい値を渡してください。
