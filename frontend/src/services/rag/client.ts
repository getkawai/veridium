import { BaseClientService } from '@/services/baseClientService';
import { SemanticSearchResult } from '@/types/rag';

import { IRAGService } from './type';

/**
 * ClientService for RAG operations using Go backend services
 * 
 * This service integrates with:
 * - VectorSearchService (Go): Semantic search using chromem (internal/services)
 */
export class ClientService extends BaseClientService implements IRAGService {
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
}
