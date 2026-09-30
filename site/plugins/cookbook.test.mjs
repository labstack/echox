import assert from 'node:assert/strict';
import { existsSync, readdirSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import remarkParse from 'remark-parse';
import { unified } from 'unified';
import { checkCookbookCode } from '../scripts/check-cookbook.mjs';
import { redirects, retiredCookbookRedirects, legacyCookbookLocales } from '../src/redirects.mjs';
import { localePrefixes } from '../src/locales.mjs';
import { parseSourceDocument } from './source-document.mjs';
import { sourceReferenceAt } from './remark-source-code.mjs';

const parser = unified().use(remarkParse);
const root = fileURLToPath(new URL('../..', import.meta.url));

test('malformed source fences report the page and fence line in both checks', () => {
  for (const fence of ['```go file=server.go\npasted\n```', '```go file=one.go file=two.go\n```']) {
    const tree = parser.parse(`Introduction\n\n${fence}`);
    assert.throws(() => checkCookbookCode(tree, 'cookbook/recipe.md'), /cookbook\/recipe.md:3: /);
    const code = tree.children.find((node) => node.type === 'code');
    assert.throws(() => sourceReferenceAt(code, 'middleware/logger.mdx'), /middleware\/logger.mdx:3: /);
  }
});

test('cookbook check rejects pasted Go programs and long excerpts', () => {
  const program = '```go\npackage main\nfunc main() {}\n```';
  const excerpt = '```go\n' + Array(16).fill('e.GET("/", handler)').join('\n') + '\n```';
  for (const markdown of [program, excerpt]) {
    assert.throws(() => checkCookbookCode(parser.parse(markdown), 'recipe.md'), /recipe.md:1: use an empty file=/);
  }
});

test('cookbook check rejects pasted complete HTML pages inside a quote', () => {
  const markdown = '> ```html\n> <!doctype html>\n> <html></html>\n> ```';
  assert.throws(() => checkCookbookCode(parser.parse(markdown), 'recipe.md'), /pasted program/);
});

test('cookbook check finds pasted programs and invalid includes inside MDX JSX', () => {
  for (const code of ['```go\npackage main\n```', '```html\n<!doctype html>\n```', '```go file=server.go\npasted\n```']) {
    const tree = parseSourceDocument(`<TabItem>\n${code}\n</TabItem>`, 'recipe.mdx');
    assert.throws(() => checkCookbookCode(tree, 'recipe.mdx'), /pasted program|must be empty/);
  }
  assert.throws(() => checkCookbookCode(parser.parse('```go title="cp file=x.go"\npackage main\n```'), 'recipe.md'), /pasted program/);
});

test('cookbook check allows short explanations, sample output, and source references', () => {
  const markdown = [
    '```go\ne.GET("/", handler)\n```',
    '```json\n{"message": "OK"}\n```',
    '```go file=cookbook/hello-world/server.go\n```',
    '```html file=cookbook/websocket/public/index.html\n```',
  ].join('\n\n');
  assert.doesNotThrow(() => checkCookbookCode(parser.parse(markdown), 'recipe.md'));
});

test('retired cookbook URLs retain their locale and channel, including legacy URLs', () => {
  const retired = { 'http2-server-push': 'http2', jsonp: 'cors', 'load-balancing': 'reverse-proxy' };
  const next = retiredCookbookRedirects('/next/');
  for (const locale of localePrefixes) {
    for (const [from, to] of Object.entries(retired)) {
      const source = `/${locale}cookbook/${from}`;
      assert.equal(redirects[source], `/${locale}cookbook/${to}/`);
      assert.equal(next[source], `/next/${locale}cookbook/${to}/`);
      assert.ok(existsSync(join(root, 'site/src/content/docs', locale, 'cookbook', `${to}.md`)));
      assert.ok(!existsSync(join(root, 'site/src/content/docs', locale, 'cookbook', `${from}.md`)));
    }
  }
  for (const recipe of Object.keys(retired)) assert.ok(!existsSync(join(root, 'cookbook', recipe)));
  for (const [legacy, current] of Object.entries(legacyCookbookLocales)) {
    const pages = readdirSync(join(root, 'site/src/content/docs/cookbook')).filter((name) => /\.mdx?$/.test(name)).map((name) => name.replace(/\.mdx?$/, ''));
    for (const recipe of [...pages, ...Object.keys(retired)]) {
      assert.equal(redirects[`/${legacy}docs/cookbook/${recipe}`], `/${current}cookbook/${retired[recipe] ?? recipe}/`);
    }
  }
  assert.equal(redirects['/zh-cn/docs/cookbook/jsonp'], undefined);
  assert.equal(redirects['/pt-br/docs/cookbook/jsonp'], undefined);
});
