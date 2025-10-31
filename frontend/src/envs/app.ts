/* eslint-disable sort-keys-fix/sort-keys-fix */

export const getAppConfig = () => ({
  NEXT_PUBLIC_ENABLE_SENTRY: false,

  ACCESS_CODES: [],

  AGENTS_INDEX_URL: 'https://dummy-agents-index.com',

  DEFAULT_AGENT_CONFIG: '',
  SYSTEM_AGENT: undefined,

  PLUGINS_INDEX_URL: 'https://dummy-plugins-index.com',
  PLUGIN_SETTINGS: undefined,

  APP_URL: 'https://dummy-app-url.com',
  VERCEL_EDGE_CONFIG: undefined,
  MIDDLEWARE_REWRITE_THROUGH_LOCAL: false,
  ENABLE_AUTH_PROTECTION: false,

  CDN_USE_GLOBAL: false,
  CUSTOM_FONT_FAMILY: undefined,
  CUSTOM_FONT_URL: undefined,

  SSRF_ALLOW_PRIVATE_IP_ADDRESS: false,
  SSRF_ALLOW_IP_ADDRESS_LIST: undefined,
  MARKET_BASE_URL: undefined,
});

export const appEnv = getAppConfig();
