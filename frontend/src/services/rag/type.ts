import { SemanticSearchSchemaType, SemanticSearchResult } from '@/types/rag';

export interface IRAGService {
  parseFileContent(id: string, skipExist?: boolean): Promise<any>;
  createParseFileTask(id: string, skipExist?: boolean): Promise<any>;
  retryParseFile(id: string): Promise<any>;
  createEmbeddingChunksTask(id: string): Promise<any>;
  semanticSearch(query: string, fileIds?: string[]): Promise<SemanticSearchResult[]>;
  semanticSearchForChat(params: SemanticSearchSchemaType): Promise<SemanticSearchResult[]>;
}

