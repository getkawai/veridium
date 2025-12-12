'use client';

import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

import FilePanel from '@/features/FileSidePanel';

import Menu from '../KnowledgeBaseDetail/menu/Menu';

/**
 * Knowledge Base Settings Page
 * Configuration page for a specific knowledge base
 */
const KnowledgeBaseSettingsPage = memo<{ id: string }>(({ id }) => {
  return (
    <>
      <FilePanel>
        <Menu id={id} />
      </FilePanel>
      <Flexbox align="center" flex={1} justify="center" style={{ overflow: 'hidden' }}>
        {/* TODO: Add settings form components here */}
        <div>Settings page for knowledge base: {id}</div>
      </Flexbox>
    </>
  );
});

KnowledgeBaseSettingsPage.displayName = 'KnowledgeBaseSettingsPage';

export default KnowledgeBaseSettingsPage;
