// Command-palette launchpad for the empty search dialog. It is built on the
// client so its markup ships once in the cached script bundle, not in every
// page's HTML. Icon paths are Starlight's built-in icons of the same names.
const icons = {
  rocket: '<path fill-rule="evenodd" d="M1.44 8.855v-.001l3.527-3.516c.34-.344.802-.541 1.285-.548h6.649l.947-.947c3.07-3.07 6.207-3.072 7.62-2.868a1.821 1.821 0 0 1 1.557 1.557c.204 1.413.203 4.55-2.868 7.62l-.946.946v6.649a1.845 1.845 0 0 1-.549 1.286l-3.516 3.528a1.844 1.844 0 0 1-3.11-.944l-.858-4.275-4.52-4.52-2.31-.463-1.964-.394A1.847 1.847 0 0 1 .98 10.693a1.843 1.843 0 0 1 .46-1.838Zm5.379 2.017-3.873-.776L6.32 6.733h4.638l-4.14 4.14Zm8.403-5.655c2.459-2.46 4.856-2.463 5.89-2.33.134 1.035.13 3.432-2.329 5.891l-6.71 6.71-3.561-3.56 6.71-6.711Zm-1.318 15.837-.776-3.873 4.14-4.14v4.639l-3.364 3.374Z" clip-rule="evenodd"/><path d="M9.318 18.345a.972.972 0 0 0-1.86-.561c-.482 1.435-1.687 2.204-2.934 2.619a8.22 8.22 0 0 1-1.23.302c.062-.365.157-.79.303-1.229.415-1.247 1.184-2.452 2.62-2.935a.971.971 0 1 0-.62-1.842c-.12.04-.236.084-.35.13-2.02.828-3.012 2.588-3.493 4.033a10.383 10.383 0 0 0-.51 2.845l-.001.016v.063c0 .536.434.972.97.972H2.24a7.21 7.21 0 0 0 .878-.065c.527-.063 1.248-.19 2.02-.447 1.445-.48 3.205-1.472 4.033-3.494a5.828 5.828 0 0 0 .147-.407Z"/>',
  random: '<path d="M8.7 10a1 1 0 0 0 1.41 0 1 1 0 0 0 0-1.41l-6.27-6.3a1 1 0 0 0-1.42 1.42ZM21 14a1 1 0 0 0-1 1v3.59L15.44 14A1 1 0 0 0 14 15.44L18.59 20H15a1 1 0 0 0 0 2h6a1 1 0 0 0 .38-.08 1 1 0 0 0 .54-.54A1 1 0 0 0 22 21v-6a1 1 0 0 0-1-1Zm.92-11.38a1 1 0 0 0-.54-.54A1 1 0 0 0 21 2h-6a1 1 0 0 0 0 2h3.59L2.29 20.29a1 1 0 0 0 0 1.42 1 1 0 0 0 1.42 0L20 5.41V9a1 1 0 0 0 2 0V3a1 1 0 0 0-.08-.38Z"/>',
  document: '<path d="M9 10h1a1 1 0 1 0 0-2H9a1 1 0 0 0 0 2Zm0 2a1 1 0 0 0 0 2h6a1 1 0 0 0 0-2H9Zm11-3.06a1.3 1.3 0 0 0-.06-.27v-.09c-.05-.1-.11-.2-.19-.28l-6-6a1.07 1.07 0 0 0-.28-.19h-.09a.88.88 0 0 0-.33-.11H7a3 3 0 0 0-3 3v14a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3V8.94Zm-6-3.53L16.59 8H15a1 1 0 0 1-1-1V5.41ZM18 19a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V5a1 1 0 0 1 1-1h5v3a3 3 0 0 0 3 3h3v9Zm-3-3H9a1 1 0 0 0 0 2h6a1 1 0 0 0 0-2Z"/>',
  padlock: '<path d="M12 13a1.5 1.5 0 0 0-1 2.6V17a1 1 0 0 0 2 0v-1.4a1.5 1.5 0 0 0-1-2.6m5-4V7A5 5 0 0 0 7 7v2a3 3 0 0 0-3 3v7a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3v-7a3 3 0 0 0-3-3M9 7a3 3 0 0 1 6 0v2H9Zm9 12a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-7a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1Z"/>',
  'open-book': '<path d="M21.17 2.06A13.1 13.1 0 0 0 19 1.87a12.94 12.94 0 0 0-7 2.05 12.94 12.94 0 0 0-7-2 13.1 13.1 0 0 0-2.17.19 1 1 0 0 0-.83 1v12a1 1 0 0 0 1.17 1 10.9 10.9 0 0 1 8.25 1.91l.12.07h.11a.91.91 0 0 0 .7 0h.11l.12-.07A10.899 10.899 0 0 1 20.83 16 1 1 0 0 0 22 15V3a1 1 0 0 0-.83-.94ZM11 15.35a12.87 12.87 0 0 0-6-1.48H4v-10c.333-.02.667-.02 1 0a10.86 10.86 0 0 1 6 1.8v9.68Zm9-1.44h-1a12.87 12.87 0 0 0-6 1.48V5.67a10.86 10.86 0 0 1 6-1.8c.333-.02.667-.02 1 0v10.04Zm1.17 4.15a13.098 13.098 0 0 0-2.17-.19 12.94 12.94 0 0 0-7 2.05 12.94 12.94 0 0 0-7-2.05c-.727.003-1.453.066-2.17.19A1 1 0 0 0 2 19.21a1 1 0 0 0 1.17.79 10.9 10.9 0 0 1 8.25 1.91 1 1 0 0 0 1.16 0A10.9 10.9 0 0 1 20.83 20a1 1 0 0 0 1.17-.79 1 1 0 0 0-.83-1.15Z"/>',
};

const en = { start: 'Start here', popular: 'Popular', quickstart: 'Quickstart', routing: 'Routing', binding: 'Binding', guide: 'Guide', middleware: 'Middleware', cookbook: 'Cookbook', open: 'open', navigate: 'navigate', close: 'close' };

export const labels: Record<string, typeof en> = {
  en,
  es: { start: 'Empieza aquí', popular: 'Popular', quickstart: 'Inicio rápido', routing: 'Rutas', binding: 'Binding', guide: 'Guía', middleware: 'Middleware', cookbook: 'Recetario', open: 'abrir', navigate: 'navegar', close: 'cerrar' },
  ja: { start: 'はじめに', popular: 'よく見るページ', quickstart: 'クイックスタート', routing: 'ルーティング', binding: 'バインディング', guide: 'ガイド', middleware: 'ミドルウェア', cookbook: 'クックブック', open: '開く', navigate: '移動', close: '閉じる' },
  'pt-br': { start: 'Comece aqui', popular: 'Populares', quickstart: 'Guia de início', routing: 'Roteamento', binding: 'Binding', guide: 'Guia', middleware: 'Middleware', cookbook: 'Receitas', open: 'abrir', navigate: 'navegar', close: 'fechar' },
  'zh-cn': { start: '从这里开始', popular: '热门内容', quickstart: '快速入门', routing: '路由', binding: '数据绑定', guide: '指南', middleware: '中间件', cookbook: '示例', open: '打开', navigate: '导航', close: '关闭' },
};

export function launchpadHTML(locale: string, docsBase: string) {
  const l = labels[locale];
  const link = (path: string, icon: keyof typeof icons, text: string, section: string) =>
    `<a class="ess-link" href="${docsBase}${path}"><svg aria-hidden="true" viewBox="0 0 24 24" fill="currentColor">${icons[icon]}</svg> ${text} <span class="ess-sec">${section}</span></a>`;
  return `<div class="echo-search-empty">
    <div class="ess-group">
      <div class="ess-label">${l.start}</div>
      ${link('/guide/quickstart/', 'rocket', l.quickstart, l.guide)}
      ${link('/guide/routing/', 'random', l.routing, l.guide)}
      ${link('/guide/binding/', 'document', l.binding, l.guide)}
    </div>
    <div class="ess-group">
      <div class="ess-label">${l.popular}</div>
      ${link('/middleware/jwt/', 'padlock', 'JWT', l.middleware)}
      ${link('/middleware/cors/', 'random', 'CORS', l.middleware)}
      ${link('/cookbook/hello-world/', 'open-book', 'Hello World', l.cookbook)}
    </div>
    <div class="ess-foot">
      <span><kbd>↵</kbd>${l.open}</span><span><kbd>Tab</kbd>${l.navigate}</span><span><kbd>esc</kbd>${l.close}</span>
    </div>
  </div>`;
}
