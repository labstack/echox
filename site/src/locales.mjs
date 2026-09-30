export const locales = {
  root: { label: 'English', lang: 'en' },
  'zh-cn': { label: '简体中文', lang: 'zh-CN' },
  ja: { label: '日本語', lang: 'ja' },
  es: { label: 'Español', lang: 'es' },
  'pt-br': { label: 'Português', lang: 'pt-BR' },
};

export const localePrefixes = Object.keys(locales).map((locale) => locale === 'root' ? '' : `${locale}/`);
