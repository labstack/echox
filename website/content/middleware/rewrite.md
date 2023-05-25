+++
title = "Rewrite Middleware"
description = "Rewrite middleware for Echo"
[menu.main]
  name = "Rewrite"
  parent = "middleware"
+++

Rewrite middleware allows to rewrite an URL path based on provided rules. It can be helpful for backward compatibility or just creating cleaner and more descriptive links.

### Usage

```go
e.Pre(middleware.Rewrite(map[string]string{
  "/old":              "/new",
  "/api/*":            "/$1",
  "/js/*":             "/public/javascripts/$1",
  "/users/*/orders/*": "/user/$1/order/$2",
}))
```

The values captured in asterisk can be retrieved by index e.g. $1, $2 and so on.
Each asterisk will be non-greedy (translated to a capture group `(.*?)`) and if using
multiple asterisk a trailing `*` will match the "rest" of the path.

> Rewrite middleware should be registered via `Echo#Pre()` to get triggered before the router.

## Custom Configuration

### Usage

```go
e := echo.New()
e.Pre(middleware.RewriteWithConfig(middleware.RewriteConfig{}))
```

### Configuration

```go
// RewriteConfig defines the config for Rewrite middleware.
  RewriteConfig struct {
    // Skipper defines a function to skip middleware.
    Skipper Skipper

    // Rules defines the URL path rewrite rules. The values captured in asterisk can be
    // retrieved by index e.g. $1, $2 and so on.
    Rules map[string]string `yaml:"rules"`

    // RegexRules defines the URL path rewrite rules using regexp.Rexexp with captures
    // Every capture group in the values can be retrieved by index e.g. $1, $2 and so on.
    RegexRules map[*regexp.Regexp]string
  }
```

Default Configuration:

| Name    | Value          |
| ------- | -------------- |
| Skipper | DefaultSkipper |

### Regex-based Rules

For advanced rewriting of paths rules may also be defined using regular expression.
Normal capture groups can be defined using `()` and referenced by index (`$1`, `$2`, ...) for the rewritten path.

`RegexRules` and normal `Rules` can be combined.

```go
  e.Pre(RewriteWithConfig(RewriteConfig{
    Rules: map[string]string{
      "^/v1/*": "/v2/$1",
    },
    RegexRules: map[*regexp.Regexp]string{
      regexp.MustCompile("^/foo/([0-9].*)"):  "/num/$1",
      regexp.MustCompile("^/bar/(.+?)/(.*)"): "/baz/$2/$1",
    },
  }))
```
