# API reference tool trial

The first trial used `gomarkdoc` v1.1.0 against the pinned Echo checkout's
`middleware` package, which includes CORS, Request Logger, and Static, and
against the external `github.com/labstack/echo-jwt/v5` package from the current
`echox` module. The package output was 2,422 lines for core middleware and 389
lines for JWT. `gomarkdoc --check --output ...` successfully detected that the
unmodified generated core file matched its source. The output included
`EnablePathUnescaping` and excluded the removed `RequestLoggerConfig.LogError`.

`--embed` supports marked regions within authored Markdown. A check against an
unmarked, fully generated file failed because embed mode would append a second
generated block. This confirms that an existing page must add embed markers
before using that mode. Templates can change generated Markdown, but this
package-wide output is too broad for individual middleware pages and is not a
machine-readable field manifest. Those two requirements motivate the small
`config-fields` extractor in `reference/`. It extracts only fields, types,
deprecation markers, and source positions; it does not infer defaults or
rewrite authored pages.

The site currently renders field tables from that extractor. A targeted
`gomarkdoc` template and marked embedding on a real page remain to be evaluated
for signatures and other reference sections. This trial does not justify
generating behavior explanations.
