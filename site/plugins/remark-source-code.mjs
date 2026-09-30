import { readFileSync, realpathSync } from 'node:fs';
import { isAbsolute, relative, resolve, sep } from 'node:path';
import { fileURLToPath } from 'node:url';
import remarkParse from 'remark-parse';
import { unified } from 'unified';

const repoRoot = fileURLToPath(new URL('../..', import.meta.url));
const parser = unified().use(remarkParse);
const regionMarker = /^\s*\/\/\s*docs:(start|end)\s+(\S+)\s*$/;

function visitCode(tree, callback) {
  if (tree.type === 'code') callback(tree);
  for (const child of tree.children ?? []) visitCode(child, callback);
}

function reference(node) {
  const matches = [...(node.meta ?? '').matchAll(/(?:^|\s)file=(\S*)/g)];
  if (!matches.length) return;
  if (matches.length !== 1 || !matches[0][1]) {
    throw new Error('Expected exactly one file=path in the code fence');
  }
  if (node.value.trim()) {
    throw new Error('A file= code fence must be empty; edit the source file instead');
  }
  return matches[0];
}

function readSnippet(specifier, root) {
  const [path, region, extra] = specifier.split('#');
  if (!path || isAbsolute(path) || /^[A-Za-z]:[\\/]/.test(path) || path.startsWith('\\')) {
    throw new Error(`Source path must be relative to the repository root: ${specifier}`);
  }
  if (extra !== undefined || region === '') {
    throw new Error(`Expected file=path or file=path#region: ${specifier}`);
  }
  const insideRoot = (file) => {
    const fromRoot = relative(root, file);
    return fromRoot !== '..' && !fromRoot.startsWith(`..${sep}`) && !isAbsolute(fromRoot);
  };
  const resolved = resolve(root, path);
  if (!insideRoot(resolved)) throw new Error(`Source path escapes the repository: ${path}`);
  let file;
  let source;
  try {
    file = realpathSync(resolved);
    if (!insideRoot(file)) throw new Error('symlink points outside the repository');
    source = readFileSync(file, 'utf8').replace(/\r\n/g, '\n');
  } catch (error) {
    throw new Error(`Cannot include ${path}: ${error.message}`);
  }
  const lines = source.split('\n');
  if (region !== undefined) {
    const starts = [];
    const ends = [];
    for (const [index, line] of lines.entries()) {
      const marker = line.match(regionMarker);
      if (marker?.[2] !== region) continue;
      (marker[1] === 'start' ? starts : ends).push(index);
    }
    if (starts.length !== 1 || ends.length !== 1 || ends[0] <= starts[0]) {
      throw new Error(`Expected one matching docs:start/docs:end pair for ${path}#${region}`);
    }
    const excerpt = lines.slice(starts[0] + 1, ends[0]).filter((line) => !regionMarker.test(line));
    const indents = excerpt.filter((line) => line.trim()).map((line) => line.match(/^[\t ]*/)[0]);
    const indent = indents.reduce((prefix, next) => {
      while (!next.startsWith(prefix)) prefix = prefix.slice(0, -1);
      return prefix;
    }, indents[0] ?? '');
    return excerpt.map((line) => line.startsWith(indent) ? line.slice(indent.length) : line).join('\n');
  }
  // A Markdown fence has no trailing newline in its parsed value.
  return lines.filter((line) => !regionMarker.test(line)).join('\n').replace(/\n$/, '');
}

/** Fill empty file= code fences before syntax highlighting, in Markdown and MDX. */
export default function remarkSourceCode({ root = repoRoot } = {}) {
  root = realpathSync(root);
  return (tree, file) => {
    visitCode(tree, (node) => {
      try {
        const match = reference(node);
        if (!match) return;
        node.value = readSnippet(match[1], root);
        node.meta = node.meta.replace(match[0], '').trim() || null;
      } catch (error) {
        file.fail(error.message, node);
      }
    });
  };
}

/** Include source contents in Astro's cache key so source-only edits rebuild pages. */
export function withSourceCode(loader, { root = repoRoot } = {}) {
  root = realpathSync(root);
  return {
    ...loader,
    name: `${loader.name}-source-code`,
    load(context) {
      return loader.load({
        ...context,
        generateDigest(contents) {
          if (typeof contents !== 'string' || !contents.includes('file=')) {
            return context.generateDigest(contents);
          }
          const snippets = [];
          visitCode(parser.parse(contents), (node) => {
            const match = reference(node);
            if (match) snippets.push(readSnippet(match[1], root));
          });
          return context.generateDigest({ contents, snippets });
        },
      });
    },
  };
}
