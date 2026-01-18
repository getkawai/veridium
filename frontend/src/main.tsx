import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import * as Sentry from "@sentry/react";

const SENTRY_DSNS = [
  "https://5fb55a3957c99692d702a33e2993cc55@o4510639245426688.ingest.us.sentry.io/4510639291498496",
  "https://f73dd13f253093e990baf69b9c652b76@o4510675714703360.ingest.us.sentry.io/4510675718832128",
  "https://b66f862d7567c075a44c697757bb8130@o4510618985758720.ingest.us.sentry.io/4510618990804992",
];

const pickRandomSentryDsn = () => {
  if (typeof window === 'undefined') return SENTRY_DSNS[0];
  const key = 'sentry.dsn';
  const cached = window.localStorage.getItem(key);
  if (cached && SENTRY_DSNS.includes(cached)) return cached;
  const randomIndex = Math.floor(Math.random() * SENTRY_DSNS.length);
  const selected = SENTRY_DSNS[randomIndex];
  window.localStorage.setItem(key, selected);
  return selected;
};

Sentry.init({
  dsn: pickRandomSentryDsn(),
  sendDefaultPii: true,
  integrations: [
    Sentry.captureConsoleIntegration({ levels: ['error'] }),
  ],
});

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
