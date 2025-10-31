import { ChunkMetadata, FileChunk } from  '@/types';
import { and, asc, count, desc, eq, inArray, isNull, sql } from 'drizzle-orm';
import { chunk } from 'lodash-es';

import {
  NewChunkItem,
  NewUnstructuredChunkItem,
  chunks,
  embeddings,
  fileChunks,
  files,
  unstructuredChunks,
} from '../schemas';
import { LobeChatDatabase } from '../type';
import { bufferToVector, cosineSimilarity } from '../utils/vectorSearch';

export class ChunkModel {
  private userId: string;

  private db: LobeChatDatabase;

  constructor(db: LobeChatDatabase, userId: string) {
    this.userId = userId;
    this.db = db;
  }

  bulkCreate = async (params: NewChunkItem[], fileId: string) => {
    return this.db.transaction(async (trx) => {
      if (params.length === 0) return [];

      const result = await trx.insert(chunks).values(params).returning();

      const fileChunksData = result.map((chunk) => ({
        chunkId: chunk.id,
        fileId,
        userId: this.userId,
      }));

      if (fileChunksData.length > 0) {
        await trx.insert(fileChunks).values(fileChunksData);
      }

      return result;
    });
  };

  bulkCreateUnstructuredChunks = async (params: NewUnstructuredChunkItem[]) => {
    return this.db.insert(unstructuredChunks).values(params);
  };

  delete = async (id: string) => {
    return this.db.delete(chunks).where(and(eq(chunks.id, id), eq(chunks.userId, this.userId)));
  };

  deleteOrphanChunks = async () => {
    const orphanedChunks = await this.db
      .select({ chunkId: chunks.id })
      .from(chunks)
      .leftJoin(fileChunks, eq(chunks.id, fileChunks.chunkId))
      .where(isNull(fileChunks.fileId));

    const ids = orphanedChunks.map((chunk) => chunk.chunkId);
    if (ids.length === 0) return;

    const list = chunk(ids, 500);

    await this.db.transaction(async (trx) => {
      await Promise.all(
        list.map(async (chunkIds) => {
          await trx.delete(chunks).where(inArray(chunks.id, chunkIds));
        }),
      );
    });
  };

  findById = async (id: string) => {
    return this.db.query.chunks.findFirst({
      where: and(eq(chunks.id, id)),
    });
  };

  findByFileId = async (id: string, page = 0) => {
    const data = await this.db
      .select({
        abstract: chunks.abstract,
        createdAt: chunks.createdAt,
        id: chunks.id,
        index: chunks.index,
        metadata: chunks.metadata,
        text: chunks.text,
        type: chunks.type,
        updatedAt: chunks.updatedAt,
      })
      .from(chunks)
      .innerJoin(fileChunks, eq(chunks.id, fileChunks.chunkId))
      .where(and(eq(fileChunks.fileId, id), eq(chunks.userId, this.userId)))
      .limit(20)
      .offset(page * 20)
      .orderBy(asc(chunks.index));

    return data.map((item) => {
      const metadata = item.metadata as ChunkMetadata;

      return { ...item, metadata, pageNumber: metadata?.pageNumber } as FileChunk;
    });
  };

  getChunksTextByFileId = async (id: string): Promise<{ id: string; text: string }[]> => {
    const data = await this.db
      .select()
      .from(chunks)
      .innerJoin(fileChunks, eq(chunks.id, fileChunks.chunkId))
      .where(eq(fileChunks.fileId, id));

    return data
      .map((item) => item.chunks)
      .map((chunk) => ({ id: chunk.id, text: this.mapChunkText(chunk) }))
      .filter((chunk) => chunk.text) as { id: string; text: string }[];
  };

  countByFileIds = async (ids: string[]) => {
    if (ids.length === 0) return [];

    return this.db
      .select({
        count: count(fileChunks.chunkId),
        id: fileChunks.fileId,
      })
      .from(fileChunks)
      .where(inArray(fileChunks.fileId, ids))
      .groupBy(fileChunks.fileId);
  };

  countByFileId = async (ids: string) => {
    const data = await this.db
      .select({
        count: count(fileChunks.chunkId),
        id: fileChunks.fileId,
      })
      .from(fileChunks)
      .where(eq(fileChunks.fileId, ids))
      .groupBy(fileChunks.fileId);

    return data[0]?.count ?? 0;
  };

  semanticSearch = async ({
    embedding,
    fileIds,
  }: {
    embedding: number[];
    fileIds: string[] | undefined;
    query: string;
  }) => {
    // Fetch all chunks with embeddings (SQLite doesn't have native vector search)
    const data = await this.db
      .select({
        chunkEmbedding: embeddings.embeddings,
        fileId: fileChunks.fileId,
        fileName: files.name,
        id: chunks.id,
        index: chunks.index,
        metadata: chunks.metadata,
        text: chunks.text,
        type: chunks.type,
      })
      .from(chunks)
      .leftJoin(embeddings, eq(chunks.id, embeddings.chunkId))
      .leftJoin(fileChunks, eq(chunks.id, fileChunks.chunkId))
      .leftJoin(files, eq(fileChunks.fileId, files.id))
      .where(fileIds ? inArray(fileChunks.fileId, fileIds) : undefined);

    // Calculate similarity in JavaScript
    const withSimilarity = data
      .filter((item) => item.chunkEmbedding) // Only items with embeddings
      .map((item) => {
        const chunkVector = bufferToVector(item.chunkEmbedding as ArrayBuffer | Uint8Array);
        const similarity = cosineSimilarity(embedding, chunkVector);
        return {
          ...item,
          similarity,
        };
      })
      .sort((a, b) => b.similarity - a.similarity)
      .slice(0, 30);

    return withSimilarity.map((item) => ({
      fileId: item.fileId,
      fileName: item.fileName,
      id: item.id,
      index: item.index,
      metadata: item.metadata as ChunkMetadata,
      similarity: item.similarity,
      text: item.text,
      type: item.type,
    }));
  };

  semanticSearchForChat = async ({
    embedding,
    fileIds,
  }: {
    embedding: number[];
    fileIds: string[] | undefined;
    query: string;
  }) => {
    const hasFiles = fileIds && fileIds.length > 0;

    if (!hasFiles) return [];

    // Fetch all chunks with embeddings
    const result = await this.db
      .select({
        chunkEmbedding: embeddings.embeddings,
        fileId: files.id,
        fileName: files.name,
        id: chunks.id,
        index: chunks.index,
        metadata: chunks.metadata,
        text: chunks.text,
        type: chunks.type,
      })
      .from(chunks)
      .leftJoin(embeddings, eq(chunks.id, embeddings.chunkId))
      .leftJoin(fileChunks, eq(chunks.id, fileChunks.chunkId))
      .leftJoin(files, eq(files.id, fileChunks.fileId))
      .where(inArray(fileChunks.fileId, fileIds));

    // Calculate similarity in JavaScript
    const withSimilarity = result
      .filter((item) => item.chunkEmbedding)
      .map((item) => {
        const chunkVector = bufferToVector(item.chunkEmbedding as ArrayBuffer | Uint8Array);
        const similarity = cosineSimilarity(embedding, chunkVector);
        return {
          ...item,
          similarity,
        };
      })
      .sort((a, b) => b.similarity - a.similarity)
      .slice(0, 15);

    return withSimilarity.map((item) => {
      return {
        fileId: item.fileId,
        fileName: item.fileName,
        id: item.id,
        index: item.index,
        similarity: item.similarity,
        text: this.mapChunkText(item),
      };
    });
  };

  private mapChunkText = (chunk: { metadata: any; text: string | null; type: string | null }) => {
    let text = chunk.text;

    if (chunk.type === 'Table') {
      text = `${chunk.text}

content in Table html is below:
${(chunk.metadata as ChunkMetadata).text_as_html}
`;
    }

    return text;
  };
}
