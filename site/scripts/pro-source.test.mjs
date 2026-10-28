import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { cpSync, linkSync, symlinkSync, unlinkSync, mkdtempSync, readFileSync, rmSync, renameSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import test from 'node:test';
import { readPro, apiSnapshot, siteDir } from './pro-source.mjs';

test('Pro snapshot has complete catalog coverage and source-owned programs', () => {
  const bundle = readPro();
  assert.equal(bundle.provenance.version, 'v0.1.0');
  for (const [name, content] of bundle.files) {
    if (name.startsWith('pages/') && content.toString().includes('package main')) assert.fail(`${name} pasted a program`);
  }
  const api = apiSnapshot(bundle);
  assert.ok(api.configs['quota.MiddlewareConfig']);
  assert.ok(api.functions['quota.MiddlewareConfig.ToMiddleware']);
  assert.ok(!Object.values(api.functions).some((signature) => signature.includes('{')));
});

test('unreviewed bundle files and modified payloads fail closed', (t) => {
  const root = mkdtempSync(join(tmpdir(), 'pro-source-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  cpSync(join(siteDir, 'pro-source'), root, { recursive: true });
  const file = join(root, 'pages/auth.mdx');
  const original = readFileSync(file);
  writeFileSync(file, 'modified');
  assert.throws(() => readPro(root), /hash mismatch/);
  writeFileSync(file, original);
  writeFileSync(join(root, 'PRIVATE.md'), 'unreviewed');
  assert.throws(() => readPro(root), /inventory differs/);
});

test('all Pro code languages require source includes even with updated hashes', (t) => {
  const root = mkdtempSync(join(tmpdir(), 'pro-source-fence-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  cpSync(join(siteDir, 'pro-source'), root, { recursive: true });
  const file = join(root, 'pages/auth.mdx');
  const provenancePath = join(root, 'pro-source.json');
  const original = readFileSync(file, 'utf8');
  const provenance = JSON.parse(readFileSync(provenancePath, 'utf8'));
  for (const lang of ['sh', 'yaml', 'json', 'go']) {
    const page = `${original}\n\x60\x60\x60${lang}\npasted content\n\x60\x60\x60\n`;
    writeFileSync(file, page);
    provenance.hashes['pages/auth.mdx'] = createHash('sha256').update(page).digest('hex');
    writeFileSync(provenancePath, JSON.stringify(provenance));
    assert.throws(() => readPro(root), /source fences must be empty/);
  }
});

test('valid hashes cannot authorize executable MDX or private references', (t) => {
  const root = mkdtempSync(join(tmpdir(), 'pro-source-policy-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  cpSync(join(siteDir, 'pro-source'), root, { recursive: true });
  const provenancePath = join(root, 'pro-source.json');
  const provenance = JSON.parse(readFileSync(provenancePath, 'utf8'));
  function replace(name, content) {
    writeFileSync(join(root, name), content);
    provenance.hashes[name] = createHash('sha256').update(content).digest('hex');
    writeFileSync(provenancePath, JSON.stringify(provenance));
  }
  const name = 'pages/auth.mdx';
  const original = readFileSync(join(root, name), 'utf8');
  for (const content of [
    "import fs from 'node:fs';\n\n{fs.readFileSync('/private/source.go', 'utf8')}",
    "import privateSource from '/private/source.go?raw';\n\n<pre>{privateSource}</pre>",
    '[source](../internal/secrets.go)', '[notes](../docs/ARCHITECTURE%2emd)',
  ]) {
    replace(name, content);
    assert.throws(() => readPro(root));
  }
  replace(name, original);
  const apiName = 'config-fields.json';
  const api = JSON.parse(readFileSync(join(root, apiName), 'utf8'));
  api.configs[0].file = '../%69nternal/private.go';
  replace(apiName, JSON.stringify(api));
  assert.throws(() => readPro(root), /Private|private|source location/);
});

test('bundle filesystem refuses symlinks and hard links before reading bytes', (t) => {
  const root = mkdtempSync(join(tmpdir(), 'pro-source-links-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  cpSync(join(siteDir, 'pro-source'), root, { recursive: true });
  const path = join(root, 'pages/auth.mdx');
  const original = readFileSync(path);
  unlinkSync(path);
  symlinkSync('tenant.mdx', path);
  assert.throws(() => readPro(root), /must not contain symlinks/);
  unlinkSync(path);
  linkSync(join(root, 'pages/tenant.mdx'), path);
  assert.throws(() => readPro(root), /must not contain hard links/);
  unlinkSync(path);
  writeFileSync(path, original);
  // Provenance is an input too.
  const provenance = join(root, 'pro-source.json');
  linkSync(provenance, join(root, 'linked-provenance'));
  assert.throws(() => readPro(root), /must not contain hard links/);
});

test('shared filesystem publication corpus', async (t) => {
  const corpus = JSON.parse(readFileSync(new URL('./testdata/publication-corpus.json', import.meta.url)));
  for (const c of corpus.inputs) await t.test(c.id, (t) => {
    const root = mkdtempSync(join(tmpdir(), 'pro-source-input-'));
    t.after(() => rmSync(root, { recursive: true, force: true }));
    cpSync(join(siteDir, 'pro-source'), root, { recursive: true });
    const target = join(root, 'pages/auth.mdx');
    if (c.kind === 'case') {
      renameSync(target, join(root, 'pages/Auth.mdx'));
      const provenancePath = join(root, 'pro-source.json');
      const provenance = JSON.parse(readFileSync(provenancePath, 'utf8'));
      provenance.hashes['pages/Auth.mdx'] = provenance.hashes['pages/auth.mdx'];
      delete provenance.hashes['pages/auth.mdx'];
      writeFileSync(provenancePath, JSON.stringify(provenance));
    } else {
      unlinkSync(target);
      if (c.kind === 'symlink') symlinkSync('tenant.mdx', target);
      else if (c.kind === 'hardlink') linkSync(join(root, 'pages/tenant.mdx'), target);
      else assert.fail(`unknown input kind ${c.kind}`);
    }
    assert.throws(() => readPro(root), (error) => error.message.startsWith(c.reason));
  });
});
