import { execFileSync } from 'node:child_process';
import { cpSync, existsSync, readFileSync, readdirSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const proposedEcho = process.env.ECHO_SOURCE_DIR || '';

function run(args, channel, sourceDir = '') {
  execFileSync('npm', args, {
    cwd: siteDir,
    env: { ...process.env, DOCS_CHANNEL: channel, ECHO_SOURCE_DIR: sourceDir },
    stdio: 'inherit',
  });
}

run(['run', 'cookbook:check'], 'stable');
run(['run', 'source:check'], 'stable');
run(['run', 'translations:status'], 'stable');
run(['run', 'security-translations:status'], 'stable');
run(['run', 'astro', '--', 'build'], 'stable');
run(['run', 'source:check'], 'next', proposedEcho);
run(['run', 'astro', '--', 'build'], 'next', proposedEcho);

const nextDist = join(siteDir, 'dist-next');
if (!existsSync(join(nextDist, 'index.html'))) throw new Error(`Next build missing: ${nextDist}`);
function rewriteNextLinks(dir) {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const path = join(dir, entry.name);
    if (entry.isDirectory()) { rewriteNextLinks(path); continue; }
    if (!path.endsWith('.html')) continue;
    const html = readFileSync(path, 'utf8');
    const updated = html.replace(/<a\b[^>]*>/g, (tag) => {
      if (tag.includes('data-version-switch')) return tag;
      return tag.replace(/\bhref="\/((?:(?:es|ja|pt-br|zh-cn)\/)?(?:guide|middleware|cookbook)\/[^\"]*)"/g, 'href="/next/$1"');
    });
    if (updated !== html) writeFileSync(path, updated);
  }
}
rewriteNextLinks(nextDist);
cpSync(nextDist, join(siteDir, 'dist', 'next'), { recursive: true });
run(['run', 'site:check'], 'stable');
run(['run', 'performance:check'], 'stable');
