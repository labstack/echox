---
title: 自動 TLS
description: Let's Encrypt から TLS 証明書を自動で取得し、更新します。
sidebar:
  order: 3
---

このレシピは、ドメインの TLS 証明書を Let's Encrypt から自動取得します。autocert manager の
`TLSConfig` を使って `StartConfig` を設定し、ポート `443` で待ち受けます。

`https://<DOMAIN>` にアクセスしてください。すべて正しく設定されていれば、TLS 経由で提供される
ウェルカムメッセージが表示されます。

:::tip
- セキュリティを高めるには、autocert manager で host policy を指定してください。
- [Let's Encrypt rate limits](https://letsencrypt.org/docs/rate-limits) に達しないよう、証明書をキャッシュしてください。
- HTTP トラフィックを HTTPS にリダイレクトするには、[リダイレクトミドルウェア](/ja/middleware/redirect/#https-redirect)を使います。
:::

## サーバー

```go file=cookbook/auto-tls/server.go#primary-server
```

## カスタム HTTP サーバーを使う

`http.Server` を完全に制御したい場合は、代わりに autocert manager をカスタム `tls.Config` に接続します。

```go file=cookbook/auto-tls/server.go#custom-server
```
