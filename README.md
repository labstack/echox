# echox

Source for the [Echo](https://github.com/labstack/echo) documentation site —
published at **[echo.labstack.com](https://echo.labstack.com)** — together with
the runnable cookbook recipes the docs reference.

## Layout

| Path        | What it is                                                                                          |
| ----------- | --------------------------------------------------------------------------------------------------- |
| `site/`     | The docs site — [Astro](https://astro.build) + [Starlight](https://starlight.astro.build). Content lives in `site/src/content/docs/`. |
| `cookbook/` | Standalone, runnable Go example apps referenced from the docs.                                       |
| `reference/` | Source-owned middleware examples and the config-field extractor used by the site build.          |
| `docs/`     | Internal design specs.                                                                               |

## Documentation site

Requires [Node.js](https://nodejs.org) (LTS), Go 1.27, and Git. The site build
fetches the Echo commit recorded in `site/echo-source.json`, compiles the
`reference/` examples against it, extracts exported middleware fields, and
checks them against `site/reference-baseline.json` before building pages.

```bash
cd site
npm install
npm run dev      # dev server at http://localhost:4321
npm run build    # production build to site/dist
npm run preview  # preview the production build
```

The first build fetches the pinned Echo source into `.cache/`. To test a proposed
Echo checkout instead, set `ECHO_SOURCE_DIR` to its absolute path when running
`npm run build`. The build reports changed fields and stops; review the affected
pages, run `npm run source:prepare` and `npm run source:accept`, then review the
baseline diff before committing it. The generated files in `site/src/generated/`
are never edited or committed.

`npm run translations:status` reports whether the reviewed Spanish, Japanese,
Portuguese, and Chinese Request Logger and Static pages still match the current
English source. After reviewing those translations, run
`npm run translations:accept` and review the hash changes. This report tracks
freshness; it does not validate translation quality.

Content is Markdown/MDX under `site/src/content/docs/` (`guide/`, `middleware/`,
`cookbook/`). To add a page, drop a file in the right folder — the sidebar is
generated from each page's `sidebar.order` frontmatter. Every page needs a
`title` and `description`.

## Cookbook recipes

Each folder under `cookbook/` is a self-contained example. Run one with:

```bash
cd cookbook/hello-world
go run .
```

## Deployment

The site auto-deploys to GitHub Pages on every push to `master` (and once daily,
to refresh build-time data such as the GitHub star count) via
[`.github/workflows/deploy.yaml`](.github/workflows/deploy.yaml). Dependencies
are installed with `npm ci --ignore-scripts` and pinned via the committed
lockfile.

## License

[MIT](LICENSE)
