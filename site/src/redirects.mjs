import { readdirSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { localePrefixes } from './locales.mjs';

// Old Docusaurus URLs -> new locations, generated so astro.config stays small.
// The old site served docs under /docs/ with guide pages flattened (slug
// overrides). On GitHub Pages each entry still emits its own static
// meta-refresh page (no server-side wildcards).
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

export function retiredCookbookRedirects(base = '/') {
  return Object.fromEntries(localePrefixes.flatMap((locale) =>
    Object.entries(retiredRecipes).map(([from, to]) =>
      [`/${locale}cookbook/${from}`, `${base}${locale}cookbook/${to}/`],
    ),
  ));
}

// Docusaurus used different locale names, and French now falls back to English.
export const legacyCookbookLocales = { '': '', 'zh-Hans/': 'zh-cn/', 'ja/': 'ja/', 'es/': 'es/', 'fr/': '' };
const legacyCookbookRedirects = Object.fromEntries(Object.entries(legacyCookbookLocales).flatMap(([fromLocale, toLocale]) =>
  [...pages('cookbook'), ...Object.keys(retiredRecipes)].map((from) =>
    [`/${fromLocale}docs/cookbook/${from}`, `/${toLocale}cookbook/${retiredRecipes[from] ?? from}/`],
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
  ...legacyCookbookRedirects,
};
