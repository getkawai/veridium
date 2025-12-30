'use client';

import { Button, Input, Alert, Table, Typography, Space, Spin, message } from 'antd';
import { PlayCircleOutlined, ClearOutlined, HistoryOutlined } from '@ant-design/icons';
import { createStyles } from 'antd-style';
import { memo, useState, useEffect } from 'react';
import { ExecuteRawQuery } from '@@/github.com/kawai-network/veridium/internal/tableviewer/service';

const { TextArea } = Input;
const { Title, Text } = Typography;

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 16px;
  `,
  editor: css`
    margin-bottom: 16px;
  `,
  toolbar: css`
    margin-bottom: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  `,
  results: css`
    flex: 1;
    overflow: auto;
  `,
  queryHistory: css`
    max-height: 200px;
    overflow-y: auto;
    border: 1px solid ${token.colorBorder};
    border-radius: 6px;
    padding: 8px;
    margin-top: 16px;
  `,
  historyItem: css`
    padding: 4px 8px;
    cursor: pointer;
    border-radius: 4px;
    margin-bottom: 4px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
    
    &:hover {
      background: ${token.colorFillTertiary};
    }
  `,
}));

interface QueryEditorProps {
  initialTable?: string;
}

const QueryEditor = memo<QueryEditorProps>(({ initialTable }) => {
  const { styles } = useStyles();

  const [query, setQuery] = useState('');
  const [results, setResults] = useState<any[]>([]);
  const [columns, setColumns] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [queryHistory, setQueryHistory] = useState<string[]>([]);

  useEffect(() => {
    if (initialTable) {
      setQuery(`SELECT * FROM ${initialTable} LIMIT 10;`);
    }
  }, [initialTable]);

  useEffect(() => {
    // Load query history from localStorage
    const saved = localStorage.getItem('sqlite-query-history');
    if (saved) {
      try {
        setQueryHistory(JSON.parse(saved));
      } catch (e) {
        console.error('Failed to parse query history:', e);
      }
    }
  }, []);

  const saveQueryToHistory = (sql: string) => {
    const trimmed = sql.trim();
    if (!trimmed) return;

    const newHistory = [trimmed, ...queryHistory.filter(q => q !== trimmed)].slice(0, 10);
    setQueryHistory(newHistory);
    localStorage.setItem('sqlite-query-history', JSON.stringify(newHistory));
  };

  const executeQuery = async () => {
    if (!query.trim()) {
      message.warning('Please enter a SQL query');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setResults([]);
      setColumns([]);

      const result = await ExecuteRawQuery(query.trim(), []);

      // Parse JSON result
      const parsedResult = JSON.parse(result);

      if (Array.isArray(parsedResult) && parsedResult.length > 0) {
        // Generate columns from first row
        const firstRow = parsedResult[0];
        const tableColumns = Object.keys(firstRow).map(key => ({
          title: key,
          dataIndex: key,
          key: key,
          ellipsis: true,
          render: (value: any) => {
            if (value === null) return <Text type="secondary">NULL</Text>;
            if (typeof value === 'boolean') return value ? 'true' : 'false';
            if (typeof value === 'object') return JSON.stringify(value);
            return String(value);
          },
        }));

        setColumns(tableColumns);
        setResults(parsedResult);
        message.success(`Query executed successfully. ${parsedResult.length} rows returned.`);
      } else {
        message.success('Query executed successfully. No rows returned.');
      }

      saveQueryToHistory(query.trim());

    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to execute query';
      setError(errorMessage);
      message.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const clearQuery = () => {
    setQuery('');
    setResults([]);
    setColumns([]);
    setError(null);
  };

  const loadFromHistory = (historyQuery: string) => {
    setQuery(historyQuery);
  };

  const sampleQueries = [
    `SELECT * FROM ${initialTable || 'table_name'} LIMIT 10;`,
    `SELECT COUNT(*) as total FROM ${initialTable || 'table_name'};`,
    `PRAGMA table_info(${initialTable || 'table_name'});`,
    'SELECT name FROM sqlite_master WHERE type="table";',
  ];

  return (
    <div className={styles.container}>
      <div className={styles.editor}>
        <Title level={5}>SQL Query Editor</Title>
        <TextArea
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Enter your SQL query here..."
          rows={6}
          style={{ fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace' }}
        />
      </div>

      <div className={styles.toolbar}>
        <Space>
          <Button
            type="primary"
            icon={<PlayCircleOutlined />}
            onClick={executeQuery}
            loading={loading}
          >
            Execute Query
          </Button>
          <Button
            icon={<ClearOutlined />}
            onClick={clearQuery}
          >
            Clear
          </Button>
        </Space>

        <Text type="secondary">
          {results.length > 0 && `${results.length} rows returned`}
        </Text>
      </div>

      {error && (
        <Alert
          message="Query Error"
          description={error}
          type="error"
          showIcon
          style={{ marginBottom: 16 }}
        />
      )}

      <div className={styles.results}>
        {loading ? (
          <div style={{ textAlign: 'center', padding: '40px 0' }}>
            <Spin size="large" />
            <div style={{ marginTop: 16 }}>Executing query...</div>
          </div>
        ) : results.length > 0 ? (
          <Table
            columns={columns}
            dataSource={results}
            pagination={{
              pageSize: 50,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} items`,
            }}
            scroll={{ x: 'max-content', y: 400 }}
            size="small"
            rowKey={(record, index) => index?.toString() || '0'}
          />
        ) : (
          <div>
            <Title level={5}>Sample Queries</Title>
            <Space direction="vertical" style={{ width: '100%' }}>
              {sampleQueries.map((sample, index) => (
                <Button
                  key={index}
                  type="dashed"
                  block
                  style={{ textAlign: 'left', fontFamily: 'Monaco, Menlo, Ubuntu Mono, monospace' }}
                  onClick={() => setQuery(sample)}
                >
                  {sample}
                </Button>
              ))}
            </Space>
          </div>
        )}
      </div>

      {queryHistory.length > 0 && (
        <div className={styles.queryHistory}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
            <HistoryOutlined style={{ marginRight: 8 }} />
            <Text strong>Query History</Text>
          </div>
          {queryHistory.map((historyQuery, index) => (
            <div
              key={index}
              className={styles.historyItem}
              onClick={() => loadFromHistory(historyQuery)}
            >
              {historyQuery}
            </div>
          ))}
        </div>
      )}
    </div>
  );
});

QueryEditor.displayName = 'QueryEditor';

export default QueryEditor;