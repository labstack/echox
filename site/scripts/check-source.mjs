import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parseSourceDocument, visitCode } from '../plugins/source-document.mjs';
import { sourceReferenceAt } from '../plugins/remark-source-code.mjs';
import { localePrefixes } from '../src/locales.mjs';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const channel = process.env.DOCS_CHANNEL === 'next' ? 'next' : 'stable';
const baselineFile = channel === 'next' ? 'next-reference-baseline.json' : 'reference-baseline.json';
const manifest = JSON.parse(readFileSync(join(siteDir, 'src/generated/config-fields.json'), 'utf8'));
const pin = JSON.parse(readFileSync(join(siteDir, channel === 'next' ? 'next-source.json' : 'echo-source.json'), 'utf8'));
if (!process.env.ECHO_SOURCE_DIR && manifest.revision !== pin.revision) {
  throw new Error(`Generated Echo source is ${manifest.revision}; expected ${channel} revision ${pin.revision}. Run source:prepare for this channel first.`);
}
const baseline = JSON.parse(readFileSync(join(siteDir, baselineFile), 'utf8'));
const pages = JSON.parse(readFileSync(join(siteDir, 'reference-pages.json'), 'utf8'));
const externalSources = JSON.parse(readFileSync(join(siteDir, 'external-sources.json'), 'utf8'));
const externalBaseline = JSON.parse(readFileSync(join(siteDir, 'external-baseline.json'), 'utf8'));

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
const functions = Object.fromEntries(manifest.functions.map((entry) => [entry.name, entry.signature]));
for (const name of new Set([...Object.keys(baseline.functions), ...Object.keys(functions)])) {
  if (baseline.functions[name] !== functions[name]) {
    differences.push(`${name}: ${JSON.stringify(baseline.functions[name] ?? null)} -> ${JSON.stringify(functions[name] ?? null)}`);
  }
}
for (const [slug, source] of Object.entries(externalSources)) {
  const external = JSON.parse(readFileSync(join(siteDir, `src/generated/${slug}.json`), 'utf8'));
  const before = externalBaseline[slug];
  if (!before || external.module !== source.module || external.revision !== source.version || before.version !== source.version) {
    differences.push(`${slug}: source module/version differs from reviewed baseline`);
    continue;
  }
  const fields = fieldsByConfig(external.configs);
  for (const config of new Set([...Object.keys(before.configs), ...Object.keys(fields)])) {
    const oldFields = before.configs[config] || {};
    const newFields = fields[config] || {};
    for (const field of new Set([...Object.keys(oldFields), ...Object.keys(newFields)])) {
      if (JSON.stringify(oldFields[field]) !== JSON.stringify(newFields[field])) differences.push(`${slug}.${config}.${field}: ${JSON.stringify(oldFields[field] ?? null)} -> ${JSON.stringify(newFields[field] ?? null)}`);
    }
  }
  const api = Object.fromEntries(external.functions.map((entry) => [entry.name, entry.signature]));
  for (const name of new Set([...Object.keys(before.functions), ...Object.keys(api)])) {
    if (before.functions[name] !== api[name]) differences.push(`${slug}.${name}: ${JSON.stringify(before.functions[name] ?? null)} -> ${JSON.stringify(api[name] ?? null)}`);
  }
}
if (differences.length) {
  throw new Error(`${channel} middleware API changed. Review the affected pages, then update ${baselineFile}:\n${differences.join('\n')}`);
}

const documented = Object.values(pages).flat();
for (const name of manifest.configs.map((entry) => entry.name)) {
  if (documented.filter((entry) => entry === name).length !== 1) {
    throw new Error(`${name} must be mapped to exactly one middleware page`);
  }
}
for (const locale of localePrefixes) {
  for (const [page, configs] of Object.entries(pages)) {
    const file = join(siteDir, 'src/content/docs', locale, 'middleware', `${page}.mdx`);
    const content = readFileSync(file, 'utf8');
    for (const config of configs) {
      if (!content.includes(`<ConfigReference name="${config}"`) || content.includes(`type ${config} struct {`)) {
        throw new Error(`${file} must display the source-backed ${config} reference without a copied struct`);
      }
    }
    const example = { logger: 'reference/request-logger/main.go', static: 'reference/static/main.go' }[page];
    if (example) {
      let included = false;
      visitCode(parseSourceDocument(content, file), (node) => {
        const reference = sourceReferenceAt(node, file);
        if (node.lang === 'go' && reference?.specifier === example) included = true;
      });
      if (!included || (page === 'logger' && content.includes('LogError'))) {
        throw new Error(`${file} must include the compiled ${example} example without removed fields`);
      }
    }
  }
  for (const [page, configs] of Object.entries({ jwt: ['Config'], prometheus: ['MiddlewareConfig', 'HandlerConfig', 'PushGatewayConfig'], 'open-telemetry': ['Config'] })) {
    const file = join(siteDir, 'src/content/docs', locale, 'middleware', `${page}.mdx`);
    const content = readFileSync(file, 'utf8');
    for (const config of configs) {
      if (!content.includes(`<ConfigReference name="${config}" source="${page}"`) || content.includes(`type ${config} struct {`)) {
        throw new Error(`${file} must show ${page}.${config} from the pinned external module`);
      }
    }
  }
}
console.log(`Source reference matches Echo ${manifest.revision}`);
