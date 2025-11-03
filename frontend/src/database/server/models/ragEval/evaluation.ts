import { RagEvalEvaluation } from '@/types/database';
import { RAGEvalEvaluationItem } from '@/types/eval';
import { DB, currentTimestampMs, parseJSON, toNullJSON, toNullString } from '@/types/database';

export class EvalEvaluationModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: {
    name: string;
    datasetId: string;
    config?: any;
    status: string;
  }) => {
    const id = crypto.randomUUID();
    const now = currentTimestampMs();

    const result = await DB.CreateRagEvalEvaluation(
      {
        id,
        name: params.name,
        datasetId: params.datasetId,
        config: toNullJSON(params.config),
        status: params.status,
        userId: this.userId,
        createdAt: now,
        updatedAt: now,
      });

    return {
      id: result.id,
      name: result.name,
      datasetId: result.datasetId,
      config: parseJSON(result.config),
      status: result.status,
      userId: result.userId,
      createdAt: new Date(result.createdAt).toISOString(),
      updatedAt: new Date(result.updatedAt).toISOString(),
    };
  };

  delete = async (id: string) => {
    return DB.DeleteRagEvalEvaluation({ id, userId: this.userId });
  };

  queryByKnowledgeBaseId = async (_knowledgeBaseId: string) => {
    // Note: Schema doesn't have knowledge_base_id, so we get all evaluations
    // In a real implementation, you'd need to add this field to the schema
    
    // Get all datasets for this user
    const datasets = await DB.ListRagEvalDatasets(this.userId);
    const datasetIds = datasets.map(d => d.id);
    
    if (datasetIds.length === 0) {
      return [];
    }

    // Get evaluations for each dataset
    const allEvaluations: RagEvalEvaluation[] = [];
    for (const datasetId of datasetIds as string[]) {
      const evals = await DB.ListRagEvalEvaluationsByDataset({ datasetId, userId: this.userId });
      allEvaluations.push(...(evals as RagEvalEvaluation[]));
    }

    // Get dataset info for each evaluation
    const evaluationsWithDataset = allEvaluations.map(async (evaluation: RagEvalEvaluation) => {
      const dataset = datasets.find(d => d.id === evaluation.datasetId);
      
      return {
        id: evaluation.id,
        name: evaluation.name,
        status: evaluation.status,
        evalRecordsUrl: null, // Schema doesn't have this field
        dataset: dataset ? {
          id: dataset.id,
          name: dataset.name,
        } : null,
        createdAt: new Date(evaluation.createdAt).toISOString(),
        updatedAt: new Date(evaluation.updatedAt).toISOString(),
      };
    });

    const evaluations = await Promise.all(evaluationsWithDataset);
    const evaluationIds = evaluations.map(e => e.id);

    // Get record stats for each evaluation
    const recordStats: Array<{ evaluationId: string; total: number; success: number }> = [];
    
    for (const evalId of evaluationIds) {
      const records = await DB.ListRagEvalEvaluationRecordsByEvaluation({ evaluationId: evalId, userId: this.userId });
      
      const total = records.length;
      const success = records.filter((record) => {
        const metrics = parseJSON(record.metrics);
        return metrics?.status === 'success';
      }).length;
      
      recordStats.push({
        evaluationId: evalId,
        total,
        success,
      });
    }

    return evaluations.map((evaluation) => {
      const stats = recordStats.find((stat) => stat.evaluationId === evaluation.id);

      return {
        ...(evaluation as unknown as RAGEvalEvaluationItem),
        recordsStats: stats
          ? { success: stats.success, total: stats.total }
          : { success: 0, total: 0 },
      } as unknown as RAGEvalEvaluationItem;
    });
  };

  findById = async (id: string) => {
    const result = await DB.GetRagEvalEvaluation({ id, userId: this.userId });
    if (!result) return undefined;

    return {
      id: result.id,
      name: result.name,
      datasetId: result.datasetId,
      config: parseJSON(result.config),
      status: result.status,
      userId: result.userId,
      createdAt: new Date(result.createdAt),
      updatedAt: new Date(result.updatedAt),
    };
  };

  update = async (
    id: string,
    value: {
      name?: string | null;
      config?: any;
      status?: string | null;
    },
  ) => {
    const now = currentTimestampMs();

    await DB.UpdateRagEvalEvaluation(
      {
        name: value.name || '',
        config: toNullString(value.config),
        status: value.status || '',
        updatedAt: now,
        id,
        userId: this.userId,
      }
    );
  };
}

