import { Generation, GenerationAsset, GenerationBatch, GenerationConfig } from  '@/types';
import debug from 'debug';
import { nanoid } from 'nanoid';

import { FileService } from '@/server/services/file';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

import { GenerationModel } from './generation';

const log = debug('lobe-image:generation-batch-model');

export class GenerationBatchModel {
  private userId: string;
  private fileService: FileService;
  private generationModel: GenerationModel;
  private logger = createModelLogger('GenerationBatch', 'GenerationBatchModel', 'database/models/generationBatch');

  constructor(_db: any, userId: string) {
    this.userId = userId;
    // Note: FileService still needs db for now, pass dummy
    this.fileService = new FileService(null as any, userId);
    this.generationModel = new GenerationModel(null as any, userId);
  }

  async create(value: any): Promise<any> {
    log('Creating generation batch: %O', {
      topicId: value.generationTopicId,
      userId: this.userId,
    });

    const id = nanoid();
    const now = currentTimestampMs();

    const result: GenerationBatch = await DB.CreateGenerationBatch({
      id,
      userId: this.userId,
      generationTopicId: toNullString(value.generationTopicId),
      provider: toNullString(value.provider),
      model: toNullString(value.model),
      prompt: toNullString(value.prompt),
      width: value.width || 0,
      height: value.height || 0,
      ratio: toNullString(value.ratio),
      config: toNullJSON(value.config),
      createdAt: now,
      updatedAt: now,
    });

    log('Generation batch created successfully: %s', result.id);
    return result;
  }

  async findById(id: string): Promise<any | undefined> {
    log('Finding generation batch by ID: %s for user: %s', id, this.userId);

    const result = await DB.GetGenerationBatch({
      id,
      userId: this.userId,
    });

    log('Generation batch %s: %s', id, result ? 'found' : 'not found');
    return result;
  }

  async findByTopicId(topicId: string): Promise<any[]> {
    log('Finding generation batches by topic ID: %s for user: %s', topicId, this.userId);

    const results = await DB.ListGenerationBatches({
      generationTopicId: toNullString(topicId),
      userId: this.userId,
    });

    log('Found %d generation batches for topic %s', results.length, topicId);
    return results;
  }

  /**
   * OPTIMIZED: Uses JOIN query (3A) to fetch batches with generations in single query
   * Much faster than N+1 approach
   */
  async findByTopicIdWithGenerations(topicId: string): Promise<any[]> {
    log(
      'Finding generation batches with generations for topic ID: %s for user: %s',
      topicId,
      this.userId,
    );

    // OPTIMIZATION 3A: Single query with JOINs
    const results = await DB.ListGenerationBatchesWithGenerations({
      generationTopicId: toNullString(topicId),
      userId: this.userId,
    });

    // Group results by batch_id
    const batchesMap = new Map<string, any>();
    
    for (const row of results) {
      const batchId = row.batchId;
      
      if (!batchesMap.has(batchId)) {
        batchesMap.set(batchId, {
          id: batchId,
          generationTopicId: getNullableString(row.generationTopicId as any),
          provider: getNullableString(row.provider as any),
          model: getNullableString(row.model as any),
          prompt: getNullableString(row.prompt as any),
          width: row.width || 0,
          height: row.height || 0,
          ratio: getNullableString(row.ratio as any),
          config: parseNullableJSON(row.config as any),
          createdAt: row.batchCreatedAt,
          updatedAt: row.batchUpdatedAt,
          userId: this.userId,
          generations: [],
        });
      }

      // Add generation if exists
      if (row.genId) {
        const batch = batchesMap.get(batchId);
        batch.generations.push({
          id: row.genId,
          generationBatchId: batchId,
          asyncTaskId: getNullableString(row.asyncTaskId as any),
          fileId: getNullableString(row.fileId as any),
          seed: row.seed || 0,
          asset: parseNullableJSON(row.asset as any),
          createdAt: row.genCreatedAt,
          updatedAt: row.genUpdatedAt,
          userId: this.userId,
          asyncTask: row.taskId ? {
            id: row.taskId,
            status: getNullableString(row.taskState as any),
            error: parseNullableJSON(row.taskError as any),
          } : null,
        });
      }
    }

    const batches = Array.from(batchesMap.values());
    log('Found %d generation batches with generations for topic %s', batches.length, topicId);
    return batches;
  }

  async queryGenerationBatchesByTopicIdWithGenerations(
    topicId: string,
  ): Promise<(GenerationBatch & { generations: Generation[] })[]> {
    log('Fetching generation batches for topic ID: %s for user: %s', topicId, this.userId);

    const batchesWithGenerations = await this.findByTopicIdWithGenerations(topicId);
    if (batchesWithGenerations.length === 0) {
      log('No batches found for topic: %s', topicId);
      return [];
    }

    // Transform the database result to match our frontend types
    const result: GenerationBatch[] = await Promise.all(
      batchesWithGenerations.map(async (batch) => {
        const [generations, config] = await Promise.all([
          // Transform generations
          Promise.all(
            batch.generations.map((gen: any) => this.generationModel.transformGeneration(gen)),
          ),
          // Transform config
          (async () => {
            const config = batch.config as GenerationConfig;

            // Handle single imageUrl
            if (config.imageUrl) {
              config.imageUrl = await this.fileService.getFullFileUrl(config.imageUrl);
            }

            // Handle imageUrls array
            if (Array.isArray(config.imageUrls)) {
              config.imageUrls = await Promise.all(
                config.imageUrls.map((url) => this.fileService.getFullFileUrl(url)),
              );
            }
            return config;
          })(),
        ]);

        return {
          config,
          createdAt: batch.createdAt,
          generations,
          height: batch.height,
          id: batch.id,
          model: batch.model,
          prompt: batch.prompt,
          provider: batch.provider,
          width: batch.width,
        };
      }),
    );

    log('Feed construction complete for topic: %s, returning %d batches', topicId, result.length);
    return result;
  }

  /**
   * OPTIMIZED: Uses query to fetch assets, then deletes
   * Returns thumbnail URLs for cleanup
   */
  async delete(
    id: string,
  ): Promise<{ deletedBatch: any; thumbnailUrls: string[] } | undefined> {
    log('Deleting generation batch: %s for user: %s', id, this.userId);

    // 1. Get batch to verify ownership
    const batch = await this.findById(id);
    if (!batch) {
      return undefined;
    }

    // 2. Get generation assets for thumbnail URLs
    const assets = await DB.GetGenerationBatchAssets({
      generationBatchId: toNullString(id),
      userId: this.userId,
    });

    // 3. Collect thumbnail URLs
    const thumbnailUrls: string[] = [];
    for (const row of assets) {
      const asset = parseNullableJSON(row.asset as any) as GenerationAsset;
      if (asset?.thumbnailUrl) {
        thumbnailUrls.push(asset.thumbnailUrl);
      }
    }

    // 4. Delete the batch (cascade will handle generations)
    await DB.DeleteGenerationBatch({
      id,
      userId: this.userId,
    });

    log(
      'Generation batch %s deleted successfully with %d thumbnails to clean',
      id,
      thumbnailUrls.length,
    );

    return {
      deletedBatch: batch,
      thumbnailUrls,
    };
  }
}

