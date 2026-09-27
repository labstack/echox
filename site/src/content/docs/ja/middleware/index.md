---
title: 目的からミドルウェアを選ぶ
description: やりたいことから始め、実行できる例と正確な Echo API を確認します。
sidebar:
  order: 0
---

Echo を初めて使う場合は、まず[クイックスタート](../guide/quickstart/)でサーバーを作ってください。各ページには用途、使い方、記録された Echo リビジョンから取得したフィールドと関数シグネチャがあります。

## 何を実現しますか？

| 目的 | ここから始める | 試す |
| --- | --- | --- |
| リクエストと失敗を記録する | [Request Logger](./logger/) と [Recover](./recover/) | [ログの例](./logger/)を実行 |
| ブラウザーから API を呼ぶ | [CORS](./cors/) | 信頼するオリジンを明示する |
| フォームを保護する | [CSRF](./csrf/) | クライアントに合うトークン保存方法を選ぶ |
| API を認証する | [JWT](./jwt/) | [JWT クックブック](../cookbook/jwt/)を実行 |
| 過剰なリクエストを制限する | [Rate Limiter](./rate-limiter/) | 識別子とストアを選ぶ |
| ファイルを配信する | [Static](./static/) | [静的ファイルの例](./static/)を実行 |
| 別のサービスへ転送する | [Proxy](./proxy/) | [リバースプロキシの例](../cookbook/reverse-proxy/)を実行 |

## コードから動作確認まで

1. 上の目的を選び、最小の使用例をコピーします。
2. 設定欄で記録された Echo リビジョンのフィールドとシグネチャを確認します。動作を変える前に既定値とセキュリティの注意を読みます。
3. 完全な例またはクックブックを実行し、`curl` でリクエストを送ります。

最初の API では、[ルーティング](../guide/routing/)、[バインディング](../guide/binding/)、[エラー処理](../guide/error-handling/)、[テスト](../guide/testing/)へ進んでください。本番トラフィックには、必要なセキュリティ用ミドルウェアを選ぶ前に [Recover](./recover/) と [Request Logger](./logger/)を追加します。

JWT などの外部連携は別のモジュールで保守され、各ページに所有するパッケージを示します。Echo 本体のミドルウェア API は [`github.com/labstack/echo/v5/middleware`](https://pkg.go.dev/github.com/labstack/echo/v5/middleware) にあります。
