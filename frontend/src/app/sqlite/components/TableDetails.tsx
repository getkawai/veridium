'use client';

import { Table, Tag, Typography, Spin, Alert, Card, Descriptions } from 'antd';
import { KeyOutlined, LinkOutlined } from '@ant-design/icons';
import { createStyles } from 'antd-style';
import { memo, useState, useEffect } from 'react';
import { GetTableDetails } from '@@/github.com/kawai-network/veridium/internal/tableviewer/service';
import type { TableColumnInfo } from '@@/github.com/kawai-network/veridium/internal/tableviewer/models';

const { Title, Text } = Typography;

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    padding: 16px;
    overflow-y: auto;
  `,
  summary: css`
    margin-bottom: 24px;
  `,
  columnsTable: css`
    .ant-table-thead > tr > th {
      background: ${token.colorFillAlter};
    }
  `,
}));

interface TableDetailsProps {
  tableName: string;
}

const TableDetails = memo<TableDetailsProps>(({ tableName }) => {
  const { styles } = useStyles();
  const [columns, setColumns] = useState<TableColumnInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (tableName) {
      loadTableDetails();
    }
  }, [tableName]);

  const loadTableDetails = async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await GetTableDetails(tableName);
      setColumns(result || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load table details');
    } finally {
      setLoading(false);
    }
  };

  const tableColumns = [
    {
      title: 'Column Name',
      dataIndex: 'name',
      key: 'name',
      width: 200,
      render: (name: string, record: TableColumnInfo) => (
        <div>
          <Text strong>{name}</Text>
          {record.isPrimaryKey && (
            <Tag icon={<KeyOutlined />} color="gold" style={{ marginLeft: 8 }}>
              PK
            </Tag>
          )}
          {record.foreignKey && (
            <Tag icon={<LinkOutlined />} color="blue" style={{ marginLeft: 8 }}>
              FK
            </Tag>
          )}
        </div>
      ),
    },
    {
      title: 'Data Type',
      dataIndex: 'type',
      key: 'type',
      width: 150,
      render: (type: string) => (
        <Tag color="blue">{type}</Tag>
      ),
    },
    {
      title: 'Nullable',
      dataIndex: 'nullable',
      key: 'nullable',
      width: 100,
      render: (nullable: boolean) => (
        <Tag color={nullable ? 'green' : 'red'}>
          {nullable ? 'YES' : 'NO'}
        </Tag>
      ),
    },
    {
      title: 'Default Value',
      dataIndex: 'defaultValue',
      key: 'defaultValue',
      width: 150,
      render: (defaultValue: string | null) => (
        defaultValue !== null ? (
          <Text code>{defaultValue}</Text>
        ) : (
          <Text type="secondary">NULL</Text>
        )
      ),
    },
    {
      title: 'Foreign Key',
      key: 'foreignKey',
      width: 200,
      render: (_, record: TableColumnInfo) => (
        record.foreignKey ? (
          <Text>
            {record.foreignKey.table}.{record.foreignKey.column}
          </Text>
        ) : (
          <Text type="secondary">-</Text>
        )
      ),
    },
  ];

  const primaryKeys = columns.filter(col => col.isPrimaryKey);
  const foreignKeys = columns.filter(col => col.foreignKey);
  const totalColumns = columns.length;

  if (loading) {
    return (
      <div className={styles.container}>
        <div style={{ textAlign: 'center', padding: '40px 0' }}>
          <Spin size="large" />
          <div style={{ marginTop: 16 }}>Loading table structure...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.container}>
        <Alert
          message="Error loading table structure"
          description={error}
          type="error"
          showIcon
        />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.summary}>
        <Title level={4}>Table Summary</Title>
        <Card>
          <Descriptions column={2} size="small">
            <Descriptions.Item label="Table Name">
              <Text strong>{tableName}</Text>
            </Descriptions.Item>
            <Descriptions.Item label="Total Columns">
              <Tag color="blue">{totalColumns}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Primary Keys">
              <Tag color="gold">{primaryKeys.length}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Foreign Keys">
              <Tag color="green">{foreignKeys.length}</Tag>
            </Descriptions.Item>
          </Descriptions>
        </Card>
      </div>

      <Title level={4}>Column Details</Title>
      <Table
        className={styles.columnsTable}
        columns={tableColumns}
        dataSource={columns}
        pagination={false}
        size="small"
        rowKey="name"
        scroll={{ x: 'max-content' }}
      />

      {primaryKeys.length > 0 && (
        <div style={{ marginTop: 24 }}>
          <Title level={5}>Primary Keys</Title>
          <div>
            {primaryKeys.map(col => (
              <Tag key={col.name} icon={<KeyOutlined />} color="gold">
                {col.name} ({col.type})
              </Tag>
            ))}
          </div>
        </div>
      )}

      {foreignKeys.length > 0 && (
        <div style={{ marginTop: 16 }}>
          <Title level={5}>Foreign Keys</Title>
          <div>
            {foreignKeys.map(col => (
              <Tag key={col.name} icon={<LinkOutlined />} color="blue">
                {col.name} → {col.foreignKey?.table}.{col.foreignKey?.column}
              </Tag>
            ))}
          </div>
        </div>
      )}
    </div>
  );
});

TableDetails.displayName = 'TableDetails';

export default TableDetails;