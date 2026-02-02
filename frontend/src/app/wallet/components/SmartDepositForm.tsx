import { Form, InputNumber, Button, Alert } from 'antd';
import { ExternalLink } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { useTheme } from 'antd-style';
import { Browser } from '@wailsio/runtime';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';

interface SmartDepositFormProps {
  onDeposit: (amount: number) => void;
  loading: boolean;
  currentNetwork: NetworkInfo | null;
}

export const SmartDepositForm = ({ onDeposit, loading, currentNetwork }: SmartDepositFormProps) => {
  const [form] = Form.useForm();
  const theme = useTheme();

  return (
    <Flexbox gap={16}>
      <Alert
        type="warning"
        showIcon
        message={
          <span style={{ fontWeight: 600 }}>
            Only deposit {currentNetwork?.stablecoinSymbol || 'USDT'} on Monad Network!
          </span>
        }
        description={
          <Flexbox gap={8} style={{ marginTop: 8 }}>
            <span>
              Don't have {currentNetwork?.stablecoinSymbol || 'USDT'} on Monad? You need to bridge from other networks first.
            </span>
            <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
              Network: <strong>{currentNetwork?.name || 'Monad Mainnet'}</strong> (Chain ID: {currentNetwork?.id || 143})
            </span>
            <Button
              type="link"
              size="small"
              icon={<ExternalLink size={14} />}
              onClick={() => Browser.OpenURL('https://getkawai.com/docs/user-guide/deposit-from-exchange')}
              style={{ padding: 0, height: 'auto' }}
            >
              Learn how to bridge from exchanges
            </Button>
          </Flexbox>
        }
        style={{ marginBottom: 8 }}
      />

      <Form form={form} layout="vertical" onFinish={(v) => onDeposit(v.amount)} initialValues={{ amount: 10 }}>
        <Form.Item
          label={`Amount (${currentNetwork?.stablecoinShort || 'USDT'})`}
          name="amount"
          rules={[{ required: true, min: 1, type: 'number', message: `Please enter at least 1 ${currentNetwork?.stablecoinShort || 'USDT'}` }]}
        >
          <InputNumber
            style={{ width: '100%' }}
            size="large"
            addonAfter={currentNetwork?.stablecoinShort || 'USDT'}
            min={1}
          />
        </Form.Item>
        <Button type="primary" htmlType="submit" block size="large" loading={loading} style={{ marginTop: 16 }}>
          {loading ? 'Processing Transaction...' : 'Confirm Deposit'}
        </Button>
      </Form>
    </Flexbox>
  );
};
