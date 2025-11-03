import { RAGEvalDataSetItem } from '@/types';
import { DB, RagEvalDataset, currentTimestampMs, getNullableString, toNullString } from '@/types/database';

export class EvalDatasetModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: {
    name: string;
    description?: string | null;
  }): Promise<RAGEvalDataSetItem> => {
    const id = crypto.randomUUID();
    const now = currentTimestampMs();
    
    const result = await DB.CreateRagEvalDataset({
      id,
      name: params.name,
      description: toNullString(params.description),
      userId: this.userId,
      createdAt: now,
      updatedAt: now,
    } as any);

    return {
      id: result.id,
      name: result.name,
      description: getNullableString(result.description),
      createdAt: new Date(result.createdAt).toISOString(),
      updatedAt: new Date(result.updatedAt).toISOString(),
    } as unknown as RAGEvalDataSetItem;
  };

  delete = async (id: string) => {
    return DB.DeleteRagEvalDataset({ id, userId: this.userId } as any);
  };

  query = async (_knowledgeBaseId: string): Promise<RAGEvalDataSetItem[]> => {
    // Note: The schema doesn't have knowledge_base_id field
    // So we query all datasets for this user
    const results = await DB.ListRagEvalDatasets(this.userId);

    return results.map((row: RagEvalDataset) => ({
      id: row.id,
      name: row.name,
      description: getNullableString(row.description),
      createdAt: new Date(row.createdAt),
      updatedAt: new Date(row.updatedAt),
    })) as unknown as RAGEvalDataSetItem[];
  };

  findById = async (id: string) => {
    const result = await DB.GetRagEvalDataset({ id, userId: this.userId } as any);
    if (!result) return undefined;

    return {
      id: result.id,
      name: result.name,
      description: getNullableString(result.description),
      createdAt: new Date(result.createdAt).toISOString(),
      updatedAt: new Date(result.updatedAt).toISOString(),
    };
  };

  update = async (
    id: string,
    value: {
      name?: string | null;
      description?: string | null;
    },
  ) => {
    const now = currentTimestampMs();
    
    await DB.UpdateRagEvalDataset({
      name: value.name || '',
      description: toNullString(value.description || null),
      updatedAt: now,
      id,
      userId: this.userId,
    } as any);
    return {
      id: id,
      name: value.name || '',
      description: value.description || null,
      createdAt: new Date(now).toISOString(),
      updatedAt: new Date(now).toISOString(),
    } as unknown as RAGEvalDataSetItem;
  };
}

