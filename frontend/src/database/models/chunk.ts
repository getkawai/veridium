import { ChunkMetadata, FileChunk } from  '@/types';
import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

import { bufferToVector, cosineSimilarity } from '../utils/vectorSearch';
import * as VectorSearch from '../../../bindings/github.com/kawai-network/veridium/internal/services/vectorsearchservice';
import { NewChunkItem } from '@/types/database-legacy';

export class ChunkModel {
  private userId: string;
  private logger = createModelLogger('Chunk', 'ChunkModel', 'database/models/chunk');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  /**
   * OPTIMIZED: Uses transaction-like behavior via sequential inserts
   * Now also adds chunks to chromem vector database for semantic search
   */
  bulkCreate = async (params: NewChunkItem[], fileId: string) => {
    if (params.length === 0) return [];

    const result: NewChunkItem[] = [];
    
    // Create chunks in SQLite
    for (const param of params) {
      const id = param.id || nanoid();
      const now = currentTimestampMs();
      
      const chunk = await DB.CreateChunk({
        id,
        text: toNullString(param.text),
        abstract: toNullString(param.abstract),
        metadata: toNullJSON(param.metadata),
        chunkIndex: { Int64: param.index || 0, Valid: true },
        type: toNullString(param.type),
        clientId: toNullString(param.clientId),
        userId: toNullString(this.userId),
        createdAt: now,
        updatedAt: now,
      });
      
      result.push({
        id: chunk.id,
        text: getNullableString(chunk.text as any) || null,
        abstract: getNullableString(chunk.abstract as any) || null,
        metadata: getNullableString(chunk.metadata as any) || null,
        index: (chunk.chunkIndex as any)?.Int64 || 0,
        type: getNullableString(chunk.type as any) || null,
        clientId: getNullableString(chunk.clientId as any) || null,
        userId: this.userId,
        createdAt: now,
        updatedAt: now,
      });
      
      // Link to file
      await DB.LinkFileToChunk({
        fileId: toNullString(fileId),
        chunkId: toNullString(id),
        createdAt: now,
        userId: this.userId,
      });
    }

    // Add chunks to chromem vector database (async, non-blocking)
    try {
      // Get file name for metadata
      const file = await DB.GetFile({ id: fileId, userId: this.userId });
      const fileName = getNullableString(file?.name as any) || 'Unknown';

      const vectorChunks = result.map((chunk) => ({
        id: chunk.id || '',
        text: chunk.text || '',
        fileId: fileId,
        fileName: fileName,
        type: chunk.type || '',
        index: chunk.index || 0,
        metadata: {},
      }));

      // Add to vector database (fire and forget, logs errors internally)
      VectorSearch.AddChunks(this.userId, vectorChunks).catch((err) => {
        this.logger.error('Failed to add chunks to vector database:', err);
      });
    } catch (err) {
      this.logger.error('Failed to prepare chunks for vector database:', err);
    }

    return result;
  };

  bulkCreateUnstructuredChunks = async (params: any[]) => {
    const results: any[] = [];
    for (const param of params) {
      const id = param.id || nanoid();
      const now = currentTimestampMs();
      
      const chunk = await DB.CreateUnstructuredChunk({
        id,
        text: toNullString(param.text),
        metadata: toNullJSON(param.metadata),
        chunkIndex: param.index || 0,
        type: toNullString(param.type),
        parentId: toNullString(param.parentId),
        compositeId: toNullString(param.compositeId),
        clientId: toNullString(param.clientId),
        userId: toNullString(this.userId),
        fileId: toNullString(param.fileId),
        createdAt: now,
        updatedAt: now,
      });
      results.push(chunk);
    }
    return results;
  };

  delete = async (id: string) => {
    // Delete from SQLite
    const result = await DB.DeleteChunk({
      id,
      userId: toNullString(this.userId),
    });

    // Delete from chromem (fire and forget)
    VectorSearch.DeleteChunks(this.userId, [id]).catch((err) => {
      this.logger.error('Failed to delete chunk from vector database:', err);
    });

    return result;
  };

  /**
   * OPTIMIZED: Uses single query to find orphans, then batch delete
   */
  deleteOrphanChunks = async () => {
    const orphanedChunks = await DB.GetOrphanedChunks();
    
    if (orphanedChunks.length === 0) return;

    // SQLite doesn't support sqlc.slice, so delete in chunks
    const chunkSize = 500;
    for (let i = 0; i < orphanedChunks.length; i += chunkSize) {
      const batch = orphanedChunks.slice(i, i + chunkSize);
      
      // Delete one by one (limitation of no slice support)
      await Promise.all(
        batch.map((chunkId) =>
          DB.DeleteChunk({
            id: chunkId,
            userId: toNullString(this.userId),
          }),
        ),
      );
    }
  };

  findById = async (id: string) => {
    return await DB.GetChunk({
      id,
      userId: toNullString(this.userId),
    });
  };

  /**
   * OPTIMIZED: Uses JOIN query (3A)
   */
  findByFileId = async (id: string, page = 0) => {
    const data = await DB.GetFileChunksWithMetadata({
      fileId: toNullString(id),
      userId: toNullString(this.userId),
      limit: 20,
      offset: page * 20,
    });

    return data.map((item) => {
      const metadata = parseNullableJSON(item.metadata as any) as ChunkMetadata;

      return { 
        ...item, 
        metadata, 
        pageNumber: metadata?.pageNumber,
        index: item.chunkIndex,
      } as unknown as FileChunk;
    });
  };

  getChunksTextByFileId = async (id: string): Promise<{ id: string; text: string }[]> => {
    const data = await DB.GetChunksTextByFileId(toNullString(id));

    return data
      .map((item) => ({
        id: item.id,
        text: this.mapChunkText({
          text: getNullableString(item.text as any) || null,
          metadata: parseNullableJSON(item.metadata as any),
          type: getNullableString(item.type as any) || null,
        }),
      }))
      .filter((chunk) => chunk.text) as { id: string; text: string }[];
  };

  /**
   * OPTIMIZED: Single query with COUNT and GROUP BY
   */
  countByFileIds = async (ids: string[]) => {
    if (ids.length === 0) return [];

    // Note: No support for IN clause with array, so fetch all and filter
    const allCounts = await DB.CountChunksByFileIds(this.userId);

    return allCounts.filter((item) =>
      ids.includes(getNullableString(item.fileId as any) || ''),
    );
  };

  countByFileId = async (id: string) => {
    const count = await DB.CountChunksByFileId({
      fileId: toNullString(id),
      userId: this.userId,
    });

    return Number(count) || 0;
  };

  /**
   * Semantic search using chromem vector database (FAST!)
   * Falls back to client-side calculation if chromem is not available
   */
  semanticSearch = async ({
    embedding,
    fileIds,
    query,
  }: {
    embedding: number[];
    fileIds: string[] | undefined;
    query: string;
  }) => {
    try {
      // Use chromem for fast semantic search
      const results = await VectorSearch.SemanticSearchMultipleFiles(
        this.userId,
        query,
        fileIds || [],
        30
      );

      return results.map((item) => ({
        fileId: item.fileId,
        fileName: item.fileName,
        id: item.id,
        index: item.index,
        metadata: {} as ChunkMetadata, // Fetch from SQLite if needed
        similarity: item.similarity,
        text: item.text,
        type: item.type,
      }));
    } catch (err) {
      this.logger.warn('Chromem search failed, falling back to client-side:', err);

      // Fallback to old client-side calculation
      const data = fileIds && fileIds.length > 0
        ? await DB.GetChunksWithEmbeddingsByFileIds({
            fileId: toNullString(fileIds[0]),
            userId: this.userId,
          })
        : await DB.GetChunksWithEmbeddings(toNullString(this.userId));

      const withSimilarity = data
        .filter((item) => item.chunkEmbedding)
        .map((item) => {
          const chunkVector = bufferToVector(item.chunkEmbedding as any);
          const similarity = cosineSimilarity(embedding, chunkVector);
          return {
            fileId: getNullableString(item.fileId as any),
            fileName: getNullableString(item.fileName as any),
            id: item.id,
            index: item.chunkIndex || 0,
            metadata: parseNullableJSON(item.metadata as any) as ChunkMetadata,
            similarity,
            text: getNullableString(item.text as any),
            type: getNullableString(item.type as any),
          };
        })
        .sort((a, b) => b.similarity - a.similarity)
        .slice(0, 30);

      return withSimilarity;
    }
  };

  semanticSearchForChat = async ({
    embedding,
    fileIds,
    query,
  }: {
    embedding: number[];
    fileIds: string[] | undefined;
    query: string;
  }) => {
    const hasFiles = fileIds && fileIds.length > 0;
    if (!hasFiles) return [];

    try {
      // Use chromem for fast semantic search
      const results = await VectorSearch.SemanticSearchMultipleFiles(
        this.userId,
        query,
        fileIds,
        15
      );

      return results.map((item) => ({
        fileId: item.fileId,
        fileName: item.fileName,
        id: item.id,
        index: item.index,
        similarity: item.similarity,
        text: item.text, // Already mapped in backend
      }));
    } catch (err) {
      this.logger.warn('Chromem search failed, falling back to client-side:', err);

      // Fallback to old client-side calculation
      const result = await DB.GetChunksWithEmbeddingsByFileIds({
        fileId: toNullString(fileIds[0]),
        userId: this.userId,
      });

      const withSimilarity = result
        .filter((item) => item.chunkEmbedding)
        .map((item) => {
          const chunkVector = bufferToVector(item.chunkEmbedding as any);
          const similarity = cosineSimilarity(embedding, chunkVector);
          return {
            fileId: getNullableString(item.fileId as any),
            fileName: getNullableString(item.fileName as any),
            id: item.id,
            index: item.chunkIndex || 0,
            similarity,
            text: this.mapChunkText({
              text: getNullableString(item.text as any) || null,
              metadata: parseNullableJSON(item.metadata as any),
              type: getNullableString(item.type as any) || null,
            }),
          };
        })
        .sort((a, b) => b.similarity - a.similarity)
        .slice(0, 15);

      return withSimilarity;
    }
  };

  private mapChunkText = (chunk: { metadata: any; text: string | null; type: string | null }) => {
    let text = chunk.text;

    if (chunk.type === 'Table') {
      text = `${chunk.text}

content in Table html is below:
${(chunk.metadata as ChunkMetadata)?.text_as_html || ''}
`;
    }

    return text;
  };
}

