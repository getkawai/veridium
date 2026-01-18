"use client";

import { SearchOutlined, TableOutlined, EyeOutlined } from "@ant-design/icons";
import { Input, List, Typography, Badge, Spin, Alert, Button } from "antd";
import { createStyles } from "antd-style";
import { memo, useState, useEffect, useCallback } from "react";
import { GetAllTables } from "@@/github.com/kawai-network/veridium/internal/tableviewer/service";
import type { TableBasicInfo } from "@@/github.com/kawai-network/veridium/internal/tableviewer/models";

const { Text } = Typography;

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 16px;
  `,
  header: css`
    margin-bottom: 16px;
  `,
  searchBox: css`
    margin-bottom: 16px;
  `,
  tableList: css`
    flex: 1;
    overflow-y: auto;

    .ant-list-item {
      padding: 8px 12px;
      cursor: pointer;
      border-radius: 6px;
      margin-bottom: 4px;
      transition: all 0.2s;

      &:hover {
        background: ${token.colorFillTertiary};
      }

      &.selected {
        background: ${token.colorPrimaryBg};
        border: 1px solid ${token.colorPrimary};
      }
    }
  `,
  tableItem: css`
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  `,
  tableInfo: css`
    display: flex;
    align-items: center;
    gap: 8px;
  `,
  tableIcon: css`
    color: ${token.colorTextSecondary};
  `,
  tableName: css`
    font-weight: 500;
  `,
  tableType: css`
    font-size: 12px;
    color: ${token.colorTextTertiary};
  `,
  rowCount: css`
    font-size: 12px;
  `,
  privacySubtitle: css`
    font-size: 12px;
    color: ${token.colorSuccessText};
    margin-top: 4px;
  `,
}));

interface TableSidebarProps {
  selectedTable: string | null;
  onTableSelect: (tableName: string) => void;
}

const TableSidebar = memo<TableSidebarProps>(
  ({ selectedTable, onTableSelect }) => {
    const { styles, theme } = useStyles();
    const [tables, setTables] = useState<TableBasicInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [searchTerm, setSearchTerm] = useState("");

    const loadTables = useCallback(async () => {
      try {
        setLoading(true);
        setError(null);
        console.log("🔍 Starting to load tables...");

        // Check if the service is available
        if (typeof GetAllTables !== "function") {
          throw new Error("GetAllTables service is not available");
        }

        const result = await GetAllTables();
        console.log("✅ Tables loaded successfully:", result);

        if (!result) {
          console.warn("⚠️ GetAllTables returned null/undefined");
          setTables([]);
          return;
        }

        if (!Array.isArray(result)) {
          console.warn(
            "⚠️ GetAllTables returned non-array:",
            typeof result,
            result,
          );
          setTables([]);
          return;
        }

        setTables(result);
        console.log(`📊 Loaded ${result.length} tables`);
      } catch (err) {
        console.error("❌ Error loading tables:", err);
        const errorMessage =
          err instanceof Error ? err.message : "Failed to load tables";
        setError(errorMessage);

        // Additional debugging info
        if (err instanceof Error) {
          console.error("Error details:", {
            name: err.name,
            message: err.message,
            stack: err.stack,
          });
        }
      } finally {
        setLoading(false);
      }
    }, []);

    useEffect(() => {
      loadTables();
    }, [loadTables]);

    const filteredTables = tables.filter((table) =>
      table.name.toLowerCase().includes(searchTerm.toLowerCase()),
    );

    const getTableIcon = (type: string) => {
      return type === "VIEW" ? <EyeOutlined /> : <TableOutlined />;
    };

    const getTableTypeColor = (type: string) => {
      return type === "VIEW" ? "blue" : "green";
    };

    if (loading) {
      return (
        <div className={styles.container}>
          <div style={{ textAlign: "center", padding: "40px 0" }}>
            <Spin size="large" />
            <div style={{ marginTop: 16 }}>Loading tables...</div>
          </div>
        </div>
      );
    }

    // Empty state when no tables found
    if (!loading && tables.length === 0) {
      return (
        <div className={styles.container}>
          <div style={{ textAlign: "center", padding: "40px 20px" }}>
            <TableOutlined
              style={{
                fontSize: 64,
                marginBottom: 16,
                color: "rgba(0, 0, 0, 0.25)",
              }}
            />
            <Typography.Title level={4} style={{ marginBottom: 8 }}>
              No Tables Found
            </Typography.Title>
            <Text type="secondary">
              The database appears to be empty or the database file is not
              accessible.
              <br />
              <br />
              Possible causes:
              <ul style={{ textAlign: "left", marginTop: 16, paddingLeft: 24 }}>
                <li>Database file does not exist at data/veridium.db</li>
                <li>Database is locked by another process</li>
                <li>Permission issues accessing the database</li>
                <li>Database is empty (no tables created yet)</li>
              </ul>
            </Text>
          </div>
        </div>
      );
    }

    if (error) {
      return (
        <div className={styles.container}>
          <Alert
            message="Error"
            description={
              <div>
                <div>{error}</div>
                <div style={{ marginTop: 8, fontSize: "12px", color: theme.colorTextSecondary }}>
                  Possible causes:
                  <ul style={{ margin: "4px 0", paddingLeft: "16px" }}>
                    <li>Database file not found at data/veridium.db</li>
                    <li>Database service not running</li>
                    <li>Permission issues accessing database</li>
                  </ul>
                </div>
              </div>
            }
            type="error"
            showIcon
            action={
              <Button type="link" onClick={loadTables}>
                Retry
              </Button>
            }
          />
        </div>
      );
    }

    return (
      <div className={styles.container}>
        <div className={styles.header}>
          <div>
            <Typography.Title level={4} style={{ margin: 0 }}>
              Database Tables
            </Typography.Title>
            <div className={styles.privacySubtitle}>
              🔒 <Text>Local Database • Privacy First</Text>
            </div>
          </div>
          <Text type="secondary">{tables.length} tables found</Text>
        </div>

        <Input
          className={styles.searchBox}
          placeholder="Search tables..."
          prefix={<SearchOutlined />}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          allowClear
        />

        <List
          className={styles.tableList}
          dataSource={filteredTables}
          renderItem={(table) => (
            <List.Item
              className={selectedTable === table.name ? "selected" : ""}
              onClick={() => onTableSelect(table.name)}
            >
              <div className={styles.tableItem}>
                <div className={styles.tableInfo}>
                  <span className={styles.tableIcon}>
                    {getTableIcon(table.type)}
                  </span>
                  <div>
                    <div className={styles.tableName}>{table.name}</div>
                    <div className={styles.tableType}>
                      <Badge
                        color={getTableTypeColor(table.type)}
                        text={table.type}
                        size="small"
                      />
                    </div>
                  </div>
                </div>
                <div className={styles.rowCount}>
                  {table.count > 0 ? (
                    <Badge count={table.count} showZero color="blue" />
                  ) : (
                    <Badge count={table.count} showZero color="gray" />
                  )}
                </div>
              </div>
            </List.Item>
          )}
        />
      </div>
    );
  },
);

TableSidebar.displayName = "TableSidebar";

export default TableSidebar;
