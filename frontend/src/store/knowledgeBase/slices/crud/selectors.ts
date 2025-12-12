import { KnowledgeBaseStoreState } from '@/store/knowledgeBase/initialState';

const activeKnowledgeBaseId = (s: KnowledgeBaseStoreState) => s.activeKnowledgeBaseId;

const getKnowledgeBaseById = (id: string) => (s: KnowledgeBaseStoreState) =>
  s.activeKnowledgeBaseItems[id];

const getKnowledgeBaseNameById = (id: string) => (s: KnowledgeBaseStoreState) =>
  getKnowledgeBaseById(id)(s)?.name;

const isFetchingList = (s: KnowledgeBaseStoreState) => s.isFetchingList;
const knowledgeBaseList = (s: KnowledgeBaseStoreState) => s.knowledgeBaseList;
const isKnowledgeBaseLoading = (s: KnowledgeBaseStoreState) => s.isKnowledgeBaseLoading;

export const knowledgeBaseSelectors = {
  activeKnowledgeBaseId,
  getKnowledgeBaseById,
  isFetchingList,
  isKnowledgeBaseLoading,
  knowledgeBaseList,
};
