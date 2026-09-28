---
title: Servindo arquivos estáticos
description: Sirva imagens, JavaScript, CSS, fontes e outros assets com Echo.
sidebar:
  order: 9
---

Echo pode servir assets estáticos como imagens, JavaScript, CSS, PDFs e fontes a partir
do filesystem ou de um filesystem incorporado.

## Filesystem padrão

Echo usa `os.DirFS(".")` como seu filesystem padrão, com raiz no diretório de trabalho
atual. Para alterá-lo, defina o campo `Echo#Filesystem`:

```go
e := echo.New()
e.Filesystem = os.DirFS("assets")
```

## Usando o middleware Static

Veja [middleware Static](/pt-br/middleware/static/).

## Usando Echo#Static()

`Echo#Static(prefix, root string)` registra uma rota que serve arquivos estáticos sob
um prefixo de caminho a partir do diretório raiz informado.

Sirva qualquer arquivo de `assets` sob `/static/*`. Um request para `/static/js/main.js`
serve `assets/js/main.js`:

```go
e := echo.New()
e.Static("/static", "assets")
```

Sirva qualquer arquivo de `assets` sob `/*`. Um request para `/js/main.js` serve
`assets/js/main.js`:

```go
e := echo.New()
e.Static("/", "assets")
```

## Usando Echo#StaticFS()

Arquivos estáticos podem ser servidos a partir de qualquer `fs.FS`, incluindo um `embed.FS`. Use
`echo.MustSubFS` para que os arquivos servidos tenham raiz no subdiretório correto — um
`embed.FS` inclui seus subdiretórios como entradas próprias.

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

## Usando Echo#File()

`Echo#File(path, file string)` registra uma rota que serve um único arquivo estático.

Sirva uma página index de `public/index.html`:

```go
e.File("/", "public/index.html")
```

Sirva um favicon de `app/assets/favicon.ico`:

```go
e := echo.New()
e.Filesystem = os.DirFS("/")
e.File("/favicon.ico", "app/assets/favicon.ico") // The file path must not have a leading slash.
```

:::caution
Um `/` inicial no caminho do arquivo não funciona com a maioria das implementações de `fs.FS`. Use
um caminho relativo.
:::

## Atualização de segurança do Echo (v5.4.0 / v4.16.0)

Caminhos com segmentos `.`, `..` ou vazios retornam 404; no modo `HTML5` o índice ainda pode ser servido. Nomes com escape não padrão exigem `StaticConfig.EnablePathUnescaping` no middleware ou `Config.EnablePathUnescapingStaticFiles` em `Echo#Static` e `Echo#StaticFS`. Essas opções também decodificam barras codificadas; não as combine com controle de acesso baseado em rotas. `e.Use(middleware.Static(...))` executa antes das proteções de rota e grupo; mantenha arquivos protegidos fora da raiz ou use uma rota `Echo#Static` protegida.
