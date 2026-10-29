import { defineRouteMiddleware } from '@astrojs/starlight/route-data';
import { isProPage } from './pro-page.mjs';
import provenance from '../pro-source/pro-source.json';

const groups = [
  ['Start here', [['install', 'Install Echo Pro'], ['getting-started', 'Getting started']]],
  ['Packages', [['auth', 'Authentication'], ['tenant', 'Tenant resolution'], ['quota', 'Quotas'], ['webhook', 'Webhooks'], ['audit', 'Audit'], ['execution', 'Execution'], ['admin', 'Admin console']]],
  ['Tasks', [['guides/auth', 'Configure authentication'], ['guides/tenant', 'Authorize resources'], ['guides/quota', 'Configure quotas'], ['guides/webhook', 'Receive and deliver events'], ['guides/audit', 'Record audit events'], ['guides/execution', 'Bind retries and recover'], ['guides/admin', 'Mount operator panels']]],
  ['Examples', [['examples/quota-otel', 'Quota telemetry'], ['examples/stripe', 'Stripe usage export'], ['examples/admin', 'Local admin'], ['examples/admin-demo', 'Full console demo']]],
  ['Reference', [['reference/auth', 'Auth configuration'], ['reference/tenant', 'Tenant configuration'], ['reference/quota', 'Quota configuration'], ['reference/webhook', 'Webhook configuration'], ['reference/audit', 'Audit configuration'], ['reference/execution', 'Execution configuration'], ['reference/admin', 'Admin configuration'], ['problems', 'Problem types']]],
] as const;

export const onRequest = defineRouteMiddleware((context) => {
  if (!isProPage(context.url.pathname)) return;
  const route = context.locals.starlightRoute;
  route.head = route.head.filter(
    ({ tag, attrs }) => !(tag === 'link' && attrs?.rel === 'alternate' && attrs?.hreflang) &&
      !(tag === 'script' && typeof attrs?.src === 'string' && attrs.src.startsWith('https://widget.kapa.ai/')),
  );
  if (context.url.pathname.replace(/\/$/, '') === '/pro') return;
  route.siteTitle = 'Echo Pro';
  route.siteTitleHref = '/pro';
  route.editUrl = undefined;
  route.lastUpdated = undefined;
  route.pagination = { prev: undefined, next: undefined };
  route.sidebar = groups.map(([label, items]) => ({
    type: 'group', label, collapsed: false, badge: undefined,
    entries: items.map(([slug, label]) => ({
      type: 'link', label, href: `/pro/${slug}/`,
      isCurrent: context.url.pathname.replace(/\/$/, '') === `/pro/${slug}`,
      badge: undefined, attrs: {},
    })),
  }));
  const title = `${route.entry.data.title} · Echo Pro`;
  for (const head of route.head) {
    if (head.tag === 'title') head.content = title;
    if (head.tag === 'meta' && (head.attrs?.property === 'og:title' || head.attrs?.name === 'twitter:title')) head.attrs.content = title;
  }
  route.head.push({ tag: 'meta', attrs: { name: 'echo-pro-version', content: provenance.version } });
});
