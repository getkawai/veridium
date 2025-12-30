'use client';

import { SearchOutlined, TableOutlined, EyeOutlined } from '@ant-design/icons';
import { Input, List, Typography, Badge, Spin, Alert } from 'antd';
import { createStyles } from 'antd-style';
import { memo, useState, useEffect } from 'react';
import { GetAllTables } from '@@/github.com/kawai-network/veridium/internal/tableviewer/service';
import type { TableBasicInfo } from '@@/github.com/kawai-network/veridium/internal/tableviewer/models';

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
}));

interface TableSidebarProps {
  selectedTable: string | null;
  onTableSelect: (tableName: string) => void;
}

const TableSidebar = memo<TableSidebarProps>(({ selectedTable, onTableSelect }) => {
  const { styles } = useStyles();
  const [tables, setTables] = useState<TableBasicInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    loadTables();
  }, []);

  const loadTables = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log('🔍 Starting to load tables...');
      
      // Check if the service is available
      if (typeof GetAllTables !== 'function') {
        throw new Error('GetAllTables service is not available');
      }
      
      const result = await GetAllTables();
      console.log('✅ Tables loaded successfully:', result);
      
      if (!result) {
        console.warn('⚠️ GetAllTables returned null/undefined');
        setTables([]);
        return;
      }
      
      if (!Array.isArray(result)) {
        console.warn('⚠️ GetAllTables returned non-array:', typeof result, result);
        setTables([]);
        return;
      }
      
      setTables(result);
      console.log(`📊 Loaded ${result.length} tables`);
    } catch (err) {
      console.error('❌ Error loading tables:', err);
      const errorMessage = err instanceof Error ? err.message : 'Failed to load tables';
      setError(errorMessage);
      
      // Additional debugging info
      if (err instanceof Error) {
        console.error('Error details:', {
          name: err.name,
          message: err.message,
          stack: err.stack
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const filteredTables = tables.filter(table =>
    table.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getTableIcon = (type: string) => {
    return type === 'VIEW' ? <EyeOutlined /> : <TableOutlined />;
  };

  const getTableTypeColor = (type: string) => {
    return type === 'VIEW' ? 'blue' : 'green';
  };

  if (loading) {
    return (
      <div className={styles.container}>
        <div style={{ textAlign: 'center', padding: '40px 0' }}>
          <Spin size="large" />
          <div style={{ marginTop: 16 }}>Loading tables...</div>
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
              <div style={{ marginTop: 8, fontSize: '12px', color: '#666' }}>
                Possible causes:
                <ul style={{ margin: '4px 0', paddingLeft: '16px' }}>
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
            <button onClick={loadTables} style={{ border: 'none', background: 'none', color: '#1890ff', cursor: 'pointer' }}>
              Retry
            </button>
          }
        />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          Database Tables
        </Typography.Title>
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
            className={selectedTable === table.name ? 'selected' : ''}
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
                  <Badge count="?" color="gray" />
                )}
              </div>
            </div>
          </List.Item>
        )}
      />
    </div>
  );
});

TableSidebar.displayName = 'TableSidebar';

export default TableSidebar;