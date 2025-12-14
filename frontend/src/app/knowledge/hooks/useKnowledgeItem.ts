import { useEffect } from 'react';

import { useKnowledgeBaseStore } from '@/store/knowledgeBase';

export const useKnowledgeBaseItem = (id: string) => {
  const fetchKnowledgeBaseItem = useKnowledgeBaseStore((s) => s.internal_fetchKnowledgeBaseItem);

  useEffect(() => {
    if (id) {
      fetchKnowledgeBaseItem(id);
    }
  }, [fetchKnowledgeBaseItem, id]);
};
