import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, readFileSync, readdirSync, rmSync, symlinkSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';
import { unified } from 'unified';
import remarkSourceCode, { sourceReference, sourceSnippet, withSourceCode } from './remark-source-code.mjs';
import { parseSourceDocument, visitCode } from './source-document.mjs';
import { localePrefixes } from '../src/locales.mjs';

function fixture(t) {
  const dir = mkdtempSync(join(tmpdir(), 'source-code-'));
  const root = join(dir, 'repo');
  mkdirSync(root);
  writeFileSync(join(root, 'server.go'), 'package main\n\n// docs:start handler\nfunc handler() {}\n// docs:end handler\n');
  t.after(() => rmSync(dir, { recursive: true, force: true }));
  return { dir, root };
}

function render(markdown, root, path = 'recipe.md') {
  const processor = unified().use(remarkSourceCode, { root });
  return processor.runSync(parseSourceDocument(markdown, path), { path });
}

function programFences(tree) {
  const fences = [];
  visitCode(tree, (node) => { if (['go', 'html'].includes(node.lang)) fences.push(node); });
  return fences;
}

function assertSourceOwned(tree) {
  const fences = programFences(tree);
  for (const fence of fences) {
    assert.ok(sourceReference(fence), 'has no pasted program');
    assert.equal(fence.value, '');
  }
  return fences;
}

test('includes a repository file and preserves highlighting metadata', (t) => {
  const { root } = fixture(t);
  const [code] = render('```go title="server.go" file=server.go {1}\n```', root).children;
  assert.equal(code.value, 'package main\n\nfunc handler() {}');
  assert.equal(code.lang, 'go');
  assert.equal(code.meta, 'title="server.go" {1}');
});

test('includes a named region inside a nested fence', (t) => {
  const { root } = fixture(t);
  const [quote] = render('> ~~~go file=server.go#handler\n> ~~~', root).children;
  assert.equal(quote.children[0].value, 'func handler() {}');
  assert.equal(quote.children[0].meta, null);
});

test('includes HTML and normalizes CRLF without dropping indentation', (t) => {
  const { root } = fixture(t);
  writeFileSync(join(root, 'index.html'), '<html>\r\n  <body></body>\r\n</html>\r\n');
  const [code] = render('```html file=index.html\n```', root).children;
  assert.equal(code.value, '<html>\n  <body></body>\n</html>');
});

test('ordinary code fences remain unchanged', (t) => {
  const { root } = fixture(t);
  const [code] = render('```go title="demo"\npackage main\n```', root).children;
  assert.equal(code.value, 'package main');
  assert.equal(code.meta, 'title="demo"');
});

test('quoted metadata never becomes a file reference, even beside a real include', (t) => {
  const { root } = fixture(t);
  for (const title of ['title="cp file=missing.go"', "title='cp file=missing.go'", 'title="cp \\"file=missing.go\\""']) {
    const parsedTitle = parseSourceDocument(`\x60\x60\x60go ${title}\n\x60\x60\x60`, 'recipe.md').children[0].meta;
    const [ordinary] = render(`\x60\x60\x60go ${title}\npackage main\n\x60\x60\x60`, root).children;
    assert.equal(ordinary.value, 'package main');
    assert.equal(ordinary.meta, parsedTitle);
    for (const meta of [`${title} file=server.go`, `file=server.go ${title}`]) {
      const [included] = render(`\x60\x60\x60go ${meta}\n\x60\x60\x60`, root).children;
      assert.equal(included.meta, parsedTitle);
      assert.match(included.value, /^package main/);
    }
  }
});

test('regions trim boundary blanks and marker removal keeps one separator between functions', (t) => {
  const { root } = fixture(t);
  writeFileSync(join(root, 'server.go'), '// docs:start first\n\nfunc first() {}\n\n// docs:end first\n\n// docs:start second\nfunc second() {}\n\n// docs:end second\n');
  assert.equal(render('```go file=server.go#first\n```', root).children[0].value, 'func first() {}');
  assert.equal(render('```go file=server.go\n```', root).children[0].value, 'func first() {}\n\nfunc second() {}');
});

test('Auto TLS and graceful shutdown render standalone programs without variant dispatch', () => {
  const repo = fileURLToPath(new URL('../..', import.meta.url));
  for (const recipe of ['auto-tls', 'graceful-shutdown']) {
    for (const path of ['server.go', 'custom-server/server.go']) {
      const code = render(`\x60\x60\x60go file=cookbook/${recipe}/${path}\n\x60\x60\x60`, repo).children[0].value;
      assert.match(code, /^package main\n/);
      assert.match(code, /func main\(\)/);
      assert.doesNotMatch(code, /flag\.|customHTTPServer|mainWithHTTPServer|\n\n\n|\n$/);
    }
  }
});

test('missing files report the Markdown filename', (t) => {
  const { root } = fixture(t);
  assert.throws(() => render('```go file=missing.go\n```', root), (error) => {
    assert.equal(error.file, 'recipe.md');
    assert.match(error.message, /Cannot include missing.go/);
    return true;
  });
});

test('rejects conflicting or malformed fences', (t) => {
  const { root } = fixture(t);
  for (const [fence, message] of [
    ['file=', /exactly one file=path/],
    ['file=server.go file=server.go', /exactly one file=path/],
    ['file=server.go#', /file=path#region/],
    ['file=server.go#handler#extra', /file=path#region/],
  ]) {
    assert.throws(() => render(`\x60\x60\x60go ${fence}\n\x60\x60\x60`, root), message);
  }
  assert.throws(() => render('```go file=server.go\npasted code\n```', root), /must be empty/);
});

test('rejects missing, reversed, and duplicate region markers', (t) => {
  const { root } = fixture(t);
  for (const source of [
    '// docs:start handler\nfunc handler() {}\n',
    '// docs:end handler\n// docs:start handler\n',
    '// docs:start handler\n// docs:start handler\n// docs:end handler\n',
  ]) {
    writeFileSync(join(root, 'server.go'), source);
    assert.throws(() => render('```go file=server.go#handler\n```', root), /matching docs:start\/docs:end pair/);
  }
});

test('rejects absolute paths, traversal, and symlinks outside the repository', (t) => {
  const { dir, root } = fixture(t);
  const outside = join(dir, 'outside.go');
  writeFileSync(outside, 'outside the repository');
  symlinkSync(outside, join(root, 'linked.go'));
  for (const path of [outside, '../outside.go', 'linked.go', 'C:\\outside.go']) {
    assert.throws(() => render(`\x60\x60\x60go file=${path}\n\x60\x60\x60`, root), /relative to|escapes|outside/);
  }
});

test('source-only changes update the loader cache key, including nested regions', (t) => {
  const { root } = fixture(t);
  const loader = withSourceCode({ name: 'test', load: (context) => context.generateDigest('> ```go file=server.go#handler\n> ```') }, { root });
  const context = { generateDigest: (value) => JSON.stringify(value) };
  const before = loader.load(context);
  writeFileSync(join(root, 'server.go'), '// docs:start handler\nfunc changed() {}\n// docs:end handler\n');
  const after = loader.load(context);
  assert.notEqual(before, after);
  assert.match(after, /func changed/);
});

test('invalid references fail in the loader before Astro can cache a failed render', (t) => {
  const { root } = fixture(t);
  const loader = withSourceCode({ name: 'test', load: (context) => context.generateDigest('```go file=missing.go\n```') }, { root });
  assert.throws(() => loader.load({ generateDigest: JSON.stringify }), /Cannot include missing.go/);
});

test('MDX fences inside JSX participate in rendering, cache invalidation, and validation', (t) => {
  const { root } = fixture(t);
  const document = (file) => `---\ntitle: Example\n---\n\n<TabItem>\n\x60\x60\x60go file=${file}\n\x60\x60\x60\n</TabItem>`;
  const markdown = document('server.go#handler');
  const code = programFences(render(markdown, root, 'recipe.mdx'));
  assert.equal(code[0].value, 'func handler() {}');
  const context = { generateDigest: JSON.stringify };
  const loader = withSourceCode({ name: 'test', load: (ctx) => ctx.generateDigest(markdown) }, { root });
  const before = loader.load(context);
  writeFileSync(join(root, 'server.go'), '// docs:start handler\nfunc changed() {}\n// docs:end handler\n');
  assert.notEqual(loader.load(context), before);
  const missing = withSourceCode({ name: 'test', load: (ctx) => ctx.generateDigest(document('missing.go')) }, { root });
  assert.throws(() => missing.load(context), /Cannot include missing.go/);
  assert.throws(() => render(document('missing.go'), root, 'recipe.mdx'), /Cannot include missing.go/);
});

test('source ownership assertions find short Go and HTML fences nested in lists and quotes', () => {
  for (const markdown of ['> ```go\n> e.GET("/", handler)\n> ```', '- Example:\n\n  ```html\n  <div>example</div>\n  ```']) {
    assert.throws(() => assertSourceOwned(parseSourceDocument(markdown, 'recipe.md')), /has no pasted program/);
  }
});

test('every cookbook page in every locale displays source-owned programs and excerpts', () => {
  const root = fileURLToPath(new URL('../..', import.meta.url));
  const pages = readdirSync(join(root, 'site/src/content/docs/cookbook')).filter((name) => /\.mdx?$/.test(name)).sort();
  const normalized = (value) => value.split('\n').filter((line) => !/^\s*\/\/\s*docs:(?:start|end)\b/.test(line)).map((line) => line.trim()).join('\n').replace(/\n{3,}/g, '\n\n').trim();
  assert.ok(pages.length);
  for (const locale of localePrefixes) {
    const dir = join(root, 'site/src/content/docs', locale, 'cookbook');
    assert.deepEqual(readdirSync(dir).filter((name) => /\.mdx?$/.test(name)).sort(), pages);
    for (const page of pages) {
      const markdown = readFileSync(join(dir, page), 'utf8');
      const before = parseSourceDocument(markdown, page);
      const fences = assertSourceOwned(before);
      assert.ok(fences.length, `${locale}${page} has a program`);
      const rendered = render(markdown, root, page);
      const code = programFences(rendered);
      for (const [index, fence] of fences.entries()) {
        const specifier = sourceReference(fence).specifier;
        const [path, region] = specifier.split('#');
        const source = normalized(readFileSync(join(root, path), 'utf8'));
        const displayed = normalized(code[index].value);
        assert.ok(displayed.length, `${locale}${page}#${region ?? 'file'} is not empty`);
        assert.ok(source.includes(displayed), `${locale}${page} displays an excerpt of ${path}`);
        if (!region) assert.equal(displayed, source);
        assert.doesNotMatch(code[index].value, /\/\/\s*docs:(?:start|end)\b/);
      }
    }
  }
});

test('raw component imports and Markdown includes share source excerpts; Casbin middleware appears once per locale', () => {
  const root = fileURLToPath(new URL('../..', import.meta.url));
  const path = 'cookbook/hello-world/server.go';
  const code = sourceSnippet(readFileSync(join(root, path), 'utf8'), 'hero', path);
  const [included] = programFences(render(`\x60\x60\x60go file=${path}#hero\n\x60\x60\x60`, root));
  assert.equal(code, included.value);
  assert.match(code, /e\.GET\("\/"/);
  assert.doesNotMatch(code, /docs:(start|end)/);
  for (const locale of localePrefixes) {
    const page = join(root, 'site/src/content/docs', locale, 'middleware/casbin-auth.md');
    const rendered = render(readFileSync(page, 'utf8'), root, page);
    const fences = programFences(rendered);
    const declarations = fences.flatMap((fence) => [...fence.value.matchAll(/^func NewCasbinMiddleware\b/gm)]);
    assert.equal(declarations.length, 1, `${locale}Casbin middleware has a single source-backed definition`);
    assert.ok(fences.some((fence) => fence.value.includes('enforcer.Enforce(')));
    assert.ok(fences.some((fence) => fence.value.includes('echojwt.JWT(')));
    assert.ok(!fences.some((fence) => /^package main\b/m.test(fence.value)), 'the full server is linked rather than duplicated');
  }
});
