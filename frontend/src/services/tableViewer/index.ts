/**
 * TableViewer Service - Direct Wails binding for database table inspection
 * Replaces the old tRPC pgTable router
 */

import { Service as TableViewerService } from '@@/github.com/kawai-network/veridium/internal/services/tableviewer/service';

export { TableViewerService };

// Re-export types from generated bindings
export type {
  TableBasicInfo,
  TableColumnInfo,
  PaginationParams,
  FilterCondition,
  TableDataResult,
  PaginationResult,
} from '@@/github.com/kawai-network/veridium/internal/services/tableviewer/models';
