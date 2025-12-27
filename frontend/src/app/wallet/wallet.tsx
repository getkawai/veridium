import { Card, Modal, QRCode, App, Table, Tag, Tabs, Button, InputNumber, Form, Input } from 'antd';
import { memo, useEffect, useState } from 'react';
import { DeAIService, WalletService } from '@@/github.com/kawai-network/veridium/internal/services';
import { ListWalletTransactions } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import type { WalletTransaction } from '@@/github.com/kawai-network/veridium/internal/database/generated/models';
import { useUserStore } from '@/store/user';
import { ArrowDownToLine, Copy, Send, Eye, EyeOff, ArrowUp, Repeat2, Wallet as WalletIcon, History, Key, Plus } from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { createStyles } from 'antd-style';

const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    width: 100%;
    background: ${token.colorBgLayout};
    padding: 24px;
    overflow-y: auto;
  `,
  balanceCard: css`
    background: linear-gradient(135deg, ${token.colorBgContainer} 0%, ${token.colorBgElevated} 100%);
    border: 1px solid ${token.colorBorderSecondary};
    border-radius: 16px;
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      right: 0;
      width: 200px;
      height: 200px;
      background: radial-gradient(circle, ${token.colorPrimary} 0%, transparent 70%);
      opacity: 0.1;
      border-radius: 50%;
      transform: translate(30%, -30%);
    }

    .ant-card-body {
      padding: 32px;
    }
  `,
  actionButton: css`
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 16px;
    border-radius: 12px;
    
    &:hover {
      background: ${token.colorFillTertiary};
      transform: translateY(-2px);
    }
  `,
  actionCircle: css`
    width: 56px;
    height: 56px;
    border-radius: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  `,
  eyeButton: css`
    position: absolute;
    top: 24px;
    right: 24px;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 0.3s ease;
    
    &:hover {
      opacity: 1;
    }
  `,
  statValue: css`
    font-size: 42px;
    font-weight: 700;
    line-height: 1.2;
    background: -webkit-linear-gradient(120deg, ${token.colorText} 30%, ${token.colorTextSecondary});
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  `,
}));

const DesktopWalletLayout = memo(() => {
  const [address, setAddress] = useState<string>('');
  const [balance, setBalance] = useState<string>('0');
  const [loading, setLoading] = useState(false);
  const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
  const [modalType, setModalType] = useState<'deposit' | 'send' | 'receive' | 'swap' | null>(null);
  const [balanceVisible, setBalanceVisible] = useState(true);
  const { message } = App.useApp();
  const { isWalletLoaded } = useUserStore();
  const { styles, theme } = useStyles();

  // Mock price change data (can be replaced with real API later)
  const priceChangePercent = 2.27;
  const priceChange24h = 1245.32;

  useEffect(() => {
    WalletService.GetCurrentAddress().then(setAddress).catch(console.error);
    loadBalance();
    loadHistory();
  }, [isWalletLoaded]);

  const loadBalance = async () => {
    try {
      const bal = await DeAIService.GetVaultBalance();
      setBalance(bal);
    } catch (e) {
      console.error("Failed to load balance", e);
    }
  };

  const loadHistory = async () => {
    try {
      const history = await ListWalletTransactions({ limit: 100, offset: 0 });
      setTransactions(history);
    } catch (e) {
      console.error("Failed to load history", e);
    }
  };

  const handleDeposit = async (amount: number) => {
    setLoading(true);
    const hide = message.loading("Processing Smart Deposit via Blockchain...", 0);

    try {
      const rawAmount = Math.floor(amount * 1_000_000).toString();
      const txHash = await DeAIService.DepositToVault(rawAmount);
      message.success(`Deposit Successful! TX: ${txHash.substring(0, 10)}...`);
      loadBalance();
      loadHistory();
      setModalType(null);
    } catch (e: any) {
      console.error(e);
      let msg = e.message || e;
      if (typeof msg === 'string' && msg.includes("insufficient funds")) {
        msg = "Insufficient BNB for gas fees! Please top up BNB.";
      }
      message.error(`Deposit Failed: ${msg}`);
    } finally {
      hide();
      setLoading(false);
    }
  };

  const handleSend = async (to: string, amount: number) => {
    setLoading(true);
    const hide = message.loading("Sending USDT...", 0);
    try {
      const rawAmount = Math.floor(amount * 1_000_000).toString();
      const tx = await DeAIService.TransferUSDT(to, rawAmount);
      message.success(`Transfer Successful! TX: ${tx.substring(0, 10)}...`);
      loadBalance();
      loadHistory();
      setModalType(null);
    } catch (e: any) {
      console.error(e);
      message.error(`Transfer Failed: ${e.message || e}`);
    } finally {
      hide();
      setLoading(false);
    }
  };

  const TransactionTable = () => (
    <Table
      dataSource={transactions}
      rowKey="id"
      pagination={{ pageSize: 5 }}
      columns={[
        {
          title: 'Type',
          dataIndex: 'txType',
          key: 'txType',
          render: (type) => (
            <Tag color={type === 'DEPOSIT' ? 'green' : 'blue'}>
              {type}
            </Tag>
          ),
        },
        {
          title: 'Amount',
          dataIndex: 'amount',
          key: 'amount',
          render: (amount, record) => (
            <span style={{
              color: record.txType === 'DEPOSIT' ? theme.colorSuccess : theme.colorText,
              fontWeight: 600
            }}>
              {record.txType === 'DEPOSIT' ? '+' : '-'}{amount} USDT
            </span>
          ),
        },
        {
          title: 'Date',
          dataIndex: 'createdAt',
          key: 'createdAt',
          render: (date) => new Date(date).toLocaleString(),
        },
        {
          title: 'Status',
          key: 'status',
          render: () => <Tag color="success">Confirmed</Tag>,
        },
      ]}
    />
  );

  return (
    <Flexbox className={styles.container} align="center" justify="flex-start">
      <Flexbox style={{ maxWidth: 1000, width: '100%' }} gap={32}>

        {/* Header Section */}
        <Flexbox justify="space-between" align="center">
          <div>
            <h1 style={{ margin: 0, fontSize: 24, fontWeight: 700 }}>Overview</h1>
            <span style={{ color: theme.colorTextSecondary }}>Manage your assets and keys</span>
          </div>
          <Tag icon={<WalletIcon size={14} />} color="default" style={{ padding: '6px 12px' }}>
            {address ? `${address.substring(0, 6)}...${address.substring(address.length - 4)}` : 'Loading...'}
          </Tag>
        </Flexbox>

        {/* Portfolio Value Card */}
        <Card className={styles.balanceCard}>
          <div className={styles.eyeButton} onClick={() => setBalanceVisible(!balanceVisible)}>
            {balanceVisible ? <Eye size={20} /> : <EyeOff size={20} />}
          </div>

          <Flexbox gap={24}>
            <Flexbox style={{ flexDirection: 'column' }} gap={8}>
              <span style={{ fontSize: 14, color: theme.colorTextSecondary }}>Total Balance</span>
              <div className={styles.statValue}>
                {balanceVisible ? `$${balance}` : '••••••'}
                <span style={{ fontSize: 20, color: theme.colorTextTertiary, marginLeft: 8 }}>USDT</span>
              </div>

              <Flexbox align="center" gap={8} style={{ color: theme.colorSuccess }}>
                <ArrowUp size={16} />
                <span style={{ fontWeight: 600 }}>{priceChangePercent}%</span>
                <span style={{ opacity: 0.7 }}>+${priceChange24h.toFixed(2)} (24h)</span>
              </Flexbox>
            </Flexbox>
          </Flexbox>
        </Card>

        {/* Quick Actions */}
        <Flexbox gap={16} justify="space-between">
          {[
            { label: 'Send', icon: Send, color: '#06b6d4', action: () => setModalType('send') },
            { label: 'Receive', icon: ArrowDownToLine, color: '#22c55e', action: () => setModalType('receive') },
            { label: 'Swap', icon: Repeat2, color: '#eab308', action: () => setModalType('swap') },
          ].map((item) => (
            <div key={item.label} className={styles.actionButton} onClick={item.action} style={{ flex: 1 }}>
              <div className={styles.actionCircle} style={{ background: `${item.color}15`, color: item.color }}>
                <item.icon size={24} />
              </div>
              <span style={{ fontWeight: 600 }}>{item.label}</span>
            </div>
          ))}
        </Flexbox>

        {/* Tabs Section */}
        <Card styles={{ body: { padding: 0 } }} bordered={false}>
          <Tabs
            defaultActiveKey="transactions"
            items={[
              {
                key: 'transactions',
                label: (
                  <Flexbox gap={8} align="center" style={{ padding: '0 16px' }}>
                    <History size={16} />
                    Transactions
                  </Flexbox>
                ),
                children: <div style={{ padding: 16 }}><TransactionTable /></div>,
              },
            ]}
          />
        </Card>
      </Flexbox>

      {/* Modals */}
      <Modal
        title="Smart Deposit"
        open={modalType === 'deposit'}
        onCancel={() => setModalType(null)}
        footer={null}
        destroyOnClose
      >
        <SmartDepositForm onDeposit={handleDeposit} loading={loading} />
      </Modal>

      <Modal
        title="Send USDT"
        open={modalType === 'send'}
        onCancel={() => setModalType(null)}
        footer={null}
        destroyOnClose
      >
        <SendForm onSend={handleSend} loading={loading} />
      </Modal>

      <Modal
        title="Receive USDT"
        open={modalType === 'receive'}
        onCancel={() => setModalType(null)}
        footer={null}
        width={400}
      >
        <Flexbox style={{ flexDirection: 'column', padding: 24 }} align="center" gap={24}>
          <div style={{ background: '#fff', padding: 16, borderRadius: 16 }}>
            <QRCode value={address || "0x"} size={200} />
          </div>
          <div style={{ background: theme.colorFillTertiary, padding: 16, borderRadius: 12, width: '100%' }}>
            <p style={{ color: theme.colorTextSecondary, fontSize: 12, marginBottom: 8, textAlign: 'center' }}>
              Your Wallet Address (Monad Testnet)
            </p>
            <Flexbox gap={10} align="center" justify="center">
              <span style={{ fontFamily: 'monospace', fontWeight: 600 }}>
                {address.substring(0, 10)}...{address.substring(address.length - 10)}
              </span>
              <CopyButton text={address} />
            </Flexbox>
          </div>
        </Flexbox>
      </Modal>

      <Modal open={modalType === 'swap'} onCancel={() => setModalType(null)} footer={null}>
        <Flexbox style={{ flexDirection: 'column', padding: 32 }} align="center" gap={16}>
          <Repeat2 size={48} color={theme.colorTextQuaternary} />
          <h3>Coming Soon</h3>
          <p style={{ textAlign: 'center', color: theme.colorTextSecondary }}>
            Token swapping will be available in the next update.
          </p>
        </Flexbox>
      </Modal>
    </Flexbox>
  );
});

const CopyButton = ({ text }: { text: string }) => {
  const { message } = App.useApp();
  return (
    <ActionIcon
      icon={Copy}
      size="small"
      onClick={() => {
        navigator.clipboard.writeText(text);
        message.success("Address copied!");
      }}
      title="Copy address"
    />
  );
};

const SmartDepositForm = ({ onDeposit, loading }: { onDeposit: (val: number) => void, loading: boolean }) => {
  const [form] = Form.useForm();

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(v) => onDeposit(v.amount)}
      initialValues={{ amount: 10 }}
    >
      <Form.Item
        label="Amount (USDT)"
        name="amount"
        rules={[{ required: true, min: 1, type: 'number', message: 'Please enter at least 1 USDT' }]}
      >
        <InputNumber
          style={{ width: '100%' }}
          size="large"
          addonAfter="USDT"
          min={1}
        />
      </Form.Item>

      <Button type="primary" htmlType="submit" block size="large" loading={loading} style={{ marginTop: 16 }}>
        {loading ? 'Processing Transaction...' : 'Confirm Deposit'}
      </Button>
    </Form>
  );
};

const SendForm = ({ onSend, loading }: { onSend: (to: string, val: number) => void, loading: boolean }) => {
  const [form] = Form.useForm();

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={(v) => onSend(v.to, v.amount)}
    >
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
        rules={[{ required: true, type: 'number', min: 0.1 }]}
      >
        <InputNumber style={{ width: '100%' }} size="large" addonAfter="USDT" min={0.1} />
      </Form.Item>

      <Button type="primary" htmlType="submit" block size="large" loading={loading} style={{ marginTop: 16 }}>
        Send USDT
      </Button>
    </Form>
  );
};

export default DesktopWalletLayout;
