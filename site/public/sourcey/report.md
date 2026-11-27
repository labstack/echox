# Echo API reference — maintainer report

Sourcey documentation site for the `github.com/labstack/echo/v5` Go package,
generated from the pinned source commit
`59aacddc4a13904b8a80a76f1ec838c495b5feca` (20 Nov 2026, release v5.4.7, MIT).
Built with Sourcey 3.9.5 using the godoc adapter in snapshot mode:
a `godoc.json` snapshot of the root package (**424 public API symbols**)
plus an overview, a quickstart, and five guides condensed from the
project's own README at the pinned commit. The static site has 10 pages.

Intended public URL (pending merge of the docs PR into `labstack/echox`):
`https://echo.labstack.com/sourcey/`

## What this documents

- The full root-package API: 38 types, 55 package-level functions,
  232 methods, 78 constants, 21 variables — including `Echo`, `Context`,
  `Config`, `Router`, `DefaultRouter`, `Group`, `RouteInfo`, `DefaultBinder`,
  `HTTPError`, `ProblemDetailsHTTPErrorHandler`, and the generic
  `ContextGet` / `ContextGetOr` helpers.
- Guides: routing and context, middleware (`Pre`/`Use`, RequestLogger,
  Recover), binding and error handling, context and responses, and the
  v4 → v5 migration (including the `Context` interface → struct change).

## Gaps this reference exposes (maintainer-facing)

1. **No generated API reference on echo.labstack.com.** The official docs are
   hand-written guides; the complete 424-symbol root-package reference with
   per-symbol source links at a pinned commit exists only on pkg.go.dev,
   outside the docs site's navigation. This Sourcey site closes that gap.
2. **The live docs site lags the current release.** The site header reads
   "Release docs · Echo v5.4.4", but v5.4.7 shipped on 20 Nov 2026 — three
   patch versions behind, and v5.4.7 is a security release (three GHSA fixes:
   proxy dot-segment resolution, static-prefix guard bypass, and the
   `middleware.Static` group-guard limitation). Readers of the docs see
   pre-fix behavior with no signal that a newer release exists. The docs
   deploy (weekly rebuild + push-to-master in `echox`) should track releases,
   or the header should link the changelog.
3. **No changelog page on the docs site.** `https://echo.labstack.com/guide/changelog/`
   returns 404. `CHANGELOG.md` exists in the repo and documents every release,
   but it is not published as a docs page — users cannot see what changed per
   release from the docs themselves.
4. **Per-API allocation and reuse behavior is undocumented.** `Echo` pools
   `Context` values (`AcquireContext` / `ReleaseContext`), and several
   response helpers commit or flush at different points, but no symbol
   documents which calls allocate, copy, or return reused buffers. Callers
   discover this only by reading source or by getting bitten.
5. **Middleware ordering and guard rules have no central page.** The `Pre`
   vs `Use` split is described, but the interaction of group-scoped `Use`
   with route matching — and the fact that `middleware.Static` attached with
   `Group.Use` is not protected by guards on routes inside the group
   (GHSA-5hp7-4fcm-57f9) — is documented only in a security advisory, never
   as a guide rule.
6. **Binder error taxonomy has no central page.** `BindError`,
   `BindingError`, and validator errors each surface differently, and the
   v5.4.7 map-binding behavior change (non-string map values now error
   instead of being skipped) has no guide explaining which failures each
   error type covers.
7. **Subpackages are not covered.** This reference covers the root package
   only. `middleware/*` and the new `echotest` test helpers have no API
   reference pages here — the same coverage boundary as the guides.
8. **No per-release API diff.** The docs site versions its guides, but there
   is no generated changelog mapping which symbols were added, changed, or
   removed per release — the v5 migration guide lists the big breaks, but
   minor-version additions are invisible.
