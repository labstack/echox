import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { createWorkspace, downloadModule, echoSource, moduleSources, run } from './module-sources.mjs';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const repoDir = dirname(siteDir);
const channel = process.env.DOCS_CHANNEL === 'next' ? 'next' : 'stable';
const sources = moduleSources(siteDir);
const override = process.env.ECHO_SOURCE_DIR;
if (override && channel !== 'next') throw new Error('Use DOCS_CHANNEL=next when testing ECHO_SOURCE_DIR');
const download = channel === 'stable' ? downloadModule(sources.echo, repoDir) : null;
const pin = download ? echoSource(download) : { ...JSON.parse(readFileSync(join(siteDir, 'next-source.json'), 'utf8')), channel };
const sourceDir = download ? download.Dir : override ? resolve(siteDir, override) : join(repoDir, '.cache', 'echo-next');
const referenceDir = join(repoDir, 'reference');
const generatedDir = join(siteDir, 'src', 'generated');

if (channel === 'next' && !override) {
  mkdirSync(dirname(sourceDir), { recursive: true });
  if (!existsSync(join(sourceDir, '.git'))) {
    run('git', ['clone', '--filter=blob:none', '--no-checkout', pin.repository, sourceDir], repoDir);
  }
  const current = run('git', ['rev-parse', 'HEAD'], sourceDir);
  if (current !== pin.revision || !existsSync(join(sourceDir, 'go.mod'))) {
    if (current !== pin.revision) run('git', ['fetch', '--depth=1', 'origin', pin.revision], sourceDir);
    run('git', ['checkout', '--detach', pin.revision], sourceDir);
  }
}

const revision = download ? pin.revision : run('git', ['rev-parse', 'HEAD'], sourceDir);
if (!override && revision !== pin.revision) {
  throw new Error(`Echo checkout is ${revision}; expected ${pin.revision}`);
}
if (!download && run('git', ['status', '--porcelain', '--untracked-files=no'], sourceDir)) {
  throw new Error(`Echo checkout has tracked changes: ${sourceDir}`);
}

const workspaceDir = join(repoDir, '.cache', 'echo-work');
const workspaceFile = createWorkspace(workspaceDir, sourceDir, referenceDir);
const goEnvironment = { GOWORK: workspaceFile };
const manifestText = run('go', ['run', './cmd/config-fields', '-root', sourceDir, '-revision', revision], referenceDir, goEnvironment);
const manifest = JSON.parse(manifestText);
if (manifest.revision !== revision || !manifest.configs.some((config) => config.name === 'RequestLoggerConfig') || !manifest.configs.some((config) => config.name === 'StaticConfig')) {
  throw new Error('Echo source extractor returned an incomplete or mismatched manifest');
}
run('go', ['vet', './...'], referenceDir, goEnvironment);
run('go', ['test', '-race', './...'], referenceDir, goEnvironment);
mkdirSync(generatedDir, { recursive: true });
writeFileSync(join(generatedDir, 'config-fields.json'), `${JSON.stringify(manifest, null, 2)}\n`);
writeFileSync(join(generatedDir, 'echo-source.json'), `${JSON.stringify({ ...pin, revision }, null, 2)}\n`);
for (const [slug, source] of Object.entries(sources.external)) {
  const download = downloadModule(source, referenceDir);
  const extracted = JSON.parse(run('go', [
    'run', './cmd/config-fields', '-root', download.Dir, '-revision', source.version,
    '-module', source.module, '-directory', '.', '-package', source.package,
  ], referenceDir, goEnvironment));
  if (extracted.module !== source.module || extracted.revision !== source.version || !extracted.configs.length || !extracted.functions.length) {
    throw new Error(`${source.module}@${source.version}: incomplete API manifest`);
  }
  writeFileSync(join(generatedDir, `${slug}.json`), `${JSON.stringify(extracted, null, 2)}\n`);
}
console.log(`Prepared Echo source ${revision}`);
