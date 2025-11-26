import { BaseClientService } from '@/services/baseClientService';
import { SemanticSearchSchemaType, SemanticSearchResult } from '@/types/rag';

import { IRAGService } from './type';

/**
 * ClientService for RAG operations using Go backend services
 * 
 * This service integrates with:
 * - FileProcessorService (Go): File parsing, document storage, RAG processing
 * - VectorSearchService (Go): Semantic search using chromem (internal/services)
 * - DocumentService (Go): Document CRUD operations
 */
export class ClientService extends BaseClientService implements IRAGService {
  /**
   * Parse file content and create document
   * Uses Go FileProcessorService.ProcessFileForStorage
   */
  parseFileContent = async (id: string, skipExist?: boolean): Promise<any> => {
    try {
      // Import Go services dynamically
      const { GetDocument } = await import('@@/github.com/kawai-network/veridium/internal/database/generated/queries');
      
      // Get document from database
      const doc = await GetDocument({
        id: id,
        userId: this.userId,
      });

      return {
        id: doc.id,
        title: this.getNullableString(doc.title),
        content: this.getNullableString(doc.content),
        fileType: doc.fileType,
        filename: this.getNullableString(doc.filename),
        totalCharCount: doc.totalCharCount || 0,
        totalLineCount: doc.totalLineCount || 0,
        metadata: doc.metadata ? JSON.parse(this.getNullableString(doc.metadata) || '{}') : undefined,
        pages: doc.pages ? JSON.parse(this.getNullableString(doc.pages) || '[]') : undefined,
        sourceType: doc.sourceType,
        source: doc.source,
        fileId: this.getNullableString(doc.fileId),
        userId: doc.userId,
        clientId: this.getNullableString(doc.clientId),
        createdAt: new Date(doc.createdAt),
        updatedAt: new Date(doc.updatedAt),
      };
    } catch (error) {
      console.error('Failed to parse file content:', error);
      throw error;
    }
  };

  /**
   * Create parse file task
   * Uses Go FileProcessorService.ProcessFileForStorage
   * Note: This processes immediately (no task queue in Wails)
   */
  createParseFileTask = async (id: string, skipExist?: boolean): Promise<any> => {
    try {
      // Import Go services dynamically
      const { GetFile } = await import('@@/github.com/kawai-network/veridium/internal/database/generated/queries');
      const { ProcessFileForStorage } = await import('@@/github.com/kawai-network/veridium/fileprocessorservice');

      // Get file info
      const file = await GetFile({
        id: id,
        userId: this.userId,
      });

      if (!file) {
        throw new Error(`File not found: ${id}`);
      }

      // Process file (parse + save + RAG)
      const result = await ProcessFileForStorage(
        file.url, // filePath
        this.getNullableString(file.name) || 'unknown',
        file.fileType,
        this.userId,
        true, // enableRAG
      );

      return {
        id: result?.documentId || id,
        status: 'completed',
        fileId: result?.fileId || id,
        documentId: result?.documentId,
        chunkIds: result?.chunkIds || [],
      };
    } catch (error) {
      console.error('Failed to create parse file task:', error);
      return {
        id,
        status: 'failed',
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  };

  /**
   * Retry parse file
   * Uses Go FileProcessorService.ProcessFileForStorage
   */
  retryParseFile = async (id: string): Promise<any> => {
    // Same as createParseFileTask but forces reprocessing
    return this.createParseFileTask(id, false);
  };

  /**
   * Create embedding chunks task
   * Uses Go RAGProcessor (automatically called by ProcessFileForStorage)
   */
  createEmbeddingChunksTask = async (id: string): Promise<any> => {
    try {
      // Import Go services dynamically
      const { ProcessFileForStorage } = await import('@@/github.com/kawai-network/veridium/fileprocessorservice');
      const { GetFile } = await import('@@/github.com/kawai-network/veridium/internal/database/generated/queries');

      // Get file info
      const file = await GetFile({
        id: id,
        userId: this.userId,
      });

      if (!file) {
        throw new Error(`File not found: ${id}`);
      }

      // Process file with RAG enabled
      const result = await ProcessFileForStorage(
        file.url,
        this.getNullableString(file.name) || 'unknown',
        file.fileType,
        this.userId,
        true, // enableRAG - this triggers embedding generation
      );

      return {
        id: result?.documentId || id,
        status: 'completed',
        chunkIds: result?.chunkIds || [],
      };
    } catch (error) {
      console.error('Failed to create embedding chunks task:', error);
      return {
        id,
        status: 'failed',
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  };

  /**
   * Semantic search using Go VectorSearchService (internal/services)
   */
  semanticSearch = async (query: string, fileIds?: string[]): Promise<SemanticSearchResult[]> => {
    try {
      // Import Go VectorSearchService dynamically (from internal/services)
      const { SemanticSearch } = await import('@@/github.com/kawai-network/veridium/internal/services/vectorsearchservice');

      // Call Go semantic search
      const results = await SemanticSearch(
        this.userId,
        query,
        fileIds || [],
        30, // limit
      );

      // Convert Go results to frontend format
      return results.map((result: any) => ({
        id: result.id,
        similarity: result.similarity,
        text: result.text,
        fileId: result.fileId,
        fileName: result.fileName,
        type: result.type,
        index: result.index,
        metadata: result.metadata,
      }));
    } catch (error) {
      console.error('Semantic search failed:', error);
      // Fallback to empty results
      return [];
    }
  };


  /**
   * Delete message RAG query
   */

  /**
   * Helper to extract string from Go NullString
   */
  private getNullableString(ns: any): string | undefined {
    if (!ns) return undefined;
    if (typeof ns === 'string') return ns;
    if (ns.String && ns.Valid) return ns.String;
    return undefined;
  }
}
