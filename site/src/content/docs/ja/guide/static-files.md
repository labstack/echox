---
title: 静的ファイルの配信
description: Echo で画像、JavaScript、CSS、フォント、その他のアセットを配信します。
sidebar:
  order: 9
---

Echo はファイルシステムまたは埋め込みファイルシステムから、画像、JavaScript、CSS、
PDF、フォントなどの静的アセットを配信できます。

## デフォルトファイルシステム

Echo は現在の作業ディレクトリをルートにした `os.DirFS(".")` をデフォルトファイルシステムとして使います。
変更するには `Echo#Filesystem` フィールドを設定します。

```go
e := echo.New()
e.Filesystem = os.DirFS("assets")
```

## Static ミドルウェアを使う

[Static ミドルウェア](/ja/middleware/static/)を参照してください。

## Echo#Static() を使う

`Echo#Static(prefix, root string)` は、指定した root ディレクトリからパス prefix 配下の
静的ファイルを配信するルートを登録します。

`assets` から `/static/*` 配下で任意のファイルを配信します。`/static/js/main.js` へのリクエストは
`assets/js/main.js` を配信します。

```go
e := echo.New()
e.Static("/static", "assets")
```

`assets` から `/*` 配下で任意のファイルを配信します。`/js/main.js` へのリクエストは
`assets/js/main.js` を配信します。

```go
e := echo.New()
e.Static("/", "assets")
```

## Echo#StaticFS() を使う

静的ファイルは `embed.FS` を含む任意の `fs.FS` から配信できます。配信されるファイルが
正しいサブディレクトリをルートにするよう、`echo.MustSubFS` を使います。`embed.FS` は
サブディレクトリもそれぞれのエントリとして含みます。

```go
//go:embed "assets/images"
var images embed.FS

func main() {
	e := echo.New()

	e.StaticFS("/images", echo.MustSubFS(images, "assets/images"))

	sc := echo.StartConfig{Address: ":1323"}
	if err := sc.Start(context.Background(), e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
```

## Echo#File() を使う

`Echo#File(path, file string)` は、単一の静的ファイルを配信するルートを登録します。

`public/index.html` からインデックスページを配信します。

```go
e.File("/", "public/index.html")
```

`app/assets/favicon.ico` から favicon を配信します。

```go
e := echo.New()
e.Filesystem = os.DirFS("/")
e.File("/favicon.ico", "app/assets/favicon.ico") // The file path must not have a leading slash.
```

:::caution
ファイルパス先頭の `/` は、ほとんどの `fs.FS` 実装では機能しません。相対パスを使ってください。
:::

## Echo のセキュリティ更新 (v5.4.0 / v4.16.0)

`.`、`..`、空のパスセグメントは 404 になり、`HTML5` モードではインデックスを返す場合があります。標準以外のエスケープを使うファイル名には、ミドルウェアで `StaticConfig.EnablePathUnescaping`、`Echo#Static` と `Echo#StaticFS` で `Config.EnablePathUnescapingStaticFiles` (v5) または `Echo#EnablePathUnescapingStaticFiles` (v4) が必要です。これらはエンコードされたスラッシュも復元するため、ルート単位のアクセス制御と併用しないでください。`e.Use(middleware.Static(...))` はルートのガードより先に実行されます。保護するファイルはルート外に置くか、ガード付きの `Echo#Static` で配信してください。
