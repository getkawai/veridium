import { KnowledgeBaseModel } from '@/database/models/knowledgeBase';
import { BaseClientService } from '@/services/baseClientService';
import { KnowledgeBaseItem } from '@/types';

import { IKnowledgeBaseService } from './type';
import { DB } from '@/types/database';

export class ClientService extends BaseClientService implements IKnowledgeBaseService {
  private get knowledgeBaseModel(): KnowledgeBaseModel {
    return new KnowledgeBaseModel(DB as any, this.userId);
  }

  getKnowledgeBaseList = async (): Promise<KnowledgeBaseItem[]> => {
    return this.knowledgeBaseModel.query();
  };

  getKnowledgeBaseById = async (id: string): Promise<KnowledgeBaseItem | undefined> => {
    return this.knowledgeBaseModel.findById(id) as Promise<KnowledgeBaseItem | undefined>;
  };

  addFilesToKnowledgeBase = async (knowledgeBaseId: string, ids: string[]): Promise<void> => {
    await this.knowledgeBaseModel.addFilesToKnowledgeBase(knowledgeBaseId, ids);
  };

  removeFilesFromKnowledgeBase = async (knowledgeBaseId: string, ids: string[]): Promise<void> => {
    await this.knowledgeBaseModel.removeFilesFromKnowledgeBase(knowledgeBaseId, ids);
  };
}

