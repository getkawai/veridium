import { memo, useEffect, useState } from 'react';
import { Spin, Empty, Button, Tooltip } from 'antd';
import { Copy } from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { useTheme } from 'antd-style';
import { App } from 'antd';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { getTokenListFromBackend } from '@/config/network';
import { StablecoinIcon } from '../StablecoinIcon';

interface AddTokenModalProps {
  currentNetwork: NetworkInfo | null;
  onClose: () => void;
}

interface ProjectToken {
  address: string;
  name: string;
  symbol: string;
}

export const AddTokenModal = memo<AddTokenModalProps>(({ currentNetwork, onClose }) => {
  const theme = useTheme();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(true);
  const [projectTokens, setProjectTokens] = useState<ProjectToken[]>([]);

  useEffect(() => {
    loadProjectTokens();
  }, [currentNetwork]);

  const loadProjectTokens = async () => {
    try {
      setLoading(true);
      const tokens = await getTokenListFromBackend();
      setProjectTokens(tokens);
    } catch (e) {
      console.error('Failed to load project tokens', e);
      setProjectTokens([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCopyAddress = async (address: string) => {
    try {
      await navigator.clipboard.writeText(address);
      message.success('Address copied!');
    } catch (err) {
      console.error('Failed to copy address:', err);
      message.error('Failed to copy address. Please copy manually.');
    }
  };

  if (loading) {
    return (
      <Flexbox align="center" justify="center" style={{ padding: 40 }}>
        <Spin />
      </Flexbox>
    );
  }

  return (
    <Flexbox gap={16}>
      <div style={{ color: theme.colorTextSecondary, fontSize: 13 }}>
        Supported tokens on {currentNetwork?.name || 'current network'}:
      </div>

      {projectTokens.length === 0 ? (
        <Empty description="No project tokens found" />
      ) : (
        <Flexbox gap={8}>
          {projectTokens.map((token) => (
            <div
              key={token.address}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '12px 16px',
                background: theme.colorFillTertiary,
                borderRadius: 12,
                border: `1px solid ${theme.colorBorderSecondary}`,
              }}
            >
              <Flexbox horizontal align="center" gap={12}>
                <div style={{
                  width: 36,
                  height: 36,
                  borderRadius: '50%',
                  background: token.symbol === 'USDT'
                    ? '#26a17b'
                    : 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontWeight: 700,
                  fontSize: 12,
                }}>
                  {token.symbol === 'USDT' || token.symbol === 'MockUSDT' || token.symbol === 'USDC'
                    ? <StablecoinIcon currentNetwork={currentNetwork} size={24} variant="branded" />
                    : token.symbol.substring(0, 2)}
                </div>
                <div>
                  <div style={{ fontWeight: 600 }}>{token.symbol}</div>
                  <div style={{ fontSize: 11, color: theme.colorTextSecondary }}>{token.name}</div>
                </div>
              </Flexbox>
              <Flexbox horizontal align="center" gap={8}>
                <Tooltip title={token.address}>
                  <span style={{ fontSize: 11, fontFamily: 'monospace', color: theme.colorTextTertiary }}>
                    {token.address.substring(0, 6)}...{token.address.substring(token.address.length - 4)}
                  </span>
                </Tooltip>
                <ActionIcon
                  icon={Copy}
                  size="small"
                  onClick={() => handleCopyAddress(token.address)}
                  title="Copy address"
                />
              </Flexbox>
            </div>
          ))}
        </Flexbox>
      )}

      <div style={{
        padding: 12,
        background: theme.colorInfoBg,
        borderRadius: 8,
        fontSize: 12,
        color: theme.colorTextSecondary
      }}>
        <strong>Note:</strong> These are the official project tokens for {currentNetwork?.name || 'the current network'}.
        Custom token import is not supported yet.
      </div>

      <Button block onClick={onClose}>Close</Button>
    </Flexbox>
  );
});

AddTokenModal.displayName = 'AddTokenModal';
