import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { copyFileSync, mkdtempSync, mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join } from 'node:path';
import { test } from 'node:test';
import { externalDomains } from './performance-origins.mjs';

function fixture(t, scripts) {
  const site = mkdtempSync(join(tmpdir(), 'echox-site-check-'));
  t.after(() => rmSync(site, { recursive: true, force: true }));
  mkdirSync(join(site, 'scripts'));
  mkdirSync(join(site, 'dist'));
  for (const name of scripts) copyFileSync(new URL(name, import.meta.url), join(site, 'scripts', name));
  return site;
}

function file(site, path, contents) {
  const target = join(site, path);
  mkdirSync(dirname(target), { recursive: true });
  writeFileSync(target, contents);
}

function check(site, script) {
  return spawnSync(process.execPath, [join(site, 'scripts', script)], { encoding: 'utf8' });
}

test('new routes pass without accepting a baseline; removals and new-page broken links still fail', (t) => {
  const site = fixture(t, ['check-site.mjs']);
  file(site, 'route-baseline.json', JSON.stringify({ routes: ['/'] }));
  file(site, 'dist/index.html', '<html lang="en"><script src="https://widget.kapa.ai/widget.js"></script><button id="echo-ask-ai"></button></html>');
  for (const locale of ['es', 'ja', 'pt-br', 'zh-cn']) file(site, `dist/${locale}/middleware/index.html`, `<html lang="${locale}"></html>`);
  for (const path of ['sitemap-index.xml', 'pagefind/pagefind.js']) file(site, `dist/${path}`, '');
  for (const path of ['llms.txt', 'next/llms.txt']) file(site, `dist/${path}`, '# Echo\n');
  file(site, 'dist/new-page/index.html', '<html lang="en"><a href="/">Home</a></html>');
  const accepted = check(site, 'check-site.mjs');
  assert.equal(accepted.status, 0, accepted.stderr);
  file(site, 'route-baseline.json', JSON.stringify({ routes: ['/', '/old-page/'] }));
  const removed = check(site, 'check-site.mjs');
  assert.notEqual(removed.status, 0);
  assert.match(removed.stderr, /Removed route: \/old-page\//);
  file(site, 'route-baseline.json', JSON.stringify({ routes: ['/'] }));
  file(site, 'dist/new-page/index.html', '<html lang="en"><a href="/#missing">Broken anchor</a></html>');
  const broken = check(site, 'check-site.mjs');
  assert.notEqual(broken.status, 0);
  assert.match(broken.stderr, /new-page\/: missing anchor/);
});

test('origin budget counts scripts and stylesheets, excluding links, preconnects and inline text', () => {
  const html = `
    <a href="https://reference.example/">Reference</a>
    <img src="https://images.example/image.png">
    <link rel="preconnect" href="https://fonts.example/">
    <script data-src="https://data-attribute.example/ignored.js"></script>
    <link rel="preconnect" title="rel='stylesheet' href='https://title.example/ignored.css'">
    <!-- <script src="https://comment.example/ignored.js"></script> -->
    <script>const text = '<link rel="stylesheet" href="https://inline.example/style.css">';</script>
    <script type="module" src='//scripts.example/main.js'></script>
    <SCRIPT SRC=https://scripts.example/other.js></SCRIPT>
    <link REL='alternate stylesheet' HREF='https://styles.example/css?x=1&amp;y=2'>
    <script src="/local.js"></script>
    <link rel="stylesheet" href="https://echo.labstack.com/local.css">
  `;
  assert.deepEqual(externalDomains(html), ['scripts.example', 'styles.example']);
});

test('external links pass the performance CLI; a new script or stylesheet origin fails', (t) => {
  const site = fixture(t, ['check-performance.mjs', 'performance-origins.mjs']);
  const routes = ['/', '/middleware/cors/', '/next/', '/next/middleware/cors/'];
  const baseline = Object.fromEntries(routes.map((route) => [route, {
    html: { gzip: 1000 }, firstLoadAssets: { gzip: 0 }, externalDomains: [],
  }]));
  file(site, 'performance-baseline.json', JSON.stringify(baseline));
  for (const route of routes) file(site, `dist${route}index.html`, '<html><a href="https://new-reference.example/">Reference</a></html>');
  const links = check(site, 'check-performance.mjs');
  assert.equal(links.status, 0, links.stderr);
  for (const resource of [
    '<script src="https://new-resource.example/main.js"></script>',
    '<link rel="stylesheet" href="https://new-resource.example/style.css">',
  ]) {
    file(site, 'dist/index.html', `<html>${resource}</html>`);
    const result = check(site, 'check-performance.mjs');
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /added external origin new-resource.example/);
  }
  assert.deepEqual(JSON.parse(readFileSync(join(site, 'performance-baseline.json'), 'utf8')), baseline);
});
