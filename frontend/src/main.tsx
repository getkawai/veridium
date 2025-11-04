import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

// Polyfill for navigator.mediaDevices in Wails desktop environment
// This prevents crashes in libraries that check for media device availability
if (typeof navigator !== 'undefined' && !navigator.mediaDevices) {
  (navigator as any).mediaDevices = {
    getUserMedia: () => Promise.reject(new Error('Media devices not available in desktop environment')),
    enumerateDevices: () => Promise.resolve([]),
    getSupportedConstraints: () => ({}),
  };
}

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
