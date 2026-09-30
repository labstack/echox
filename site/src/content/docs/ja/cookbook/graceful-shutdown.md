---
title: グレースフルシャットダウン
description: 割り込み信号でサーバーを停止する前に、処理中のリクエストを完了させます。
sidebar:
  order: 8
---

グレースフルシャットダウンは、プロセスが終了する前に処理中のリクエストを完了させます。
もっとも簡単な方法は、キャンセル可能なコンテキストを `StartConfig.Start` に渡し、
`GracefulTimeout` を設定することです。割り込み信号でコンテキストがキャンセルされると、
Echo は新しい接続の受け付けを停止し、アクティブなリクエストが完了するまで最大タイムアウトまで待ちます。

## サーバー

```go file=cookbook/graceful-shutdown/server.go#primary-server
```

## カスタム HTTP サーバーを使う

`http.Server` を自分で管理する場合は、goroutine で起動し、signal context を待ってから、
タイムアウト付きで `Shutdown` を呼び出します。

```go file=cookbook/graceful-shutdown/server.go#custom-server
```

```sh
go run . -custom-server
```
