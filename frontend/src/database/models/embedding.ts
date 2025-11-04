import { nanoid } from 'nanoid';

import { NewEmbeddingsItem } from '../schemas';
import {
  DB,
  toNullString,
  getNullableString,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

export class EmbeddingModel {
  private userId: string;
  private logger = createModelLogger('Embedding', 'EmbeddingModel', 'database/models/embedding');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (value: Omit<NewEmbeddingsItem, 'userId'>) => {
    const id = nanoid();

    await DB.CreateEmbeddingsItem({
      id,
      chunkId: toNullString(value.chunkId as any),
      embeddings: value.embeddings as any,
      model: toNullString(value.model as any),
      clientId: toNullString(value.clientId as any),
      userId: toNullString(this.userId as any),
    });

    return id;
  };

  bulkCreate = async (values: Omit<NewEmbeddingsItem, 'userId'>[]) => {
    await Promise.all(
      values.map((value) =>
        DB.BulkCreateEmbeddingsItems({
          id: nanoid(),
          chunkId: toNullString(value.chunkId as any),
          embeddings: value.embeddings as any,
          model: toNullString(value.model as any),
          clientId: toNullString(value.clientId as any),
          userId: toNullString(this.userId as any),
        }),
      ),
    );
  };

  delete = async (id: string) => {
    await DB.DeleteEmbeddingsItem({
      id,
      userId: toNullString(this.userId as any),
    });
  };

  query = async () => {
    const results = await DB.ListEmbeddingsItems(toNullString(this.userId as any));
    return results.map((r) => this.mapEmbedding(r));
  };

  findById = async (id: string) => {
    try {
      const result = await DB.GetEmbeddingsItem({
        id,
        userId: toNullString(this.userId as any),
      });
      return this.mapEmbedding(result);
    } catch {
      return undefined;
    }
  };

  countUsage = async (): Promise<number> => {
    const result = await DB.CountEmbeddingsItems(toNullString(this.userId as any));
    return Number(result) || 0;
  };

  // **************** Helper *************** //

  private mapEmbedding = (emb: any) => {
    return {
      id: emb.id,
      chunkId: emb.chunkId,
      embeddings: emb.embeddings,
      model: getNullableString(emb.model as any),
      clientId: getNullableString(emb.clientId as any),
      userId: emb.userId,
      createdAt: new Date(emb.createdAt),
      updatedAt: new Date(emb.updatedAt),
    };
  };
}
