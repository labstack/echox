import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
import remarkParse from 'remark-parse';
import { unified } from 'unified';
import remarkSourceCode from '../plugins/remark-source-code.mjs';

const docsDir = fileURLToPath(new URL('../src/content/docs', import.meta.url));
const parser = unified().use(remarkParse);
const processor = unified().use(remarkParse).use(remarkSourceCode);

export function checkCookbookCode(tree, path) {
  function visit(node) {
    if (node.type === 'code' && !/(?:^|\s)file=/.test(node.meta ?? '')) {
      const program = node.lang === 'go' && (node.value.split('\n').length > 15 || /^\s*package\s/.test(node.value));
      const html = node.lang === 'html' && /^\s*(?:<!doctype\b|<html\b)/i.test(node.value);
      if (program || html) {
        throw new Error(`${path}:${node.position.start.line}: use an empty file= fence instead of a pasted program`);
      }
    }
    for (const child of node.children ?? []) visit(child);
  }
  visit(tree);
}

export function checkCookbook() {
  let count = 0;
  for (const locale of ['', 'es/', 'ja/', 'pt-br/', 'zh-cn/']) {
    const dir = join(docsDir, locale, 'cookbook');
    for (const page of readdirSync(dir).filter((name) => /\.mdx?$/.test(name))) {
      const path = join(dir, page);
      const markdown = readFileSync(path, 'utf8');
      const tree = parser.parse(markdown);
      checkCookbookCode(tree, path);
      // Validate every referenced file and region even when Astro has cached a page.
      processor.runSync(tree, { path });
      count++;
    }
  }
  console.log(`Checked ${count} cookbook pages: source includes resolve and no programs are pasted`);
}

if (process.argv[1] === fileURLToPath(import.meta.url)) checkCookbook();
