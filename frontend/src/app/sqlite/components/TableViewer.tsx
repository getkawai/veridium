'use client';

import { Tabs, Typography, Empty } from 'antd';
import { createStyles } from 'antd-style';
import { memo } from 'react';
import { DataTable, TableDetails, QueryEditor } from './index';

const { Title } = Typography;

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 16px;
  `,
  header: css`
    margin-bottom: 16px;
    border-bottom: 1px solid ${token.colorBorder};
    padding-bottom: 16px;
  `,
  content: css`
    flex: 1;
    overflow: hidden;
    
    .ant-tabs {
      height: 100%;
      
      .ant-tabs-content-holder {
        height: calc(100% - 46px);
        
        .ant-tabs-tabpane {
          height: 100%;
          overflow: hidden;
        }
      }
    }
  `,
  emptyState: css`
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  `,
}));

interface TableViewerProps {
  selectedTable: string | null;
}

const TableViewer = memo<TableViewerProps>(({ selectedTable }) => {
  const { styles } = useStyles();

  if (!selectedTable) {
    return (
      <div className={styles.emptyState}>
        <Empty
          description="Select a table from the sidebar to view its data"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      </div>
    );
  }

  const tabItems = [
    {
      key: 'data',
      label: 'Data',
      children: <DataTable tableName={selectedTable} />,
    },
    {
      key: 'structure',
      label: 'Structure',
      children: <TableDetails tableName={selectedTable} />,
    },
    {
      key: 'query',
      label: 'Query',
      children: <QueryEditor initialTable={selectedTable} />,
    },
  ];

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title level={3} style={{ margin: 0 }}>
          {selectedTable}
        </Title>
      </div>

      <div className={styles.content}>
        <Tabs
          defaultActiveKey="data"
          items={tabItems}
          size="large"
        />
      </div>
    </div>
  );
});

TableViewer.displayName = 'TableViewer';

export default TableViewer;