import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import * as Sentry from "@sentry/react";

// "https://5fb55a3957c99692d702a33e2993cc55@o4510639245426688.ingest.us.sentry.io/4510639291498496"
Sentry.init({
  dsn: "https://b66f862d7567c075a44c697757bb8130@o4510618985758720.ingest.us.sentry.io/4510618990804992",
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
