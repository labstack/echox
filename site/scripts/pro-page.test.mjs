import assert from 'node:assert/strict';
import { test } from 'node:test';
import { isProPage } from '../src/pro-page.mjs';

test('isProPage matches /pro on both channels, with or without a trailing slash', () => {
  for (const path of ['/pro', '/pro/', '/next/pro', '/next/pro/']) assert.equal(isProPage(path), true, path);
  for (const path of ['/', '/next/', '/product/', '/pro/x/', '/ja/pro/', '/nextpro/', '/guide/pro/']) assert.equal(isProPage(path), false, path);
});
