import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightDotMd from 'starlight-dot-md';
import starlightLlmsTxt from 'starlight-llms-txt';
import { redirects } from './src/redirects.mjs';

const next = process.env.DOCS_CHANNEL === 'next';

// https://astro.build/config
export default defineConfig({
  site: 'https://echo.labstack.com',
  base: next ? '/next/' : '/',
  outDir: next ? './dist-next' : './dist',
  // Preserve every live Docusaurus /docs/* URL at cutover (generated — see ./src/redirects.mjs).
  redirects: next ? {} : redirects,
  integrations: [
    starlight({
      title: 'Echo',
      defaultLocale: 'root',
      locales: {
        root: { label: 'English', lang: 'en' },
        'zh-cn': { label: '简体中文', lang: 'zh-CN' },
        ja: { label: '日本語', lang: 'ja' },
        es: { label: 'Español', lang: 'es' },
        'pt-br': { label: 'Português', lang: 'pt-BR' },
      },
      logo: {
        light: './src/assets/logo-light.svg',
        dark: './src/assets/logo-dark.svg',
        replacesTitle: true,
      },
      customCss: ['./src/styles/terminal.css'],
      // Machine-readable surfaces for coding agents. Two separate concerns:
      //
      //   starlight-dot-md   — a raw Markdown twin of every page at <path>.md,
      //                        so an agent fetching a doc gets the source
      //                        instead of parsing rendered HTML.
      //   starlight-llms-txt — /llms.txt (a link manifest), /llms-full.txt
      //                        (everything concatenated) and /llms-small.txt
      //                        (trimmed for small context windows).
      //
      // These do nothing for search ranking: AI search crawlers overwhelmingly
      // ignore llms.txt and Google states it has no effect on Search or AI
      // Overviews. They are here for the audience that does read them — Claude
      // Code, Cursor, Copilot, Cline and friends, writing Echo code.
      plugins: [
        starlightDotMd(),
        starlightLlmsTxt({
          projectName: 'Echo',
          description:
            'High performance, minimalist Go web framework. The current major version is v5, imported as `github.com/labstack/echo/v5`.',
          // The single most useful thing this file can carry. Models trained
          // before v5 default to v4 APIs, and v4 code fails to compile against
          // v5 in ways that look like user error rather than a version skew.
          details: [
            '## Versions',
            '',
            'This documentation describes **Echo v5**, the current major version.',
            '',
            '- Module path: `github.com/labstack/echo/v5` — the `/v5` suffix is required.',
            '- Echo v4 (`github.com/labstack/echo/v4`) is in long-term support and is a',
            '  different API. Code written for v4 does not compile against v5 unchanged.',
            '- If you are asked for Echo code without a version being stated, use v5 and',
            '  say which version you used.',
            '',
            '## Using this documentation',
            '',
            'Every page is also served as raw Markdown by appending `.md` to its URL,',
            'e.g. `https://echo.labstack.com/guide/routing.md`. Prefer those over the',
            'HTML when reading a single page.',
          ].join('\n'),
          optionalLinks: [
            {
              label: 'API reference (pkg.go.dev)',
              url: 'https://pkg.go.dev/github.com/labstack/echo/v5',
              description: 'Generated godoc for every exported symbol in v5.',
            },
            {
              label: 'Source (GitHub)',
              url: 'https://github.com/labstack/echo',
              description: 'v5 lives on `master`; v4 lives on the `v4` branch.',
            },
          ],
          // Guide before middleware before cookbook — the order someone
          // learning the framework needs, not alphabetical.
          promote: ['index*', 'guide/**'],
          demote: ['cookbook/**'],
          // llms-small.txt is for small context windows, so it has to be
          // meaningfully smaller than llms-full.txt to be worth generating.
          // The cookbook is 20 pages of standalone recipes — the least useful
          // thing to spend a constrained budget on. Code blocks stay: this is
          // a Go framework, and prose without the snippets is not usable.
          exclude: ['cookbook/**'],
          minify: {
            note: true,
            tip: true,
            details: true,
            whitespace: true,
            customSelectors: [],
          },
        }),
      ],
      social: [
        { icon: 'github', label: 'GitHub', href: 'https://github.com/labstack/echo' },
      ],
      editLink: {
        baseUrl: 'https://github.com/labstack/echox/edit/master/site/',
      },
      lastUpdated: true,
      // Echo "E" cube mark; .ico kept as legacy fallback, apple-touch-icon added in head.
      favicon: '/favicon.svg',
      // Keep Starlight's built-in Pagefind ⌘K search; Search override adds the
      // empty-state launchpad. "Ask AI" is the kapa.ai widget (see head).
      components: {
        Banner: './src/components/VersionBanner.astro',
        Footer: './src/components/Footer.astro',
        Search: './src/components/Search.astro',
      },
      head: [
        // Google Analytics (carried over from the Docusaurus site, anonymized IP).
        { tag: 'script', attrs: { async: true, src: 'https://www.googletagmanager.com/gtag/js?id=G-H19TMZLQFN' } },
        {
          tag: 'script',
          content:
            "window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','G-H19TMZLQFN',{anonymize_ip:true});",
        },
        // kapa.ai "Ask AI" widget — real answers from the indexed Echo docs.
        // We hide kapa's default floating launcher and open the modal from our
        // own "Ask AI" pill in the header (see Search.astro) so the trigger
        // matches the Terminal chrome instead of kapa's stock button.
        ...(!next ? [{
          tag: 'script',
          attrs: {
            async: true,
            src: 'https://widget.kapa.ai/kapa-widget.bundle.js',
            'data-website-id': '3ff47090-a571-4c4a-bec0-c0a377028db5',
            'data-project-name': 'Echo',
            // Modal header title (defaults to "Echo Docs AI" from the project name).
            'data-modal-title': 'Ask AI',
            'data-project-color': '#00afd1',
            // Modal header logo — the same sparkle as the "Ask AI" header pill
            // (not the Echo cube). The floating launcher that also used this is hidden.
            'data-project-logo': '/ask-ai.svg',
            // Sync the widget's light/dark with the site (we set data-theme on <html>).
            'data-color-scheme-selector': "[data-theme='dark']",
            // Hide the stock floating button; kapa wires clicks on our header pill.
            'data-button-hide': 'true',
            'data-modal-override-open-selector': '#echo-ask-ai',
          },
        }] : []),
        // Dark-first: default new visitors to dark unless they've chosen otherwise.
        {
          tag: 'script',
          content: "try{if(!localStorage.getItem('starlight-theme')){localStorage.setItem('starlight-theme','dark');document.documentElement.dataset.theme='dark';}}catch(e){document.documentElement.dataset.theme='dark';}",
        },
        { tag: 'link', attrs: { rel: 'preconnect', href: 'https://fonts.googleapis.com' } },
        { tag: 'link', attrs: { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: true } },
        {
          tag: 'link',
          attrs: {
            rel: 'stylesheet',
            href: 'https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700;800&family=Fragment+Mono&display=swap',
          },
        },
        {
          tag: 'link',
          attrs: {
            rel: 'stylesheet',
            href: 'https://unpkg.com/@phosphor-icons/web@2.1.1/src/regular/style.css',
          },
        },
        // Default social card. Starlight already emits og:title/description/url
        // and twitter:card=summary_large_image, but no image — add a site-wide
        // default so shares aren't imageless. Absolute URLs are required by
        // social scrapers. Per-page overrides can set their own og:image later.
        { tag: 'meta', attrs: { property: 'og:image', content: 'https://echo.labstack.com/og.png' } },
        { tag: 'meta', attrs: { property: 'og:image:width', content: '1200' } },
        { tag: 'meta', attrs: { property: 'og:image:height', content: '630' } },
        { tag: 'meta', attrs: { name: 'twitter:image', content: 'https://echo.labstack.com/og.png' } },
        { tag: 'link', attrs: { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' } },
      ],
      // Autogenerated from the content dirs — new pages appear automatically,
      // ordered by each page's `sidebar.order` frontmatter.
      sidebar: [
        { label: 'Guide', translations: { 'zh-CN': '指南', ja: 'ガイド', es: 'Guía', 'pt-BR': 'Guia' }, items: [{ autogenerate: { directory: 'guide' } }] },
        { label: 'Middleware', translations: { 'zh-CN': '中间件', ja: 'ミドルウェア', es: 'Middleware', 'pt-BR': 'Middleware' }, items: [{ autogenerate: { directory: 'middleware' } }] },
        { label: 'Cookbook', translations: { 'zh-CN': '示例', ja: 'クックブック', es: 'Recetario', 'pt-BR': 'Receitas' }, items: [{ autogenerate: { directory: 'cookbook' } }] },
      ],
      // tune the built-in code theme toward our terminal palette
      expressiveCode: { themes: ['github-dark', 'github-light'] },
    }),
  ],
});
