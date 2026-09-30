import assert from 'node:assert/strict';
import { existsSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import remarkParse from 'remark-parse';
import { unified } from 'unified';
import { checkCookbookCode } from '../scripts/check-cookbook.mjs';
import { redirects, retiredCookbookRedirects } from '../src/redirects.mjs';

const parser = unified().use(remarkParse);
const root = fileURLToPath(new URL('../..', import.meta.url));

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
  for (const locale of ['', 'es/', 'ja/', 'pt-br/', 'zh-cn/']) {
    for (const [from, to] of Object.entries(retired)) {
      const source = `/${locale}cookbook/${from}`;
      assert.equal(redirects[source], `/${locale}cookbook/${to}/`);
      assert.equal(redirects[`/${locale}docs/cookbook/${from}`], `/${locale}cookbook/${to}/`);
      assert.equal(next[source], `/next/${locale}cookbook/${to}/`);
      assert.ok(existsSync(join(root, 'site/src/content/docs', locale, 'cookbook', `${to}.md`)));
      assert.ok(!existsSync(join(root, 'site/src/content/docs', locale, 'cookbook', `${from}.md`)));
    }
  }
  for (const recipe of Object.keys(retired)) assert.ok(!existsSync(join(root, 'cookbook', recipe)));
});
