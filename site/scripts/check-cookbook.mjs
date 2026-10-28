import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { unified } from 'unified';
import remarkSourceCode, { sourceReferenceAt } from '../plugins/remark-source-code.mjs';
import { parseSourceDocument, visitCode } from '../plugins/source-document.mjs';
import { localePrefixes } from '../src/locales.mjs';

const docsDir = fileURLToPath(new URL('../src/content/docs', import.meta.url));
const siteDir = fileURLToPath(new URL('..', import.meta.url));
const processor = unified().use(remarkSourceCode);

export function checkCookbookCode(tree, path) {
  visitCode(tree, (node) => {
    if (!sourceReferenceAt(node, path)) {
      const program = node.lang === 'go' && (node.value.split('\n').length > 15 || /^\s*package\s/.test(node.value));
      const html = node.lang === 'html' && /^\s*(?:<!doctype\b|<html\b)/i.test(node.value);
      if (program || html) {
        throw new Error(`${path}:${node.position.start.line}: use an empty file= fence instead of a pasted program`);
      }
    }
  });
}

export function checkCookbook() {
  let count = 0;
  for (const locale of localePrefixes) {
    const dir = join(docsDir, locale, 'cookbook');
    for (const page of readdirSync(dir).filter((name) => /\.mdx?$/.test(name))) {
      const path = join(dir, page);
      const markdown = readFileSync(path, 'utf8');
      const tree = parseSourceDocument(markdown, path);
      checkCookbookCode(tree, path);
      // Validate every referenced file and region even when Astro has cached a page.
      processor.runSync(tree, { path });
      count++;
    }
  }
  const pro = join(siteDir, 'pro-source/pages');
  function checkPro(dir) {
    for (const entry of readdirSync(dir, { withFileTypes: true })) {
      const path = join(dir, entry.name);
      if (entry.isDirectory()) { checkPro(path); continue; }
      if (!/\.mdx?$/.test(path)) continue;
      const tree = parseSourceDocument(readFileSync(path, 'utf8'), path);
      checkCookbookCode(tree, path);
      processor.runSync(tree, { path });
      count++;
    }
  }
  checkPro(pro);
  console.log(`Checked ${count} cookbook pages: source includes resolve and no programs are pasted`);
}

if (process.argv[1] === fileURLToPath(import.meta.url)) checkCookbook();
