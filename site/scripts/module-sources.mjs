import { execFileSync } from 'node:child_process';
import { mkdirSync, readFileSync, rmSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';

export const echoModule = 'github.com/labstack/echo/v5';

export function run(command, args, cwd, environment = {}) {
  try {
    return execFileSync(command, args, { cwd, env: { ...process.env, ...environment }, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] }).trim();
  } catch (error) {
    const detail = error.stderr?.toString().trim() || error.message;
    throw new Error(`${command} ${args.join(' ')} failed: ${detail}`, { cause: error });
  }
}

function selectedModule(module, cwd) {
  const selected = JSON.parse(run('go', ['list', '-m', '-json', module], cwd, { GOWORK: 'off' }));
  if (selected.Path !== module || !selected.Version || selected.Replace) {
    throw new Error(`${cwd}/go.mod must select a versioned, unreplaced ${module}`);
  }
  return { module, version: selected.Version };
}

export function moduleSources(siteDir) {
  const repoDir = dirname(siteDir);
  const referenceDir = join(repoDir, 'reference');
  const echo = selectedModule(echoModule, repoDir);
  const referenceEcho = selectedModule(echoModule, referenceDir);
  if (referenceEcho.version !== echo.version) {
    throw new Error(`Echo versions differ: go.mod selects ${echo.version}, reference/go.mod selects ${referenceEcho.version}. Update both modules in the same PR.`);
  }
  const rootRequirements = JSON.parse(run('go', ['mod', 'edit', '-json'], repoDir, { GOWORK: 'off' })).Require || [];
  const external = JSON.parse(readFileSync(join(siteDir, 'external-sources.json'), 'utf8'));
  for (const [slug, source] of Object.entries(external)) {
    const selected = selectedModule(source.module, referenceDir);
    if (rootRequirements.some((entry) => entry.Path === source.module)) {
      const root = selectedModule(source.module, repoDir);
      if (root.version !== selected.version) {
        throw new Error(`${source.module} versions differ between go.mod and reference/go.mod (${root.version} / ${selected.version})`);
      }
    }
    external[slug] = { ...source, version: selected.version };
  }
  return { echo, external };
}

export function downloadModule(source, cwd) {
  // An unqualified module argument uses the version selected by this go.mod.
  const download = JSON.parse(run('go', ['mod', 'download', '-json', source.module], cwd, { GOWORK: 'off' }));
  if (download.Error || !download.Dir || download.Path !== source.module || download.Version !== source.version) {
    throw new Error(`${source.module}@${source.version}: ${download.Error || 'download does not match the selected module'}`);
  }
  return download;
}

export function echoSource(download) {
  return {
    channel: 'stable',
    repository: 'https://github.com/labstack/echo.git',
    // Origin metadata is optional and cached downloads may omit it. Tags and
    // pseudo-version commit suffixes are deterministic GitHub source refs.
    // Go verifies the downloaded module via go.sum in either case.
    revision: download.Version.match(/[-.]\d{14}-([0-9a-f]{12})$/)?.[1] || download.Version,
    release: download.Version,
  };
}

export function createWorkspace(workspaceDir, sourceDir, referenceDir) {
  mkdirSync(workspaceDir, { recursive: true });
  const workspaceFile = join(workspaceDir, 'go.work');
  rmSync(workspaceFile, { force: true });
  // Let Go choose a compatible directive for the toolchain and both modules.
  run('go', ['work', 'init', sourceDir, referenceDir], workspaceDir, { GOWORK: 'off' });
  return workspaceFile;
}

export function checkReference(channel, referenceDir, environment) {
  const tags = channel === 'stable' ? ['-tags=docs_external'] : [];
  run('go', ['vet', ...tags, './...'], referenceDir, environment);
  run('go', ['test', '-race', ...tags, './...'], referenceDir, environment);
}

export function preparedSource(siteDir) {
  const channel = process.env.DOCS_CHANNEL === 'next' ? 'next' : 'stable';
  const sources = moduleSources(siteDir);
  const source = JSON.parse(readFileSync(join(siteDir, 'src/generated/echo-source.json'), 'utf8'));
  const manifest = JSON.parse(readFileSync(join(siteDir, 'src/generated/config-fields.json'), 'utf8'));
  const nextPin = channel === 'next' ? JSON.parse(readFileSync(join(siteDir, 'next-source.json'), 'utf8')) : null;
  const revision = channel === 'stable' ? echoSource({ Version: sources.echo.version }).revision : process.env.ECHO_SOURCE_DIR
    ? run('git', ['rev-parse', 'HEAD'], resolve(siteDir, process.env.ECHO_SOURCE_DIR))
    : nextPin.revision;
  if (source.channel !== channel || manifest.module !== echoModule || manifest.revision !== source.revision ||
      (channel === 'stable' && source.release !== sources.echo.version) ||
      (revision && source.revision !== revision)) {
    throw new Error(`Generated Echo source does not match the ${channel} channel and selected module. Run source:prepare for this channel first.`);
  }
  return { ...sources, source, manifest };
}
