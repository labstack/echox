// Count resources the page loads, rather than sites readers may visit.
export function externalDomains(html, ownHost = 'echo.labstack.com') {
  const markup = html.replace(/<!--[\s\S]*?-->/g, '')
    .replace(/(<script\b[^>]*>)[\s\S]*?<\/script>/gi, '$1</script>');
  const domains = new Set();
  for (const [, tag, attributes] of markup.matchAll(/<(script|link)\b([^>]*)>/gi)) {
    const attrs = Object.fromEntries([...attributes.matchAll(/([^\s=/'">]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]+)))?/g)]
      .map(([, key, double, single, unquoted]) => [key.toLowerCase(), double ?? single ?? unquoted]));
    const script = tag.toLowerCase() === 'script';
    if (!script && !attrs.rel?.toLowerCase().split(/\s+/).includes('stylesheet')) continue;
    const value = script ? attrs.src : attrs.href;
    if (!value) continue;
    let url;
    try { url = new URL(value.replaceAll('&amp;', '&'), `https://${ownHost}/`); } catch { continue; }
    if (['http:', 'https:'].includes(url.protocol) && url.hostname !== ownHost) domains.add(url.hostname);
  }
  return [...domains].sort();
}
