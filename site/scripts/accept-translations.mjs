import { createHash } from 'node:crypto';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const hash = (page) => createHash('sha256').update(readFileSync(join(siteDir, 'src/content/docs/middleware', `${page}.mdx`))).digest('hex');
const pages = { logger: hash('logger'), static: hash('static') };
const baseline = Object.fromEntries(['es', 'ja', 'pt-br', 'zh-cn'].map((locale) => [locale, pages]));
writeFileSync(join(siteDir, 'translation-baseline.json'), `${JSON.stringify(baseline, null, 2)}\n`);
console.log('Accepted Request Logger and Static translations for review');
