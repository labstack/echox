// /pro is an English-only product page on the release site; the next preview
// (base /next/) only redirects there. One check for every place that treats it
// differently: header, footer, banner, route data and the sitemap filter.
export function isProPage(pathname) {
  return pathname.replace(/^\/next(?=\/|$)/, '').replace(/\/$/, '') === '/pro';
}
