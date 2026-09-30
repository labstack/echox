import { readFileSync, realpathSync } from 'node:fs';
import { isAbsolute, relative, resolve, sep } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parseSourceDocument, visitCode } from './source-document.mjs';

const repoRoot = fileURLToPath(new URL('../..', import.meta.url));
const regionMarker = /^\s*\/\/\s*docs:(start|end)\s+(\S+)\s*$/;

export function sourceReference(node) {
  // Tokenize whole attributes so file= inside a quoted title is ordinary text.
  const matches = [...(node.meta ?? '').matchAll(/(?:[^\s"']|"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*')+/g)]
    .filter((match) => match[0].startsWith('file='));
  if (!matches.length) return;
  const specifier = matches[0][0].slice(5);
  if (matches.length !== 1 || !specifier || /["']/.test(specifier)) {
    throw new Error('Expected exactly one file=path in the code fence');
  }
  if (node.value.trim()) {
    throw new Error('A file= code fence must be empty; edit the source file instead');
  }
  return { specifier, start: matches[0].index, end: matches[0].index + matches[0][0].length };
}

function withoutMarkers(lines) {
  const result = [];
  let removed = false;
  for (const line of lines) {
    if (regionMarker.test(line)) { removed = true; continue; }
    if (!line.trim() && removed && !result.at(-1)?.trim()) continue;
    result.push(line);
    if (line.trim()) removed = false;
  }
  while (result.length && !result[0].trim()) result.shift();
  while (result.length && !result.at(-1).trim()) result.pop();
  return result;
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
    const excerpt = withoutMarkers(lines.slice(starts[0] + 1, ends[0]));
    const indents = excerpt.filter((line) => line.trim()).map((line) => line.match(/^[\t ]*/)[0]);
    const indent = indents.reduce((prefix, next) => {
      while (!next.startsWith(prefix)) prefix = prefix.slice(0, -1);
      return prefix;
    }, indents[0] ?? '');
    return excerpt.map((line) => line.startsWith(indent) ? line.slice(indent.length) : line).join('\n');
  }
  // A Markdown fence has no trailing newline in its parsed value.
  return withoutMarkers(lines).join('\n');
}

/** Fill empty file= code fences before syntax highlighting, in Markdown and MDX. */
export default function remarkSourceCode({ root = repoRoot } = {}) {
  root = realpathSync(root);
  return (tree, file) => {
    visitCode(tree, (node) => {
      try {
        const match = sourceReference(node);
        if (!match) return;
        node.value = readSnippet(match.specifier, root);
        node.meta = [node.meta.slice(0, match.start).trim(), node.meta.slice(match.end).trim()].filter(Boolean).join(' ') || null;
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
          visitCode(parseSourceDocument(contents), (node) => {
            const match = sourceReference(node);
            if (match) snippets.push(readSnippet(match.specifier, root));
          });
          return context.generateDigest({ contents, snippets });
        },
      });
    },
  };
}
