import { ShoppingCart } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import type { OTCContentProps } from './types';

const OTCContent = ({ styles, theme }: OTCContentProps) => {
  return (
    <Flexbox style={{ maxWidth: 700 }} gap={20}>
      <div>
        <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>OTC Market</h2>
        <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>P2P trading for KAWAI tokens</span>
      </div>

      <div className={styles.placeholderCard}>
        <ShoppingCart size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
        <h3 style={{ margin: '0 0 8px' }}>Coming Soon</h3>
        <p style={{ color: theme.colorTextSecondary, margin: 0 }}>
          Buy and sell KAWAI tokens directly with other users.<br />
          No slippage, atomic swaps via smart contract.
        </p>
      </div>
    </Flexbox>
  );
};

export default OTCContent;

