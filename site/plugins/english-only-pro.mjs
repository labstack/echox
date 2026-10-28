/** Keep Starlight's renderer while filtering automatic locale fallback paths.
 * Intercept its Astro route registration, using the exported static route's
 * getStaticPaths in our adapter. Fail explicitly if upstream routing changes.
 */
export function englishOnlyPro(integration) {
  const setup = integration.hooks['astro:config:setup'];
  integration.hooks['astro:config:setup'] = async (context) => {
    let replaced = false;
    await setup({
      ...context,
      injectRoute(route) {
        if (route.pattern === '[...slug]') {
          if (route.entrypoint !== '@astrojs/starlight/routes/static/index.astro') throw new Error('Review Pro routing adapter after Starlight routing change');
          route = { ...route, entrypoint: './src/starlight-docs.astro' };
          replaced = true;
        }
        context.injectRoute(route);
      },
    });
    if (!replaced) throw new Error('Starlight docs route missing; review Pro routing adapter');
  };
  return integration;
}
