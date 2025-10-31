/* eslint-disable sort-keys-fix/sort-keys-fix */

export const getAnalyticsConfig = () => ({
  ENABLED_PLAUSIBLE_ANALYTICS: false,
  PLAUSIBLE_SCRIPT_BASE_URL: 'https://dummy-plausible.com',
  PLAUSIBLE_DOMAIN: undefined,

  ENABLED_POSTHOG_ANALYTICS: false,
  POSTHOG_KEY: undefined,
  POSTHOG_HOST: 'https://dummy-posthog.com',
  DEBUG_POSTHOG_ANALYTICS: false,

  ENABLED_UMAMI_ANALYTICS: false,
  UMAMI_WEBSITE_ID: undefined,
  UMAMI_SCRIPT_URL: 'https://dummy-umami.com/script.js',

  ENABLED_CLARITY_ANALYTICS: false,
  CLARITY_PROJECT_ID: undefined,

  ENABLE_VERCEL_ANALYTICS: false,
  DEBUG_VERCEL_ANALYTICS: false,

  ENABLE_GOOGLE_ANALYTICS: false,
  GOOGLE_ANALYTICS_MEASUREMENT_ID: undefined,

  REACT_SCAN_MONITOR_API_KEY: undefined,
});

export const analyticsEnv = getAnalyticsConfig();
