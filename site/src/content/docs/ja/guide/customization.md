---
title: カスタマイズ
description: Echo の logger、validator、binder、renderer、serializer、エラー処理をカスタマイズします。
sidebar:
  order: 12
---

Echo は `Echo` インスタンス上に一連のフィールドを公開しており、組み込みの挙動を
独自実装に置き換えられます。

## ログ

`Echo#Logger` は構造化ログを書き込みます。デフォルトのハンドラは JSON を
`os.Stdout` に出力します。

### カスタム logger

logger は `*slog.Logger` なので、任意の `slog` ハンドラを登録できます。

```go
e.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

## Validator

`Echo#Validator` はリクエストペイロード検証用の validator を登録します。

[詳しく見る](/ja/guide/request/#validate-data)

## カスタム binder

`Echo#Binder` はリクエストペイロードをバインドするカスタム binder を登録します。

[詳しく見る](/ja/guide/binding/#custom-binder)

## カスタム JSON serializer

`Echo#JSONSerializer` はカスタム JSON serializer を登録します。
[json.go](https://github.com/labstack/echo/blob/master/json.go) の
`DefaultJSONSerializer` を参照してください。

## Renderer

`Echo#Renderer` はテンプレートレンダリング用の renderer を登録します。

[詳しく見る](/ja/guide/templates/)

<a id="http-error-handler"></a>

## HTTP エラーハンドラ

`Echo#HTTPErrorHandler` はカスタム HTTP エラーハンドラを登録します。

[詳しく見る](/ja/guide/error-handling/)

## ルートコールバック

`Echo#OnAddRoute` は、新しいルートがルーターに追加されるたびに呼び出される
コールバックを登録します。

## IP 抽出器

`Echo#IPExtractor` は実際のクライアント IP アドレスをどう判定するかを制御します。
信頼性と安全性を保って取得するには、アプリケーションがインフラ全体を把握している必要があります。

[詳しく見る](/ja/guide/ip-address/)
