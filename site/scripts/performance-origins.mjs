// Count resource URLs declared by the page, rather than sites readers may visit.
export function externalDomains(html, ownHost = 'echo.labstack.com') {
  const markup = html.replace(/<!--[\s\S]*?-->/g, '')
    .replace(/(<script\b[^>]*>)[\s\S]*?<\/script>/gi, '$1</script>');
  const domains = new Set();
  const add = (value) => {
    if (!value) return;
    let url;
    try { url = new URL(value.replaceAll('&amp;', '&'), `https://${ownHost}/`); } catch { return; }
    if (['http:', 'https:'].includes(url.protocol) && url.hostname !== ownHost) domains.add(url.hostname);
  };
  for (const [, tag, attributes] of markup.matchAll(/<(script|link|img|iframe|source|video|audio|embed|object|input)\b([^>]*)>/gi)) {
    const attrs = Object.fromEntries([...attributes.matchAll(/([^\s=/'">]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]+)))?/g)]
      .map(([, key, double, single, unquoted]) => [key.toLowerCase(), double ?? single ?? unquoted]));
    const name = tag.toLowerCase();
    if (name === 'link') {
      if (attrs.rel?.toLowerCase().split(/\s+/).some((rel) => ['stylesheet', 'preload', 'modulepreload', 'icon', 'apple-touch-icon', 'mask-icon'].includes(rel))) add(attrs.href);
      continue;
    }
    if (name === 'input' && attrs.type?.toLowerCase() !== 'image') continue;
    add(name === 'object' ? attrs.data : attrs.src);
    if (name === 'video') add(attrs.poster);
    if (name === 'img' || name === 'source') {
      for (const [, value] of (attrs.srcset || '').matchAll(/(?:^|,)\s*((?:https?:)?\/\/[^\s,]+)/gi)) add(value);
    }
  }
  return [...domains].sort();
}
