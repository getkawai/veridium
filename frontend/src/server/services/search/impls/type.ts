import { SearchParams, UniformSearchResponse } from '@/types';

/**
 * Search service implementation interface
 */
export interface SearchServiceImpl {
  /**
   * Query for search results
   */
  query(query: string, params?: SearchParams): Promise<UniformSearchResponse>;
}
