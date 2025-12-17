import { memo } from 'react';

interface HTMLRendererProps {
  height?: string;
  htmlContent: string;
  width?: string;
}
const HTMLRenderer = memo<HTMLRendererProps>(({ htmlContent, width = '100%', height = '100%' }) => {

  // Use key with hash of htmlContent to force re-render when content changes significantly?
  // Actually, srcdoc is better but we need to ensure it updates.
  // The issue with document.write in StrictMode is double execution on the same window.
  // Switching to srcdoc is cleaner and handles isolation correctly.

  return (
    <iframe
      srcDoc={htmlContent}
      style={{ border: 'none', height, width }}
      title="html-renderer"
    />
  );
}
);

export default HTMLRenderer;
