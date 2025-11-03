/**
 * Search Service - Direct Wails binding for web search and crawling
 * Replaces the old tRPC search router
 */

import * as SearchServiceNamespace from '@@/github.com/kawai-network/veridium/internal/services/search';

// Export the service
export const SearchService = SearchServiceNamespace.Service;

// Re-export types from generated bindings
export type {
	SearchQuery,
	SearchParams,
	UniformSearchResponse,
	CrawlPagesRequest,
	CrawlPagesResponse,
	CrawlResult,
	CrawlSuccessResult,
	CrawlErrorResult,
	CrawlImplType,
} from '@@/github.com/kawai-network/veridium/internal/services/search';

