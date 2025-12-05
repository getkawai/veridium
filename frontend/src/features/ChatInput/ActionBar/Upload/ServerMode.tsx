import { MenuProps, Tooltip } from '@lobehub/ui';
import { css, cx } from 'antd-style';
import { FileUp, FolderUp, ImageUp, Paperclip } from 'lucide-react';
import { memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Dialogs } from '@wailsio/runtime';

import { message } from '@/components/AntdStaticMethods';
import { useModelSupportVision } from '@/hooks/useModelSupportVision';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/selectors';
import { useFileStore } from '@/store/file';
import { ProcessFileFromPath } from '@@/github.com/kawai-network/veridium/fileprocessorservice';

import Action from '../components/Action';
import { getUserId } from '@/store/user/helpers';

const hotArea = css`
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-color: transparent;
  }
`;

// Helper to get MIME type from extension (for UI display only)
const getMimeType = (ext: string): string => {
  const mimeTypes: Record<string, string> = {
    png: 'image/png',
    jpg: 'image/jpeg',
    jpeg: 'image/jpeg',
    gif: 'image/gif',
    webp: 'image/webp',
    svg: 'image/svg+xml',
    mp4: 'video/mp4',
    webm: 'video/webm',
    pdf: 'application/pdf',
    txt: 'text/plain',
    json: 'application/json',
    xml: 'application/xml',
  };
  return mimeTypes[ext] || 'application/octet-stream';
};

const FileUpload = memo(() => {
  const { t } = useTranslation('chat');

  const model = useAgentStore(agentSelectors.currentAgentModel);
  const provider = useAgentStore(agentSelectors.currentAgentModelProvider);

  const canUploadImage = useModelSupportVision(model, provider);

  // Handle native file dialog for images
  const handleImageUpload = useCallback(async () => {
    try {
      const result = await Dialogs.OpenFile({
        CanChooseFiles: true,
        CanChooseDirectories: false,
        AllowsMultipleSelection: true,
        Filters: [
          {
            DisplayName: 'Images',
            Pattern: '*.png;*.jpg;*.jpeg;*.gif;*.webp;*.svg',
          },
        ],
        Title: 'Select Images',
      });

      if (!result) return;

      const filePaths = Array.isArray(result) ? result : [result];

      // Process files via backend (copy + parse + RAG in one call)
      const processedFiles = await Promise.all(
        filePaths.map(async (filePath) => {
          const userId = getUserId();

          // Process file from absolute path (copies to local storage and processes)
          const result = await ProcessFileFromPath(filePath, userId);
          if (!result) return null;

          const ext = result.filename.split('.').pop()?.toLowerCase() || '';
          const mimeType = getMimeType(ext); // For UI display only

          return {
            fileId: result.fileId,  // Include file ID from backend
            name: result.filename,
            type: mimeType,
            url: result.relativeUrl,
          };
        }),
      );

      const validFiles = processedFiles.filter(Boolean);

      if (validFiles.length === 0) {
        message.warning('No valid files to upload');
        return;
      }

      // Create upload items and add directly to upload list (files already saved by backend)
      const uploadItems = validFiles.map((info) => ({
        id: info!.fileId,  // Use actual file ID from backend, not filename
        file: { name: info!.name, type: info!.type, size: 0 } as File,
        previewUrl: info!.url,
        base64Url: undefined,
        status: 'success' as const,
      }));

      useFileStore.getState().dispatchChatUploadFileList({
        files: uploadItems as any,
        type: 'addFiles',
      });

      message.success(`Successfully uploaded ${uploadItems.length} image(s)`);
    } catch (error) {
      console.error('Image upload error:', error);
      message.error('Failed to upload images');
    }
  }, [canUploadImage]);

  // Handle native file dialog for files
  const handleFileUpload = useCallback(async () => {
    try {
      const result = await Dialogs.OpenFile({
        CanChooseFiles: true,
        CanChooseDirectories: false,
        AllowsMultipleSelection: true,
        Title: 'Select Files',
      });

      if (!result) return;

      const filePaths = Array.isArray(result) ? result : [result];

      // Show loading message
      const hideLoading = message.loading('Processing files...', 0);

      try {
        // Process files via backend (copy + parse + RAG in one call)
        const processedFiles = await Promise.all(
          filePaths.map(async (filePath) => {
            const fileName = filePath.split('/').pop() || filePath.split('\\').pop() || 'file';
            const ext = fileName.split('.').pop()?.toLowerCase() || '';
            const mimeType = getMimeType(ext); // For UI display only

            // Skip video if model doesn't support vision
            if (!canUploadImage && mimeType.startsWith('video')) {
              return null;
            }

            const userId = getUserId();

            // Process file from absolute path (copies to local storage and processes)
            const result = await ProcessFileFromPath(filePath, userId);
            if (!result) return null;

            return {
              fileId: result.fileId,  // Include file ID from backend
              name: result.filename,
              type: mimeType,
              url: result.relativeUrl,
            };
          }),
        );

        const validFiles = processedFiles.filter(Boolean);

        if (validFiles.length === 0) {
          hideLoading();
          message.warning('No valid files to upload');
          return;
        }

        // Create upload items and add directly to upload list (files already saved by backend)
        const uploadItems = validFiles.map((info) => ({
          id: info!.fileId,  // Use actual file ID from backend, not filename
          file: { name: info!.name, type: info!.type, size: 0 } as File,
          previewUrl: info!.url,
          base64Url: undefined,
          status: 'success' as const,
        }));

        useFileStore.getState().dispatchChatUploadFileList({
          files: uploadItems as any,
          type: 'addFiles',
        });

        hideLoading();
        message.success(`Successfully uploaded ${uploadItems.length} file(s)`);
      } catch (error) {
        hideLoading();
        throw error;
      }
    } catch (error) {
      console.error('File upload error:', error);
      message.error('Failed to upload files');
    }
  }, [canUploadImage]);

  // Handle native folder dialog
  const handleFolderUpload = useCallback(async () => {
    try {
      const result = await Dialogs.OpenFile({
        CanChooseFiles: false,
        CanChooseDirectories: true,
        AllowsMultipleSelection: false,
        Title: 'Select Folder',
      });

      if (!result) return;

      message.info('Folder upload not yet implemented');
      // TODO: Implement recursive folder scanning in backend
    } catch (error) {
      console.error('Folder upload error:', error);
      message.error('Failed to upload folder');
    }
  }, []);

  const items: MenuProps['items'] = [
    {
      disabled: !canUploadImage,
      icon: ImageUp,
      key: 'upload-image',
      label: canUploadImage ? (
        <div className={cx(hotArea)} onClick={handleImageUpload}>
          {t('upload.action.imageUpload')}
        </div>
      ) : (
        <Tooltip placement={'right'} title={t('upload.action.imageDisabled')}>
          <div className={cx(hotArea)}>{t('upload.action.imageUpload')}</div>
        </Tooltip>
      ),
    },
    {
      icon: FileUp,
      key: 'upload-file',
      label: (
        <div className={cx(hotArea)} onClick={handleFileUpload}>
          {t('upload.action.fileUpload')}
        </div>
      ),
    },
    {
      icon: FolderUp,
      key: 'upload-folder',
      label: (
        <div className={cx(hotArea)} onClick={handleFolderUpload}>
          {t('upload.action.folderUpload')}
        </div>
      ),
    },
  ];

  return (
    <Action
      dropdown={{
        menu: { items },
      }}
      icon={Paperclip}
      showTooltip={false}
      title={t('upload.action.tooltip')}
    />
  );
});

export default FileUpload;
