import { FileTypeIcon, Icon, Text } from '@lobehub/ui';
import { Dialogs } from '@wailsio/runtime';
import { createStyles, useTheme } from 'antd-style';

import { message } from '@/components/AntdStaticMethods';
import { ArrowUpIcon, PlusIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Center, Flexbox } from 'react-layout-kit';

import { useCreateNewModal } from '@/features/KnowledgeBaseModal';
import { useFileStore } from '@/store/file';
import { ProcessFileFromPath } from '@@/github.com/kawai-network/veridium/fileprocessorservice';

const ICON_SIZE = 80;

const useStyles = createStyles(({ css, token }) => ({
  actionTitle: css`
    margin-block-start: 12px;
    font-size: 16px;
    color: ${token.colorTextSecondary};
  `,
  card: css`
    cursor: pointer;

    position: relative;

    overflow: hidden;

    width: 200px;
    height: 140px;
    border-radius: ${token.borderRadiusLG}px;

    font-weight: 500;
    text-align: center;

    background: ${token.colorFillTertiary};
    box-shadow: 0 0 0 1px ${token.colorFillTertiary} inset;

    transition: background 0.3s ease-in-out;

    &:hover {
      background: ${token.colorFillSecondary};
    }
  `,
  glow: css`
    position: absolute;
    inset-block-end: -12px;
    inset-inline-end: 0;

    width: 48px;
    height: 48px;

    opacity: 0.5;
    filter: blur(24px);
  `,
  icon: css`
    position: absolute;
    z-index: 1;
    inset-block-end: -24px;
    inset-inline-end: 8px;

    flex: none;
  `,
}));

interface EmptyStatusProps {
  knowledgeBaseId?: string;
  showKnowledgeBase: boolean;
}
const EmptyStatus = ({ showKnowledgeBase, knowledgeBaseId }: EmptyStatusProps) => {
  const { t } = useTranslation('components');
  const theme = useTheme();
  const { styles } = useStyles();



  const { open } = useCreateNewModal();
  const refreshFileList = useFileStore((s) => s.refreshFileList);

  // Helper from ServerMode.tsx
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
        const processedFiles = await Promise.all(
          filePaths.map(async (filePath) => {
            const fileName = filePath.split('/').pop() || filePath.split('\\').pop() || 'file';
            const ext = fileName.split('.').pop()?.toLowerCase() || '';
            const mimeType = getMimeType(ext);

            const result = await ProcessFileFromPath(filePath);

            if (!result) return null;

            return {
              fileId: result.fileId,
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

        // TODO: ServerMode.tsx adds to chat list, here we want to refresh the file manager list
        // Note: files are not linked to Knowledge Base yet because ProcessFileFromPath doesn't support it

        await refreshFileList();

        hideLoading();
        message.success(`Successfully uploaded ${validFiles.length} file(s)`);
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

  return (
    <Center gap={24} height={'100%'} style={{ paddingBottom: 100 }} width={'100%'}>
      <Flexbox justify={'center'} style={{ textAlign: 'center' }}>
        <Text as={'h4'}>{t('FileManager.emptyStatus.title')}</Text>
        <Text type={'secondary'}>{t('FileManager.emptyStatus.or')}</Text>
      </Flexbox>
      <Flexbox gap={12} horizontal>
        {showKnowledgeBase && (
          <Flexbox
            className={styles.card}
            onClick={() => {
              open();
            }}
            padding={16}
          >
            <span className={styles.actionTitle}>
              {t('FileManager.emptyStatus.actions.knowledgeBase')}
            </span>
            <div className={styles.glow} style={{ background: theme.purple }} />
            <FileTypeIcon
              className={styles.icon}
              color={theme.purple}
              icon={<Icon color={'#fff'} icon={PlusIcon} />}
              size={ICON_SIZE}
              type={'folder'}
            />
          </Flexbox>
        )}
        <Flexbox
          className={styles.card}
          onClick={handleFileUpload}
          padding={16}
        >
          <span className={styles.actionTitle}>{t('FileManager.emptyStatus.actions.file')}</span>
          <div className={styles.glow} style={{ background: theme.gold }} />
          <FileTypeIcon
            className={styles.icon}
            color={theme.gold}
            icon={<Icon color={'#fff'} icon={ArrowUpIcon} />}
            size={ICON_SIZE}
          />
        </Flexbox>

        <Flexbox
          className={styles.card}
          onClick={handleFolderUpload}
          padding={16}
        >
          <span className={styles.actionTitle}>
            {t('FileManager.emptyStatus.actions.folder')}
          </span>
          <div className={styles.glow} style={{ background: theme.geekblue }} />
          <FileTypeIcon
            className={styles.icon}
            color={theme.geekblue}
            icon={<Icon color={'#fff'} icon={ArrowUpIcon} />}
            size={ICON_SIZE}
            type={'folder'}
          />
        </Flexbox>
      </Flexbox>
    </Center>
  );
};

export default EmptyStatus;
