---
title: Static Files
menu:
  side:
    parent: guide
    weight: 3
---

Images, JavaScript, CSS, PDF, Fonts and so on...

### Using Static Middleware

Static middleware can be used to serve static files from the provided root directory.

#### From Default Configuration

```go
DefaultStaticConfig = StaticConfig{
  Root:   "",
  Index:  []string{"index.html"},
  Browse: false,
}
```

##### Usage

```go
e := echo.New()
e.Use(middleware.Static("public"))
```

This will serve any file from the public directory. For example, a request `/js/main.js`
will fetch and serve `public/js/main.js` file.

##### Caveat

With this setup, each and every request goes through this middleware which might
affect the performance. To overcome that, you can use `Echo#Group()` or [`Echo#Static()`]({{< relref "#using-echo-static">}})
APIs.

##### Using `Echo#Group()`

```go
e := echo.New()
e.Group("/static*", middleware.Static("public"))
```

This will serve any file from the public directory. For example, a request `/static/js/main.js`
will fetch and serve `public/js/main.js` file.

#### Custom Configuration

##### Usage

```go
e := echo.New()
e.Use(middleware.StaticFromConfig(middleware.StaticConfig{
  Root:   "public",
  Browse: true,
}))
```

This uses `Root` as public directory to serve the static files and sets `Browse`
to true, enabling directory browsing.

##### Configuration

```go
// StaticConfig defines the config for static middleware.
StaticConfig struct {
  // Root is the directory from where the static content is served.
  // Optional with default value as `DefaultStaticConfig.Root`.
  Root string `json:"root"`

  // Index is the list of index files to be searched and used when serving
  // a directory.
  // Optional with default value as `DefaultStaticConfig.Index`.
  Index []string `json:"index"`

  // Browse is a flag to enable/disable directory browsing.
  // Required.
  Browse bool `json:"browse"`
}
```

### Using `Echo#Static()`

`Echo#Static(prefix, root string)` serves static files from the provided `root` directory for path `/prefix*`.

##### Example 2

```go
e.Static("/static", "assets")
```

This will serve any file from the assets directory for path `/static/*`. For example,
a request `/static/js/main.js` will fetch and serve `assets/js/main.js` file.

##### Example 2

```go
e.Static("/", "assets")
```

This will serve any file from the assets directory for path `/*`. For example,
a request `/js/main.js` will fetch and serve `assets/js/main.js` file.

### Using `Echo#File()`

`Echo#File(path, file string)` serves static file from the provided `file` for `/path`.

##### Example 1

Serving an index page from `public/index.html`

```go
e.File("/", "public/index.html")
```

##### Example 2

Serving a favicon from `images/facicon.ico`

```go
e.File("/favicon.ico", "images/facicon.ico")
```
