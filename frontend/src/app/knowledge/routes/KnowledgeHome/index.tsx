'use client';

import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import NProgress from '@/components/NProgress';
import PanelTitle from '@/components/PanelTitle';
import FileManager from '@/features/FileManager';
import FilePanel from '@/features/FileSidePanel';
import { FilesTabs } from '@/types/files';

import { useFileCategory } from '../../hooks/useFileCategory';
import FileModalQueryRoute from '../../shared/FileModalQueryRoute';
import { useSetFileModalId } from '../../shared/useFileQueryParam';
import Container from './layout/Container';
import RegisterHotkeys from './layout/RegisterHotkeys';
import FileMenu from './menu/FileMenu';
import KnowledgeBase from './menu/KnowledgeBase';

// Menu content component
const MenuContent = memo(() => {
  const { t } = useTranslation('file');

  return (
    <Flexbox gap={16} height={'100%'}>
      <Flexbox paddingInline={8}>
        <PanelTitle desc={t('desc')} title={t('title')} />
        <FileMenu />
      </Flexbox>
      <KnowledgeBase />
    </Flexbox>
  );
});

MenuContent.displayName = 'MenuContent';

// Main files list component
const FilesListPage = memo(() => {
  const [category] = useFileCategory();
  const setFileModalId = useSetFileModalId();

  return (
    <FileManager
      category={category}
      onOpenFile={setFileModalId}
      title={`${category as FilesTabs}`}
    />
  );
});

FilesListPage.displayName = 'FilesListPage';

// Main Knowledge Home Page
const KnowledgeHomePage = memo(() => {
  return (
    <>
      <NProgress />
      <Flexbox
        height={'100%'}
        horizontal
        style={{ maxWidth: '100%', overflow: 'hidden', position: 'relative' }}
        width={'100%'}
      >
        <FilePanel>
          <MenuContent />
        </FilePanel>
        <Container>
          <FilesListPage />
        </Container>
      </Flexbox>
      <RegisterHotkeys />
      <FileModalQueryRoute />
    </>
  );
});

KnowledgeHomePage.displayName = 'KnowledgeHomePage';

export default KnowledgeHomePage;
