import { KnowledgeBaseItem } from '@/types/knowledgeBase';

export interface KnowledgeBaseState {
  activeKnowledgeBaseId: string | null;
  activeKnowledgeBaseItems: Record<string, KnowledgeBaseItem>;
  initKnowledgeBaseList: boolean;
  isFetchingList: boolean;
  knowledgeBaseList: KnowledgeBaseItem[];
  knowledgeBaseLoadingIds: string[];
  knowledgeBaseRenamingId?: string | null;
}

export const initialKnowledgeBaseState: KnowledgeBaseState = {
  activeKnowledgeBaseId: null,
  activeKnowledgeBaseItems: {},
  initKnowledgeBaseList: false,
  isFetchingList: false,
  knowledgeBaseList: [],
  knowledgeBaseLoadingIds: [],
};
