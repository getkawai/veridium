import { useState } from 'react';
import { Form, Input, InputNumber, Button, Select } from 'antd';
import { Send } from 'lucide-react';
import { Icon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { useTheme } from 'antd-style';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';

interface SendFormProps {
  onSend: (to: string, amount: number, asset: string, customAddr?: string) => void;
  loading: boolean;
  currentNetwork?: NetworkInfo | null;
}

interface PendingTx {
  to: string;
  amount: number;
  assetType: string;
  customAddr?: string;
}

export const SendForm = ({ onSend, loading, currentNetwork }: SendFormProps) => {
  const [form] = Form.useForm();
  const [selectedAsset, setSelectedAsset] = useState('usdt');
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [pendingTx, setPendingTx] = useState<PendingTx | null>(null);
  const theme = useTheme();

  const assetOptions = [
    { label: `Native Token (${currentNetwork?.nativeTokenSymbol || 'ETH'})`, value: 'native' },
    { label: currentNetwork?.stablecoinSymbol === 'USDC' ? 'USDC (USD Coin)' : 'USDT (Tether)', value: 'usdt' },
    { label: 'KAWAI (Kawai Token)', value: 'kawai' },
  ];

  const getAssetLabel = (assetType: string) => {
    if (assetType === 'native') return currentNetwork?.nativeTokenSymbol || 'ETH';
    if (assetType === 'usdt') return currentNetwork?.stablecoinSymbol || 'USDT';
    if (assetType === 'kawai') return 'KAWAI';
    return assetType.toUpperCase();
  };

  const handleFinish = (values: any) => {
    let assetType = values.asset || selectedAsset;
    let customAddr: string | undefined = undefined;

    if (values.asset === 'custom') {
      customAddr = values.customAddr;
      assetType = 'custom';
    }

    setPendingTx({ to: values.to, amount: values.amount, assetType, customAddr });
    setShowConfirmation(true);
  };

  const handleConfirm = () => {
    if (pendingTx) {
      onSend(pendingTx.to, pendingTx.amount, pendingTx.assetType, pendingTx.customAddr);
      setShowConfirmation(false);
    }
  };

  if (showConfirmation && pendingTx) {
    return (
      <Flexbox gap={16}>
        <div style={{ textAlign: 'center', marginBottom: 16 }}>
          <Icon icon={Send} size={48} style={{ color: theme.colorPrimary, marginBottom: 8 }} />
          <div style={{ fontSize: 18, fontWeight: 600 }}>Confirm Transaction</div>
        </div>

        <div style={{ background: theme.colorFillTertiary, borderRadius: 12, padding: 16 }}>
          <Flexbox gap={12}>
            <Flexbox horizontal justify="space-between">
              <span style={{ color: theme.colorTextSecondary }}>To</span>
              <span style={{ fontFamily: 'monospace', fontSize: 12 }}>
                {pendingTx.to.substring(0, 10)}...{pendingTx.to.substring(pendingTx.to.length - 8)}
              </span>
            </Flexbox>
            <Flexbox horizontal justify="space-between">
              <span style={{ color: theme.colorTextSecondary }}>Amount</span>
              <span style={{ fontWeight: 600 }}>{pendingTx.amount} {getAssetLabel(pendingTx.assetType)}</span>
            </Flexbox>
            <Flexbox horizontal justify="space-between">
              <span style={{ color: theme.colorTextSecondary }}>Network</span>
              <span>{currentNetwork?.name || 'Monad Testnet'}</span>
            </Flexbox>
          </Flexbox>
        </div>

        <div style={{ fontSize: 12, color: theme.colorTextTertiary, textAlign: 'center' }}>
          Please review transaction details before confirming.
        </div>

        <Flexbox horizontal gap={8}>
          <Button block size="large" onClick={() => setShowConfirmation(false)}>
            Cancel
          </Button>
          <Button type="primary" block size="large" loading={loading} onClick={handleConfirm}>
            Confirm & Send
          </Button>
        </Flexbox>
      </Flexbox>
    );
  }

  return (
    <Form form={form} layout="vertical" onFinish={handleFinish} initialValues={{ asset: 'usdt' }}>
      <Form.Item label="Asset" name="asset">
        <Select
          options={[
            ...assetOptions,
            { label: 'Custom Token', value: 'custom' }
          ]}
          value={selectedAsset}
          onChange={setSelectedAsset}
          size="large"
        />
      </Form.Item>

      {selectedAsset === 'custom' && (
        <Form.Item
          label="Token Contract Address"
          name="customAddr"
          rules={[
            { required: true, message: 'Token address is required' },
            { pattern: /^0x[a-fA-F0-9]{40}$/, message: 'Invalid EVM address' }
          ]}
        >
          <Input placeholder="0x..." size="large" />
        </Form.Item>
      )}

      <Form.Item
        label="Recipient Address"
        name="to"
        rules={[
          { required: true, message: 'Address is required' },
          { pattern: /^0x[a-fA-F0-9]{40}$/, message: 'Invalid EVM address' }
        ]}
      >
        <Input placeholder="0x..." size="large" />
      </Form.Item>

      <Form.Item
        label="Amount"
        name="amount"
        rules={[{ required: true, type: 'number', min: 0.000001 }]}
      >
        <InputNumber style={{ width: '100%' }} size="large" min={0.000001} />
      </Form.Item>

      <Button type="primary" htmlType="submit" block size="large" style={{ marginTop: 16 }}>
        Review Transaction
      </Button>
    </Form>
  );
};
