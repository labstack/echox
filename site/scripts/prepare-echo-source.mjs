import { execFileSync } from 'node:child_process';
import { copyFileSync, existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const siteDir = fileURLToPath(new URL('..', import.meta.url));
const repoDir = dirname(siteDir);
const pin = JSON.parse(readFileSync(join(siteDir, 'echo-source.json'), 'utf8'));
const override = process.env.ECHO_SOURCE_DIR;
const sourceDir = override ? resolve(siteDir, override) : join(repoDir, '.cache', 'echo-source');
const referenceDir = join(repoDir, 'reference');
const generatedDir = join(siteDir, 'src', 'generated');

function run(command, args, cwd = repoDir, environment = {}) {
  try {
    return execFileSync(command, args, { cwd, env: { ...process.env, ...environment }, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();
  } catch (error) {
    const detail = error.stderr?.toString().trim() || error.message;
    throw new Error(`${command} ${args.join(' ')} failed: ${detail}`);
  }
}

if (!override) {
  mkdirSync(dirname(sourceDir), { recursive: true });
  if (!existsSync(join(sourceDir, '.git'))) {
    run('git', ['clone', '--filter=blob:none', '--no-checkout', pin.repository, sourceDir]);
  }
  const current = run('git', ['rev-parse', 'HEAD'], sourceDir);
  if (current !== pin.revision || !existsSync(join(sourceDir, 'go.mod'))) {
    if (current !== pin.revision) run('git', ['fetch', '--depth=1', 'origin', pin.revision], sourceDir);
    run('git', ['checkout', '--detach', pin.revision], sourceDir);
  }
}

const revision = run('git', ['rev-parse', 'HEAD'], sourceDir);
if (!override && revision !== pin.revision) {
  throw new Error(`Echo checkout is ${revision}; expected ${pin.revision}`);
}
if (run('git', ['status', '--porcelain', '--untracked-files=no'], sourceDir)) {
  throw new Error(`Echo checkout has tracked changes: ${sourceDir}`);
}

const workspaceDir = join(repoDir, '.cache', 'echo-work');
mkdirSync(workspaceDir, { recursive: true });
const workspaceFile = join(workspaceDir, 'go.work');
writeFileSync(workspaceFile, `go 1.27.0\n\nuse (\n\t${JSON.stringify(sourceDir)}\n\t${JSON.stringify(referenceDir)}\n)\n`);
const goEnvironment = { GOWORK: workspaceFile };
const manifestText = run('go', ['run', './cmd/config-fields', '-root', sourceDir, '-revision', revision], referenceDir, goEnvironment);
const manifest = JSON.parse(manifestText);
if (manifest.revision !== revision || !manifest.configs.some((config) => config.name === 'RequestLoggerConfig') || !manifest.configs.some((config) => config.name === 'StaticConfig')) {
  throw new Error('Echo source extractor returned an incomplete or mismatched manifest');
}
run('go', ['test', './...'], referenceDir, goEnvironment);
mkdirSync(generatedDir, { recursive: true });
writeFileSync(join(generatedDir, 'config-fields.json'), `${JSON.stringify(manifest, null, 2)}\n`);
copyFileSync(join(referenceDir, 'request-logger', 'main.go'), join(generatedDir, 'request-logger.go.txt'));
copyFileSync(join(referenceDir, 'static', 'main.go'), join(generatedDir, 'static.go.txt'));
console.log(`Prepared Echo source ${revision}`);
