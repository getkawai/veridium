'use client';

import { useState, useEffect, useCallback } from 'react';
import {
  GetAllTables,
  GetTableData,
  GetTableDetails,
  ExecuteRawQuery
} from '@@/github.com/kawai-network/veridium/internal/tableviewer/service';
import type {
  TableBasicInfo,
  TableDataResult,
  TableColumnInfo,
  PaginationParams,
  FilterCondition
} from '@@/github.com/kawai-network/veridium/internal/tableviewer/models';

export interface UseSQLiteDataReturn {
  // Tables
  tables: TableBasicInfo[];
  tablesLoading: boolean;
  tablesError: string | null;
  loadTables: () => Promise<void>;

  // Table Data
  tableData: TableDataResult | null;
  tableDataLoading: boolean;
  tableDataError: string | null;
  loadTableData: (tableName: string, pagination: PaginationParams, filters?: FilterCondition[]) => Promise<void>;

  // Table Structure
  tableColumns: TableColumnInfo[];
  tableColumnsLoading: boolean;
  tableColumnsError: string | null;
  loadTableColumns: (tableName: string) => Promise<void>;

  // Raw Query
  queryResult: any[];
  queryLoading: boolean;
  queryError: string | null;
  executeQuery: (query: string, args?: any[]) => Promise<void>;
}

export const useSQLiteData = (): UseSQLiteDataReturn => {
  // Tables state
  const [tables, setTables] = useState<TableBasicInfo[]>([]);
  const [tablesLoading, setTablesLoading] = useState(false);
  const [tablesError, setTablesError] = useState<string | null>(null);

  // Table data state
  const [tableData, setTableData] = useState<TableDataResult | null>(null);
  const [tableDataLoading, setTableDataLoading] = useState(false);
  const [tableDataError, setTableDataError] = useState<string | null>(null);

  // Table columns state
  const [tableColumns, setTableColumns] = useState<TableColumnInfo[]>([]);
  const [tableColumnsLoading, setTableColumnsLoading] = useState(false);
  const [tableColumnsError, setTableColumnsError] = useState<string | null>(null);

  // Query state
  const [queryResult, setQueryResult] = useState<any[]>([]);
  const [queryLoading, setQueryLoading] = useState(false);
  const [queryError, setQueryError] = useState<string | null>(null);

  const loadTables = useCallback(async () => {
    try {
      setTablesLoading(true);
      setTablesError(null);
      const result = await GetAllTables();
      setTables(result || []);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load tables';
      setTablesError(errorMessage);
      console.error('Error loading tables:', err);
    } finally {
      setTablesLoading(false);
    }
  }, []);

  const loadTableData = useCallback(async (
    tableName: string,
    pagination: PaginationParams,
    filters: FilterCondition[] = []
  ) => {
    try {
      setTableDataLoading(true);
      setTableDataError(null);
      const result = await GetTableData(tableName, pagination, filters);
      setTableData(result);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load table data';
      setTableDataError(errorMessage);
      console.error('Error loading table data:', err);
    } finally {
      setTableDataLoading(false);
    }
  }, []);

  const loadTableColumns = useCallback(async (tableName: string) => {
    try {
      setTableColumnsLoading(true);
      setTableColumnsError(null);
      const result = await GetTableDetails(tableName);
      setTableColumns(result || []);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load table columns';
      setTableColumnsError(errorMessage);
      console.error('Error loading table columns:', err);
    } finally {
      setTableColumnsLoading(false);
    }
  }, []);

  const executeQuery = useCallback(async (query: string, args: any[] = []) => {
    try {
      setQueryLoading(true);
      setQueryError(null);
      setQueryResult([]);

      const result = await ExecuteRawQuery(query, args);
      const parsedResult = JSON.parse(result);

      if (Array.isArray(parsedResult)) {
        setQueryResult(parsedResult);
      } else {
        setQueryResult([]);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to execute query';
      setQueryError(errorMessage);
      console.error('Error executing query:', err);
    } finally {
      setQueryLoading(false);
    }
  }, []);

  // Auto-load tables on mount
  useEffect(() => {
    loadTables();
  }, [loadTables]);

  return {
    // Tables
    tables,
    tablesLoading,
    tablesError,
    loadTables,

    // Table Data
    tableData,
    tableDataLoading,
    tableDataError,
    loadTableData,

    // Table Structure
    tableColumns,
    tableColumnsLoading,
    tableColumnsError,
    loadTableColumns,

    // Raw Query
    queryResult,
    queryLoading,
    queryError,
    executeQuery,
  };
};

export default useSQLiteData;