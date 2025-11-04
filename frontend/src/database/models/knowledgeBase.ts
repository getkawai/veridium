import { KnowledgeBaseItem, NewKnowledgeBase } from '@/types';
import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  getNullableString,
  parseNullableJSON,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';

export class KnowledgeBaseModel {
  private userId: string;
  private logger = createModelLogger('KnowledgeBase', 'KnowledgeBaseModel', 'database/models/knowledgeBase');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  /**
   * Show error notification to user
   */
  private async showErrorNotification(title: string, message: string) {
    try {
      await NotificationService.SendNotification(
        new NotificationOptions({
          id: `knowledge-base-error-${Date.now()}`,
          title: `Knowledge Base Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  // create

  create = async (params: Omit<NewKnowledgeBase, 'userId'>) => {
    await this.logger.methodEntry('create', { name: params.name, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      const result = await DB.CreateKnowledgeBase({
        id: nanoid(),
        name: params.name,
        description: toNullString(params.description as any),
        avatar: toNullString(params.avatar as any),
        type: toNullString(params.type as any),
        userId: this.userId,
        clientId: toNullString(params.clientId as any),
        isPublic: boolToInt(params.isPublic || false),
        settings: toNullJSON(params.settings) as any,
        createdAt: now,
        updatedAt: now,
      });

      await this.logger.methodExit('create', { id: result.id });
      return this.mapKnowledgeBase(result);
    } catch (error) {
      await this.logger.error('Failed to create knowledge base', { error, params });
      await this.showErrorNotification(
        'Create Failed',
        `Failed to create knowledge base "${params.name}". Please try again.`
      );
      throw error;
    }
  };

  addFilesToKnowledgeBase = async (id: string, fileIds: string[]) => {
    await this.logger.methodEntry('addFilesToKnowledgeBase', { id, count: fileIds.length });
    
    try {
      const now = currentTimestampMs();

      await Promise.all(
        fileIds.map((fileId) =>
          DB.BatchLinkKnowledgeBaseToFiles({
            knowledgeBaseId: id,
            fileId,
            userId: this.userId,
            createdAt: now,
          }),
        ),
      );
      
      await this.logger.methodExit('addFilesToKnowledgeBase', { id, count: fileIds.length });
    } catch (error) {
      await this.logger.error('Failed to add files to knowledge base', { error, id, fileIds });
      await this.showErrorNotification(
        'Add Files Failed',
        `Failed to add ${fileIds.length} files to knowledge base. Please try again.`
      );
      throw error;
    }
  };

  // delete
  delete = async (id: string) => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    try {
      await DB.DeleteKnowledgeBase({
        id,
        userId: this.userId,
      });
      
      await this.logger.methodExit('delete', { id });
    } catch (error) {
      await this.logger.error('Failed to delete knowledge base', { error, id });
      await this.showErrorNotification(
        'Delete Failed',
        `Failed to delete knowledge base. Please try again.`
      );
      throw error;
    }
  };

  deleteAll = async () => {
    try {
      await DB.DeleteAllKnowledgeBases(this.userId);
    } catch (error) {
      await this.logger.error('Failed to delete all knowledge bases', { error });
      await this.showErrorNotification(
        'Delete All Failed',
        `Failed to delete all knowledge bases. Please try again.`
      );
      throw error;
    }
  };

  removeFilesFromKnowledgeBase = async (knowledgeBaseId: string, ids: string[]) => {
    await this.logger.methodEntry('removeFilesFromKnowledgeBase', { knowledgeBaseId, count: ids.length });
    
    try {
      await Promise.all(
        ids.map((fileId) =>
          DB.BatchUnlinkKnowledgeBaseFromFiles({
            knowledgeBaseId,
            fileId,
          }),
        ),
      );
      
      await this.logger.methodExit('removeFilesFromKnowledgeBase', { knowledgeBaseId, count: ids.length });
    } catch (error) {
      await this.logger.error('Failed to remove files from knowledge base', { error, knowledgeBaseId, ids });
      await this.showErrorNotification(
        'Remove Files Failed',
        `Failed to remove ${ids.length} files from knowledge base. Please try again.`
      );
      throw error;
    }
  };

  // query
  query = async () => {
    try {
      const results = await DB.ListKnowledgeBases(this.userId);
      return results.map((r) => this.mapKnowledgeBase(r)) as KnowledgeBaseItem[];
    } catch (error) {
      await this.logger.error('Failed to query knowledge bases', { error });
      throw error;
    }
  };

  findById = async (id: string) => {
    try {
      const result = await DB.GetKnowledgeBase({
        id,
        userId: this.userId,
      });
      return this.mapKnowledgeBase(result);
    } catch (error) {
      await this.logger.warn('Knowledge base not found', { id, error });
      return undefined;
    }
  };

  // update
  update = async (id: string, value: Partial<KnowledgeBaseItem>) => {
    await this.logger.methodEntry('update', { id, value, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      await DB.UpdateKnowledgeBase({
        id,
        userId: this.userId,
        name: value.name || '',
        description: toNullString(value.description as any),
        avatar: toNullString(value.avatar as any),
        settings: toNullJSON(value.settings) as any,
        updatedAt: now,
      });
      
      await this.logger.methodExit('update', { id });
    } catch (error) {
      await this.logger.error('Failed to update knowledge base', { error, id, value });
      await this.showErrorNotification(
        'Update Failed',
        `Failed to update knowledge base "${value.name || ''}". Please try again.`
      );
      throw error;
    }
  };

  static findById = async (_db: any, id: string) => {
    try {
      const result = await DB.GetKnowledgeBase({
        id,
        userId: '', // Static method doesn't have userId context
      });
      return {
        id: result.id,
        name: result.name,
        description: getNullableString(result.description as any),
        avatar: getNullableString(result.avatar as any),
        type: getNullableString(result.type as any),
        userId: result.userId,
        clientId: getNullableString(result.clientId as any),
        isPublic: intToBool(result.isPublic),
        settings: parseNullableJSON(result.settings as any),
        createdAt: new Date(result.createdAt),
        updatedAt: new Date(result.updatedAt),
      };
    } catch {
      return undefined;
    }
  };

  // **************** Helper *************** //

  private mapKnowledgeBase = (kb: any) => {
    return {
      id: kb.id,
      name: kb.name,
      description: getNullableString(kb.description as any),
      avatar: getNullableString(kb.avatar as any),
      type: getNullableString(kb.type as any),
      userId: kb.userId,
      clientId: getNullableString(kb.clientId as any),
      isPublic: intToBool(kb.isPublic),
      settings: parseNullableJSON(kb.settings as any),
      createdAt: new Date(kb.createdAt),
      updatedAt: new Date(kb.updatedAt),
    };
  };
}

