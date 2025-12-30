import { Suspense } from 'react';
import SQLiteViewer from './components/SQLiteViewer';

const DesktopSQLiteLayout = () => {
  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Suspense fallback={<div>Loading SQLite Viewer...</div>}>
        <SQLiteViewer />
      </Suspense>
    </div>
  );
};

export default DesktopSQLiteLayout;