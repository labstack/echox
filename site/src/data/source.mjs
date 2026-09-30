export function formatRevision(revision) {
  return /^[0-9a-f]{7,40}$/i.test(revision) ? revision.slice(0, 7) : revision;
}
