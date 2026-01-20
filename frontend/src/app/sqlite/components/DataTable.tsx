"use client";

import {
  Table,
  Pagination,
  Input,
  Select,
  Button,
  Space,
  Spin,
  Alert,
  Tag,
  Typography,
} from "antd";
import {
  InboxOutlined,
  SearchOutlined,
  ReloadOutlined,
  FilterOutlined,
} from "@ant-design/icons";

const { Title, Text } = Typography;
import { createStyles } from "antd-style";
import { Flexbox } from "react-layout-kit";
import { memo, useState, useEffect, useMemo } from "react";
import {
  GetTableData,
  GetTableDetails,
} from "@@/github.com/kawai-network/veridium/internal/tableviewer/service";
import type {
  TableDataResult,
  PaginationParams,
  FilterCondition,
  TableColumnInfo,
} from "@@/github.com/kawai-network/veridium/internal/tableviewer/models";

const { Option } = Select;

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    display: flex;
    flex-direction: column;
  `,
  toolbar: css`
    padding: 16px;
    border-bottom: 1px solid ${token.colorBorder};
    background: ${token.colorBgContainer};
  `,
  filters: css`
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;
  `,
  tableContainer: css`
    flex: 1;
    overflow: auto;
    padding: 16px;
  `,
  pagination: css`
    padding: 16px;
    border-top: 1px solid ${token.colorBorder};
    text-align: center;
    background: ${token.colorBgContainer};
  `,
}));

interface DataTableProps {
  tableName: string;
}

const DataTable = memo<DataTableProps>(({ tableName }) => {
  const { styles } = useStyles();

  const [data, setData] = useState<TableDataResult | null>(null);
  const [columns, setColumns] = useState<TableColumnInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(50);

  // Filter state
  const [searchColumn, setSearchColumn] = useState<string>("");
  const [searchValue, setSearchValue] = useState<string>("");
  const [searchOperator, setSearchOperator] = useState<string>("contains");

  useEffect(() => {
    loadTableStructure();
  }, [tableName]);

  useEffect(() => {
    if (tableName) {
      loadTableData();
    }
  }, [tableName, currentPage, pageSize, searchColumn, searchValue, searchOperator]);

  useEffect(() => {
    if (tableName) {
      loadTableData();
    }
  }, [currentPage, pageSize, searchColumn, searchValue, searchOperator]);

  const loadTableStructure = async () => {
    try {
      const result = await GetTableDetails(tableName);
      setColumns(result || []);
    } catch (err) {
      console.error("Failed to load table structure:", err);
    }
  };

  const loadTableData = async () => {
    try {
      setLoading(true);
      setError(null);

      const pagination: PaginationParams = {
        page: currentPage,
        pageSize: pageSize,
        sortBy: null,
        sortOrder: null,
      };

      const filters: FilterCondition[] = [];
      if (searchColumn && searchValue) {
        filters.push({
          column: searchColumn,
          operator: searchOperator,
          value: searchValue,
        });
      }

      const result = await GetTableData(tableName, pagination, filters);
      setData(result);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load table data",
      );
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    setCurrentPage(1);
    loadTableData();
  };

  const handleClearSearch = () => {
    setSearchColumn("");
    setSearchValue("");
    setCurrentPage(1);
  };

  const tableColumns = useMemo(() => {
    if (!data?.data || data.data.length === 0) return [];

    const firstRow = data.data[0];
    return Object.keys(firstRow).map((key) => {
      const columnInfo = columns.find((col) => col.name === key);

      return {
        title: (
          <div>
            <div>{key}</div>
            {columnInfo && (
              <div
                style={{
                  fontSize: "12px",
                  fontWeight: "normal",
                  color: token.colorTextSecondary,
                }}
              >
                <Tag color={columnInfo.isPrimaryKey ? "gold" : "blue"}>
                  {columnInfo.type}
                </Tag>
                {columnInfo.isPrimaryKey && <Tag color="red">PK</Tag>}
                {!columnInfo.nullable && <Tag color="orange">NOT NULL</Tag>}
              </div>
            )}
          </div>
        ),
        dataIndex: key,
        key: key,
        width: 150,
        ellipsis: true,
        render: (value: any) => {
          if (value === null) return <Tag color="default">NULL</Tag>;
          if (typeof value === "boolean") return value ? "true" : "false";
          if (typeof value === "object") return JSON.stringify(value);
          return String(value);
        },
      };
    });
  }, [data, columns]);

  if (error) {
    return (
      <div className={styles.container}>
        <Alert
          message="Error loading table data"
          description={error}
          type="error"
          showIcon
          action={
            <Button onClick={loadTableData} icon={<ReloadOutlined />}>
              Retry
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.toolbar}>
        <div className={styles.filters}>
          <FilterOutlined />
          <Select
            placeholder="Select column"
            style={{ width: 150 }}
            value={searchColumn}
            onChange={setSearchColumn}
            allowClear
          >
            {columns.map((col) => (
              <Option key={col.name} value={col.name}>
                {col.name}
              </Option>
            ))}
          </Select>

          <Select
            value={searchOperator}
            onChange={setSearchOperator}
            style={{ width: 120 }}
          >
            <Option value="equals">Equals</Option>
            <Option value="contains">Contains</Option>
            <Option value="startsWith">Starts with</Option>
            <Option value="endsWith">Ends with</Option>
          </Select>

          <Input
            placeholder="Search value..."
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            onPressEnter={handleSearch}
            style={{ width: 200 }}
            suffix={<SearchOutlined />}
          />

          <Space>
            <Button onClick={handleSearch} type="primary">
              Search
            </Button>
            <Button onClick={handleClearSearch}>Clear</Button>
            <Button onClick={loadTableData} icon={<ReloadOutlined />}>
              Refresh
            </Button>
          </Space>
        </div>
      </div>

      <div className={styles.tableContainer}>
        <Spin spinning={loading}>
          {(!data || !data.data || data.data.length === 0) && !loading ? (
            <Flexbox
              align="center"
              justify="center"
              style={{ height: "300px", color: "rgba(0, 0, 0, 0.45)" }}
            >
              <div style={{ textAlign: "center" }}>
                <InboxOutlined style={{ fontSize: 48, marginBottom: 16 }} />
                <Title level={4} style={{ marginBottom: 8 }}>
                  No Data Found
                </Title>
                <Text type="secondary">
                  Try selecting a different table or adjusting your filters
                </Text>
              </div>
            </Flexbox>
          ) : (
            <Table
              columns={tableColumns}
              dataSource={data?.data || []}
              pagination={false}
              scroll={{ x: "max-content", y: "calc(100vh - 300px)" }}
              size="small"
              rowKey={(record, index) => index?.toString() || "0"}
            />
          )}
        </Spin>
      </div>

      {data?.pagination && (
        <div className={styles.pagination}>
          <Pagination
            current={currentPage}
            pageSize={pageSize}
            total={data.pagination.total}
            showSizeChanger
            showQuickJumper
            showTotal={(total, range) =>
              `${range[0]}-${range[1]} of ${total} items`
            }
            onChange={(page, size) => {
              setCurrentPage(page);
              if (size !== pageSize) {
                setPageSize(size);
              }
            }}
            pageSizeOptions={["10", "25", "50", "100", "200"]}
          />
        </div>
      )}
    </div>
  );
});

DataTable.displayName = "DataTable";

export default DataTable;
