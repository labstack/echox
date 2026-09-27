import { createHash } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const baseline = JSON.parse(readFileSync(join(siteDir, 'translation-baseline.json'), 'utf8'));
const hash = (page) => createHash('sha256').update(readFileSync(join(siteDir, 'src/content/docs/middleware', `${page}.mdx`))).digest('hex');
let stale = 0;
for (const [locale, pages] of Object.entries(baseline)) {
  for (const [page, reviewedHash] of Object.entries(pages)) {
    const status = reviewedHash === hash(page) ? 'current' : 'needs review';
    if (status !== 'current') stale++;
    console.log(`${locale}/middleware/${page}: ${status}`);
  }
}
console.log(`${stale} translation(s) need review; this report does not validate the prose`);
