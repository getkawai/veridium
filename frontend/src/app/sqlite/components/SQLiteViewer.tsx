'use client';

import { useSize } from 'ahooks';
import { DraggablePanel, DraggablePanelContainer } from '@lobehub/ui';
import { createStyles, useResponsive } from 'antd-style';
import isEqual from 'fast-deep-equal';
import { memo, useEffect, useRef, useState } from 'react';
import TableSidebar from './TableSidebar';
import TableViewer from './TableViewer';

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100vh;
    display: flex;
    background: ${token.colorBgContainer};
  `,
  panel: css`
    height: 100%;
    background: ${token.colorBgContainerSecondary};
    border-right: 1px solid ${token.colorBorder};
  `,
  mainContent: css`
    flex: 1;
    height: 100%;
    background: ${token.colorBgContainer};
    overflow: hidden;
  `,
}));

const SQLiteViewer = memo(() => {
  const { md = true } = useResponsive();
  const { styles } = useStyles();

  const [panelWidth, setPanelWidth] = useState(280);
  const [expand, setExpand] = useState(true);
  const [selectedTable, setSelectedTable] = useState<string | null>(null);

  const innerRef = useRef(null);
  useSize(innerRef);

  const handleExpand = (newExpand: boolean) => {
    setExpand(newExpand);
  };

  useEffect(() => {
    if (!md) setExpand(false);
  }, [md]);

  const handleSizeChange = (_, s) => {
    if (!s) return;
    const nextWidth = typeof s.width === 'string' ? Number.parseInt(s.width) : s.width;
    if (!nextWidth || isEqual(nextWidth, panelWidth)) return;
    setPanelWidth(nextWidth);
  };

  return (
    <div className={styles.container}>
      <DraggablePanel
        className={styles.panel}
        defaultSize={{ width: panelWidth }}
        expand={expand}
        maxWidth={400}
        minWidth={200}
        mode={md ? 'fixed' : 'float'}
        onExpandChange={handleExpand}
        onSizeChange={handleSizeChange}
        placement="left"
        size={{ height: '100%', width: panelWidth }}
      >
        <div ref={innerRef} style={{ height: '100%', width: '100%' }}>
          <DraggablePanelContainer
            style={{
              flex: 'none',
              height: '100%',
              minWidth: 200,
              overflow: 'hidden',
            }}
          >
            <TableSidebar 
              selectedTable={selectedTable}
              onTableSelect={setSelectedTable}
            />
          </DraggablePanelContainer>
        </div>
      </DraggablePanel>
      
      <div className={styles.mainContent}>
        <TableViewer selectedTable={selectedTable} />
      </div>
    </div>
  );
});

SQLiteViewer.displayName = 'SQLiteViewer';

export default SQLiteViewer;