import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { locales, pages, sections, siteDir } from './translation-sections.mjs';

const baseline = JSON.parse(readFileSync(join(siteDir, 'translation-baseline.json'), 'utf8'));
let stale = 0;
for (const locale of locales) {
  for (const page of pages) {
    const english = sections('', page);
    const translated = sections(locale, page);
    const accepted = baseline[locale]?.[page];
    if (!Array.isArray(accepted)) {
      console.log(`${locale}/middleware/${page}: needs review (no section baseline)`);
      stale++;
      continue;
    }
    const count = Math.max(english.length, translated.length, accepted.length);
    for (let index = 0; index < count; index++) {
      const source = english[index];
      const translation = translated[index];
      const recorded = accepted[index];
      const changes = [];
      if (!source || !translation || !recorded) changes.push('section count changed');
      if (source && recorded && source.hash !== recorded.source) changes.push('English changed');
      if (translation && recorded && translation.hash !== recorded.translation) changes.push('translation changed');
      if (changes.length) {
        console.log(`${locale}/middleware/${page} §${index + 1} ${source?.heading ?? recorded?.heading ?? 'missing'}: ${changes.join(', ')}`);
        stale++;
      }
    }
  }
}
console.log(`${stale} section(s) need review; matching hashes track edits, not translation quality`);
if (stale) process.exitCode = 1;
