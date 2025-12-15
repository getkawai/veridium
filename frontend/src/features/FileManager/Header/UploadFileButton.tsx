'use client';

import { Button, Dropdown, Icon, MenuProps } from '@lobehub/ui';
import { css, cx } from 'antd-style';
import { FileUp, FolderUp, UploadIcon } from 'lucide-react';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Dialogs } from '@wailsio/runtime';

import { message } from '@/components/AntdStaticMethods';
import { AddFileToKnowledgeBase } from '@@/github.com/kawai-network/veridium/internal/services/knowledgebaseservice';
import { ProcessFileFromPath } from '@@/github.com/kawai-network/veridium/fileprocessorservice';

import DragUpload from '@/components/DragUpload';
import { useFileStore } from '@/store/file';

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

  const pushDockFileList = useFileStore((s) => s.pushDockFileList);
  const refreshFileList = useFileStore((s) => s.refreshFileList);

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
    <>
      <Dropdown menu={{ items }} placement="bottomRight">
        <Button icon={UploadIcon}>{t('header.uploadButton')}</Button>
      </Dropdown>
      <DragUpload
        enabledFiles
        onUploadFiles={(files) => pushDockFileList(files, knowledgeBaseId)}
      />
    </>
  );
};

export default UploadFileButton;
