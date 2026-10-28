import assert from 'node:assert/strict';
import { existsSync, readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';
import { readPro, siteDir } from './pro-source.mjs';

const { files } = readPro();
const dist = join(siteDir, 'dist');
const sitemap = readdirSync(dist).filter((n) => /^sitemap-\d+\.xml$/.test(n)).map((n) => readFileSync(join(dist, n), 'utf8')).join('');
let count = 0;
for (const name of files.keys()) {
  if (!name.startsWith('pages/')) continue;
  const slug = name.slice(6).replace(/\.mdx$/, '').replace(/\/index$/, '');
  const route = `/pro/${slug}/`;
  const html = readFileSync(join(dist, route, 'index.html'), 'utf8');
  assert.match(html, /<h1\b/);
  assert.match(html, /<title>[^<]+ · Echo Pro<\/title>/);
  assert.match(html, /data-pagefind-filter="product\[content\]" content="Echo Pro"/);
  assert.match(html, /Echo Pro is commercial software/);
  assert.doesNotMatch(html, /hreflang=|widget\.kapa\.ai|github\.com\/labstack\/echo-pro|Edit page/);
  assert.ok(sitemap.includes(`https://echo.labstack.com${route}`), `${route} absent from release sitemap`);
  assert.ok(!existsSync(join(dist, 'next', route, 'index.html')), `${route} has next copy`);
  for (const locale of ['es', 'ja', 'pt-br', 'zh-cn']) assert.ok(!existsSync(join(dist, locale, route, 'index.html')), `${route} has locale fallback`);
  count++;
}
for (const pkg of ['auth', 'tenant', 'quota', 'webhook', 'audit', 'execution', 'admin']) {
  const overview = readFileSync(join(dist, 'pro', pkg, 'index.html'), 'utf8');
  assert.ok(overview.includes(`/pro/reference/${pkg}/`), `${pkg} has no configuration reference`);
  assert.ok(overview.includes(`/pro/guides/${pkg}/`), `${pkg} has no task guide`);
  const html = readFileSync(join(dist, 'pro/reference', pkg, 'index.html'), 'utf8');
  assert.match(html, /echo-config-reference/);
  assert.match(html, /Fields from package source/);
  assert.match(html, /Defaults \/ requirements/);
  assert.match(html, /\.go:\d+/);
}
assert.match(readFileSync(join(dist, 'next/pro/index.html'), 'utf8'), /http-equiv="refresh"/);
const oss = readFileSync(join(dist, 'guide/installation/index.html'), 'utf8');
assert.ok(!oss.includes('href="/pro/install/"'), 'Pro sidebar leaked into OSS');
console.log(`Verified ${count} rendered Pro docs, release sitemap coverage, package tables and English/release-only routes`);
