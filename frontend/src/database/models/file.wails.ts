import { FilesTabs, QueryFileListParams, SortType } from  '@/types';
import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  toNullInt,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';

export class FileModel {
  private readonly userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
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
    const fileId = nanoid();
    const now = currentTimestampMs();

    // Use backend transaction method for atomic operations
    const result = await DBService.CreateFileWithLinks({
      file: {
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
      globalFile: insertToGlobalFiles && params.fileHash ? {
        hashId: toNullString(params.fileHash) as any,
        fileType: toNullString(params.fileType) as any,
        size: params.size || 0,
        url: toNullString(params.url) as any,
        metadata: toNullJSON(params.metadata) as any,
        creator: toNullString(this.userId) as any,
        createdAt: now,
        accessedAt: now,
      } : null,
      knowledgeBase: params.knowledgeBaseId || null,
    });

    return { id: result?.id || fileId };
  };

  createGlobalFile = async (file: any) => {
    const now = currentTimestampMs();
    return await DB.CreateGlobalFile({
      hashId: toNullString(file.hashId),
      fileType: toNullString(file.fileType),
      size: file.size || 0,
      url: toNullString(file.url),
      metadata: toNullJSON(file.metadata),
      creator: toNullString(file.creator || this.userId),
      createdAt: now,
      accessedAt: now,
    });
  };

  checkHash = async (hash: string) => {
    const item = await DB.GetGlobalFile({
      hashId: toNullString(hash),
    });
    
    if (!item) return { isExist: false };

    return {
      fileType: getNullableString(item.fileType as any),
      isExist: true,
      metadata: parseNullableJSON(item.metadata as any),
      size: item.size,
      url: getNullableString(item.url as any),
    };
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
    // 1. Get file first
    const file = await this.findById(id);
    if (!file) return;

    const fileHash = getNullableString(file.fileHash as any);

    // 2. Use backend transaction method for atomic delete
    await DBService.DeleteFileWithCascade({
      fileId: id,
      userId: this.userId,
      removeGlobalFile: removeGlobalFile,
      fileHash: fileHash || '',
    });

    return file;
  };

  deleteGlobalFile = async (hashId: string) => {
    return await DB.DeleteGlobalFile({
      hashId: toNullString(hashId),
    });
  };

  countUsage = async () => {
    const result = await DB.CountFilesUsage({
      userId: this.userId,
    });

    return Number(result.totalSize) || 0;
  };

  /**
   * COMPLEX: Batch delete with cascading cleanup
   * Similar to single delete but for multiple files
   *
   * LIMITATION: No transaction - partial failure possible
   * TODO: Create backend transaction method
   */
  deleteMany = async (ids: string[], removeGlobalFile: boolean = true) => {
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
          DB.DeleteGlobalFile({
            hashId: toNullString(hash),
          }),
        ),
      );
    }

    return fileList;
  };

  clear = async () => {
    return await DB.DeleteAllFiles({
      userId: this.userId,
    });
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
    // If filtering by knowledge base, use JOIN query
    if (knowledgeBaseId) {
      const files = await DB.QueryFilesByKnowledgeBase({
        knowledgeBaseId: toNullString(knowledgeBaseId),
        userId: this.userId,
      });

      return this.filterAndSortFiles(files, { category, q, sortType, sorter });
    }

    // Otherwise, get all files and filter client-side
    const allFiles = await DB.QueryFiles({
      userId: this.userId,
    });

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
  };

  findByIds = async (ids: string[]) => {
    // No batch query available, fetch one by one
    const results = await Promise.all(
      ids.map((id) => this.findById(id)),
    );
    return results.filter(Boolean) as any[];
  };

  findById = async (id: string) => {
    return await DB.GetFile({
      id,
      userId: this.userId,
    });
  };

  countFilesByHash = async (hash: string) => {
    const result = await DB.CountFilesByHash({
      fileHash: toNullString(hash),
    });

    return Number(result.count) || 0;
  };

  update = async (id: string, value: any) => {
    return await DB.UpdateFile({
      id,
      userId: this.userId,
      name: toNullString(value.name),
      metadata: toNullJSON(value.metadata),
      updatedAt: currentTimestampMs(),
    });
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
    // Get all files and filter client-side
    const allFiles = await DB.GetFilesByNames({
      userId: this.userId,
    });

    return allFiles.filter((f) =>
      fileNames.some((name) =>
        getNullableString(f.name as any)?.startsWith(name),
      ),
    );
  };

  /**
   * Complex delete operation for file chunks
   * Deletes in order: embeddings -> documentChunks -> chunks -> fileChunks
   *
   * LIMITATION: No transaction support - partial cleanup possible
   */
  private deleteFileChunks = async (fileIds: string[]) => {
    if (fileIds.length === 0) return;

    // Get all chunk IDs for these files
    const allChunkIds: string[] = [];
    for (const fileId of fileIds) {
      const chunks = await DB.GetFileChunkIds({
        fileId,
      });
      allChunkIds.push(...chunks.map((c) => c.chunkId));
    }

    if (allChunkIds.length === 0) return;

    // Batch delete in chunks of 500
    const BATCH_SIZE = 500;
    for (let i = 0; i < allChunkIds.length; i += BATCH_SIZE) {
      const batchIds = allChunkIds.slice(i, i + BATCH_SIZE);

      // Delete embeddings
      await Promise.all(
        batchIds.map((chunkId) =>
          DB.DeleteEmbedding({
            id: chunkId,
            userId: this.userId,
          }).catch(() => {}), // Ignore errors
        ),
      );

      // Delete chunks
      await Promise.all(
        batchIds.map((chunkId) =>
          DB.DeleteChunk({
            id: chunkId,
            userId: this.userId,
          }).catch(() => {}), // Ignore errors
        ),
      );
    }

    // Delete file_chunks links
    for (const fileId of fileIds) {
      // TODO: Need batch unlink query
      // For now, individual deletes handled by CASCADE
    }

    return allChunkIds;
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

