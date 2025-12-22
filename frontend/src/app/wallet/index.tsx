import { Card, Modal, QRCode, App, List, Avatar } from 'antd';
import { memo, useEffect, useState } from 'react';
import { DeAIService, WalletService, HistoryService } from '@@/github.com/kawai-network/veridium/internal/services';
import { useUserStore } from '@/store/user';
import { QrCode, ArrowDownToLine, Copy, Send } from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { createStyles } from 'antd-style';
import type { TransactionRecord } from '@@/github.com/kawai-network/veridium/internal/services/models';

const useStyles = createStyles(({ css, token }) => ({
  container: css`
    height: 100%;
    width: 100%;
    background: ${token.colorBgLayout};
    padding: 24px;
    overflow-y: auto;
  `,
  balanceCard: css`
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    border-radius: 16px;
    color: white;
    
    .ant-card-body {
      padding: 32px;
    }
  `,
  actionButton: css`
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 16px;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
    }
  `,
}));

const DesktopWalletLayout = memo(() => {
  const [address, setAddress] = useState<string>('');
  const [balance, setBalance] = useState<string>('0');
  const [loading, setLoading] = useState(false);
  const [transactions, setTransactions] = useState<TransactionRecord[]>([]);
  const [modalType, setModalType] = useState<'deposit' | 'send' | 'receive' | null>(null);
  const { message } = App.useApp();
  const { isWalletLoaded } = useUserStore();
  const { styles } = useStyles();

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
      const history = await HistoryService.GetTransactions();
      setTransactions(history);
    } catch (e) {
      console.error("Failed to load history", e);
    }
  };

  const handleDeposit = async (amount: number) => {
    if (!amount || amount <= 0) {
      message.error("Please enter a valid amount");
      return;
    }

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

  return (
    <Flexbox className={styles.container} align="center" justify="flex-start">
      <Flexbox style={{ maxWidth: 1200, width: '100%', flexDirection: 'column' }} gap={24}>
        <h1 style={{ margin: 0 }}>My Wallet</h1>

        {/* Balance Card */}
        <Card className={styles.balanceCard}>
          <Flexbox style={{ flexDirection: 'column' }} gap={8}>
            <div style={{ fontSize: 14, opacity: 0.9 }}>Total Balance</div>
            <div style={{ fontSize: 48, fontWeight: 'bold' }}>{balance} USDT</div>
            <div style={{ fontSize: 12, opacity: 0.8 }}>{address}</div>
          </Flexbox>
        </Card>

        {/* Action Grid */}
        <Flexbox gap={16} justify="space-around">
          <Flexbox className={styles.actionButton} onClick={() => setModalType('deposit')}>
            <ActionIcon
              icon={ArrowDownToLine}
              size={{ blockSize: 64, size: 32 }}
              style={{ background: '#1677ff', color: 'white' }}
            />
            <span style={{ fontWeight: 600, fontSize: 14 }}>Deposit</span>
          </Flexbox>

          <Flexbox className={styles.actionButton} onClick={() => setModalType('send')}>
            <ActionIcon
              icon={Send}
              size={{ blockSize: 64, size: 32 }}
              style={{ background: '#52c41a', color: 'white' }}
            />
            <span style={{ fontWeight: 600, fontSize: 14 }}>Send</span>
          </Flexbox>

          <Flexbox className={styles.actionButton} onClick={() => setModalType('receive')}>
            <ActionIcon
              icon={QrCode}
              size={{ blockSize: 64, size: 32 }}
              style={{ background: '#722ed1', color: 'white' }}
            />
            <span style={{ fontWeight: 600, fontSize: 14 }}>Receive</span>
          </Flexbox>
        </Flexbox>

        {/* Transaction History */}
        <div>
          <h2>Recent Transactions</h2>
          <List
            dataSource={transactions}
            locale={{ emptyText: 'No transactions yet' }}
            renderItem={(tx) => (
              <List.Item
                extra={
                  <Flexbox style={{ flexDirection: 'column' }} align="flex-end" gap={4}>
                    <div style={{ fontWeight: 600 }}>{tx.amount}</div>
                    <div style={{ fontSize: 12, color: '#888' }}>
                      {tx.hash.substring(0, 10)}...
                    </div>
                  </Flexbox>
                }
              >
                <List.Item.Meta
                  avatar={
                    <Avatar
                      style={{
                        background: tx.type === 'DEPOSIT' ? '#1677ff' : '#52c41a'
                      }}
                      icon={tx.type === 'DEPOSIT' ? <ArrowDownToLine size={20} /> : <Send size={20} />}
                    />
                  }
                  title={tx.description}
                  description={new Date(tx.timestamp * 1000).toLocaleString()}
                />
              </List.Item>
            )}
          />
        </div>
      </Flexbox>

      {/* Modals */}
      <Modal
        title="Smart Deposit"
        open={modalType === 'deposit'}
        onCancel={() => setModalType(null)}
        footer={null}
      >
        <SmartDepositForm onDeposit={handleDeposit} loading={loading} />
      </Modal>

      <Modal
        title="Send USDT"
        open={modalType === 'send'}
        onCancel={() => setModalType(null)}
        footer={null}
      >
        <SendForm onSend={handleSend} loading={loading} />
      </Modal>

      <Modal
        title="Receive USDT"
        open={modalType === 'receive'}
        onCancel={() => setModalType(null)}
        footer={null}
      >
        <Flexbox style={{ flexDirection: 'column' }} align="center" gap={24}>
          <p style={{ color: '#aaa', textAlign: 'center' }}>
            Scan this QR code to deposit USDT (BEP20) or BNB to your wallet.
          </p>
          <div style={{ background: '#fff', padding: 16, borderRadius: 12 }}>
            <QRCode value={address || "0x"} size={200} color="#000000" />
          </div>
          <div style={{ background: '#2c2c2c', padding: 16, borderRadius: 8, width: '100%' }}>
            <p style={{ color: '#888', fontSize: 12, marginBottom: 4, textAlign: 'center' }}>Your Wallet Address</p>
            <Flexbox gap={10} align="center" justify="center">
              <span style={{ fontFamily: 'monospace', fontSize: 14 }}>{address}</span>
              <CopyButton text={address} />
            </Flexbox>
          </div>
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
      size={{ blockSize: 24, size: 16 }}
      onClick={() => {
        navigator.clipboard.writeText(text);
        message.success("Address copied!");
      }}
      title="Copy address"
    />
  );
};

const SmartDepositForm = ({ onDeposit, loading }: { onDeposit: (val: number) => void, loading: boolean }) => {
  const [val, setVal] = useState(10);
  return (
    <Flexbox style={{ flexDirection: 'column' }} gap={16}>
      <div>
        <p style={{ marginBottom: 8, color: '#aaa' }}>Amount (USDT)</p>
        <input
          type="number"
          value={val}
          onChange={e => setVal(Number(e.target.value))}
          style={{ width: '100%', padding: '12px', borderRadius: 6, border: '1px solid #444', background: '#333', color: '#fff' }}
        />
      </div>
      <button
        disabled={loading}
        onClick={() => onDeposit(val)}
        style={{
          padding: '12px',
          borderRadius: 6,
          background: '#1677ff',
          color: '#fff',
          border: 'none',
          cursor: loading ? 'not-allowed' : 'pointer',
          opacity: loading ? 0.7 : 1,
          fontWeight: 600
        }}
      >
        {loading ? 'Processing...' : 'Deposit'}
      </button>
    </Flexbox>
  );
};

const SendForm = ({ onSend, loading }: { onSend: (to: string, val: number) => void, loading: boolean }) => {
  const [to, setTo] = useState('');
  const [val, setVal] = useState(0);

  return (
    <Flexbox style={{ flexDirection: 'column' }} gap={16}>
      <div>
        <p style={{ marginBottom: 8, color: '#aaa' }}>Recipient Address</p>
        <input
          value={to}
          onChange={e => setTo(e.target.value)}
          placeholder="0x..."
          style={{ width: '100%', padding: '12px', borderRadius: 6, border: '1px solid #444', background: '#333', color: '#fff' }}
        />
      </div>
      <div>
        <p style={{ marginBottom: 8, color: '#aaa' }}>Amount (USDT)</p>
        <input
          type="number"
          value={val}
          onChange={e => setVal(Number(e.target.value))}
          style={{ width: '100%', padding: '12px', borderRadius: 6, border: '1px solid #444', background: '#333', color: '#fff' }}
        />
      </div>
      <button
        disabled={loading || !to || val <= 0}
        onClick={() => onSend(to, val)}
        style={{
          padding: '12px',
          borderRadius: 6,
          background: '#52c41a',
          color: '#fff',
          border: 'none',
          cursor: (loading || !to || val <= 0) ? 'not-allowed' : 'pointer',
          opacity: (loading || !to || val <= 0) ? 0.7 : 1,
          fontWeight: 600
        }}
      >
        {loading ? 'Sending...' : 'Send USDT'}
      </button>
    </Flexbox>
  );
};

export default DesktopWalletLayout;
