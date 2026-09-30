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

function check(site, script, ...args) {
  return spawnSync(process.execPath, [join(site, 'scripts', script), ...args], { encoding: 'utf8' });
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
  assert.match(accepted.stderr, /new route\(s\) are not yet protected[\s\S]*\/new-page\/[\s\S]*npm run site:accept/);
  assert.deepEqual(JSON.parse(readFileSync(join(site, 'route-baseline.json'), 'utf8')).routes, ['/']);
  assert.equal(check(site, 'check-site.mjs', '--accept').status, 0);
  rmSync(join(site, 'dist/new-page'), { recursive: true });
  const removedNewPage = check(site, 'check-site.mjs');
  assert.notEqual(removedNewPage.status, 0);
  assert.match(removedNewPage.stderr, /Removed route: \/new-page\//);
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

test('origin budget counts resource URLs, excluding links, preconnects and inline text', () => {
  const html = `
    <a href="https://reference.example/">Reference</a>
    <img src="https://images.example/image.png">
    <link rel="preload" as="font" href="https://font-preloads.example/font.woff2">
    <link rel="modulepreload" href="https://modules.example/module.js">
    <link rel="icon" href="https://icons.example/favicon.ico">
    <iframe src="https://frames.example/embed"></iframe>
    <source src="https://media.example/video.webm" srcset="https://responsive.example/small.webp 1x, //retina.example/large.webp 2x">
    <video poster="https://posters.example/video.jpg"></video>
    <object data="https://objects.example/document.pdf"></object>
    <input type="image" src="https://buttons.example/submit.png">
    <input type="text" src="https://ignored-input.example/unused.png">
    <img src="data:image/png;base64,AAAA">
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
  assert.deepEqual(externalDomains(html), ['buttons.example', 'font-preloads.example', 'frames.example', 'icons.example', 'images.example', 'media.example', 'modules.example', 'objects.example', 'posters.example', 'responsive.example', 'retina.example', 'scripts.example', 'styles.example']);
});

test('external links pass the performance CLI; new resource origins fail', (t) => {
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
    '<link rel="modulepreload" href="https://new-resource.example/module.js">',
    '<link rel="preload" as="font" href="https://new-resource.example/font.woff2">',
    '<link rel="icon" href="https://new-resource.example/favicon.ico">',
    '<img src="https://new-resource.example/image.png">',
    '<iframe src="https://new-resource.example/embed"></iframe>',
    '<source src="https://new-resource.example/video.webm">',
  ]) {
    file(site, 'dist/index.html', `<html>${resource}</html>`);
    const result = check(site, 'check-performance.mjs');
    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /added external origin new-resource.example/);
  }
  assert.deepEqual(JSON.parse(readFileSync(join(site, 'performance-baseline.json'), 'utf8')), baseline);
});
