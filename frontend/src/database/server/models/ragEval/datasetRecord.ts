import { EvalDatasetRecordRefFile } from '@/types';
import {
  DB,
  currentTimestampMs,
  getNullableString,
  parseJSON,
  toNullJSON,
  toNullString,
  RagEvalDatasetRecord,
} from '@/types/database';

export class EvalDatasetRecordModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: {
    datasetId: string;
    query: string;
    referenceAnswer?: string | null;
    referenceContexts?: string | null;
    metadata?: any;
  }) => {
    const id = crypto.randomUUID();
    const now = currentTimestampMs();

    const result: RagEvalDatasetRecord = await DB.CreateRagEvalDatasetRecord(
      {
        id,
        datasetId: params.datasetId,
        query: params.query,
        referenceAnswer: toNullString(params.referenceAnswer),
        referenceContexts: toNullString(params.referenceContexts),
        metadata: toNullJSON(params.metadata),
        userId: this.userId,
        createdAt: now,
        updatedAt: now,
      }
    );

    return {
      id: result.id,
      datasetId: result.datasetId,
      query: result.query,
      referenceAnswer: getNullableString(result.referenceAnswer),
      referenceContexts: getNullableString(result.referenceContexts),
      metadata: parseJSON(result.metadata),
      userId: result.userId,
      createdAt: new Date(result.createdAt).toISOString(),
      updatedAt: new Date(result.updatedAt).toISOString(),
    };
  };

  batchCreate = async (
    params: Array<{
      datasetId: string;
      query: string;
      referenceAnswer?: string | null;
      referenceContexts?: string | null;
      metadata?: any;
    }>,
  ) => {
    const now = currentTimestampMs();
    
    // Insert records one by one since SQLite batch insert is complex
    const results: RagEvalDatasetRecord[] = [];
    for (const param of params) {
      const id = crypto.randomUUID();
      const result: RagEvalDatasetRecord = await DB.CreateRagEvalDatasetRecord(
        {
          id,
          datasetId: param.datasetId,
          query: param.query,
          referenceAnswer: toNullString(param.referenceAnswer),
          referenceContexts: toNullString(param.referenceContexts),
          metadata: toNullJSON(param.metadata),
          userId: this.userId,
          createdAt: now,
          updatedAt: now,
        }
      );
      results.push(result);
    }

    return results[0]; // Return first result to match Drizzle behavior
  };

  delete = async (id: string) => {
    return DB.DeleteRagEvalDatasetRecord({ id, userId: this.userId });
  };

  query = async (datasetId: string) => {
    const list = await DB.ListRagEvalDatasetRecords({ datasetId, userId: this.userId });
    
    // Get file IDs from reference_files metadata
    const fileList: string[] = list
      .flatMap((item: RagEvalDatasetRecord) => {
        const metadata = parseJSON(item.metadata);
        return metadata?.referenceFiles || [];
      })
      .filter(Boolean);

    // Get file details if there are any
    let fileItems: EvalDatasetRecordRefFile[] = [];
    if (fileList.length > 0) {
      const files = await Promise.all(
        fileList.map(async (fileId) => {
          try {
            const file = await DB.GetFile({ id: fileId, userId: this.userId });
            if (file) {
              return {
                id: file.id,
                name: file.name,
                fileType: file.fileType,
              } as EvalDatasetRecordRefFile;
            }
          } catch (e) {
            console.error('Error fetching file:', e);
          }
          return null;
        }),
      );
      fileItems = files.filter((f): f is EvalDatasetRecordRefFile => f !== null);
    }

    return list.map((item) => {
      const metadata = parseJSON(item.metadata);
      const refFileIds = (metadata?.referenceFiles as string[]) || [];
      
      return {
        ...item,
        referenceContexts: getNullableString(item.referenceContexts),
        referenceAnswer: getNullableString(item.referenceAnswer),
        metadata,
        referenceFiles: refFileIds
          .map((fileId) => fileItems.find((file) => file.id === fileId))
          .filter(Boolean) as EvalDatasetRecordRefFile[],
      };
    });
  };

  findByDatasetId = async (datasetId: string) => {
    const results = await DB.ListRagEvalDatasetRecords({ datasetId, userId: this.userId });
    
    return results.map((row: RagEvalDatasetRecord) => ({
      id: row.id,
      datasetId: row.datasetId,
      query: row.query,
      referenceAnswer: getNullableString(row.referenceAnswer),
      referenceContexts: getNullableString(row.referenceContexts),
      metadata: parseJSON(row.metadata),
      userId: row.userId,
      createdAt: new Date(row.createdAt),
      updatedAt: new Date(row.updatedAt),
    }));
  };

  findById = async (id: string) => {
    const result = await DB.GetRagEvalDatasetRecord({ id, userId: this.userId });
    if (!result) return undefined;

    return {
      id: result.id,
      datasetId: result.datasetId,
      query: result.query,
      referenceAnswer: getNullableString(result.referenceAnswer),
      referenceContexts: getNullableString(result.referenceContexts),
      metadata: parseJSON(result.metadata),
      userId: result.userId,
      createdAt: new Date(result.createdAt),
      updatedAt: new Date(result.updatedAt),
    };
  };

  update = async (
    id: string,
    value: {
      query?: string | null;
      referenceAnswer?: string | null;
      referenceContexts?: string | null;
      metadata?: any;
    },
  ) => {
    const now = currentTimestampMs();

    await DB.UpdateRagEvalDatasetRecord(
      {
        query: value.query || '',
        referenceAnswer: toNullString(value.referenceAnswer),
        referenceContexts: toNullString(value.referenceContexts),
        metadata: toNullJSON(value.metadata),
        updatedAt: now,
        id,
        userId: this.userId,
      }
    );
  };
}

