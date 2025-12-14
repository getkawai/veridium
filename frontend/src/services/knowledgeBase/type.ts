import { KnowledgeBaseItem } from '@/types';

export interface IKnowledgeBaseService {
  getKnowledgeBaseList(): Promise<KnowledgeBaseItem[]>;
  getKnowledgeBaseById(id: string): Promise<KnowledgeBaseItem | undefined>;
  addFilesToKnowledgeBase(knowledgeBaseId: string, ids: string[]): Promise<void>;
  removeFilesFromKnowledgeBase(knowledgeBaseId: string, ids: string[]): Promise<void>;
}

