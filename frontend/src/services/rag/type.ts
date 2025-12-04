import { SemanticSearchResult } from '@/types/rag';

export interface IRAGService {
  semanticSearch(query: string, fileIds?: string[]): Promise<SemanticSearchResult[]>;
}
