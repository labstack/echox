import assert from 'node:assert/strict';
import { copyFileSync, mkdtempSync, mkdirSync, readFileSync, rmSync, symlinkSync, writeFileSync } from 'node:fs';
import { spawnSync } from 'node:child_process';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { test } from 'node:test';
import { fileURLToPath } from 'node:url';
import { echoSource, moduleSources, preparedSource, createWorkspace, run } from './module-sources.mjs';
import { reportApiChanges } from './report-api.mjs';

function fixture(t, referenceVersion = 'v5.4.0') {
  const repo = mkdtempSync(join(tmpdir(), 'echox-modules-'));
  t.after(() => rmSync(repo, { recursive: true, force: true }));
  const site = join(repo, 'site');
  const reference = join(repo, 'reference');
  mkdirSync(join(site, 'src/generated'), { recursive: true });
  mkdirSync(reference);
  writeFileSync(join(repo, 'go.mod'), 'module example.test/cookbook\n\ngo 1.25.0\n\nrequire github.com/labstack/echo/v5 v5.4.0\n');
  writeFileSync(join(reference, 'go.mod'), `module example.test/reference\n\ngo 1.25.0\n\nrequire github.com/labstack/echo/v5 ${referenceVersion}\n`);
  const sums = readFileSync(new URL('../../go.sum', import.meta.url));
  writeFileSync(join(repo, 'go.sum'), sums);
  writeFileSync(join(reference, 'go.sum'), sums);
  writeFileSync(join(site, 'external-sources.json'), '{}\n');
  return { repo, site, reference };
}

test('source selection follows the two Go modules and rejects a mismatched Echo bump', (t) => {
  const { site } = fixture(t, 'v5.0.0');
  assert.throws(() => moduleSources(site), /Echo versions differ: go.mod selects v5.4.0, reference\/go.mod selects v5.0.0/);
});

test('external version comes from reference/go.mod; conflicting cookbook versions fail', (t) => {
  const { repo, site, reference } = fixture(t);
  writeFileSync(join(site, 'external-sources.json'), JSON.stringify({ jwt: { module: 'github.com/labstack/echo-jwt/v5', package: 'echojwt' } }));
  run('go', ['mod', 'edit', '-require=github.com/labstack/echo-jwt/v5@v5.0.2'], reference);
  assert.equal(moduleSources(site).external.jwt.version, 'v5.0.2');
  run('go', ['mod', 'edit', '-require=github.com/labstack/echo-jwt/v5@v5.0.0'], repo);
  assert.throws(() => moduleSources(site), /echo-jwt\/v5 versions differ/);
});

test('stable source metadata supports module proxies with and without origin information', () => {
  assert.deepEqual(echoSource({ Version: 'v5.4.0' }), {
    channel: 'stable', repository: 'https://github.com/labstack/echo.git', revision: 'v5.4.0', release: 'v5.4.0',
  });
  assert.equal(echoSource({ Version: 'v5.4.0', Origin: { Hash: 'abc123' } }).revision, 'v5.4.0');
  assert.equal(echoSource({ Version: 'v5.4.1-0.20260930160000-5196b9b0ad8f' }).revision, '5196b9b0ad8f');
});

test('a source prepared before a module bump is rejected', (t) => {
  const { repo, site, reference } = fixture(t);
  const source = echoSource({ Version: 'v5.4.0' });
  writeFileSync(join(site, 'src/generated/echo-source.json'), JSON.stringify(source));
  writeFileSync(join(site, 'src/generated/config-fields.json'), JSON.stringify({ module: 'github.com/labstack/echo/v5', revision: source.revision }));
  const channel = process.env.DOCS_CHANNEL;
  const override = process.env.ECHO_SOURCE_DIR;
  delete process.env.DOCS_CHANNEL;
  delete process.env.ECHO_SOURCE_DIR;
  t.after(() => {
    if (channel === undefined) delete process.env.DOCS_CHANNEL; else process.env.DOCS_CHANNEL = channel;
    if (override === undefined) delete process.env.ECHO_SOURCE_DIR; else process.env.ECHO_SOURCE_DIR = override;
  });
  assert.equal(preparedSource(site).source.release, 'v5.4.0');
  for (const directory of [repo, reference]) run('go', ['mod', 'edit', '-require=github.com/labstack/echo/v5@v5.0.0'], directory);
  assert.throws(() => preparedSource(site), /Run source:prepare/);
});

test('Go creates a workspace for the actual modules, replacing an old workspace', (t) => {
  const { repo, reference } = fixture(t);
  const workspace = join(repo, '.cache/work');
  mkdirSync(workspace, { recursive: true });
  writeFileSync(join(workspace, 'go.work'), 'invalid old workspace');
  const path = createWorkspace(workspace, repo, reference);
  const data = JSON.parse(run('go', ['work', 'edit', '-json'], workspace, { GOWORK: path }));
  assert.deepEqual(data.Use.map((entry) => entry.DiskPath).sort(), [repo, reference].sort());
  assert.ok(data.Go);
});

test('API summary reports success and retains the actionable diff on failure', (t) => {
  const { repo } = fixture(t);
  const summary = join(repo, 'summary.md');
  reportApiChanges('stable', 'v5.4.0', [], summary);
  reportApiChanges('next', 'abc123', ['Config.Field: null -> {"type":"bool"}'], summary);
  const text = readFileSync(summary, 'utf8');
  assert.match(text, /stable middleware API \(v5.4.0\)[\s\S]*No differences/);
  assert.match(text, /next middleware API \(abc123\)[\s\S]*Review required: 1 API difference/);
  assert.match(text, /Config.Field: null ->/);
});

test('accepting stale external source fails without changing either reviewed baseline', (t) => {
  const { site, reference } = fixture(t);
  run('go', ['mod', 'edit', '-require=github.com/labstack/echo-jwt/v5@v5.0.2'], reference);
  writeFileSync(join(site, 'external-sources.json'), JSON.stringify({ jwt: { module: 'github.com/labstack/echo-jwt/v5' } }));
  writeFileSync(join(site, 'src/generated/echo-source.json'), JSON.stringify(echoSource({ Version: 'v5.4.0' })));
  writeFileSync(join(site, 'src/generated/config-fields.json'), JSON.stringify({ module: 'github.com/labstack/echo/v5', revision: 'v5.4.0', configs: [], functions: [] }));
  writeFileSync(join(site, 'src/generated/jwt.json'), JSON.stringify({ module: 'github.com/labstack/echo-jwt/v5', revision: 'v5.0.0', configs: [], functions: [] }));
  for (const name of ['reference-baseline.json', 'external-baseline.json']) writeFileSync(join(site, name), '{}\n');
  mkdirSync(join(site, 'scripts'));
  for (const name of ['accept-source.mjs', 'module-sources.mjs']) copyFileSync(new URL(name, import.meta.url), join(site, 'scripts', name));
  const result = spawnSync(process.execPath, [join(site, 'scripts/accept-source.mjs')], {
    encoding: 'utf8', env: { ...process.env, DOCS_CHANNEL: 'stable', ECHO_SOURCE_DIR: '' },
  });
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /jwt: generated source is stale/);
  for (const name of ['reference-baseline.json', 'external-baseline.json']) assert.equal(readFileSync(join(site, name), 'utf8'), '{}\n');
});

test('source-check CLI reports API field changes even when an external version also changed', (t) => {
  const { site, reference } = fixture(t);
  run('go', ['mod', 'edit', '-require=github.com/labstack/echo-jwt/v5@v5.0.2'], reference);
  writeFileSync(join(site, 'external-sources.json'), JSON.stringify({ jwt: { module: 'github.com/labstack/echo-jwt/v5' } }));
  writeFileSync(join(site, 'src/generated/echo-source.json'), JSON.stringify(echoSource({ Version: 'v5.4.0' })));
  writeFileSync(join(site, 'src/generated/config-fields.json'), JSON.stringify({ module: 'github.com/labstack/echo/v5', revision: 'v5.4.0', configs: [], functions: [] }));
  writeFileSync(join(site, 'src/generated/jwt.json'), JSON.stringify({
    module: 'github.com/labstack/echo-jwt/v5', revision: 'v5.0.2',
    configs: [{ name: 'Config', fields: [{ name: 'AddedField', type: 'bool' }] }], functions: [],
  }));
  writeFileSync(join(site, 'reference-baseline.json'), JSON.stringify({ configs: {}, functions: {} }));
  writeFileSync(join(site, 'external-baseline.json'), JSON.stringify({ jwt: {
    module: 'github.com/labstack/echo-jwt/v5', version: 'v5.0.0', configs: {}, functions: {},
  } }));
  writeFileSync(join(site, 'reference-pages.json'), '{}\n');
  mkdirSync(join(site, 'scripts'));
  for (const name of ['check-source.mjs', 'module-sources.mjs', 'report-api.mjs']) copyFileSync(new URL(name, import.meta.url), join(site, 'scripts', name));
  symlinkSync(fileURLToPath(new URL('../plugins', import.meta.url)), join(site, 'plugins'), 'dir');
  copyFileSync(new URL('../src/locales.mjs', import.meta.url), join(site, 'src/locales.mjs'));
  const summary = join(site, 'summary.md');
  const result = spawnSync(process.execPath, [join(site, 'scripts/check-source.mjs')], {
    encoding: 'utf8', env: { ...process.env, DOCS_CHANNEL: 'stable', ECHO_SOURCE_DIR: '', GITHUB_STEP_SUMMARY: summary },
  });
  assert.notEqual(result.status, 0);
  for (const output of [result.stderr, readFileSync(summary, 'utf8')]) {
    assert.match(output, /echo-jwt\/v5@v5.0.0 -> github.com\/labstack\/echo-jwt\/v5@v5.0.2/);
    assert.match(output, /jwt.Config.AddedField: null ->/);
  }
});
