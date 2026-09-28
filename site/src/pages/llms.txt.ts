import { getCollection } from 'astro:content';
import type { APIRoute } from 'astro';
import stableSource from '../../echo-source.json';
import nextSource from '../../next-source.json';

export const prerender = true;

const origin = 'https://echo.labstack.com';
const sections = [
  { title: 'Guide', prefix: 'guide/' },
  { title: 'Middleware', prefix: 'middleware/' },
  { title: 'Cookbook', prefix: 'cookbook/' },
];

export const GET: APIRoute = async () => {
  const base = import.meta.env.BASE_URL;
  const source = base === '/next/' ? nextSource : stableSource;
  const pages = await getCollection('docs');
  const englishPages = pages.filter(({ id }) =>
    !/^(?:es|ja|pt-br|zh-cn)\//.test(id),
  );
  const lines = [
    '# Echo',
    '',
    '> Echo is a Go web framework. These docs cover setup, routing, requests, responses, middleware, and runnable examples.',
    '',
    `This index describes the ${source.release} docs, built against [Echo source](${source.repository.replace(/\.git$/, '')}/tree/${source.revision}). For package API signatures, see [pkg.go.dev](https://pkg.go.dev/github.com/labstack/echo/v5).`,
    '',
  ];

  for (const { title, prefix } of sections) {
    lines.push(`## ${title}`, '');
    const entries = englishPages
      .filter(({ id }) => id.startsWith(prefix))
      .sort((a, b) =>
        (a.data.sidebar?.order ?? 999) - (b.data.sidebar?.order ?? 999) ||
        a.data.title.localeCompare(b.data.title),
      );
    for (const { id, data } of entries) {
      const route = id.endsWith('/index') ? id.slice(0, -'/index'.length) : id;
      lines.push(`- [${data.title}](${origin}${base}${route}/): ${data.description}`);
    }
    lines.push('');
  }

  lines.push('## Other versions and languages', '');
  if (base === '/') lines.push(`- [Next version](${origin}/next/llms.txt): preview docs for the next Echo release.`);
  else lines.push(`- [Stable version](${origin}/llms.txt): current release docs.`);
  for (const [locale, label] of [
    ['es', 'Español'],
    ['ja', '日本語'],
    ['pt-br', 'Português'],
    ['zh-cn', '简体中文'],
  ]) {
    lines.push(`- [${label}](${origin}${base}${locale}/): localized documentation.`);
  }

  return new Response(`${lines.join('\n')}\n`, {
    headers: { 'Content-Type': 'text/plain; charset=utf-8' },
  });
};
