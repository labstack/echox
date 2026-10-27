import { defineRouteMiddleware } from '@astrojs/starlight/route-data';

// /pro exists only in English and only on the release site, so drop the
// hreflang alternates Starlight adds for the other locales.
export const onRequest = defineRouteMiddleware((context) => {
  if (context.url.pathname.replace(/\/$/, '') !== '/pro') return;
  const route = context.locals.starlightRoute;
  route.head = route.head.filter(
    ({ tag, attrs }) => !(tag === 'link' && attrs?.rel === 'alternate' && attrs?.hreflang),
  );
});
