import {
  AsyncTaskError,
  AsyncTaskStatus,
  FileSource,
  Generation,
  ImageGenerationAsset,
} from '@/types';
import debug from 'debug';
import { nanoid } from 'nanoid';

// import { FileService } from '@/server/services/file';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  toNullInt,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';

// Local type definitions since schemas are missing
export interface NewFile {
  fileType: string;
  fileHash?: string;
  name: string;
  size?: number;
  url: string;
  metadata?: any;
  knowledgeBaseId?: string;
}

export interface GenerationItem {
  id: string;
  userId: string;
  generationBatchId: string;
  asyncTaskId?: string;
  fileId?: string;
  seed?: number;
  asset?: any;
  createdAt: Date;
  updatedAt: Date;
}

export interface NewGeneration {
  generationBatchId: string;
  asyncTaskId?: string;
  fileId?: string;
  seed?: number;
  asset?: any;
  userId?: string;
}

export interface GenerationWithAsyncTask extends GenerationItem {
  asyncTask?: {
    id: string;
    status: AsyncTaskStatus;
    error?: any;
  };
}

// Create debug logger
const log = debug('lobe-image:generation-model');

export class GenerationModel {
  private userId: string;
  private fileService: FileService;

  constructor(_db: any, userId: string) {
    this.userId = userId;
    this.fileService = new FileService(_db, userId);
  }

  async create(value: Omit<NewGeneration, 'userId'>): Promise<GenerationItem> {
    log('Creating generation: %O', {
      generationBatchId: value.generationBatchId,
      userId: this.userId,
    });

    const now = currentTimestampMs();
    const result = await DB.CreateGeneration({
      id: nanoid(),
      userId: this.userId,
      generationBatchId: value.generationBatchId,
      asyncTaskId: toNullString(value.asyncTaskId as any),
      fileId: toNullString(value.fileId as any),
      seed: toNullInt(value.seed as any),
      asset: toNullJSON(value.asset),
      createdAt: now,
      updatedAt: now,
    });

    log('Generation created successfully: %s', result.id);
    return this.mapGeneration(result);
  }

  async findById(id: string): Promise<GenerationItem | undefined> {
    log('Finding generation by ID: %s for user: %s', id, this.userId);

    try {
      const result = await DB.GetGeneration({
        id,
        userId: this.userId,
      });

      log('Generation %s: found', id);
      return this.mapGeneration(result);
    } catch {
      log('Generation %s: not found', id);
      return undefined;
    }
  }

  async findByIdWithAsyncTask(id: string): Promise<GenerationWithAsyncTask | undefined> {
    log('Finding generation with async task by ID: %s for user: %s', id, this.userId);

    try {
      const result = await DB.GetGenerationWithAsyncTask({
        id,
        userId: this.userId,
      });

      log('Generation %s: found', id);

      // Map the joined result to GenerationWithAsyncTask
      return {
        ...this.mapGeneration(result),
        asyncTask: result.asyncTaskStatus ? {
          id: getNullableString(result.asyncTaskId as any) || '',
          status: result.asyncTaskStatus,
          error: parseNullableJSON(result.asyncTaskError as any),
        } : undefined,
      } as GenerationWithAsyncTask;
    } catch {
      log('Generation %s: not found', id);
      return undefined;
    }
  }

  async update(id: string, value: Partial<NewGeneration>, _trx?: any) {
    log('Updating generation: %s with values: %O', id, {
      asyncTaskId: value.asyncTaskId,
      hasAsset: !!value.asset,
    });

    // Note: No transaction support in Wails!
    // The trx parameter is ignored

    const now = currentTimestampMs();
    const result = await DB.UpdateGeneration({
      id,
      userId: this.userId,
      asyncTaskId: toNullString(value.asyncTaskId as any),
      fileId: toNullString(value.fileId as any),
      asset: toNullJSON(value.asset),
      updatedAt: now,
    });

    log('Generation %s updated successfully', id);
    return this.mapGeneration(result);
  }

  async createAssetAndFile(
    id: string,
    asset: ImageGenerationAsset,
    file: Omit<NewFile, 'id' | 'userId'>,
  ) {
    log('Creating generation asset and file: %s', id);

    // Note: No transaction support in Wails!
    // This is a potential data consistency issue

    // Create file first
    const fileId = nanoid();
    const now = currentTimestampMs();
    const fileHash = file.fileHash || `gen-${fileId}`; // Fallback hash

    // Check if global file exists (unlikely for new generation but good practice)
    const existingGlobalFile = await DB.GetGlobalFile(fileHash);
    const isExist = !!existingGlobalFile;

    await DBService.CreateFileWithLinks({
      File: {
        id: fileId,
        userId: this.userId,
        fileType: toNullString(file.fileType) as any,
        fileHash: toNullString(fileHash) as any,
        name: toNullString(file.name) as any,
        size: file.size || 0,
        url: toNullString(file.url) as any,
        source: toNullString(FileSource.ImageGeneration) as any,
        clientId: toNullString(this.userId) as any,
        metadata: toNullJSON(file.metadata) as any,
        chunkTaskId: toNullString('') as any,
        embeddingTaskId: toNullString('') as any,
        createdAt: now,
        updatedAt: now,
      },
      GlobalFile: !isExist ? {
        hashId: toNullString(fileHash) as any,
        fileType: toNullString(file.fileType) as any,
        size: file.size || 0,
        url: toNullString(file.url) as any,
        metadata: toNullJSON(file.metadata) as any,
        creator: toNullString(this.userId) as any,
        createdAt: now,
      } : null,
      KnowledgeBase: null,
    });

    const newFile = {
      ...file,
      id: fileId,
      userId: this.userId,
      createdAt: now,
      updatedAt: now,
    };

    // Update generation with asset and fileId
    await this.update(id, {
      asset,
      fileId: newFile.id,
    });

    log('Generation %s updated with asset and file %s successfully', id, newFile.id);

    return {
      file: newFile,
    };
  }

  async delete(id: string, _trx?: any) {
    log('Deleting generation: %s for user: %s', id, this.userId);

    // Note: No transaction support in Wails!
    // The trx parameter is ignored

    await DB.DeleteGeneration({
      id,
      userId: this.userId,
    });

    log('Generation %s deleted successfully', id);

    // Note: Drizzle returns the deleted item, but Wails DELETE doesn't return anything
    // We need to fetch it first if we need to return it
    return undefined;
  }

  /**
   * Find generation by ID and transform it to frontend type
   * This method uses findByIdWithAsyncTask and applies transformation
   */
  async findByIdAndTransform(id: string): Promise<Generation | null> {
    log('Finding and transforming generation: %s', id);

    const generation = await this.findByIdWithAsyncTask(id);
    if (!generation) {
      log('Generation %s not found', id);
      return null;
    }

    return await this.transformGeneration(generation);
  }

  /**
   * Transform a GenerationItem (database type) to Generation (frontend type)
   * This method processes asset URLs and async task information
   */
  async transformGeneration(generation: GenerationWithAsyncTask): Promise<Generation> {
    // Process asset URLs if they exist, following the same logic as in generationBatch.ts
    const asset = generation.asset as ImageGenerationAsset | null;
    if (asset && asset.url && asset.thumbnailUrl) {
      const [url, thumbnailUrl] = await Promise.all([
        this.fileService.getFullFileUrl(asset.url),
        this.fileService.getFullFileUrl(asset.thumbnailUrl),
      ]);
      asset.url = url;
      asset.thumbnailUrl = thumbnailUrl;
    }

    // Build the Generation object following the same structure as in generationBatch.ts
    const result: Generation = {
      asset,
      asyncTaskId: generation.asyncTaskId || null,
      createdAt: generation.createdAt,
      id: generation.id,
      seed: generation.seed,
      task: {
        error: generation.asyncTask?.error
          ? (generation.asyncTask.error as unknown as AsyncTaskError)
          : undefined,
        id: generation.asyncTaskId || '',
        status: (generation.asyncTask?.status as AsyncTaskStatus) || 'pending',
      },
    };
    return result;
  }

  // **************** Helper *************** //

  private mapGeneration = (gen: any): GenerationItem => {
    return {
      id: gen.id,
      userId: gen.userId,
      generationBatchId: gen.generationBatchId,
      asyncTaskId: getNullableString(gen.asyncTaskId as any),
      fileId: getNullableString(gen.fileId as any),
      seed: gen.seed,
      asset: parseNullableJSON(gen.asset as any),
      createdAt: new Date(gen.createdAt),
      updatedAt: new Date(gen.updatedAt),
    } as GenerationItem;
  };
}
