import { Card, Modal, QRCode, App, Table, Tag, Tabs, Button, InputNumber, Form, Input, Menu, Empty } from 'antd';
import { memo, useEffect, useState } from 'react';
import { DeAIService, WalletService } from '@@/github.com/kawai-network/veridium/internal/services';
import { ListWalletTransactions } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import type { WalletTransaction } from '@@/github.com/kawai-network/veridium/internal/database/generated/models';
import { useUserStore } from '@/store/user';
import { ArrowDownToLine, Copy, Send, Eye, EyeOff, ArrowUp, Repeat2, Wallet as WalletIcon, History, Home, ShoppingCart, Gift, Settings, Coins, ExternalLink } from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { createStyles } from 'antd-style';

type MenuKey = 'home' | 'otc' | 'rewards' | 'settings';

const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    width: 100%;
    background: ${token.colorBgLayout};
    display: flex;
  `,
  sidebar: css`
    width: 200px;
    min-width: 200px;
    background: ${token.colorBgContainer};
    border-right: 1px solid ${token.colorBorderSecondary};
    padding: 16px 0;
    display: flex;
    flex-direction: column;
  `,
  sidebarMenu: css`
    border: none !important;
    background: transparent !important;
    
    .ant-menu-item {
      margin: 4px 8px !important;
      border-radius: 8px !important;
      height: 44px !important;
      line-height: 44px !important;
      
      &:hover {
        background: ${token.colorFillTertiary} !important;
      }
    }
    
    .ant-menu-item-selected {
      background: ${token.colorPrimaryBg} !important;
      color: ${token.colorPrimary} !important;
    }
  `,
  content: css`
    flex: 1;
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
      padding: 24px;
    }
  `,
  actionButton: css`
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 12px 16px;
    border-radius: 12px;
    
    &:hover {
      background: ${token.colorFillTertiary};
      transform: translateY(-2px);
    }
  `,
  actionCircle: css`
    width: 48px;
    height: 48px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
  `,
  eyeButton: css`
    position: absolute;
    top: 20px;
    right: 20px;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 0.3s ease;
    
    &:hover {
      opacity: 1;
    }
  `,
  statValue: css`
    font-size: 36px;
    font-weight: 700;
    line-height: 1.2;
    background: -webkit-linear-gradient(120deg, ${token.colorText} 30%, ${token.colorTextSecondary});
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  `,
  tokenRow: css`
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-radius: 12px;
    background: ${token.colorBgContainer};
    border: 1px solid ${token.colorBorderSecondary};
    
    &:hover {
      border-color: ${token.colorPrimary};
    }
  `,
  placeholderCard: css`
    background: ${token.colorBgContainer};
    border: 1px solid ${token.colorBorderSecondary};
    border-radius: 16px;
    padding: 48px;
    text-align: center;
  `,
}));

const DesktopWalletLayout = memo(() => {
  const [activeMenu, setActiveMenu] = useState<MenuKey>('home');
  const [address, setAddress] = useState<string>('');
  const [balance, setBalance] = useState<string>('0');
  const [loading, setLoading] = useState(false);
  const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
  const [modalType, setModalType] = useState<'deposit' | 'send' | 'receive' | 'swap' | null>(null);
  const [balanceVisible, setBalanceVisible] = useState(true);
  const { message } = App.useApp();
  const { isWalletLoaded } = useUserStore();
  const { styles, theme } = useStyles();

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
        msg = "Insufficient MON for gas fees!";
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

  const menuItems = [
    { key: 'home', icon: <Home size={18} />, label: 'Home' },
    { key: 'otc', icon: <ShoppingCart size={18} />, label: 'OTC Market' },
    { key: 'rewards', icon: <Gift size={18} />, label: 'Rewards' },
    { key: 'settings', icon: <Settings size={18} />, label: 'Settings' },
  ];

  const renderContent = () => {
    switch (activeMenu) {
      case 'home':
        return <HomeContent
          address={address}
          balance={balance}
          balanceVisible={balanceVisible}
          setBalanceVisible={setBalanceVisible}
          setModalType={setModalType}
          transactions={transactions}
          styles={styles}
          theme={theme}
        />;
      case 'otc':
        return <OTCContent styles={styles} theme={theme} />;
      case 'rewards':
        return <RewardsContent styles={styles} theme={theme} />;
      case 'settings':
        return <SettingsContent address={address} styles={styles} theme={theme} />;
      default:
        return null;
    }
  };

  return (
    <Flexbox className={styles.container}>
      {/* Sidebar */}
      <div className={styles.sidebar}>
        <Flexbox align="center" gap={8} style={{ padding: '0 16px 16px' }}>
          <WalletIcon size={20} />
          <span style={{ fontWeight: 600 }}>Kawai Wallet</span>
        </Flexbox>

        <Menu
          className={styles.sidebarMenu}
          mode="inline"
          selectedKeys={[activeMenu]}
          items={menuItems}
          onClick={({ key }) => setActiveMenu(key as MenuKey)}
        />

        <div style={{ flex: 1 }} />

        <Flexbox style={{ padding: '16px' }} gap={4}>
          <Tag color="green" style={{ margin: 0 }}>Monad Testnet</Tag>
          <span style={{ fontSize: 11, color: theme.colorTextTertiary }}>Chain ID: 10143</span>
        </Flexbox>
      </div>

      {/* Content */}
      <div className={styles.content}>
        {renderContent()}
      </div>

      {/* Modals */}
      <Modal title="Smart Deposit" open={modalType === 'deposit'} onCancel={() => setModalType(null)} footer={null} destroyOnClose>
        <SmartDepositForm onDeposit={handleDeposit} loading={loading} />
      </Modal>

      <Modal title="Send USDT" open={modalType === 'send'} onCancel={() => setModalType(null)} footer={null} destroyOnClose>
        <SendForm onSend={handleSend} loading={loading} />
      </Modal>

      <Modal title="Receive" open={modalType === 'receive'} onCancel={() => setModalType(null)} footer={null} width={400}>
        <Flexbox style={{ flexDirection: 'column', padding: 24 }} align="center" gap={24}>
          <div style={{ background: '#fff', padding: 16, borderRadius: 16 }}>
            <QRCode value={address || "0x"} size={200} />
          </div>
          <div style={{ background: theme.colorFillTertiary, padding: 16, borderRadius: 12, width: '100%' }}>
            <p style={{ color: theme.colorTextSecondary, fontSize: 12, marginBottom: 8, textAlign: 'center' }}>
              Your Wallet Address (Monad Testnet)
            </p>
            <Flexbox gap={10} align="center" justify="center">
              <span style={{ fontFamily: 'monospace', fontWeight: 600, fontSize: 12 }}>
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

// ============ HOME CONTENT ============
const HomeContent = ({ address, balance, balanceVisible, setBalanceVisible, setModalType, transactions, styles, theme }: any) => {
  return (
    <Flexbox style={{ maxWidth: 700 }} gap={20}>
      {/* Header */}
      <Flexbox justify="space-between" align="center">
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Overview</h2>
          <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>Manage your assets</span>
        </div>
        <Tag icon={<WalletIcon size={12} />} style={{ padding: '4px 10px', cursor: 'pointer' }}>
          {address ? `${address.substring(0, 6)}...${address.substring(address.length - 4)}` : 'Loading...'}
        </Tag>
      </Flexbox>

      {/* Balance Card */}
      <Card className={styles.balanceCard}>
        <div className={styles.eyeButton} onClick={() => setBalanceVisible(!balanceVisible)}>
          {balanceVisible ? <Eye size={16} /> : <EyeOff size={16} />}
        </div>
        <Flexbox horizontal justify="space-between" align="center">
          <Flexbox style={{ flexDirection: 'column' }} gap={4}>
            <span style={{ fontSize: 11, color: theme.colorTextSecondary, textTransform: 'uppercase', letterSpacing: '0.5px' }}>Total Balance</span>
            <div className={styles.statValue}>
              {balanceVisible ? `$${balance}` : '••••••'}
              <span style={{ fontSize: 16, color: theme.colorTextTertiary, marginLeft: 6, fontWeight: 500 }}>USDT</span>
            </div>
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Quick Actions */}
      <Flexbox horizontal gap={8}>
        {[
          { label: 'Send', icon: Send, color: '#06b6d4', action: () => setModalType('send') },
          { label: 'Receive', icon: ArrowDownToLine, color: '#22c55e', action: () => setModalType('receive') },
          { label: 'Swap', icon: Repeat2, color: '#eab308', action: () => setModalType('swap') },
        ].map((item) => (
          <div key={item.label} className={styles.actionButton} onClick={item.action}>
            <div className={styles.actionCircle} style={{ background: `${item.color}20`, color: item.color }}>
              <item.icon size={20} />
            </div>
            <span style={{ fontWeight: 500, fontSize: 12 }}>{item.label}</span>
          </div>
        ))}
      </Flexbox>

      {/* Token List */}
      <Card title={<Flexbox horizontal align="center" gap={8}><Coins size={16} /> Tokens</Flexbox>} size="small">
        <Flexbox gap={8}>
          <div className={styles.tokenRow}>
            <Flexbox horizontal align="center" gap={12}>
              <div style={{ width: 36, height: 36, borderRadius: '50%', background: '#26a17b', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 700, fontSize: 12 }}>U</div>
              <div>
                <div style={{ fontWeight: 600 }}>USDT</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Tether USD</div>
              </div>
            </Flexbox>
            <div style={{ textAlign: 'right' }}>
              <div style={{ fontWeight: 600 }}>{balanceVisible ? balance : '••••'}</div>
              <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>${balanceVisible ? balance : '••••'}</div>
            </div>
          </div>

          <div className={styles.tokenRow}>
            <Flexbox horizontal align="center" gap={12}>
              <div style={{ width: 36, height: 36, borderRadius: '50%', background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 700, fontSize: 12 }}>K</div>
              <div>
                <div style={{ fontWeight: 600 }}>KAWAI</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Kawai Token</div>
              </div>
            </Flexbox>
            <div style={{ textAlign: 'right' }}>
              <div style={{ fontWeight: 600 }}>0</div>
              <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>$0.00</div>
            </div>
          </div>
        </Flexbox>
      </Card>

      {/* Activity */}
      <Card title={<Flexbox horizontal align="center" gap={8}><History size={16} /> Recent Activity</Flexbox>} size="small">
        {transactions.length > 0 ? (
          <Table
            dataSource={transactions.slice(0, 5)}
            rowKey="id"
            pagination={false}
            size="small"
            columns={[
              { title: 'Type', dataIndex: 'txType', key: 'txType', render: (type) => <Tag color={type === 'DEPOSIT' ? 'green' : 'blue'}>{type}</Tag> },
              { title: 'Amount', dataIndex: 'amount', key: 'amount', render: (amount, record: any) => <span style={{ color: record.txType === 'DEPOSIT' ? theme.colorSuccess : theme.colorText, fontWeight: 600 }}>{record.txType === 'DEPOSIT' ? '+' : '-'}{amount} USDT</span> },
              { title: 'Date', dataIndex: 'createdAt', key: 'createdAt', render: (date) => new Date(date).toLocaleDateString() },
            ]}
          />
        ) : (
          <Empty description="No transactions yet" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        )}
      </Card>
    </Flexbox>
  );
};

// ============ OTC MARKET CONTENT ============
const OTCContent = ({ styles, theme }: any) => {
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

// ============ REWARDS CONTENT ============
const RewardsContent = ({ styles, theme }: any) => {
  return (
    <Flexbox style={{ maxWidth: 700 }} gap={20}>
      <div>
        <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Rewards</h2>
        <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>Claim your KAWAI mining rewards</span>
      </div>

      <div className={styles.placeholderCard}>
        <Gift size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
        <h3 style={{ margin: '0 0 8px' }}>Coming Soon</h3>
        <p style={{ color: theme.colorTextSecondary, margin: 0 }}>
          Contributors can claim weekly KAWAI rewards here.<br />
          Merkle-based distribution for gas efficiency.
        </p>
      </div>
    </Flexbox>
  );
};

// ============ SETTINGS CONTENT ============
const SettingsContent = ({ address, styles, theme }: any) => {
  return (
    <Flexbox style={{ maxWidth: 700 }} gap={20}>
      <div>
        <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Settings</h2>
        <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>Network and wallet configuration</span>
      </div>

      <Card title="Network Info" size="small">
        <Flexbox gap={12}>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Network</span>
            <Tag color="green">Monad Testnet</Tag>
          </Flexbox>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Chain ID</span>
            <span style={{ fontFamily: 'monospace' }}>10143</span>
          </Flexbox>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>RPC</span>
            <span style={{ fontFamily: 'monospace', fontSize: 12 }}>testnet-rpc.monad.xyz</span>
          </Flexbox>
          <Flexbox horizontal justify="space-between">
            <span style={{ color: theme.colorTextSecondary }}>Explorer</span>
            <a href="https://testnet.monad.xyz" target="_blank" rel="noopener noreferrer" style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
              testnet.monad.xyz <ExternalLink size={12} />
            </a>
          </Flexbox>
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

// ============ HELPER COMPONENTS ============
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

const SmartDepositForm = ({ onDeposit, loading }: { onDeposit: (val: number) => void, loading: boolean }) => {
  const [form] = Form.useForm();
  return (
    <Form form={form} layout="vertical" onFinish={(v) => onDeposit(v.amount)} initialValues={{ amount: 10 }}>
      <Form.Item label="Amount (USDT)" name="amount" rules={[{ required: true, min: 1, type: 'number', message: 'Please enter at least 1 USDT' }]}>
        <InputNumber style={{ width: '100%' }} size="large" addonAfter="USDT" min={1} />
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
    <Form form={form} layout="vertical" onFinish={(v) => onSend(v.to, v.amount)}>
      <Form.Item label="Recipient Address" name="to" rules={[{ required: true, message: 'Address is required' }, { pattern: /^0x[a-fA-F0-9]{40}$/, message: 'Invalid EVM address' }]}>
        <Input placeholder="0x..." size="large" />
      </Form.Item>
      <Form.Item label="Amount" name="amount" rules={[{ required: true, type: 'number', min: 0.1 }]}>
        <InputNumber style={{ width: '100%' }} size="large" addonAfter="USDT" min={0.1} />
      </Form.Item>
      <Button type="primary" htmlType="submit" block size="large" loading={loading} style={{ marginTop: 16 }}>
        Send USDT
      </Button>
    </Form>
  );
};

export default DesktopWalletLayout;
