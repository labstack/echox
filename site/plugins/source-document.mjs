import remarkParse from 'remark-parse';
import remarkMdx from 'remark-mdx';
import remarkFrontmatter from 'remark-frontmatter';
import { unified } from 'unified';

const markdown = unified().use(remarkParse).use(remarkFrontmatter);
const mdx = unified().use(remarkParse).use(remarkFrontmatter).use(remarkMdx);

export function visitCode(tree, callback) {
  if (tree.type === 'code') callback(tree);
  for (const child of tree.children ?? []) visitCode(child, callback);
}

export function parseSourceDocument(contents, path) {
  if (path?.endsWith('.mdx')) return mdx.parse(contents);
  if (path) return markdown.parse(contents);
  // Astro's digest callback receives contents without the file extension. A
  // valid MDX parse finds fences inside JSX; ordinary HTML Markdown may need
  // the CommonMark parser instead. Invalid MDX still fails during rendering.
  try { return mdx.parse(contents); } catch { return markdown.parse(contents); }
}
