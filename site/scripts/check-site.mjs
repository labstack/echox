import { existsSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { extname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const distDir = join(siteDir, 'dist');
const baselineFile = join(siteDir, 'route-baseline.json');
const accept = process.argv.includes('--accept');
const origin = 'https://echo.labstack.com';

function walk(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const path = join(dir, entry.name);
    return entry.isDirectory() ? walk(path) : [path];
  });
}

const files = walk(distDir);
const htmlFiles = files.filter((file) => file.endsWith('.html'));
const routes = htmlFiles.map((file) => {
  const path = `/${relative(distDir, file).replaceAll('\\', '/')}`;
  return path.endsWith('/index.html') ? path.slice(0, -10) : path;
}).sort();

if (accept) {
  writeFileSync(baselineFile, `${JSON.stringify({ routes }, null, 2)}\n`);
  console.log(`Accepted ${routes.length} routes for review`);
  process.exit(0);
}

const expected = JSON.parse(readFileSync(baselineFile, 'utf8')).routes;
const problems = [];
for (const route of expected.filter((item) => !routes.includes(item))) problems.push(`Removed route: ${route}`);
for (const route of routes.filter((item) => !expected.includes(item))) problems.push(`New route needs baseline review: ${route}`);
const ids = new Map();
function idsIn(file) {
  if (!ids.has(file)) {
    const html = readFileSync(file, 'utf8');
    ids.set(file, new Set([...html.matchAll(/\bid="([^"]+)"/g)].map((match) => match[1])));
  }
  return ids.get(file);
}

let checked = 0;
for (const file of htmlFiles) {
  const html = readFileSync(file, 'utf8');
  const markup = html.replace(/<script\b[^>]*>[\s\S]*?<\/script>/gi, '').replace(/<pre\b[^>]*>[\s\S]*?<\/pre>/gi, '');
  const route = routes.find((item) => file === join(distDir, item, 'index.html') || file === join(distDir, item));
  const redirect = route.startsWith('/docs/') || route === '/docs/' || route === '/search/';
  if (!redirect && !/\<html\b[^>]*\blang="[^"]+"/.test(html)) problems.push(`${route}: missing document language`);
  if (route.endsWith('/404.html')) continue;
  for (const [tag] of markup.matchAll(/<img\b[^>]*>/g)) {
    if (!/\balt(?:\s|=)/.test(tag)) problems.push(`${route}: image missing alt text`);
  }
  for (const [, names] of markup.matchAll(/\baria-labelledby="([^"]+)"/g)) {
    for (const id of names.split(/\s+/)) {
      if (id && !idsIn(file).has(id)) problems.push(`${route}: aria-labelledby references missing #${id}`);
    }
  }
  if (route.startsWith('/next/')) {
    for (const [tag] of markup.matchAll(/<a\b[^>]*>/g)) {
      if (tag.includes('data-version-switch')) continue;
      const href = tag.match(/\bhref="(\/[^"]+)"/)?.[1];
      if (href && /^\/(?:es\/|ja\/|pt-br\/|zh-cn\/)?(?:guide|middleware|cookbook)\//.test(href)) {
        problems.push(`${route}: next content links to stable route ${href}`);
      }
    }
  }
  for (const match of markup.matchAll(/\b(?:href|src)="([^"]+)"/g)) {
    const value = match[1].replaceAll('&amp;', '&');
    if (!value || value.startsWith('data:') || value.startsWith('mailto:') || value.startsWith('javascript:')) continue;
    let target;
    try { target = new URL(value, `${origin}${route}`); } catch { continue; }
    if (target.origin !== origin) continue;
    let pathname;
    try { pathname = decodeURIComponent(target.pathname); } catch { problems.push(`${route}: malformed URL ${value}`); continue; }
    const resolved = resolve(distDir, `.${pathname}`);
    if (!resolved.startsWith(`${distDir}/`) && resolved !== distDir) { problems.push(`${route}: path escapes site ${value}`); continue; }
    const targetFile = extname(resolved) ? resolved : join(resolved, 'index.html');
    if (!existsSync(targetFile) || !statSync(targetFile).isFile()) { problems.push(`${route}: missing ${value}`); continue; }
    if (target.hash && targetFile.endsWith('.html') && !target.hash.startsWith('#:~:text=')) {
      const fragment = decodeURIComponent(target.hash.slice(1));
      if (fragment && !idsIn(targetFile).has(fragment)) problems.push(`${route}: missing anchor ${value}`);
    }
    checked++;
  }
}

for (const locale of ['es', 'ja', 'pt-br', 'zh-cn']) {
  if (!routes.includes(`/${locale}/middleware/`)) problems.push(`Missing ${locale} middleware task page`);
}
for (const path of ['sitemap-index.xml', 'pagefind/pagefind.js']) {
  if (!existsSync(join(distDir, path))) problems.push(`Missing ${path}`);
}
if (problems.length) throw new Error(`${problems.length} site check failure(s):\n${problems.join('\n')}`);
console.log(`Checked ${routes.length} routes and ${checked} internal links/assets; all resolve`);
