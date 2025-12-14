'use client';

import { memo } from 'react';

import { useSetFileModalId } from '@/app/knowledge/shared/useFileQueryParam';
import FileManager from '@/features/FileManager';
import { knowledgeBaseSelectors, useKnowledgeBaseStore } from '@/store/knowledgeBase';

import { useKnowledgeBaseItem } from '../../hooks/useKnowledgeItem';


const KnowledgeBaseDetailPage = memo<{ id: string }>(({ id }) => {
  const setFileModalId = useSetFileModalId();

  useKnowledgeBaseItem(id!);
  const item = useKnowledgeBaseStore(knowledgeBaseSelectors.getKnowledgeBaseById(id!));
  const name = item?.name || '';

  if (!id) {
    return <div>Knowledge base ID is required</div>;
  }

  return <FileManager knowledgeBaseId={id} onOpenFile={setFileModalId} title={name} />;
});

KnowledgeBaseDetailPage.displayName = 'KnowledgeBaseDetailPage';

export default KnowledgeBaseDetailPage;
