import { existsSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { locales, pages, sections, siteDir } from './translation-sections.mjs';

const [locale, page] = process.argv.slice(2);
if (!locales.includes(locale) || !pages.includes(page)) {
  console.error('Usage: npm run translations:accept -- <es|ja|pt-br|zh-cn> <logger|static|index>');
  process.exit(1);
}
const path = join(siteDir, 'translation-baseline.json');
const baseline = existsSync(path) ? JSON.parse(readFileSync(path, 'utf8')) : {};
const english = sections('', page);
const translated = sections(locale, page);
if (english.length !== translated.length) {
  throw new Error(`${locale}/middleware/${page}: ${english.length} English sections and ${translated.length} translated sections`);
}
baseline[locale] ??= {};
baseline[locale][page] = english.map((section, index) => ({
  heading: section.heading,
  source: section.hash,
  translation: translated[index].hash,
}));
writeFileSync(path, `${JSON.stringify(baseline, null, 2)}\n`);
console.log(`Recorded reviewed sections for ${locale}/middleware/${page}`);
