import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { preparedSource } from './module-sources.mjs';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const channel = process.env.DOCS_CHANNEL === 'next' ? 'next' : 'stable';
const baselineFile = channel === 'next' ? 'next-reference-baseline.json' : 'reference-baseline.json';
const { manifest, external: sources } = preparedSource(siteDir);
const configs = Object.fromEntries(manifest.configs.map((config) => [config.name, Object.fromEntries(config.fields.map((field) => [field.name, { type: field.type, deprecated: field.deprecated ?? false }]))]));
const functions = Object.fromEntries(manifest.functions.map((entry) => [entry.name, entry.signature]));
const external = Object.fromEntries(Object.keys(sources).map((slug) => {
  const data = JSON.parse(readFileSync(join(siteDir, `src/generated/${slug}.json`), 'utf8'));
  if (data.module !== sources[slug].module || data.revision !== sources[slug].version) {
    throw new Error(`${slug}: generated source is stale. Run source:prepare before accepting the baseline.`);
  }
  return [slug, {
    module: data.module,
    version: data.revision,
    configs: Object.fromEntries(data.configs.map((config) => [config.name, Object.fromEntries(config.fields.map((field) => [field.name, { type: field.type, deprecated: field.deprecated ?? false }]))])),
    functions: Object.fromEntries(data.functions.map((entry) => [entry.name, entry.signature])),
  }];
}));
writeFileSync(join(siteDir, baselineFile), `${JSON.stringify({ configs, functions }, null, 2)}\n`);
writeFileSync(join(siteDir, 'external-baseline.json'), `${JSON.stringify(external, null, 2)}\n`);
console.log(`Accepted ${channel} ${manifest.configs.length} Echo config types, ${manifest.functions.length} functions, and ${Object.keys(external).length} external modules for review`);
