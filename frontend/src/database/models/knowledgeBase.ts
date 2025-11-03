import { KnowledgeBaseItem } from '@/types';
import { nanoid } from 'nanoid';

import { NewKnowledgeBase } from '../schemas';
import {
  DB,
  toNullString,
  toNullJSON,
  toNullInt,
  getNullableString,
  parseNullableJSON,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';

export class KnowledgeBaseModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // create

  create = async (params: Omit<NewKnowledgeBase, 'userId'>) => {
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

    return this.mapKnowledgeBase(result);
  };

  addFilesToKnowledgeBase = async (id: string, fileIds: string[]) => {
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
  };

  // delete
  delete = async (id: string) => {
    await DB.DeleteKnowledgeBase({
      id,
      userId: this.userId,
    });
  };

  deleteAll = async () => {
    await DB.DeleteAllKnowledgeBases(this.userId);
  };

  removeFilesFromKnowledgeBase = async (knowledgeBaseId: string, ids: string[]) => {
    await Promise.all(
      ids.map((fileId) =>
        DB.BatchUnlinkKnowledgeBaseFromFiles({
          knowledgeBaseId,
          fileId,
        }),
      ),
    );
  };

  // query
  query = async () => {
    const results = await DB.ListKnowledgeBases(this.userId);
    return results.map((r) => this.mapKnowledgeBase(r)) as KnowledgeBaseItem[];
  };

  findById = async (id: string) => {
    try {
      const result = await DB.GetKnowledgeBase({
        id,
        userId: this.userId,
      });
      return this.mapKnowledgeBase(result);
    } catch {
      return undefined;
    }
  };

  // update
  update = async (id: string, value: Partial<KnowledgeBaseItem>) => {
    const now = currentTimestampMs();

    await DB.UpdateKnowledgeBase({
      id,
      userId: this.userId,
      name: toNullString(value.name as any),
      description: toNullString(value.description as any),
      avatar: toNullString(value.avatar as any),
      settings: toNullJSON(value.settings) as any,
      updatedAt: now,
    });
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

