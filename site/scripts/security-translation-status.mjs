import { createHash } from 'node:crypto';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const docsDir = join(siteDir, 'src/content/docs');
const baselinePath = join(siteDir, 'security-translation-baseline.json');
const locales = ['es', 'ja', 'pt-br', 'zh-cn'];
const pages = [
  'guide/ip-address.md',
  'guide/request-scheme.md',
  'guide/response.md',
  'guide/static-files.md',
  'guide/testing.md',
  'middleware/redirect.mdx',
  'middleware/secure.mdx',
  'middleware/proxy.mdx',
  'middleware/method-override.mdx',
  'middleware/static.mdx',
  'middleware/trailing-slash.mdx',
  'cookbook/reverse-proxy.md',
];

function hash(path) {
  return createHash('sha256').update(readFileSync(path, 'utf8')).digest('hex');
}

const actual = Object.fromEntries(locales.map((locale) => [
  locale,
  Object.fromEntries(pages.map((page) => [page, {
    source: hash(join(docsDir, page)),
    translation: hash(join(docsDir, locale, page)),
  }])),
]));

if (process.argv.includes('--accept')) {
  writeFileSync(baselinePath, `${JSON.stringify(actual, null, 2)}\n`);
  console.log(`Recorded ${pages.length} reviewed security pages in ${locales.length} locales`);
} else {
  const baseline = JSON.parse(readFileSync(baselinePath, 'utf8'));
  const stale = [];
  for (const locale of locales) {
    for (const page of pages) {
      const before = baseline[locale]?.[page];
      const now = actual[locale][page];
      if (!before || before.source !== now.source || before.translation !== now.translation) {
        stale.push(`${locale}/${page}: ${!before ? 'not reviewed' : before.source !== now.source ? 'English changed' : 'translation changed'}`);
      }
    }
  }
  if (stale.length) throw new Error(`${stale.length} security translations need review:\n${stale.join('\n')}`);
  console.log(`${pages.length} security pages match reviewed translations in ${locales.length} locales`);
}
