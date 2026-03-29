---
description: Serving static files
slug: /static-files
sidebar_position: 11
---

# Serving Static Files

Images, JavaScript, CSS, PDF, Fonts and so on...

## Default file system

Echo uses `os.DirFS(".")` as default a file system which is set to current working directory. 
To change the default file system, use `Echo#Filesystem` field.

```go
e := echo.New()
e.Filesystem = os.DirFS("assets")
```

## Using Static Middleware

[See ](/middleware/static.md)

## Using Echo#Static()

`Echo#Static(prefix, root string)` registers a new route with path prefix to serve
static files from the provided root directory.

*Usage 1*

```go
e := echo.New()
e.Static("/static", "assets")
```

Example above will serve any file from the assets directory for path `/static/*`. For example,
a request to `/static/js/main.js` will fetch and serve `assets/js/main.js` file.

*Usage 2*

```go
e := echo.New()
e.Static("/", "assets")
```

Example above will serve any file from the assets directory for path `/*`. For example,
a request to `/js/main.js` will fetch and serve `assets/js/main.js` file.

## Using Echo#StaticFS()

Static files can be served from an `embed.FS` instance. Be sure to use `echo.MustSubFS` as embed.FS includes
subdirectories as their own entries  and staticFS needs to server files from the correct root directory.

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

## Using Echo#File()

`Echo#File(path, file string)` registers a new route with path to serve a static
file.

*Usage 1*

Serving an index page from `public/index.html`

```go
e.File("/", "public/index.html")
```

*Usage 2*

:::important

Leading `/` in the file path is will not work with most of the fs.FS implementations.

:::

Serving a favicon from `/app/assets/favicon.ico`

```go
e := echo.New()
e.Filesystem = os.DirFS("/")
e.File("/favicon.ico", "app/assets/favicon.ico") // <--- file path must not have a leading slash
```
