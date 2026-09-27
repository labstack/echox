import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const manifest = JSON.parse(readFileSync(join(siteDir, 'src/generated/config-fields.json'), 'utf8'));
const baseline = JSON.parse(readFileSync(join(siteDir, 'reference-baseline.json'), 'utf8'));

function fieldsByConfig(configs) {
  return Object.fromEntries(configs.map((config) => [config.name, Object.fromEntries(config.fields.map((field) => [field.name, { type: field.type, deprecated: field.deprecated ?? false }]))]));
}

const actual = fieldsByConfig(manifest.configs);
const differences = [];
for (const name of new Set([...Object.keys(baseline.configs), ...Object.keys(actual)])) {
  const before = baseline.configs[name] || {};
  const after = actual[name] || {};
  for (const field of new Set([...Object.keys(before), ...Object.keys(after)])) {
    if (JSON.stringify(before[field]) !== JSON.stringify(after[field])) {
      differences.push(`${name}.${field}: ${JSON.stringify(before[field] ?? null)} -> ${JSON.stringify(after[field] ?? null)}`);
    }
  }
}
if (differences.length) {
  throw new Error(`Echo config fields changed. Review the affected pages, then update reference-baseline.json:\n${differences.join('\n')}`);
}

for (const locale of ['', 'es/', 'ja/', 'pt-br/', 'zh-cn/']) {
  for (const [page, config, example] of [['logger', 'RequestLoggerConfig', 'request-logger.go'], ['static', 'StaticConfig', 'static.go']]) {
    const file = join(siteDir, 'src/content/docs', locale, 'middleware', `${page}.mdx`);
    const content = readFileSync(file, 'utf8');
    if (!content.includes(`name="${config}"`) || !content.includes(example) || (page === 'logger' && content.includes('LogError'))) {
      throw new Error(`${file} must display the source-backed ${config} reference and ${example} example without removed fields`);
    }
  }
}
console.log(`Source reference matches Echo ${manifest.revision}`);
