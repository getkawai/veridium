import { Events } from '@wailsio/runtime';
import { memo, useEffect } from 'react';

import { message } from '@/components/AntdStaticMethods';
import { useModelSupportFiles } from '@/hooks/useModelSupportFiles';
import { useModelSupportVision } from '@/hooks/useModelSupportVision';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/slices/chat';
import { useFileStore } from '@/store/file';

import FileItemList from './FileList';

interface WailsDropEvent {
  files: string[];
  elementId?: string;
  classList?: string[];
  x?: number;
  y?: number;
  attributes?: Record<string, string>;
}

const FilePreview = memo(() => {
  const model = useAgentStore(agentSelectors.currentAgentModel);
  const provider = useAgentStore(agentSelectors.currentAgentModelProvider);

  const enabledFiles = useModelSupportFiles(model, provider);
  const supportVision = useModelSupportVision(model, provider);
  const canUpload = enabledFiles || supportVision;

  const [uploadFiles] = useFileStore((s) => [s.uploadChatFiles]);

  useEffect(() => {
    // Listen to Wails drag & drop event
    const unsubscribe = Events.On('files:dropped', async (event: any) => {
      console.log('[Wails D&D] Event received:', event);

      // Extract data from Wails event
      const data: WailsDropEvent = event.data?.[0] || event;
      console.log('[Wails D&D] Files dropped:', data);

      if (!canUpload) {
        console.log('[Wails D&D] Upload not supported for current model');
        message.warning('File upload is not supported for the current model');
        return;
      }

      const filePaths = data.files;
      if (!filePaths || filePaths.length === 0) {
        console.log('[Wails D&D] No files in drop event');
        return;
      }

      console.log('[Wails D&D] Processing', filePaths.length, 'file(s)...');

      try {
        // Convert file paths to File objects
        const files = await Promise.all(
          filePaths.map(async (filePath) => {
            console.log('[Wails D&D] Reading file:', filePath);

            // Extract filename from path
            const fileName = filePath.split('/').pop() || filePath.split('\\').pop() || 'file';

            // Read file using fetch with file:// protocol
            // In Wails, we can access local files directly
            const response = await fetch(`file://${filePath}`);
            const blob = await response.blob();

            // Determine MIME type from extension
            const ext = fileName.split('.').pop()?.toLowerCase();
            const mimeType = getMimeType(ext || '');

            return new File([blob], fileName, { type: mimeType });
          }),
        );

        console.log('[Wails D&D] Files prepared:', files.length);

        // Check if trying to upload non-image files when only vision is supported
        if (!enabledFiles) {
          const nonImageFiles = files.filter((file) => !file.type.startsWith('image'));
          if (nonImageFiles.length > 0) {
            console.warn('[Wails D&D] Non-image files not supported:', nonImageFiles);
            message.warning('Only image files are supported for this model');
            return;
          }
        }

        // Upload files
        console.log('[Wails D&D] Uploading files...');
        await uploadFiles(files);
        console.log('[Wails D&D] Upload completed');
        message.success(`Successfully uploaded ${files.length} file(s)`);
      } catch (error) {
        console.error('[Wails D&D] Error:', error);
        message.error('Failed to upload files');
      }
    });

    return () => {
      unsubscribe();
    };
  }, [canUpload, enabledFiles, uploadFiles]);

  return <FileItemList />;
});

// Helper function to determine MIME type from extension
function getMimeType(ext: string): string {
  const mimeTypes: Record<string, string> = {
    // Images
    jpg: 'image/jpeg',
    jpeg: 'image/jpeg',
    png: 'image/png',
    gif: 'image/gif',
    webp: 'image/webp',
    bmp: 'image/bmp',
    svg: 'image/svg+xml',
    // Documents
    pdf: 'application/pdf',
    doc: 'application/msword',
    docx: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    xls: 'application/vnd.ms-excel',
    xlsx: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    ppt: 'application/vnd.ms-powerpoint',
    pptx: 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    txt: 'text/plain',
    md: 'text/markdown',
    csv: 'text/csv',
    json: 'application/json',
    xml: 'application/xml',
    // Archives
    zip: 'application/zip',
    rar: 'application/x-rar-compressed',
    '7z': 'application/x-7z-compressed',
    tar: 'application/x-tar',
    gz: 'application/gzip',
  };

  return mimeTypes[ext] || 'application/octet-stream';
}

export default FilePreview;
