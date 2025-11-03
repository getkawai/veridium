import { clientDB } from '@/database/client/db';
import { KnowledgeBaseModel } from '@/database/models/knowledgeBase';
import { BaseClientService } from '@/services/baseClientService';
import { CreateKnowledgeBaseParams, KnowledgeBaseItem } from '@/types';

import { IKnowledgeBaseService } from './type';

export class ClientService extends BaseClientService implements IKnowledgeBaseService {
  private get knowledgeBaseModel(): KnowledgeBaseModel {
    return new KnowledgeBaseModel(clientDB as any, this.userId);
  }

  createKnowledgeBase = async (params: CreateKnowledgeBaseParams): Promise<KnowledgeBaseItem> => {
    return this.knowledgeBaseModel.create(params) as Promise<KnowledgeBaseItem>;
  };

  getKnowledgeBaseList = async (): Promise<KnowledgeBaseItem[]> => {
    return this.knowledgeBaseModel.query();
  };

  getKnowledgeBaseById = async (id: string): Promise<KnowledgeBaseItem | undefined> => {
    return this.knowledgeBaseModel.findById(id) as Promise<KnowledgeBaseItem | undefined>;
  };

  updateKnowledgeBaseList = async (id: string, value: any): Promise<KnowledgeBaseItem | undefined> => {
    await this.knowledgeBaseModel.update(id, value);
    return this.knowledgeBaseModel.findById(id) as Promise<KnowledgeBaseItem | undefined>;
  };

  deleteKnowledgeBase = async (id: string): Promise<void> => {
    await this.knowledgeBaseModel.delete(id);
  };

  addFilesToKnowledgeBase = async (knowledgeBaseId: string, ids: string[]): Promise<void> => {
    await this.knowledgeBaseModel.addFilesToKnowledgeBase(knowledgeBaseId, ids);
  };

  removeFilesFromKnowledgeBase = async (knowledgeBaseId: string, ids: string[]): Promise<void> => {
    await this.knowledgeBaseModel.removeFilesFromKnowledgeBase(knowledgeBaseId, ids);
  };
}

