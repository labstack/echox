---
title: Servir archivos estáticos
description: Sirve imágenes, JavaScript, CSS, fuentes y otros assets con Echo.
sidebar:
  order: 9
---

Echo puede servir assets estáticos como imágenes, JavaScript, CSS, PDFs y fuentes desde
el filesystem o desde un filesystem embebido.

## Filesystem por defecto

Echo usa `os.DirFS(".")` como filesystem por defecto, con raíz en el directorio de trabajo
actual. Para cambiarlo, establece el campo `Echo#Filesystem`:

```go
e := echo.New()
e.Filesystem = os.DirFS("assets")
```

## Usar el middleware Static

Consulta [middleware Static](/es/middleware/static/).

## Usar Echo#Static()

`Echo#Static(prefix, root string)` registra una ruta que sirve archivos estáticos bajo
un prefijo de path desde el directorio raíz dado.

Sirve cualquier archivo de `assets` bajo `/static/*`. Un request a `/static/js/main.js`
sirve `assets/js/main.js`:

```go
e := echo.New()
e.Static("/static", "assets")
```

Sirve cualquier archivo de `assets` bajo `/*`. Un request a `/js/main.js` sirve
`assets/js/main.js`:

```go
e := echo.New()
e.Static("/", "assets")
```

## Usar Echo#StaticFS()

Los archivos estáticos se pueden servir desde cualquier `fs.FS`, incluido un `embed.FS`. Usa
`echo.MustSubFS` para que los archivos servidos tengan raíz en el subdirectorio correcto: un
`embed.FS` incluye sus subdirectorios como entradas propias.

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

## Usar Echo#File()

`Echo#File(path, file string)` registra una ruta que sirve un único archivo estático.

Sirve una página index desde `public/index.html`:

```go
e.File("/", "public/index.html")
```

Sirve un favicon desde `app/assets/favicon.ico`:

```go
e := echo.New()
e.Filesystem = os.DirFS("/")
e.File("/favicon.ico", "app/assets/favicon.ico") // The file path must not have a leading slash.
```

:::caution
Un `/` inicial en el path del archivo no funciona con la mayoría de implementaciones de `fs.FS`.
Usa un path relativo.
:::

## Actualización de seguridad de Echo (v5.4.0 / v4.16.0)

Las rutas estáticas con segmentos `.`, `..` o vacíos devuelven 404; en modo `HTML5` puede servirse el índice. La codificación no estándar de nombres requiere `StaticConfig.EnablePathUnescaping` para el middleware o `Config.EnablePathUnescapingStaticFiles` (v5) o `Echo#EnablePathUnescapingStaticFiles` (v4) para `Echo#Static` y `Echo#StaticFS`. Estas opciones también decodifican barras codificadas: no las combines con control de acceso basado en rutas. `e.Use(middleware.Static(...))` se ejecuta antes que las guardas de rutas y grupos; guarda los archivos protegidos fuera de su raíz o usa una ruta `Echo#Static` protegida.
