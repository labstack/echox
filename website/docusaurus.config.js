// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const {themes} = require('prism-react-renderer');
const lightTheme = themes.github;
const darkTheme = themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Echo',
  tagline: 'High performance, extensible, minimalist Go web framework',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://echo.labstack.com/',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'labstack', // Usually your GitHub org/user name.
  projectName: 'echo', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  trailingSlash: false,

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/labstack/echox/blob/master/website/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        gtag: {
          trackingID: 'G-H19TMZLQFN',
          anonymizeIP: true,
        }
      }),
    ]
  ],

  plugins: [
    [
      '@docusaurus/plugin-client-redirects',
      {
        redirects: [
          {
            to: '/docs/start-server',
            from: '/guide/http_server'
          },
          {
            to: '/docs/cookbook/subdomain',
            from: '/cookbook/subdomains'
          },
          {
            to: '/docs/cookbook/timeout',
            from: '/cookbook/timeouts'
          }
        ],
        createRedirects(existingPath) {
          if (existingPath.includes('/docs/middleware')) {
            return [
              existingPath.replace('/docs/middleware', '/middleware'),
            ];
          }
          if (existingPath.includes('/docs/cookbook')) {
            return [
              existingPath.replace('/docs/cookbook', '/cookbook'),
            ];
          }
          if (existingPath.includes('/docs')) {
            return [
              existingPath.replace('/docs', '/guide'),
            ];
          }
          return undefined; // Return a falsy value: no redirect created
        },
      },
    ],
  ],

  themes: ['docusaurus-theme-github-codeblock'],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.jpg',
      navbar: {
        logo: {
          alt: 'Echo',
          src: 'img/logo-light.svg',
          srcDark: 'img/logo-dark.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'docsSidebar',
            position: 'left',
            label: 'Docs',
          },
          // {to: '/blog', label: 'Blog', position: 'left'},
          {
            href: 'https://github.com/labstack/echo/discussions',
            label: 'Discussions',
            position: 'right',
          },
          {
            href: 'https://github.com/labstack/echo',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        copyright: `Copyright © ${new Date().getFullYear()} Echo, LabStack LLC.`,
      },
      prism: {
        theme: lightTheme,
        darkTheme: darkTheme,
      },
      codeblock: {
        showGithubLink: false,
        showRunmeLink: false
      },
      algolia: {
        appId: '4AFTIR2HGK',
        apiKey: 'b92d819e8e37b59dc5d564d763542682',
        indexName: 'labstack_echo',
        contextualSearch: false,
        facetFilters: []
      }
    }),
};

module.exports = config;
