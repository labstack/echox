import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const manifest = JSON.parse(readFileSync(join(siteDir, 'src/generated/config-fields.json'), 'utf8'));
const configs = Object.fromEntries(manifest.configs.map((config) => [config.name, Object.fromEntries(config.fields.map((field) => [field.name, { type: field.type, deprecated: field.deprecated ?? false }]))]));
writeFileSync(join(siteDir, 'reference-baseline.json'), `${JSON.stringify({ configs }, null, 2)}\n`);
console.log(`Accepted ${manifest.configs.length} Echo config types for review`);
