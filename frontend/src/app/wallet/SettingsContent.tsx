import { Card, Tag, App } from 'antd';
import { ExternalLink, Copy } from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import type { SettingsContentProps } from './types';

const CopyButton = ({ text }: { text: string }) => {
  const { message } = App.useApp();
  return (
    <ActionIcon
      icon={Copy}
      size="small"
      onClick={() => {
        navigator.clipboard.writeText(text);
        message.success("Copied!");
      }}
      title="Copy"
    />
  );
};

const SettingsContent = ({ address, styles, theme, currentNetwork }: SettingsContentProps) => {
  return (
    <Flexbox style={{ maxWidth: 700 }} gap={20}>
      <div>
        <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Settings</h2>
        <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>Network and wallet configuration</span>
      </div>

      <Card title="Active Session" size="small">
        <Flexbox gap={12}>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Network</span>
            <Tag color={currentNetwork?.isTestnet ? 'orange' : 'green'}>
              {currentNetwork?.name || 'Not Connected'}
            </Tag>
          </Flexbox>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Status</span>
            <span style={{ color: theme.colorSuccess, fontWeight: 600, fontSize: 12 }}>Connected</span>
          </Flexbox>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Native Token</span>
            <span style={{ fontWeight: 600 }}>{currentNetwork?.nativeTokenSymbol || '-'}</span>
          </Flexbox>
        </Flexbox>
      </Card>

      <Card title="Technical Details" size="small">
        <Flexbox gap={12}>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Chain ID</span>
            <span style={{ fontFamily: 'monospace' }}>{currentNetwork?.id || '-'}</span>
          </Flexbox>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Native Token Decimals</span>
            <span style={{ fontFamily: 'monospace' }}>{currentNetwork?.nativeTokenDecimal || '-'}</span>
          </Flexbox>
          {currentNetwork?.explorerURL && (
            <Flexbox horizontal justify="space-between">
              <span style={{ color: theme.colorTextSecondary }}>Explorer API</span>
              <a href={currentNetwork.explorerURL} target="_blank" rel="noopener noreferrer" style={{ display: 'flex', alignItems: 'center', gap: 4, fontSize: 12 }}>
                View <ExternalLink size={12} />
              </a>
            </Flexbox>
          )}
        </Flexbox>
      </Card>

      <Card title="Wallet" size="small">
        <Flexbox gap={12}>
          <Flexbox horizontal justify="space-between" align="center">
            <span style={{ color: theme.colorTextSecondary }}>Address</span>
            <Flexbox horizontal align="center" gap={8}>
              <span style={{ fontFamily: 'monospace', fontSize: 12 }}>
                {address ? `${address.substring(0, 10)}...${address.substring(address.length - 8)}` : 'Loading...'}
              </span>
              <CopyButton text={address} />
            </Flexbox>
          </Flexbox>
        </Flexbox>
      </Card>
    </Flexbox>
  );
};

export default SettingsContent;

