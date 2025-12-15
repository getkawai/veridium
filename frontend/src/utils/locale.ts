import { resolveAcceptLanguage } from 'resolve-accept-language';

import { DEFAULT_LANG } from '@/const/locale';
import { Locales, locales, normalizeLocale } from '@/locales/resources';

export const getAntdLocale = async (lang?: string) => {
  let normalLang: any = normalizeLocale(lang);

  // due to antd only have ar-EG locale, we need to convert ar to ar-EG
  // refs: https://ant.design/docs/react/i18n

  // And we don't want to handle it in `normalizeLocale` function
  // because of other locale files are all `ar` not `ar-EG`
  if (normalLang === 'ar') normalLang = 'ar-EG';

  const { default: locale } = await import(/* @vite-ignore */ `antd/locale/${normalLang.replace('-', '_')}.js`);

  return locale;
};

/**
 * Parse the browser language and return the fallback language
 */
export const parseBrowserLanguage = (headers?: Headers | string, defaultLang: string = DEFAULT_LANG) => {
  // if the default language is not 'en-US', just return the default language as fallback lang
  if (defaultLang !== 'en-US') return defaultLang;

  let acceptLanguage = '';

  // Handle different input types for browser compatibility
  if (typeof headers === 'string') {
    acceptLanguage = headers;
  } else if (headers instanceof Headers) {
    acceptLanguage = headers.get('accept-language') || '';
  } else {
    // Browser environment - use navigator.languages
    acceptLanguage = navigator.languages?.join(',') || navigator.language || '';
  }

  /**
   * The arguments are as follows:
   *
   * 1) The HTTP accept-language header or browser languages.
   * 2) The available locales (they must contain the default locale).
   * 3) The default locale.
   */
  let browserLang: string = resolveAcceptLanguage(
    acceptLanguage,
    //  Invalid locale identifier 'ar'. A valid locale should follow the BCP 47 'language-country' format.
    locales.map((locale) => (locale === 'ar' ? 'ar-EG' : locale)),
    defaultLang,
  );

  // if match the ar-EG then fallback to ar
  if (browserLang === 'ar-EG') browserLang = 'ar';

  return browserLang;
};

/**
 * Parse the page locale from the URL and search
 * Browser-compatible version
 */
export const parsePageLocale = (url?: string): Locales => {
  // Use provided URL or current browser location
  const currentUrl = url || (typeof window !== 'undefined' ? window.location.href : '');
  const urlObj = new URL(currentUrl);

  // Get locale from URL search params (hl parameter)
  const searchParams = urlObj.searchParams;
  const hlParam = searchParams.get('hl');

  // Get locale from URL path if not in search params
  // This handles cases like /en/page or /zh-CN/page
  const pathSegments = urlObj.pathname.split('/').filter(Boolean);
  const pathLocale = pathSegments[0];

  // Try to get browser language as fallback
  const browserLocale = parseBrowserLanguage();

  // Priority: URL search param > URL path > browser language > default
  const detectedLocale = hlParam || pathLocale || browserLocale;

  return normalizeLocale(detectedLocale) as Locales;
};

// Note: Server-side locale parsing is not available in browser environment
// Use parsePageLocale() for browser-compatible locale detection
