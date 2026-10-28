// Pro is English-only and independent of the OSS preview channel.
export function isProPage(pathname) {
  return /^\/pro(?:\/|$)/.test(pathname.replace(/^\/next(?=\/|$)/, ''));
}

// Starlight synthesizes locale fallbacks. Filter them before static rendering.
export function includeDocsRoute(slug = '', next = false) {
  if (/^(?:es|ja|pt-br|zh-cn)\/pro(?:\/|$)/.test(slug)) return false;
  return !(next && /^pro(?:\/|$)/.test(slug));
}
