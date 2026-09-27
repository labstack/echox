import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const manifest = JSON.parse(readFileSync(join(siteDir, 'src/generated/config-fields.json'), 'utf8'));
const configs = Object.fromEntries(manifest.configs.map((config) => [config.name, Object.fromEntries(config.fields.map((field) => [field.name, { type: field.type, deprecated: field.deprecated ?? false }]))]));
const functions = Object.fromEntries(manifest.functions.map((entry) => [entry.name, entry.signature]));
writeFileSync(join(siteDir, 'reference-baseline.json'), `${JSON.stringify({ configs, functions }, null, 2)}\n`);
const sources = JSON.parse(readFileSync(join(siteDir, 'external-sources.json'), 'utf8'));
const external = Object.fromEntries(Object.keys(sources).map((slug) => {
  const data = JSON.parse(readFileSync(join(siteDir, `src/generated/${slug}.json`), 'utf8'));
  return [slug, {
    module: data.module,
    version: data.revision,
    configs: Object.fromEntries(data.configs.map((config) => [config.name, Object.fromEntries(config.fields.map((field) => [field.name, { type: field.type, deprecated: field.deprecated ?? false }]))])),
    functions: Object.fromEntries(data.functions.map((entry) => [entry.name, entry.signature])),
  }];
}));
writeFileSync(join(siteDir, 'external-baseline.json'), `${JSON.stringify(external, null, 2)}\n`);
console.log(`Accepted ${manifest.configs.length} Echo config types, ${manifest.functions.length} functions, and ${Object.keys(external).length} external modules for review`);
