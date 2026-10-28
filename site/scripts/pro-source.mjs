import { createHash } from 'node:crypto';
import { existsSync, lstatSync, mkdirSync, readFileSync, readdirSync, rmSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parseSourceDocument, visitCode } from '../plugins/source-document.mjs';
import { sourceReference, sourceSnippet } from '../plugins/remark-source-code.mjs';

export const siteDir = fileURLToPath(new URL('..', import.meta.url));
const page = /^(?:install|getting-started|auth|tenant|quota|webhook|audit|execution|admin|recipes\/(?:index|quota-otel|stripe|admin|admin-demo)|(?:guides|reference)\/(?:index|auth|tenant|quota|webhook|audit|execution|admin)|problems\/(?:index|[a-z]+(?:-[a-z]+)*))\.mdx$/;
const example = /^(?:getting-started\/main|quota-otel\/main|stripe\/main|admin\/main|admin-demo\/(?:main|auth|server|quota|audit|execution|webhook|reports|throughput))\.go$/;

const problemExample = /^examples\/problems\/[a-z]+(?:-[a-z]+)*-[45]\d\d\.json$/;

const sourceExamples = new Set(["auth/quickstart_example_test.go", "auth/example_test.go", "auth/site_example_test.go", "auth/keystore_example_test.go", "auth/authbind/example_test.go", "tenant/example_test.go", "tenant/site_example_test.go", "quota/example_test.go", "quota/site_example_test.go", "quota/middleware_example_test.go", "quota/auth_example_test.go", "quota/plans_example_test.go", "quota/accounting_example_test.go", "quota/reset_example_test.go", "quota/observations_example_test.go", "quota/shadow_admin_example_test.go", "quota/redisstore/example_test.go", "quota/quotaadmin/credits_example_test.go", "quota/quotaotel/tracing_example_test.go", "quota/quotaotel/lifecycle_example_test.go", "quota/quotaprom/lifecycle_example_test.go", "quota/quotaclient/example_test.go", "quota/quotamcp/example_test.go", "quota/quotastripe/example_test.go", "webhook/example_test.go", "webhook/http_example_test.go", "webhook/webhookmem/example_test.go", "webhook/webhookpg/example_test.go", "webhook/webhookpg/recipe_example_test.go", "webhook/webhookadmin/example_test.go", "audit/quickstart_example_test.go", "audit/business_example_test.go", "audit/reader_example_test.go", "audit/middleware_example_test.go", "audit/edges_example_test.go", "audit/auditpg/example_test.go", "audit/auditadmin/example_test.go", "execution/example_test.go", "execution/http_example_test.go", "execution/problem_example_test.go", "execution/call_example_test.go", "execution/scheduler_example_test.go", "execution/pgstore/example_test.go", "execution/pgstore/transactional_example_test.go", "execution/executionadmin/example_test.go", "execution/executiontest/example_test.go", "admin/quickstart_example_test.go", "admin/example_test.go", "admin/admintest/example_test.go"]);
const installExamples = new Set(["site-docs/commands.sh", "site-docs/netrc.example", "site-docs/github-actions.yml", "site-docs/Dockerfile"]);

function filesIn(dir, prefix = '') {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    if (entry.isSymbolicLink()) throw new Error('Pro bundle must not contain symlinks');
    const name = prefix + entry.name;
    return entry.isDirectory() ? filesIn(join(dir, entry.name), `${name}/`) : [name];
  }).sort();
}

/** Validate only the committed public snapshot. Never access the private repo. */
export function readPro(dir = join(siteDir, 'pro-source')) {
  if (lstatSync(dir).isSymbolicLink()) throw new Error('Pro bundle root must not be a symlink');
  const provenance = JSON.parse(readFileSync(join(dir, 'pro-source.json'), 'utf8'));
  if (provenance.module !== 'echo.labstack.com/pro' || !/^v\d+\.\d+\.\d+$/.test(provenance.version) || !/^[a-f0-9]{40}$/.test(provenance.commit)) {
    throw new Error('Invalid Pro provenance');
  }
  const names = filesIn(dir).filter((name) => name !== 'pro-source.json');
  if (JSON.stringify(names) !== JSON.stringify(Object.keys(provenance.hashes).sort())) throw new Error('Pro bundle inventory differs from provenance');
  const files = new Map();
  for (const name of names) {
    const isExample = name.startsWith('examples/') && (example.test(name.slice(9)) || installExamples.has(name.slice(9)) || (name.startsWith('examples/source/') && sourceExamples.has(name.slice(16))));
    if (!['go.mod', 'config-fields.json', 'problems.json'].includes(name) && !(name.startsWith('pages/') && page.test(name.slice(6))) && !problemExample.test(name) && !isExample) {
      throw new Error(`Unreviewed Pro bundle path: ${name}`);
    }
    const content = readFileSync(join(dir, name));
    if (createHash('sha256').update(content).digest('hex') !== provenance.hashes[name]) throw new Error(`Pro content hash mismatch: ${name}`);
    if (/github\.com\/labstack\/echo-pro|echo\.labstack\.com\/pro\/internal\/|docs\/(?:private|[a-z0-9_/-]+\.md)/i.test(content.toString())) throw new Error(`Private reference in Pro bundle: ${name}`);
    files.set(name, content);
  }
  if (files.get('go.mod')?.toString() !== 'module echo.labstack.com/pro-site-snapshot\n\ngo 1.26.0\n') throw new Error('Pro examples must be isolated from public root Go tests');
  const api = JSON.parse(files.get('config-fields.json'));
  const problems = JSON.parse(files.get('problems.json'));
  if (api.module !== provenance.module || api.revision !== provenance.version) throw new Error('Pro API version differs from provenance');
  const expected = problems.filter((d) => d.Type.startsWith('https://echo.labstack.com/pro/problems/')).map((d) => `pages/problems/${d.Type.split('/').at(-1)}.mdx`).sort();
  const actual = names.filter((name) => name.startsWith('pages/problems/') && name !== 'pages/problems/index.mdx');
  if (JSON.stringify(actual) !== JSON.stringify(expected)) throw new Error('Pro problem pages differ from the real catalog');
  const localProblems = problems.filter((d) => d.Type.startsWith('https://echo.labstack.com/pro/problems/'));
  const expectedExamples = localProblems.flatMap((d) => d.Messages.map((m) => `examples/problems/${d.Type.split('/').at(-1)}-${m.Status}.json`)).sort();
  const actualExamples = names.filter((name) => name.startsWith('examples/problems/'));
  if (JSON.stringify(actualExamples) !== JSON.stringify(expectedExamples)) throw new Error('Pro problem examples differ from the real catalog statuses');
  for (const d of localProblems) {
    for (const m of d.Messages) {
      const sample = JSON.parse(files.get(`examples/problems/${d.Type.split('/').at(-1)}-${m.Status}.json`));
      if (sample.type !== d.Type || sample.code !== d.Code || sample.status !== m.Status) throw new Error('Pro problem example identity mismatch');
    }
  }
  for (const [name, content] of files) {
    if (!name.startsWith('pages/')) continue;
    visitCode(parseSourceDocument(content.toString(), name), (node) => {
      const reference = sourceReference(node);
      if (!reference || !reference.specifier.startsWith('site/pro-source/examples/')) throw new Error(`${name}: every Pro code block must include an exported source`);
      const [path, region, extra] = reference.specifier.slice('site/pro-source/'.length).split('#');
      if (extra !== undefined || region === '' || !files.has(path)) throw new Error(`${name}: invalid Pro source include`);
      sourceSnippet(files.get(path).toString(), region, path);
    });
  }
  return { provenance, files, api, problems };
}

export function apiSnapshot({ api, problems }) {
  return {
    module: api.module, version: api.revision,
    configs: Object.fromEntries(api.configs.map((c) => [c.name, c.fields.map(({ name, type, doc, default: defaults, deprecated }) => ({ name, type, doc, default: defaults, deprecated }))])),
    functions: Object.fromEntries(api.functions.map((f) => [f.name, f.signature])),
    problems,
  };
}

/** Materialize ephemeral collection pages; the bundle remains the source. */
export function preparePro() {
  const bundle = readPro();
  const target = join(siteDir, 'src/content/docs/pro');
  if (existsSync(target) && lstatSync(target).isSymbolicLink()) throw new Error('Pro content target must not be a symlink');
  rmSync(target, { recursive: true, force: true });
  if (process.env.DOCS_CHANNEL === 'next') return;
  for (const [name, content] of bundle.files) {
    if (!name.startsWith('pages/')) continue;
    const file = join(target, name.slice(6));
    mkdirSync(dirname(file), { recursive: true });
    writeFileSync(file, content);
  }
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  if (process.argv.includes('--prepare')) preparePro();
  else {
    const snapshot = `${JSON.stringify(apiSnapshot(readPro()), null, 2)}\n`;
    const baseline = join(siteDir, 'pro-reference-baseline.json');
    if (process.argv.includes('--accept')) {
      writeFileSync(baseline, snapshot);
      console.log('Accepted Pro declarations and problem catalog for review');
    } else if (!existsSync(baseline) || readFileSync(baseline, 'utf8') !== snapshot) {
      throw new Error('Pro API/catalog changed. Review pages and run npm run pro:accept.');
    } else console.log('Pro bundle hashes, publication boundary, declarations and problem catalog match');
  }
}
