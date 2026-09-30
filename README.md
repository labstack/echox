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

Requires [Node.js](https://nodejs.org) (LTS), the Go version in `go.mod`, and Git.
The `toolchain` directive selects the preferred compiler; the `go` directive
remains the minimum. CI and publishing use that file and enable automatic
toolchain upgrades when reference dependencies or next Echo need a newer Go.
The build publishes the stable docs at `/` and a next preview at `/next/`, each
with its own search index and source revision. Stable source comes from the Echo version
selected by `go.mod`, downloaded through Go's module cache and verified against
`go.sum`. The preview revision remains pinned in `site/next-source.json`.
The build compiles the `reference/` examples against both, extracts middleware fields and function
signatures, and checks stable against `site/reference-baseline.json` and next
against `site/next-reference-baseline.json`. JWT,
Prometheus, and OpenTelemetry use the versions selected by `reference/go.mod`;
`site/external-sources.json` only maps packages and repository URLs. Their APIs
are checked against `site/external-baseline.json`. Shared Echo/JWT versions
must agree between the two Go modules. Dependabot groups Go dependency updates
across both directories into one PR.

Released external middleware packages are compiled against stable Echo with
the `docs_external` build tag. Next and `ECHO_SOURCE_DIR` workspaces validate
the Echo examples and extract the external API facts from their selected module
versions. Dependabot has separate groups for Go version and security updates
across both directories.

```bash
cd site
npm install
npm run dev      # dev server at http://localhost:4321
npm run build    # production build to site/dist
npm run preview  # preview the production build
```

To update stable Echo, update the dependency in both Go modules in the same PR.
Run `npm run source:prepare` from `site/`, review the affected docs and API
differences, then run `npm run source:accept` and review the baseline diff.
Rerun the full build. External middleware updates follow the same review flow;
the baselines record reviewed API facts and versions, rather than selecting
which source is downloaded.

The first build downloads stable/external modules and fetches preview source
into `.cache/`. To test a proposed next Echo checkout, set `ECHO_SOURCE_DIR` to its absolute path when running
`npm run build`. The build reports changed API facts and stops. Review the
affected pages and behavior, then prepare and accept the **next** baseline:

```bash
DOCS_CHANNEL=next ECHO_SOURCE_DIR=/absolute/path/to/echo npm run source:prepare
DOCS_CHANNEL=next ECHO_SOURCE_DIR=/absolute/path/to/echo npm run source:accept
```

Review the baseline diff, update `site/next-source.json` to the proposed
revision, and rerun the full build. The release baseline remains independent.
The generated files in `site/src/generated/` are never edited or committed.
`npm run site:check` checks routes, local links and fragments, image text,
search assets, and locale coverage. `npm run performance:check` catches large
HTML or first-load asset growth on representative stable and next pages, plus
new external resource origins, including scripts, stylesheets, preloads, icons,
images and embedded content. Plain external links and preconnect hints do not
count as loaded resources. New routes pass without baseline acceptance and
receive the same link/accessibility checks. The check warns about new routes
that are not yet in `route-baseline.json`; run `npm run site:accept` and commit
that baseline to protect them from later removal. Removing a baseline route
still requires review.

The single PR workflow runs `go vet`, Go tests with race detection, Node tests,
and the complete site build/checks. API checks also write their result and any
differences to the GitHub Actions step summary. The publishing workflow remains
separate.

Generated field descriptions come from the pinned Go source comments and remain
in English on localized pages; the surrounding task guidance is authored per
locale. `npm run translations:status` identifies changed sections in the Spanish,
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

Use an empty `file=` code fence to display a runnable example directly from its
source file. Paths are relative to the repository root, so the same fence works
in every locale:

````markdown
```go file=cookbook/jwt/custom-claims/server.go
```
````

This works in Markdown and MDX, for Go, HTML, and other code languages. Keep
explanations in the page and edit the program in its source file. The next build
refreshes the displayed code even if only that source file changed. A missing
file or an invalid region fails the build.

For an excerpt, use `file=path#region` and surround the source lines with
`// docs:start region` and `// docs:end region`. Region markers are omitted from
the excerpt. Existing code-fence options such as `title="server.go"` and line
highlights can be combined with `file=`. Excerpts have their common indentation
removed and boundary blank lines trimmed. Region markers are hidden in both
excerpts and full-file examples, with surrounding blank separators kept to one.
Run `cd site && npm test` to check the include behavior and all cookbook pages.
The build runs `npm run cookbook:check`, which rejects pasted Go programs and
complete HTML examples and validates every included file and region.

The build also publishes `/llms.txt` and `/next/llms.txt`. Their page links and
descriptions come from that same content, and each index identifies its pinned
Echo source revision. `site:check` verifies both files and their site links.

## Cookbook recipes

Each folder under `cookbook/` is a self-contained example. Run one with:

```bash
cd cookbook/hello-world
go run .
```

## Deployment

The site auto-deploys to GitHub Pages on every push to `master` (and once weekly,
to refresh build-time data such as the GitHub star count) via
[`.github/workflows/deploy.yaml`](.github/workflows/deploy.yaml). Dependencies
are installed with `npm ci --ignore-scripts` and pinned via the committed
lockfile.

## License

[MIT](LICENSE)
