import { appendFileSync } from 'node:fs';

export function reportApiChanges(channel, revision, differences, summaryFile = process.env.GITHUB_STEP_SUMMARY) {
  if (!summaryFile) return;
  const result = differences.length
    ? `Review required: ${differences.length} API difference(s).\n\n\`\`\`text\n${differences.join('\n')}\n\`\`\`\n`
    : 'No differences from the reviewed API baseline.\n';
  appendFileSync(summaryFile, `## ${channel} middleware API (${revision})\n\n${result}\n`);
}
