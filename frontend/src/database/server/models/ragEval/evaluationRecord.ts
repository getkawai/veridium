import { DB, RagEvalEvaluationRecord, currentTimestampMs, getNullableString, parseJSON, toNullJSON, toNullString } from '@/types/database';

export class EvaluationRecordModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: {
    evaluationId: string;
    datasetRecordId: string;
    retrievedContexts?: string | null;
    generatedAnswer?: string | null;
    metrics?: any;
  }) => {
    const id = crypto.randomUUID();
    const now = currentTimestampMs();

    const result: RagEvalEvaluationRecord = await DB.CreateRagEvalEvaluationRecord({
      id,
      evaluationId: params.evaluationId,
      datasetRecordId: params.datasetRecordId,
      retrievedContexts: toNullString(params.retrievedContexts || null),
      generatedAnswer: toNullString(params.generatedAnswer || null),
      metrics: toNullJSON(params.metrics || null),
      userId: this.userId,
      createdAt: now,
      updatedAt: now,
    });

    return {
      id: result.id,
      evaluationId: result.evaluationId,
      datasetRecordId: result.datasetRecordId,
      retrievedContexts: getNullableString(result.retrievedContexts),
      generatedAnswer: getNullableString(result.generatedAnswer),
      metrics: parseJSON(result.metrics),
      userId: result.userId,
      createdAt: new Date(result.createdAt).toISOString(),
      updatedAt: new Date(result.updatedAt).toISOString(),
    };
  };

  batchCreate = async (
    params: Array<{
      evaluationId: string;
      datasetRecordId: string;
      retrievedContexts?: string | null;
      generatedAnswer?: string | null;
      metrics?: any;
    }>,
  ) => {
    const now = currentTimestampMs();
    
    // Insert records one by one
    const results: any[] = [];
    for (const param of params) {
      const id = crypto.randomUUID();
      const result: RagEvalEvaluationRecord = await DB.CreateRagEvalEvaluationRecord({
        id,
        evaluationId: param.evaluationId,
        datasetRecordId: param.datasetRecordId,
        retrievedContexts: toNullString(param.retrievedContexts || null),
        generatedAnswer: toNullString(param.generatedAnswer || null),
        metrics: toNullJSON(param.metrics || null),
        userId: this.userId,
        createdAt: now,
        updatedAt: now,
      } as any);
      results.push(result);
    }

    return results.map((r: any) => ({
      id: r.id,
      evaluationId: r.evaluationId,
      datasetRecordId: r.datasetRecordId,
      retrievedContexts: getNullableString(r.retrievedContexts),
      generatedAnswer: getNullableString(r.generatedAnswer),
      metrics: parseJSON(r.metrics),
      userId: r.userId,
      createdAt: new Date(r.createdAt).toISOString(),
      updatedAt: new Date(r.updatedAt).toISOString(),
    }));
  };

  delete = async (id: string) => {
    return DB.DeleteRagEvalEvaluationRecord({ id, userId: this.userId } as any);
  };

  query = async (evaluationId: string) => {
    const results: RagEvalEvaluationRecord[] = await DB.ListRagEvalEvaluationRecordsByEvaluation({ evaluationId, userId: this.userId } as any);

    return results.map((row: RagEvalEvaluationRecord) => ({
      id: row.id,
      evaluationId: row.evaluationId,
      datasetRecordId: row.datasetRecordId,
      retrievedContexts: getNullableString(row.retrievedContexts),
      generatedAnswer: getNullableString(row.generatedAnswer),
      metrics: parseJSON(row.metrics),
      userId: row.userId,
      createdAt: new Date(row.createdAt).toISOString(),
      updatedAt: new Date(row.updatedAt).toISOString(),
    }));
  };

  findById = async (id: string) => {
    const result: RagEvalEvaluationRecord | undefined = await DB.GetRagEvalEvaluationRecord({ id, userId: this.userId } as any);
    if (!result) return undefined;

    return {
      id: result.id,
      evaluationId: result.evaluationId,
      datasetRecordId: result.datasetRecordId,
      retrievedContexts: getNullableString(result.retrievedContexts),
      generatedAnswer: getNullableString(result.generatedAnswer),
      metrics: parseJSON(result.metrics),
      userId: result.userId,
      createdAt: new Date(result.createdAt).toISOString(),
      updatedAt: new Date(result.updatedAt).toISOString(),
    };
  };

  findByEvaluationId = async (evaluationId: string) => {
    const results: RagEvalEvaluationRecord[] = await DB.ListRagEvalEvaluationRecordsByEvaluation({ evaluationId, userId: this.userId } as any);

    return results.map((row: RagEvalEvaluationRecord) => ({
      id: row.id,
      evaluationId: row.evaluationId,
      datasetRecordId: row.datasetRecordId,
      retrievedContexts: getNullableString(row.retrievedContexts),
      generatedAnswer: getNullableString(row.generatedAnswer),
      metrics: parseJSON(row.metrics),
      userId: row.userId,
      createdAt: new Date(row.createdAt).toISOString(),
      updatedAt: new Date(row.updatedAt).toISOString(),
    }));
  };

  update = async (
    id: string,
    value: {
      retrievedContexts?: string | null;
      generatedAnswer?: string | null;
      metrics?: any;
    },
  ) => {
    const now = currentTimestampMs();

    await DB.UpdateRagEvalEvaluationRecord({
      retrievedContexts: toNullString(value.retrievedContexts || null),
      generatedAnswer: toNullString(value.generatedAnswer || null),
      metrics: toNullJSON(value.metrics || null),
      updatedAt: now,
      id,
      userId: this.userId,
    } as any);
  };
}

