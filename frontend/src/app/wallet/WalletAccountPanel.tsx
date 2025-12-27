import { useSize } from 'ahooks';
import { DraggablePanel, DraggablePanelContainer, type DraggablePanelProps } from '@lobehub/ui';
import { createStyles, useResponsive } from 'antd-style';
import isEqual from 'fast-deep-equal';
import { PropsWithChildren, memo, useEffect, useRef, useState } from 'react';

export const useStyles = createStyles(({ css, token }) => ({
  panel: css`
    height: 100%;
    background: ${token.colorBgContainerSecondary};
  `,
}));

const WalletAccountPanel = memo<PropsWithChildren>(({ children }) => {
  const { md = true } = useResponsive();
  const { styles } = useStyles();

  const [panelWidth, setPanelWidth] = useState(80);
  const [expand, setExpand] = useState(true);

  const innerRef = useRef(null);
  useSize(innerRef);

  const handleExpand = (newExpand: boolean) => {
    setExpand(newExpand);
  };

  useEffect(() => {
    if (!md) setExpand(false);
  }, [md]);

  const handleSizeChange: DraggablePanelProps['onSizeChange'] = (_, s) => {
    if (!s) return;
    const nextWidth = typeof s.width === 'string' ? Number.parseInt(s.width) : s.width;
    if (!nextWidth || isEqual(nextWidth, panelWidth)) return;
    setPanelWidth(nextWidth);
  };

  return (
    <DraggablePanel
      className={styles.panel}
      defaultSize={{ width: panelWidth }}
      expand={expand}
      maxWidth={320}
      minWidth={80}
      mode={md ? 'fixed' : 'float'}
      onExpandChange={handleExpand}
      onSizeChange={handleSizeChange}
      placement="right"
      size={{ height: '100%', width: panelWidth }}
    >
      <div ref={innerRef} style={{ height: '100%', width: '100%' }}>
        <DraggablePanelContainer
          style={{
            flex: 'none',
            height: '100%',
            minWidth: 80,
            overflow: 'hidden',
          }}
        >
          {children}
        </DraggablePanelContainer>
      </div>
    </DraggablePanel>
  );
});

WalletAccountPanel.displayName = 'WalletAccountPanel';

export default WalletAccountPanel;
