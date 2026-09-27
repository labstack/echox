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

Requires [Node.js](https://nodejs.org) (LTS), Go 1.27, and Git. The build
publishes the stable docs at `/` and a next preview at `/next/`, each with its
own search index and source revision. Those revisions are pinned in
`site/echo-source.json` and `site/next-source.json`. The build compiles the
`reference/` examples against both, extracts middleware fields and function
signatures, and checks them against `site/reference-baseline.json`. JWT,
Prometheus, and OpenTelemetry are extracted from their own pinned modules in
`site/external-sources.json` and checked against `site/external-baseline.json`.

```bash
cd site
npm install
npm run dev      # dev server at http://localhost:4321
npm run build    # production build to site/dist
npm run preview  # preview the production build
```

The first build fetches pinned source into `.cache/`. To test a proposed next
Echo checkout, set `ECHO_SOURCE_DIR` to its absolute path when running
`npm run build`. The build reports changed API facts and stops. Review the
affected pages and behavior, run `npm run source:prepare` and
`npm run source:accept`, then review the baseline diff before committing it.
The generated files in `site/src/generated/` are never edited or committed.
`npm run site:check` checks routes, local links and fragments, image text,
search assets, and locale coverage. `npm run performance:check` catches large
HTML or first-load asset growth on representative stable and next pages.

`npm run translations:status` identifies changed sections in the Spanish,
Japanese, Portuguese, and Chinese Request Logger, Static, and middleware task
pages. Translate and review the affected section, then record that page with
`npm run translations:accept -- es logger` (replace locale and page) and review
the baseline diff. This tracks edits to both English and localized text; it
does not judge translation quality. The API tables are generated from source,
so translators focus on explanations, task guidance, examples, and safety notes.

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
