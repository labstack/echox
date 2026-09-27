import { createHash } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

export const siteDir = fileURLToPath(new URL('..', import.meta.url));
export const locales = ['es', 'ja', 'pt-br', 'zh-cn'];
export const pages = ['logger', 'static', 'index'];

export function sections(locale, page) {
  const extension = page === 'index' ? 'md' : 'mdx';
  const path = join(siteDir, 'src/content/docs', locale ? `${locale}/middleware` : 'middleware', `${page}.${extension}`);
  const content = readFileSync(path, 'utf8');
  const blocks = [];
  let heading = 'Introduction';
  let lines = [];
  for (const line of content.split('\n')) {
    if (/^#{2,3} /.test(line)) {
      blocks.push({ heading, text: lines.join('\n').trim() });
      heading = line.replace(/^#{2,3} /, '');
      lines = [];
    } else {
      lines.push(line);
    }
  }
  blocks.push({ heading, text: lines.join('\n').trim() });
  return blocks.map(({ heading, text }) => ({
    heading,
    hash: createHash('sha256').update(`${heading}\n${text}`).digest('hex'),
  }));
}
