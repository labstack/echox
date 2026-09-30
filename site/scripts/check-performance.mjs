import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { gzipSync } from 'node:zlib';
import { externalDomains } from './performance-origins.mjs';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const distDir = join(siteDir, 'dist');
const baselineFile = join(siteDir, 'performance-baseline.json');
const routes = ['/', '/middleware/cors/', '/next/', '/next/middleware/cors/'];

function bytes(path) {
  const data = readFileSync(path);
  return { raw: data.length, gzip: gzipSync(data).length };
}

function metrics(route) {
  const htmlFile = join(distDir, route, 'index.html');
  const html = readFileSync(htmlFile, 'utf8');
  const assetPaths = [...html.matchAll(/(?:src|href)="(\/(?:next\/)?_astro\/[^"?#]+)"/g)]
    .map((match) => match[1]);
  const assets = [...new Set(assetPaths)].filter((path) => /\.(?:js|css)$/.test(path));
  const assetBytes = assets.reduce((sum, path) => {
    const size = bytes(join(distDir, path));
    return { raw: sum.raw + size.raw, gzip: sum.gzip + size.gzip };
  }, { raw: 0, gzip: 0 });
  return { html: bytes(htmlFile), firstLoadAssets: assetBytes, externalDomains: externalDomains(html) };
}

const current = Object.fromEntries(routes.map((route) => [route, metrics(route)]));
if (process.argv.includes('--accept')) {
  writeFileSync(baselineFile, `${JSON.stringify(current, null, 2)}\n`);
  console.log('Accepted stable and next HTML/asset sizes for review');
  process.exit(0);
}
const baseline = JSON.parse(readFileSync(baselineFile, 'utf8'));
const problems = [];
for (const route of routes) {
  const before = baseline[route];
  const after = current[route];
  for (const [part, metric] of [['html', 'gzip'], ['firstLoadAssets', 'gzip']]) {
    if (after[part][metric] > before[part][metric] * 1.15 + 1024) {
      problems.push(`${route} ${part} gzip grew from ${before[part][metric]} to ${after[part][metric]} bytes`);
    }
  }
  for (const domain of after.externalDomains) {
    if (!before.externalDomains.includes(domain)) problems.push(`${route} added external origin ${domain}`);
  }
}
if (problems.length) throw new Error(`Review performance baseline changes:\n${problems.join('\n')}`);
console.log('Stable and next HTML/asset sizes are within the recorded baseline');
