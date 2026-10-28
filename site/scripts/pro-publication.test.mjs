import assert from 'node:assert/strict';
import test from 'node:test';
import { createHash } from 'node:crypto';
import { mkdtempSync, mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { tmpdir } from 'node:os';
import { parseFrontmatter } from '@astrojs/internal-helpers/frontmatter';
import { readPro } from './pro-source.mjs';
import { checkMDX, checkFrontmatter, checkJSONReferences, checkDestination } from './pro-publication.mjs';

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

const corpus = JSON.parse(readFileSync(new URL('./testdata/publication-corpus.json', import.meta.url)));
test('shared publication corpus', async (t) => {
  const current = readPro();
  assert.deepEqual(corpus.cases.filter((c) => c.currentPage).map((c) => c.name).sort(), [...current.files.keys()].filter((name) => name.startsWith('pages/')).sort(), 'current pages must all appear in the shared corpus');
  for (const c of corpus.cases) await t.test(c.id, () => {
    const files = c.currentPage ? new Map(current.files) : new Map(Object.entries(corpus.files));
    const content = c.currentPage ? files.get(c.name)?.toString() : c.content;
    assert.equal(typeof content, 'string', 'missing corpus page');
    files.set(c.name, content);
    const check = () => checkMDX(c.name, content, files, c.currentPage ? current.api : api);
    if (c.allowed) assert.doesNotThrow(check);
    else assert.throws(check, (error) => error.message === `${c.name}: ${c.reason}`, c.id);
  });
});

test('Astro frontmatter agrees with every accepted corpus case', async (t) => {
  const current = readPro();
  for (const c of corpus.cases) await t.test(c.id, () => {
    const files = c.currentPage ? new Map(current.files) : new Map(Object.entries(corpus.files));
    const content = c.currentPage ? files.get(c.name).toString() : c.content;
    files.set(c.name, content);
    const validate = () => checkMDX(c.name, content, files, c.currentPage ? current.api : api);
    let astro;
    try { astro = parseFrontmatter(content); }
    catch {
      assert.throws(validate, undefined, 'Astro parse failures must be rejected');
      return;
    }
    if (!astro.rawFrontmatter) return;
    try { validate(); } catch { return; }
    const { frontmatter } = checkFrontmatter(c.name, content);
    assert.deepEqual(frontmatter, astro.frontmatter);
    assert.ok(Object.keys(frontmatter).every((key) => ['title', 'description', 'sidebar'].includes(key)));
    if (frontmatter.sidebar) assert.ok(Object.keys(frontmatter.sidebar).every((key) => key === 'order'));
  });
});

test('synthetic corpus traverses readPro and its source include validation', async (t) => {
  for (const c of corpus.cases.filter((c) => !c.currentPage)) await t.test(c.id, (t) => {
    const root = mkdtempSync(join(tmpdir(), 'pro-publication-corpus-'));
    t.after(() => rmSync(root, { recursive: true, force: true }));
    const payloads = {
      ...corpus.files,
      'go.mod': 'module echo.labstack.com/pro-site-snapshot\n\ngo 1.26.0\n',
      'config-fields.json': JSON.stringify({ ...api, module: 'echo.labstack.com/pro', revision: 'v0.1.0' }),
      'problems.json': '[]',
      [c.name]: c.content,
    };
    const provenance = { module: 'echo.labstack.com/pro', version: 'v0.1.0', commit: 'a'.repeat(40), hashes: {} };
    for (const [name, content] of Object.entries(payloads)) {
      mkdirSync(dirname(join(root, name)), { recursive: true });
      writeFileSync(join(root, name), content);
      provenance.hashes[name] = createHash('sha256').update(content).digest('hex');
    }
    writeFileSync(join(root, 'pro-source.json'), JSON.stringify(provenance));
    const check = () => readPro(root);
    if (c.allowed) assert.doesNotThrow(check);
    // The bundle's raw private-reference check may reject before MDX parsing.
    else assert.throws(check, (error) => error.message === `${c.name}: ${c.reason}` ||
      c.reason === 'private reference' && error.message === `Private reference in Pro bundle: ${c.name}`, c.id);
  });
});
