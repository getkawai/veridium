import { clientDB } from '@/database/client/db';
import { ChunkModel } from '@/database/models/chunk';
import { DocumentModel } from '@/database/models/document';
import { MessageModel } from '@/database/models/message';
import { BaseClientService } from '@/services/baseClientService';
import { SemanticSearchSchemaType, SemanticSearchResult } from '@/types/rag';

import { IRAGService } from './type';

export class ClientService extends BaseClientService implements IRAGService {
  private get chunkModel(): ChunkModel {
    return new ChunkModel(clientDB as any, this.userId);
  }

  private get documentModel(): DocumentModel {
    return new DocumentModel(clientDB as any, this.userId);
  }

  private get messageModel(): MessageModel {
    return new MessageModel(clientDB as any, this.userId);
  }

  parseFileContent = async (id: string, skipExist?: boolean): Promise<any> => {
    // This would typically call backend API for file parsing
    // For now, return document content
    const doc = await this.documentModel.findById(id);
    return doc;
  };

  createParseFileTask = async (id: string, skipExist?: boolean): Promise<any> => {
    // Task creation would go to backend
    // For Wails, we might process directly or queue
    console.warn('createParseFileTask not implemented for Wails - would need backend task queue');
    return { id, status: 'pending' };
  };

  retryParseFile = async (id: string): Promise<any> => {
    console.warn('retryParseFile not implemented for Wails - would need backend task queue');
    return { id, status: 'pending' };
  };

  createEmbeddingChunksTask = async (id: string): Promise<any> => {
    console.warn('createEmbeddingChunksTask not implemented for Wails - would need backend embedding service');
    return { id, status: 'pending' };
  };

  semanticSearch = async (query: string, fileIds?: string[]): Promise<SemanticSearchResult[]> => {
    // Use chunk model's semantic search
    return this.chunkModel.semanticSearch(query, fileIds);
  };

  semanticSearchForChat = async (params: SemanticSearchSchemaType): Promise<SemanticSearchResult[]> => {
    // Use chunk model's semantic search with chat context
    return this.chunkModel.semanticSearch(params.query, params.fileIds);
  };

  deleteMessageRagQuery = async (id: string): Promise<void> => {
    await this.messageModel.deleteMessageQuery(id);
  };
}

