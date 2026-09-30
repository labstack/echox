import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, readFileSync, readdirSync, rmSync, symlinkSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';
import remarkParse from 'remark-parse';
import { unified } from 'unified';
import remarkSourceCode, { withSourceCode } from './remark-source-code.mjs';

function fixture(t) {
  const dir = mkdtempSync(join(tmpdir(), 'source-code-'));
  const root = join(dir, 'repo');
  mkdirSync(root);
  writeFileSync(join(root, 'server.go'), 'package main\n\n// docs:start handler\nfunc handler() {}\n// docs:end handler\n');
  t.after(() => rmSync(dir, { recursive: true, force: true }));
  return { dir, root };
}

function render(markdown, root) {
  const processor = unified().use(remarkParse).use(remarkSourceCode, { root });
  return processor.runSync(processor.parse(markdown), { path: 'recipe.md' });
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

test('every cookbook page in every locale displays source-owned programs and excerpts', () => {
  const root = fileURLToPath(new URL('../..', import.meta.url));
  const pages = readdirSync(join(root, 'site/src/content/docs/cookbook')).filter((name) => /\.mdx?$/.test(name)).sort();
  const normalized = (value) => value.split('\n').filter((line) => !/^\s*\/\/\s*docs:(?:start|end)\b/.test(line)).map((line) => line.trim()).join('\n').trim();
  assert.ok(pages.length);
  for (const locale of ['', 'es/', 'ja/', 'pt-br/', 'zh-cn/']) {
    const dir = join(root, 'site/src/content/docs', locale, 'cookbook');
    assert.deepEqual(readdirSync(dir).filter((name) => /\.mdx?$/.test(name)).sort(), pages);
    for (const page of pages) {
      const markdown = readFileSync(join(dir, page), 'utf8');
      const before = unified().use(remarkParse).parse(markdown);
      const fences = before.children.filter((node) => node.type === 'code' && ['go', 'html'].includes(node.lang));
      assert.ok(fences.length, `${locale}${page} has a program`);
      for (const fence of fences) {
        assert.match(fence.meta ?? '', /(?:^|\s)file=/, `${locale}${page} has no pasted program`);
        assert.equal(fence.value, '');
      }
      const rendered = render(markdown, root);
      const code = rendered.children.filter((node) => node.type === 'code' && ['go', 'html'].includes(node.lang));
      for (const [index, fence] of fences.entries()) {
        const specifier = fence.meta.match(/(?:^|\s)file=(\S+)/)[1];
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
