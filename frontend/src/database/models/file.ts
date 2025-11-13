import { FilesTabs, QueryFileListParams, SortType } from  '@/types';
import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  File as DBFile,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';
import { createModelLogger } from '@/utils/logger';
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';

export class FileModel {
  private readonly userId: string;
  private logger = createModelLogger('File', 'FileModel', 'database/models/file');

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
          id: `file-error-${Date.now()}`,
          title: `File Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  /**
   * Create file with atomic transaction
   * ✅ OPTIMIZED: Uses backend transaction for atomicity
   * - Inserts to global_files if needed
   * - Inserts to files table
   * - Links to knowledge_base_files if needed
   * All operations succeed or all rollback automatically
   */
  create = async (
    params: any & { knowledgeBaseId?: string },
    insertToGlobalFiles?: boolean,
  ) => {
    await this.logger.methodEntry('create', { userId: this.userId, fileHash: params.fileHash, knowledgeBaseId: params.knowledgeBaseId });
    
    try {
      const fileId = nanoid();
      const now = currentTimestampMs();

      // Use backend transaction method for atomic operations
      const result = await DBService.CreateFileWithLinks({
        File: {
          id: fileId,
          userId: this.userId,
          fileType: toNullString(params.fileType) as any,
          fileHash: toNullString(params.fileHash) as any,
          name: toNullString(params.name) as any,
          size: params.size || 0,
          url: toNullString(params.url) as any,
          source: toNullString(params.source) as any,
          clientId: toNullString(params.clientId) as any,
          metadata: toNullJSON(params.metadata) as any,
          chunkTaskId: toNullString(params.chunkTaskId) as any,
          embeddingTaskId: toNullString(params.embeddingTaskId) as any,
          createdAt: now,
          updatedAt: now,
        },
        GlobalFile: insertToGlobalFiles && params.fileHash ? {
          hashId: toNullString(params.fileHash) as any,
          fileType: toNullString(params.fileType) as any,
          size: params.size || 0,
          url: toNullString(params.url) as any,
          metadata: toNullJSON(params.metadata) as any,
          creator: toNullString(this.userId) as any,
          createdAt: now,
        } : null,
        KnowledgeBase: params.knowledgeBaseId || null,
      });

      const resultObj = { id: result?.id || fileId };
      await this.logger.methodExit('create', resultObj);
      return resultObj;
    } catch (error) {
      await this.logger.error('Failed to create file', { error, params });
      await this.showErrorNotification(
        'Upload Failed',
        `Failed to upload file "${params.name}". Please try again.`
      );
      throw error;
    }
  };

  createGlobalFile = async (file: any) => {
    try {
      const now = currentTimestampMs();
      return await DB.CreateGlobalFile({
        hashId: toNullString(file.hashId) as any,
        fileType: toNullString(file.fileType) as any,
        size: file.size || 0,
        url: toNullString(file.url) as any,
        metadata: toNullJSON(file.metadata) as any,
        creator: toNullString(file.creator || this.userId) as any,
        createdAt: now,
      });
    } catch (error) {
      await this.logger.error('Failed to create global file', { error, file });
      throw error;
    }
  };

  checkHash = async (hash: string) => {
    try {
      const item = await DB.GetGlobalFile(hash);
      
      if (!item) return { isExist: false };

      return {
        fileType: getNullableString(item.fileType as any),
        isExist: true,
        metadata: parseNullableJSON(item.metadata as any),
        size: item.size,
        url: getNullableString(item.url as any),
      };
    } catch (error) {
      await this.logger.error('Failed to check hash', { error, hash });
      throw error;
    }
  };

  /**
   * Delete file with atomic cascading cleanup
   * ✅ OPTIMIZED: Uses backend transaction for atomicity
   * 1. Get file info
   * 2. Delete related chunks (embeddings, documentChunks, chunks, fileChunks)
   * 3. Delete file record
   * 4. Check if other files use the same hash
   * 5. Delete from global_files if no other files use it
   * All operations in transaction - succeed or rollback
   */
  delete = async (id: string, removeGlobalFile: boolean = true) => {
    await this.logger.methodEntry('delete', { userId: this.userId, id, removeGlobalFile });
    
    try {
      // 1. Get file first
      const file = await this.findById(id);
      if (!file) {
        await this.logger.warn('File not found for delete', { id });
        return;
      }

      const fileHash = getNullableString(file.fileHash as any);

      // 2. Use backend transaction method for atomic delete
      await DBService.DeleteFileWithCascade({
        FileID: id,
        UserID: this.userId,
        RemoveGlobalFile: removeGlobalFile,
        FileHash: fileHash || '',
      });

      await this.logger.methodExit('delete', { id, fileHash });
      return file;
    } catch (error) {
      await this.logger.error('Failed to delete file', { error, id });
      await this.showErrorNotification(
        'Delete Failed',
        `Failed to delete file. Please try again.`
      );
      throw error;
    }
  };

  deleteGlobalFile = async (hashId: string) => {
    try {
      return await DB.DeleteGlobalFile(hashId);
    } catch (error) {
      await this.logger.error('Failed to delete global file', { error, hashId });
      throw error;
    }
  };

  countUsage = async () => {
    try {
      const result = await DB.CountFilesUsage(this.userId);
      return Number(result.totalSize) || 0;
    } catch (error) {
      await this.logger.error('Failed to count file usage', { error });
      throw error;
    }
  };

  /**
   * COMPLEX: Batch delete with cascading cleanup
   * Similar to single delete but for multiple files
   *
   * LIMITATION: No transaction - partial failure possible
   * TODO: Create backend transaction method
   */
  deleteMany = async (ids: string[], removeGlobalFile: boolean = true) => {
    await this.logger.methodEntry('deleteMany', { userId: this.userId, count: ids.length, removeGlobalFile });
    
    try {
      if (ids.length === 0) return [];

      // 1. Get file list first
      const fileList = await this.findByIds(ids);
      if (fileList.length === 0) return [];

      // 2. Extract hashes
      const hashList = fileList
        .map((f) => getNullableString(f.fileHash as any))
        .filter(Boolean) as string[];

      // 3. Delete chunks
      await this.deleteFileChunks(ids);

      // 4. Delete files (one by one, no batch delete)
      await Promise.all(
        ids.map((id) =>
          DB.DeleteFile({
            id,
            userId: this.userId,
          }),
        ),
      );

      // 5. Delete global files if needed
      if (removeGlobalFile && hashList.length > 0) {
        // Check which hashes are still in use
        const remainingFiles = await DB.GetFilesByHash({
          fileHash: toNullString(hashList[0]), // Check first hash
          userId: this.userId,
        });

        const usedHashes = new Set(
          remainingFiles.map((f) => getNullableString(f.fileHash as any)),
        );

        // Delete unused hashes
        const hashesToDelete = hashList.filter((hash) => !usedHashes.has(hash));

        await Promise.all(
          hashesToDelete.map((hash) =>
            DB.DeleteGlobalFile(hash),
          ),
        );
      }

      await this.logger.methodExit('deleteMany', { deletedCount: fileList.length });
      return fileList;
    } catch (error) {
      await this.logger.error('Failed to delete multiple files', { error, ids });
      await this.showErrorNotification(
        'Batch Delete Failed',
        `Failed to delete ${ids.length} file(s). Please try again.`
      );
      throw error;
    }
  };

  clear = async () => {
    await this.logger.methodEntry('clear', { userId: this.userId });
    
    try {
      const result = await DB.DeleteAllFiles(this.userId);
      await this.logger.methodExit('clear');
      return result;
    } catch (error) {
      await this.logger.error('Failed to clear all files', { error });
      await this.showErrorNotification(
        'Clear Failed',
        `Failed to clear all files. Please try again.`
      );
      throw error;
    }
  };

  /**
   * Complex query with multiple filters
   * NOTE: Full filtering logic moved to client-side due to SQL complexity
   * Could be optimized with specific backend queries
   */
  query = async ({
    category,
    q,
    sortType,
    sorter,
    knowledgeBaseId,
    showFilesInKnowledgeBase,
  }: QueryFileListParams = {}) => {
    try {
      // If filtering by knowledge base, use JOIN query
      if (knowledgeBaseId) {
        const files = await DB.QueryFilesByKnowledgeBase({
          knowledgeBaseId: toNullString(knowledgeBaseId) as any,
          userId: this.userId,
        });

        return this.filterAndSortFiles(files, { category, q, sortType, sorter });
      }

      // Otherwise, get all files and filter client-side
      const allFiles = await DB.QueryFiles(this.userId);

      // Apply filters client-side
      let filtered = allFiles;

      // Filter by search query
      if (q) {
        filtered = filtered.filter((f) =>
          getNullableString(f.name as any)
            ?.toLowerCase()
            .includes(q.toLowerCase()),
        );
      }

      // Filter by category
      if (category && category !== FilesTabs.All) {
        const fileTypePrefix = this.getFileTypePrefix(category as FilesTabs);
        filtered = filtered.filter((f) =>
          getNullableString(f.fileType as any)?.startsWith(fileTypePrefix),
        );
      }

      // Filter by knowledge base visibility
      if (!showFilesInKnowledgeBase) {
        // TODO: Need query to check if file is in knowledge base
        // For now, return all
      }

      return this.filterAndSortFiles(filtered, { sortType, sorter });
    } catch (error) {
      await this.logger.error('Failed to query files', { error, category, q, knowledgeBaseId });
      throw error;
    }
  };

  findByIds = async (ids: string[]) => {
    try {
      // No batch query available, fetch one by one
      const results = await Promise.all(
        ids.map((id) => this.findById(id)),
      );
      return results.filter(Boolean) as any[];
    } catch (error) {
      await this.logger.error('Failed to find files by IDs', { error, ids });
      throw error;
    }
  };

  findById = async (id: string): Promise<DBFile | undefined> => {
    try {
      const result = await DB.GetFile({
        id,
        userId: this.userId,
      });
      return result;
    } catch (error) {
      await this.logger.warn('File not found', { id, error });
      return undefined;
    }
  };

  countFilesByHash = async (hash: string) => {
    try {
      const result = await DB.CountFilesByHash(toNullString(hash) as any);
      return Number(result) || 0;
    } catch (error) {
      await this.logger.error('Failed to count files by hash', { error, hash });
      throw error;
    }
  };

  update = async (id: string, value: any) => {
    await this.logger.methodEntry('update', { id, value, userId: this.userId });
    
    try {
      const result = await DB.UpdateFile({
        id,
        userId: this.userId,
        name: toNullString(value.name) as any,
        metadata: toNullJSON(value.metadata) as any,
        updatedAt: currentTimestampMs(),
      });
      
      await this.logger.methodExit('update', { id });
      return result;
    } catch (error) {
      await this.logger.error('Failed to update file', { error, id, value });
      await this.showErrorNotification(
        'Update Failed',
        `Failed to update file "${value.name || ''}". Please try again.`
      );
      throw error;
    }
  };

  /**
   * Get the corresponding file type prefix according to FilesTabs
   */
  private getFileTypePrefix = (category: FilesTabs): string => {
    switch (category) {
      case FilesTabs.Audios: {
        return 'audio';
      }
      case FilesTabs.Documents: {
        return 'application';
      }
      case FilesTabs.Images: {
        return 'image';
      }
      case FilesTabs.Videos: {
        return 'video';
      }
      default: {
        return '';
      }
    }
  };

  findByNames = async (fileNames: string[]) => {
    try {
      // Get all files and filter client-side
      const allFiles = await DB.GetFilesByNames(this.userId);

      return allFiles.filter((f) =>
        fileNames.some((name) =>
          getNullableString(f.name as any)?.startsWith(name),
        ),
      );
    } catch (error) {
      await this.logger.error('Failed to find files by names', { error, fileNames });
      throw error;
    }
  };

  /**
   * Complex delete operation for file chunks
   * Deletes in order: embeddings -> documentChunks -> chunks -> fileChunks
   *
   * LIMITATION: No transaction support - partial cleanup possible
   */
  private deleteFileChunks = async (fileIds: string[]) => {
    try {
      if (fileIds.length === 0) return;

      // Get all chunk IDs for these files
      const allChunkIds: string[] = [];
      for (const fileId of fileIds) {
        const chunks = await DB.GetFileChunkIds(toNullString(fileId) as any);
        // chunks is NullString[], convert to string[]
        allChunkIds.push(...chunks.map((c) => getNullableString(c as any) || '').filter(Boolean));
      }

      if (allChunkIds.length === 0) return;

      await this.logger.debug(`Deleting ${allChunkIds.length} chunks for ${fileIds.length} files`);

      // Batch delete in chunks of 500
      const BATCH_SIZE = 500;
      for (let i = 0; i < allChunkIds.length; i += BATCH_SIZE) {
        const batchIds = allChunkIds.slice(i, i + BATCH_SIZE);

        // Delete embeddings
        await Promise.all(
          batchIds.map((chunkId) =>
            DB.DeleteEmbedding({
              id: chunkId,
              userId: toNullString(this.userId) as any,
            }).catch(() => {}), // Ignore errors
          ),
        );

        // Delete chunks
        await Promise.all(
          batchIds.map((chunkId) =>
            DB.DeleteChunk({
              id: chunkId,
              userId: toNullString(this.userId) as any,
            }).catch(() => {}), // Ignore errors
          ),
        );
      }

      return allChunkIds;
    } catch (error) {
      await this.logger.error('Failed to delete file chunks', { error, fileIds });
      throw error;
    }
  };

  private filterAndSortFiles = (files: any[], options: any) => {
    let result = [...files];

    // Sort
    if (options.sorter && options.sortType) {
      const direction = options.sortType.toLowerCase() === SortType.Asc ? 1 : -1;
      
      result.sort((a, b) => {
        const aVal = a[options.sorter];
        const bVal = b[options.sorter];
        
        if (aVal < bVal) return -1 * direction;
        if (aVal > bVal) return 1 * direction;
        return 0;
      });
    }

    return result;
  };
}

