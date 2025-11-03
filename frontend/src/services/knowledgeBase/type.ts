import { KnowledgeBaseItem } from '@/types';
import { CreateKnowledgeBaseParams } from '@/types/knowledgeBase';

export interface IKnowledgeBaseService {
  createKnowledgeBase(params: CreateKnowledgeBaseParams): Promise<KnowledgeBaseItem>;
  getKnowledgeBaseList(): Promise<KnowledgeBaseItem[]>;
  getKnowledgeBaseById(id: string): Promise<KnowledgeBaseItem | undefined>;
  updateKnowledgeBaseList(id: string, value: any): Promise<KnowledgeBaseItem | undefined>;
  deleteKnowledgeBase(id: string): Promise<void>;
  addFilesToKnowledgeBase(knowledgeBaseId: string, ids: string[]): Promise<void>;
  removeFilesFromKnowledgeBase(knowledgeBaseId: string, ids: string[]): Promise<void>;
}

