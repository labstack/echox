import assert from 'node:assert/strict';
import test from 'node:test';
import { englishOnlyPro } from './english-only-pro.mjs';

test('filtering adapter replaces only the Starlight docs route registration', async () => {
  const routes = [];
  const integration = englishOnlyPro({ hooks: { 'astro:config:setup': ({ injectRoute }) => {
    injectRoute({ pattern: '404', entrypoint: '404.astro' });
    injectRoute({ pattern: '[...slug]', entrypoint: '@astrojs/starlight/routes/static/index.astro' });
  } } });
  await integration.hooks['astro:config:setup']({ injectRoute: (route) => routes.push(route) });
  assert.deepEqual(routes.map((r) => r.entrypoint), ['404.astro', './src/starlight-docs.astro']);
});

test('routing changes require explicit review', async () => {
  const integration = englishOnlyPro({ hooks: { 'astro:config:setup': () => {} } });
  await assert.rejects(() => integration.hooks['astro:config:setup']({}), /route missing/);
});
