import assert from 'node:assert/strict';
import { test } from 'node:test';
import { includeDocsRoute, isProPage } from '../src/pro-page.mjs';

test('Pro prefix matching is bounded and includes package and problem pages', () => {
  for (const path of ['/pro', '/pro/', '/pro/quota/', '/pro/problems/admin-busy', '/next/pro', '/next/pro/quota/']) assert.equal(isProPage(path), true, path);
  for (const path of ['/', '/next/', '/product/', '/ja/pro/', '/nextpro/', '/guide/pro/', '/problems/', '/professional/']) assert.equal(isProPage(path), false, path);
});

test('Pro routes have no locale fallback or next-channel copy', () => {
  for (const slug of ['pro/install', 'pro/quota', 'pro/problems/admin-busy']) {
    assert.equal(includeDocsRoute(slug), true);
    assert.equal(includeDocsRoute(slug, true), false);
    for (const locale of ['es', 'ja', 'pt-br', 'zh-cn']) assert.equal(includeDocsRoute(`${locale}/${slug}`), false);
  }
  for (const slug of ['', 'guide/installation', 'ja/middleware/jwt']) assert.equal(includeDocsRoute(slug, true), true);
});
