import { defineRouteMiddleware } from '@astrojs/starlight/route-data';
import { isProPage } from './pro-page.mjs';

// /pro exists only in English and only on the release site, so drop the
// hreflang alternates Starlight adds for the other locales.
export const onRequest = defineRouteMiddleware((context) => {
  if (!isProPage(context.url.pathname)) return;
  const route = context.locals.starlightRoute;
  route.head = route.head.filter(
    ({ tag, attrs }) => !(tag === 'link' && attrs?.rel === 'alternate' && attrs?.hreflang),
  );
});
