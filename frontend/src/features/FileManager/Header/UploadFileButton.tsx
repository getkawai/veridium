'use client';

import { Button, Dropdown, Icon, MenuProps } from '@lobehub/ui';
import { Events } from '@wailsio/runtime';
import { css, cx } from 'antd-style';
import { FileUp, FolderUp, UploadIcon } from 'lucide-react';
import { useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Dialogs } from '@wailsio/runtime';

import { message } from '@/components/AntdStaticMethods';
import { AddFileToKnowledgeBase, LinkFileToKnowledgeBase } from '@@/github.com/kawai-network/veridium/internal/services/knowledgebaseservice';
import { ProcessFileFromPath } from '@@/github.com/kawai-network/veridium/fileprocessorservice';

import { useFileStore } from '@/store/file';

interface WailsDropEvent {
  files: Array<{
    name: string;
    fileId?: string;
    url: string;
    processing: boolean;
  }>;
  elementId?: string;
  x?: number;
  y?: number;
}

const hotArea = css`
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-color: transparent;
  }
`;

const UploadFileButton = ({ knowledgeBaseId }: { knowledgeBaseId?: string }) => {
  const { t } = useTranslation('file');

  const refreshFileList = useFileStore((s) => s.refreshFileList);

  // Handle native Wails Drag & Drop events
  useEffect(() => {
    // Listen for files:dropped event from main.go
    const unsubscribe = Events.On('files:dropped', async (event: any) => {
      // Extract data from Wails event
      const data: WailsDropEvent = event.data || event;
      const processedFiles = data.files;

      console.log('[Native DnD] Files dropped:', processedFiles);

      if (!processedFiles || processedFiles.length === 0) return;

      // If we are in a specific Knowledge Base context (knowledgeBaseId exists)
      if (knowledgeBaseId) {
        const fileIdsToLink = processedFiles
          .map((f) => f.fileId)
          .filter((id): id is string => !!id);

        if (fileIdsToLink.length > 0) {
          try {
            message.loading(`Linking ${fileIdsToLink.length} file(s) to Knowledge Base...`, 1);

            // Link each file to the KB (files are already processed by main.go)
            await Promise.all(
              fileIdsToLink.map((fileId) => LinkFileToKnowledgeBase(knowledgeBaseId, fileId))
            );

            await refreshFileList();
            message.success(`Linked ${fileIdsToLink.length} file(s) to Knowledge Base`);
          } catch (error) {
            console.error('Failed to link files to KB:', error);
            message.error('Failed to link files to Knowledge Base');
          }
        }
      } else {
        // If global context (no KB ID), just refresh list to show new files
        // main.go already processed them into Global/Generic files
        await refreshFileList();
        message.success(`Uploaded ${processedFiles.length} file(s)`);
      }
    });

    return () => {
      unsubscribe();
    };
  }, [knowledgeBaseId, refreshFileList]);

  const handleFileUpload = async () => {
    try {
      const result = await Dialogs.OpenFile({
        CanChooseFiles: true,
        CanChooseDirectories: false,
        AllowsMultipleSelection: true,
        Title: 'Select Files',
      });

      if (!result) return;

      const filePaths = Array.isArray(result) ? result : [result];
      const hideLoading = message.loading('Processing files...', 0);

      try {
        await Promise.all(
          filePaths.map(async (filePath) => {
            if (knowledgeBaseId) {
              // If KB ID is present, use AddFileToKnowledgeBase (Backend: Load -> DB -> RAG -> Link)
              await AddFileToKnowledgeBase(knowledgeBaseId, filePath, {});
            } else {
              // If no KB ID, just process file (Backend: Load -> DB -> RAG)
              await ProcessFileFromPath(filePath);
            }
          })
        );

        await refreshFileList();
        hideLoading();
        message.success(`Successfully uploaded ${filePaths.length} file(s)`);
      } catch (error) {
        hideLoading();
        throw error;
      }
    } catch (error) {
      console.error('File upload error:', error);
      message.error('Failed to upload files');
    }
  };

  const handleFolderUpload = async () => {
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
  };

  const items = useMemo<MenuProps['items']>(
    () => [
      {
        icon: <Icon icon={FileUp} />,
        key: 'upload-file',
        label: (
          <div className={cx(hotArea)} onClick={handleFileUpload}>
            {t('header.actions.uploadFile')}
          </div>
        ),
      },
      {
        icon: <Icon icon={FolderUp} />,
        key: 'upload-folder',
        label: (
          <div className={cx(hotArea)} onClick={handleFolderUpload}>
            {t('header.actions.uploadFolder')}
          </div>
        ),
      },
    ],
    [handleFileUpload, handleFolderUpload, t],
  );

  return (
    <Dropdown menu={{ items }} placement="bottomRight">
      <Button icon={UploadIcon}>{t('header.uploadButton')}</Button>
    </Dropdown>
  );
};

export default UploadFileButton;
