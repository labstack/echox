import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import { unified } from '@astrojs/markdown-remark';
import { redirects, retiredCookbookRedirects } from './src/redirects.mjs';
import remarkSourceCode from './plugins/remark-source-code.mjs';
import { locales } from './src/locales.mjs';

const next = process.env.DOCS_CHANNEL === 'next';

// https://astro.build/config
export default defineConfig({
  site: 'https://echo.labstack.com',
  base: next ? '/next/' : '/',
  outDir: next ? './dist-next' : './dist',
  markdown: { processor: unified({ remarkPlugins: [remarkSourceCode] }) },
  // Sponsor avatars: github.com/<org>.png redirects to avatars.githubusercontent.com.
  image: { domains: ['github.com', 'avatars.githubusercontent.com'] },
  // Preserve every live Docusaurus /docs/* URL at cutover (generated — see ./src/redirects.mjs).
  redirects: next ? retiredCookbookRedirects('/next/') : redirects,
  integrations: [
    starlight({
      title: 'Echo',
      defaultLocale: 'root',
      locales,
      logo: {
        light: './src/assets/logo-light.svg',
        dark: './src/assets/logo-dark.svg',
        replacesTitle: true,
      },
      customCss: ['./src/styles/terminal.css'],
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
            // Modal header logo — the same star as the "Ask AI" header pill
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
