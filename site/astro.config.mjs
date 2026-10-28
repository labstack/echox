import { defineConfig, fontProviders } from 'astro/config';
import starlightBase from '@astrojs/starlight';
import { englishOnlyPro } from './plugins/english-only-pro.mjs';
const starlight = (options) => englishOnlyPro(starlightBase(options));
import sitemap from '@astrojs/sitemap';
import { unified } from '@astrojs/markdown-remark';
import { redirects, retiredCookbookRedirects } from './src/redirects.mjs';
import remarkSourceCode from './plugins/remark-source-code.mjs';
import { locales } from './src/locales.mjs';
import { preparePro } from './scripts/pro-source.mjs';
import { isProPage } from './src/pro-page.mjs';

preparePro();

const next = process.env.DOCS_CHANNEL === 'next';

// https://astro.build/config
export default defineConfig({
  site: 'https://echo.labstack.com',
  base: next ? '/next/' : '/',
  outDir: next ? './dist-next' : './dist',
  // Self-hosted fonts (no third-party font requests). Latin subsets only; other
  // scripts fall back to the system stack. Astro generates size-adjusted
  // fallbacks and the preload links emitted by ./src/components/Head.astro.
  fonts: [
    {
      provider: fontProviders.local(),
      name: 'DM Sans',
      cssVariable: '--font-dm-sans',
      fallbacks: ['ui-sans-serif', 'system-ui', 'sans-serif'],
      display: 'swap',
      options: {
        variants: [{ src: ['@fontsource-variable/dm-sans/files/dm-sans-latin-wght-normal.woff2'], weight: '400 700', style: 'normal' }],
      },
    },
    {
      provider: fontProviders.local(),
      name: 'Fragment Mono',
      cssVariable: '--font-fragment-mono',
      fallbacks: ['ui-monospace', 'monospace'],
      display: 'swap',
      options: {
        variants: [{ src: ['@fontsource/fragment-mono/files/fragment-mono-latin-400-normal.woff2'], weight: '400', style: 'normal' }],
      },
    },
  ],
  markdown: { processor: unified({ remarkPlugins: [remarkSourceCode] }) },
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
      routeMiddleware: './src/routeData.ts',
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
        Head: './src/components/Head.astro',
        Header: './src/components/Header.astro',
        Footer: './src/components/Footer.astro',
        LanguageSelect: './src/components/LanguageSelect.astro',
        Search: './src/components/Search.astro',
        SocialIcons: './src/components/SocialIcons.astro',
      },
      head: [
        // Google Analytics (carried over from the Docusaurus site, anonymized IP).
        // Links with data-ga-event send that GA4 event on click; data-ga-plan
        // and data-ga-value add the plan name and USD value (begin_checkout);
        // data-ga-placement names where the link sits (nav, home_hero, ...).
        { tag: 'script', attrs: { async: true, src: 'https://www.googletagmanager.com/gtag/js?id=G-H19TMZLQFN' } },
        {
          tag: 'script',
          content:
            "window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','G-H19TMZLQFN',{anonymize_ip:true});" +
            "document.addEventListener('click',function(e){var a=e.target.closest&&e.target.closest('[data-ga-event]');if(!a)return;var d=a.dataset,p={link_text:a.textContent.trim(),link_url:a.getAttribute('href')};" +
            "if(d.gaPlacement){p.placement=d.gaPlacement;}if(d.gaPlan){p.plan=d.gaPlan;}if(d.gaValue){var v=Number(d.gaValue);p.currency='USD';p.value=v;p.items=[{item_id:'echo_pro_'+d.gaPlan.toLowerCase(),item_name:'Echo Pro '+d.gaPlan,price:v,quantity:1}];}" +
            "gtag('event',d.gaEvent,p);});",
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
    // Starlight adds @astrojs/sitemap only when it is missing; add it here with
    // the same i18n mapping and keep Pro only on the release site.
    sitemap({
      i18n: {
        defaultLocale: 'root',
        locales: Object.fromEntries(Object.entries(locales).map(([locale, { lang }]) => [locale, lang])),
      },
      // List the launched /pro page; skip only the /next/pro redirect stub.
      filter: (page) => { const path = new URL(page).pathname; return !(isProPage(path) && path.startsWith('/next/')); },
    }),
  ],
});
