import assert from 'node:assert/strict';
import test from 'node:test';
import { checkMDX, checkJSONReferences, checkDestination } from './pro-publication.mjs';

const reviewedImport = "import ConfigReference from '../../../../components/ConfigReference.astro';";
const api = { configs: [{ name: 'auth.Config' }] };
const files = new Map();

test('only reviewed MDX imports, components and literal attributes are accepted', () => {
  for (const content of [
    "import fs from 'node:fs';", "> import fs from 'node:fs';", "- import fs from 'node:fs';", "import secret from '/private/source.go?raw';",
    "import secret from '../private/source.go';", "import secret from './local.mdx';",
    'export const secret = 1;', '{process.env.SECRET}', "{await import('node:fs')}",
    '<pre>{privateSource}</pre>', '<script src="https://example.org/script.js" />',
    `${reviewedImport}\n<ConfigReference source={process.env.SECRET} name="auth.Config" />`,
    `${reviewedImport}\n<ConfigReference source="pro" name="auth.Config" {...props} />`,
    `${reviewedImport}\n<ConfigReference source="pro" name="unknown.Config" />`,
    "import\nfs from 'node:fs';",
  ]) assert.throws(() => checkMDX('pages/reference/auth.mdx', content, files, api), undefined, content);
  checkMDX('pages/reference/auth.mdx', `${reviewedImport}\n\n<ConfigReference source="pro" name="auth.Config" />\n\n{/* Diataxis: reference */}`, files, api);
});

test('links, images, definitions, HTML and JSON reject decoded private destinations', () => {
  for (const destination of [
    '/pro/auth/private.go', '/pro/../../outside.md', 'https://github.com/labstack/./echo-pro/blob/main/private.go', 'https://github.com./labstack/echo-pro/blob/main/private.go', 'https://github。com/labstack/echo-pro/blob/main/private.go', '../internal/secrets.go', '../docs/ARCHITECTURE%2emd', '../docs/private/plan', '../.tmp/log', '../fixtures/data.json',
    '/pro/%69nternal/key', '/pro/docs%252fTESTING%252emd', 'file:///private/source.go', '../outside.go',
    'https://github.com/labstack/echo-pro/blob/main/auth/key.go', 'https://raw.githubusercontent.com/labstack/echo-pro/main/auth/key.go',
    'https://example.org/a/../internal/key', 'https://example.org/docs/ARCHITECTURE.md', '../&#105;nternal/key',
  ]) {
    for (const content of [
      `[source](${destination})`, `![image](${destination})`, `[ref][x]\n\n[x]: ${destination}`,
      `<a href="${destination}">link</a>`, `<img src="${destination}" />`,
    ]) assert.throws(() => checkMDX('pages/auth.mdx', content, files, api), undefined, content);
    for (const key of ['href', 'src', 'url', 'path', 'source', 'file']) assert.throws(() => checkJSONReferences('config-fields.json', JSON.stringify({ nested: [{ [key]: destination }] }), files));
  }
  checkMDX('pages/auth.mdx', '[bundle file](AuthGuide.mdx)', new Map([['pages/AuthGuide.mdx', '']]), api);
  for (const content of ['[quota](/pro/quota/)', '[Go](https://go.dev/doc/)', '`{notExecutable}`', '[part](#limits)']) {
    checkMDX('pages/auth.mdx', content, files, api);
  }
});

test('JSON source locations name reviewed public packages and type URLs use HTTP', () => {
  for (const value of [
    { file: 'outside/source.go' }, { file: '../auth/config.go' },
    { type: 'javascript:alert(1)' }, { type: 'file:/private/source.go' },
  ]) assert.throws(() => checkJSONReferences('config-fields.json', JSON.stringify(value), files));
  checkJSONReferences('config-fields.json', JSON.stringify({ file: 'auth/middleware.go', type: '*slog.Logger' }), files);
});

test('unframed URLs and JSON documentation cannot point at private repositories', () => {
  assert.throws(() => checkDestination('pages/auth.mdx', 'https://user:test-marker@example.org/pro/', files), (error) => !String(error).includes('test-marker'));
  for (const target of [
    'https://github%2ecom/labstack/echo-pro/blob/main/private.go',
    'www.github.com./labstack/echo-pro/blob/main/private.go',
  ]) {
    assert.throws(() => checkMDX('pages/auth.mdx', `Source: ${target}.`, files, api));
    assert.throws(() => checkJSONReferences('config-fields.json', JSON.stringify({ doc: `Source: ${target}.` }), files));
  }
});
