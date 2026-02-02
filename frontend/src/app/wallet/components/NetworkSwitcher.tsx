import { memo, useEffect, useState } from 'react';
import { Popover, Spin } from 'antd';
import { Check, Globe } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { useTheme } from 'antd-style';
import { JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { NetworkInfo, BackendConfig } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { NetworkIcon } from '../NetworkIcons';

interface NetworkSwitcherProps {
  currentNetwork: NetworkInfo | null;
  onNetworkChange: (network: NetworkInfo) => void;
  backendConfig: BackendConfig | null;
}

export const NetworkSwitcher = memo<NetworkSwitcherProps>(({ currentNetwork, onNetworkChange, backendConfig }) => {
  const theme = useTheme();
  const [networks, setNetworks] = useState<NetworkInfo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadNetworks();
  }, [backendConfig]);

  const loadNetworks = async () => {
    try {
      const supportedNetworks = await JarvisService.GetSupportedNetworks();
      if (backendConfig) {
        const targetChainId = backendConfig.network.chainId;
        const monadNetworks = supportedNetworks.filter(n => n.id === targetChainId);

        console.log(`Filtering networks for ${backendConfig.environment}:`, {
          targetChainId,
          filtered: monadNetworks.map(n => ({ id: n.id, name: n.name })),
        });

        setNetworks(monadNetworks);
      } else {
        setNetworks([]);
      }
    } catch (e) {
      console.error('Failed to load networks', e);
      setNetworks([]);
    } finally {
      setLoading(false);
    }
  };

  const displayName = currentNetwork?.name || 'Select Network';

  return (
    <Popover
      arrow={false}
      content={
        <div style={{ width: 220, maxHeight: 300, overflowY: 'auto' }}>
          <div style={{
            padding: '4px 8px',
            fontSize: 11,
            color: theme.colorTextTertiary,
            textTransform: 'uppercase',
            position: 'sticky',
            top: 0,
            background: theme.colorBgElevated
          }}>
            Switch Network
          </div>
          {loading ? (
            <Flexbox align="center" justify="center" style={{ padding: 20 }}>
              <Spin size="small" />
            </Flexbox>
          ) : (
            networks.map((network) => (
              <Flexbox
                key={network.id}
                horizontal
                align="center"
                gap={12}
                onClick={() => onNetworkChange(network)}
                style={{
                  padding: '8px 12px',
                  cursor: 'pointer',
                  borderRadius: 8,
                  background: currentNetwork?.id === network.id ? theme.colorFillSecondary : 'transparent',
                  transition: 'background 0.2s'
                }}
              >
                <span>
                  <NetworkIcon name={network.icon || 'ethereum'} size={24} variant="branded" />
                </span>
                <Flexbox flex={1}>
                  <span style={{
                    fontSize: 13,
                    fontWeight: currentNetwork?.id === network.id ? 600 : 400
                  }}>
                    {network.name}
                  </span>
                  <span style={{ fontSize: 10, color: theme.colorTextTertiary }}>
                    {network.nativeTokenSymbol}
                  </span>
                </Flexbox>
                {currentNetwork?.id === network.id && <Check size={14} color={theme.colorSuccess} />}
              </Flexbox>
            ))
          )}
        </div>
      }
      placement="rightBottom"
      trigger="click"
    >
      <div style={{
        padding: '8px 12px',
        background: theme.colorFillSecondary,
        borderRadius: 12,
        cursor: 'pointer',
        border: `1px solid ${theme.colorBorderSecondary}`,
        display: 'flex',
        alignItems: 'center',
        gap: 10,
        transition: 'all 0.2s ease'
      }}
        onMouseEnter={(e) => e.currentTarget.style.borderColor = theme.colorPrimary}
        onMouseLeave={(e) => e.currentTarget.style.borderColor = theme.colorBorderSecondary}
      >
        <div style={{
          width: 8,
          height: 8,
          borderRadius: '50%',
          background: theme.colorSuccess,
          boxShadow: `0 0 8px ${theme.colorSuccess}80`
        }} />
        <Flexbox flex={1}>
          <div style={{ fontSize: 10, color: theme.colorTextTertiary, lineHeight: 1 }}>Network</div>
          <div style={{
            fontSize: 12,
            fontWeight: 600,
            display: 'flex',
            alignItems: 'center',
            gap: 4
          }}>
            {currentNetwork ? (
              <NetworkIcon name={currentNetwork.icon || 'ethereum'} size={20} variant="branded" />
            ) : (
              <NetworkIcon name="wallet-connect" size={20} />
            )} {displayName}
          </div>
        </Flexbox>
        <Globe size={14} style={{ opacity: 0.5 }} />
      </div>
    </Popover>
  );
});

NetworkSwitcher.displayName = 'NetworkSwitcher';
