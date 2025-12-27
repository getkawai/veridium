'use client';

import { DraggablePanel, DraggablePanelContainer, type DraggablePanelProps } from '@lobehub/ui';
import { createStyles, useResponsive } from 'antd-style';
import isEqual from 'fast-deep-equal';
import { PropsWithChildren, memo, useEffect, useState } from 'react';

const WALLET_PANEL_WIDTH = 220;
const WALLET_PANEL_MIN_WIDTH = 180;

export const useStyles = createStyles(({ css, token }) => ({
  panel: css`
    height: 100%;
    background: ${token.colorBgLayout};
  `,
}));

interface WalletSidePanelProps extends PropsWithChildren {
  defaultWidth?: number;
  showPanel?: boolean;
  onExpandChange?: (expand: boolean) => void;
}

const WalletSidePanel = memo<WalletSidePanelProps>(({
  children,
  defaultWidth = WALLET_PANEL_WIDTH,
  showPanel = true,
  onExpandChange
}) => {
  const { md = true } = useResponsive();
  const { styles } = useStyles();

  const [panelWidth, setPanelWidth] = useState(defaultWidth);
  const [expand, setExpand] = useState(showPanel);

  const handleExpand = (newExpand: boolean) => {
    if (isEqual(newExpand, expand)) return;
    setExpand(newExpand);
    onExpandChange?.(newExpand);
  };

  useEffect(() => {
    // Auto-collapse on mobile
    if (!md) {
      setExpand(false);
    }
  }, [md]);

  const handleSizeChange: DraggablePanelProps['onSizeChange'] = (_, size) => {
    if (!size) return;
    const nextWidth = typeof size.width === 'string' ? Number.parseInt(size.width) : size.width;
    if (!nextWidth || isEqual(nextWidth, panelWidth)) return;
    setPanelWidth(nextWidth);
  };

  return (
    <DraggablePanel
      className={styles.panel}
      defaultSize={{ width: panelWidth }}
      expand={expand}
      maxWidth={280}
      minWidth={WALLET_PANEL_MIN_WIDTH}
      mode={md ? 'fixed' : 'float'}
      onExpandChange={handleExpand}
      onSizeChange={handleSizeChange}
      placement="left"
      size={{ height: '100%', width: panelWidth }}
    >
      <DraggablePanelContainer
        style={{
          flex: 'none',
          height: '100%',
          minWidth: WALLET_PANEL_MIN_WIDTH,
          overflow: 'hidden',
        }}
      >
        {children}
      </DraggablePanelContainer>
    </DraggablePanel>
  );
});

WalletSidePanel.displayName = 'WalletSidePanel';

export default WalletSidePanel;
