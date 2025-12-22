
import { Flex, App, message, Tabs, QRCode } from 'antd';
import { memo, useEffect, useState } from 'react';
import styled from 'styled-components';
import { DeAIService, WalletService } from '../../../bindings/github.com/kawai-network/veridium/internal/services';
import { useUserStore } from '@/store/user';
import { QrCode, ArrowDownToLine, Copy } from 'lucide-react';


const Container = styled(Flex)`
  height: 100%;
  width: 100%;
  position: relative;
  background: ${({ theme }) => theme.colorBgLayout};
  padding: 24px;
`;

const ContentWrapper = styled.div`
  max-width: 800px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const BSC_CHAIN: Chain = {
  id: 97,
  name: 'BSC Testnet',
  type: ChainType.EVM,
};

const SUPPORTED_CHAINS = [
  {
    chain: BSC_CHAIN,
  }
]

const MOCK_USDT_TOKEN = {
  decimal: 6,
  symbol: "USDT",
  name: "Tether USD",
  icon: "https://cryptologos.cc/logos/tether-usdt-logo.svg?v=032",
  availableChains: [
    {
      chain: BSC_CHAIN,
      contract: "0x312C4fC3598AC9B54375eD12BbF55af83f86f862"
    }
  ]
}

const DesktopWalletLayout = memo(() => {
  const [address, setAddress] = useState<string>('');
  const [balance, setBalance] = useState<string>('0');
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();
  const { isWalletLoaded } = useUserStore();

  // Fetch initial data
  useEffect(() => {
    WalletService.GetCurrentAddress().then(setAddress).catch(console.error);
    loadBalance();
  }, [isWalletLoaded]);

  const loadBalance = async () => {
    try {
      const bal = await DeAIService.GetVaultBalance();
      setBalance(bal);
    } catch (e) {
      console.error("Failed to load balance", e);
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

  const items = [
    {
      key: 'smart-deposit',
      label: (
        <Flex gap={8} align="center">
          <ArrowDownToLine size={16} />
          Topup Credits (Smart Deposit)
        </Flex>
      ),
      children: (
        <div style={{ background: '#1e1e1e', padding: 32, borderRadius: 12 }}>
          <h2>Smart Deposit</h2>
          <p style={{ color: '#aaa', marginBottom: 24, lineHeight: '1.6' }}>
            <b>One-Click Deposit</b> allows you to instantly convert your internal wallet's USDT balance into Service Credits.
            <br />
            <br />
            ℹ️ <b>Note:</b> You must have a small amount of <b>BNB</b> in your wallet to pay for gas fees.
          </p>
          <SmartDepositForm onDeposit={handleDeposit} loading={loading} />
        </div>
      ),
    },
    {
      key: 'receive',
      label: (
        <Flex gap={8} align="center">
          <QrCode size={16} />
          Receive / External Deposit
        </Flex>
      ),
      children: (
        children: (
          <div style = {{ background: '#1e1e1e', padding: 32, borderRadius: 12, textAlign: 'center' }} >
            <h2>Receive USDT</h2>
            <p style={{ color: '#aaa', marginBottom: 32 }}>
                Scan this QR code to deposit USDT (BEP20) or BNB to your wallet.
            </p>
            
            <Flex justify="center" style={{ marginBottom: 24 }}>
                <div style={{ background: '#fff', padding: 16, borderRadius: 12 }}>
                    <QRCode value={address || "0x"} size={200} />
                </div>
            </Flex>

            <div style={{ background: '#2c2c2c', padding: '16px', borderRadius: 8, display: 'inline-block', maxWidth: '100%', wordBreak: 'break-all' }}>
                <p style={{ color: '#888', fontSize: 12, marginBottom: 4 }}>Your Wallet Address (BSC Testnet)</p>
                <Flex gap={10} align="center" justify="center">
                    <span style={{ fontFamily: 'monospace', fontSize: 16, color: '#fff' }}>{address}</span>
                    <CopyButton text={address} />
                </Flex>
            </div>
        </div >
      ),
    },
  ];

return (
  <Container justify={'center'} align={'center'}>
    <ContentWrapper>
      <h1>My Wallet</h1>
      <Tabs defaultActiveKey="smart-deposit" items={items} />
    </ContentWrapper>
  </Container>
);
});

const CopyButton = ({ text }: { text: string }) => {
  const { message } = App.useApp();
  return (
    <div
      style={{ cursor: 'pointer', color: '#1677ff' }}
      onClick={() => {
        navigator.clipboard.writeText(text);
        message.success("Address copied to clipboard!");
      }}
    >
      <Copy size={16} />
    </div>
  )
}

// Helper component for Smart Deposit input
const SmartDepositForm = ({ onDeposit, loading }: { onDeposit: (val: number) => void, loading: boolean }) => {
  const [val, setVal] = useState(10);
  return (
    <Flex gap={10} align='center'>
      <input
        type="number"
        value={val}
        onChange={e => setVal(Number(e.target.value))}
        style={{ padding: '8px 12px', borderRadius: 6, border: '1px solid #444', background: '#333', color: '#fff', fontSize: 16, width: 120 }}
      />
      <span style={{ fontSize: 16, fontWeight: 500 }}>USDT</span>
      <button
        disabled={loading}
        onClick={() => onDeposit(val)}
        style={{
          padding: '8px 24px',
          borderRadius: 6,
          background: '#1677ff',
          color: '#fff',
          border: 'none',
          cursor: loading ? 'not-allowed' : 'pointer',
          opacity: loading ? 0.7 : 1,
          fontSize: 16,
          fontWeight: 500,
          marginLeft: 16
        }}
      >
        {loading ? 'Processing...' : 'One-Click Deposit'}
      </button>
    </Flex>
  )
}

export default DesktopWalletLayout;
