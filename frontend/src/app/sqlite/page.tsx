import { Suspense } from 'react';
import SQLiteViewer from './components/SQLiteViewer';
import PrivacyBanner from './components/PrivacyBanner';

export default function SQLitePage() {
  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <PrivacyBanner />
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <Suspense fallback={<div role="status" aria-live="polite">Loading SQLite Viewer...</div>}>
          <SQLiteViewer />
        </Suspense>
      </div>
    </div>
  );
}