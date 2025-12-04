import { Events } from '@wailsio/runtime';
import { memo, useEffect } from 'react';

import { message } from '@/components/AntdStaticMethods';
import { useModelSupportFiles } from '@/hooks/useModelSupportFiles';
import { useModelSupportVision } from '@/hooks/useModelSupportVision';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/slices/chat';
import { useFileStore } from '@/store/file';

import FileItemList from './FileList';

interface ProcessedFile {
  originalPath: string;
  savedKey: string;
  url: string;
  name: string;
}

interface WailsDropEvent {
  files: ProcessedFile[];
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

  useEffect(() => {
    // Listen to Wails drag & drop event
    const unsubscribe = Events.On('files:dropped', async (event: any) => {
      // Extract data from Wails event
      const data: WailsDropEvent = event.data || event;

      if (!canUpload) {
        message.warning('File upload is not supported for the current model');
        return;
      }

      const processedFiles = data.files;
      
      if (!processedFiles || processedFiles.length === 0) {
        return;
      }

      try {
        // Files are already saved to local storage by backend
        // Just create upload items with local URLs
        const uploadItems = processedFiles.map((fileInfo) => {
          // Determine MIME type from extension
          const ext = fileInfo.name.split('.').pop()?.toLowerCase();
          const mimeType = getMimeType(ext || '');

          // Check if trying to upload non-image files when only vision is supported
          if (!enabledFiles && !mimeType.startsWith('image')) {
            return null;
          }

          return {
            id: fileInfo.name,
            file: { name: fileInfo.name, type: mimeType, size: 0 } as File,
            previewUrl: fileInfo.url,
            base64Url: undefined,
            status: 'success' as const,
          };
        }).filter(Boolean);

        if (uploadItems.length === 0) {
          message.warning('No supported files to upload');
          return;
        }

        // Add files directly to upload list (already saved by backend)
        useFileStore.getState().dispatchChatUploadFileList({ 
          files: uploadItems as any, 
          type: 'addFiles' 
        });
        
        message.success(`Successfully uploaded ${uploadItems.length} file(s)`);
      } catch (error) {
        console.error('[Wails D&D] Error:', error);
        message.error('Failed to upload files');
      }
    });

    return () => {
      unsubscribe();
    };
  }, [canUpload, enabledFiles]);

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
