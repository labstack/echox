import { readdirSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

// Old Docusaurus URLs -> new locations, generated so astro.config stays small.
// The old site served docs under /docs/ with guide pages flattened (slug
// overrides). On GitHub Pages each entry still emits its own static
// meta-refresh page (no server-side wildcards) — that's why there are ~60.
// On a wildcard-capable host (Cloudflare/Netlify) these collapse to a handful
// of `_redirects` lines.

const docsDir = fileURLToPath(new URL('./content/docs', import.meta.url));
const pages = (sub) =>
  readdirSync(`${docsDir}/${sub}`)
    .filter((f) => /\.mdx?$/.test(f))
    .map((f) => f.replace(/\.mdx?$/, ''));

// Guide pages were flattened to /docs/<name>. Exclude pages that have no old
// URL (new in v5) or are handled explicitly below.
const GUIDE_SKIP = new Set(['context', 'installation', 'quickstart']);

const fromList = (subPrefix, toPrefix, names) =>
  Object.fromEntries(names.map((n) => [`/docs/${subPrefix}${n}`, `/${toPrefix}${n}/`]));

// Keep retired recipes independent of the content tree so deleting a page
// never deletes its redirects. Preserve the locale and release channel.
const retiredRecipes = {
  'http2-server-push': 'http2',
  jsonp: 'cors',
  'load-balancing': 'reverse-proxy',
};
const localePrefixes = ['', 'es/', 'ja/', 'pt-br/', 'zh-cn/'];

export function retiredCookbookRedirects(base = '/') {
  return Object.fromEntries(localePrefixes.flatMap((locale) =>
    Object.entries(retiredRecipes).map(([from, to]) =>
      [`/${locale}cookbook/${from}`, `${base}${locale}cookbook/${to}/`],
    ),
  ));
}

const retiredLegacyRedirects = Object.fromEntries(localePrefixes.flatMap((locale) =>
  Object.entries(retiredRecipes).map(([from, to]) =>
    [`/${locale}docs/cookbook/${from}`, `/${locale}cookbook/${to}/`],
  ),
));

export const redirects = {
  // Renamed, removed, index, Docusaurus category pages, and the old search page.
  '/docs': '/guide/quickstart/',
  '/docs/quick-start': '/guide/quickstart/',
  '/docs/start-server': '/guide/customization/',
  '/docs/category/guide': '/guide/quickstart/',
  '/docs/category/middleware': '/middleware/',
  '/docs/middleware/index': '/middleware/',
  '/docs/category/cookbook': '/cookbook/hello-world/',
  '/search': '/',
  // Bulk, generated from the current content tree.
  ...fromList('', 'guide/', pages('guide').filter((n) => !GUIDE_SKIP.has(n))),
  ...fromList('middleware/', 'middleware/', pages('middleware').filter((name) => name !== 'index')),
  ...fromList('cookbook/', 'cookbook/', pages('cookbook')),
  ...retiredCookbookRedirects(),
  ...retiredLegacyRedirects,
};
