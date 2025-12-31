import { Modal, QRCode, App, Button, InputNumber, Form, Input, Empty, Popover, Spin, Tooltip, Select } from 'antd';
import { memo, useEffect, useState, useCallback } from 'react';
import { DeAIService, WalletService, JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { ListWalletTransactions } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import type { WalletTransaction } from '@@/github.com/kawai-network/veridium/internal/database/generated/models';
import { useUserStore } from '@/store/user';
import {
  Copy,
  Send,
  Settings,
  Check,
  Gift,
  Home,
  ShoppingCart,
  Globe,
  Repeat2,
} from 'lucide-react';
import { ActionIcon, Icon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { createStyles, useTheme } from 'antd-style';
import { NetworkIcon } from './NetworkIcons';
import { TokenUSDT } from '@web3icons/react';

import Menu from '@/components/Menu';
import PanelTitle from '@/components/PanelTitle';
import WalletSidePanel from './WalletSidePanel';
import WalletAccountPanel from './WalletAccountPanel';
import AccountList from './AccountList';

// Import content components from separate files
import HomeContent from './HomeContent';
import OTCContent from './OTCContent';
import RewardsContent from './RewardsContent';
import SettingsContent from './SettingsContent';

type MenuKey = 'home' | 'otc' | 'rewards' | 'settings';

const useStyles = createStyles(({ css, token }) => ({
  container: css`
    flex: 1;
    display: flex;
    flex-direction: row;
    overflow: hidden;
  `,
  content: css`
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    background: ${token.colorBgLayout};
  `,
  balanceCard: css`
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
    border: 1px solid rgba(102, 126, 234, 0.3);
    border-radius: 20px;
    position: relative;
    overflow: hidden;
    
    /* Subtle grid pattern */
    background-image: 
      linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%),
      repeating-linear-gradient(
        45deg,
        transparent,
        transparent 20px,
        rgba(255,255,255,0.01) 20px,
        rgba(255,255,255,0.01) 40px
      );
    
    &::before {
      content: '';
      position: absolute;
      top: -50%;
      right: -20%;
      width: 300px;
      height: 300px;
      background: radial-gradient(circle, rgba(102, 126, 234, 0.4) 0%, transparent 60%);
      border-radius: 50%;
    }
    
    &::after {
      content: '';
      position: absolute;
      bottom: -30%;
      left: -10%;
      width: 200px;
      height: 200px;
      background: radial-gradient(circle, rgba(118, 75, 162, 0.25) 0%, transparent 60%);
      border-radius: 50%;
    }

    .ant-card-body {
      padding: 28px;
      position: relative;
      z-index: 1;
    }
  `,
  actionButton: css`
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    transition: all 0.3s ease;
    padding: 16px 24px;
    border-radius: 16px;
    min-width: 80px;
    
    &:hover {
      background: ${token.colorFillTertiary};
      transform: translateY(-3px);
    }
  `,
  actionCircle: css`
    width: 56px;
    height: 56px;
    border-radius: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
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

const MenuContent = memo<{
  activeMenu: MenuKey;
  setActiveMenu: (key: MenuKey) => void;
  styles: any;
  currentNetwork: NetworkInfo | null;
  onNetworkChange: (network: NetworkInfo) => void;
}>(({ activeMenu, setActiveMenu, styles, currentNetwork, onNetworkChange }) => {
  const menuItems = [
    { key: 'home', icon: <Icon icon={Home} />, label: 'Home' },
    { key: 'otc', icon: <Icon icon={ShoppingCart} />, label: 'OTC Market' },
    { key: 'rewards', icon: <Icon icon={Gift} />, label: 'Rewards' },
    { key: 'settings', icon: <Icon icon={Settings} />, label: 'Settings' },
  ];

  return (
    <Flexbox gap={16} height={'100%'}>
      <Flexbox paddingInline={8}>
        <PanelTitle desc="Manage your digital assets and transactions" title="Wallet" />
        <Menu
          compact
          selectable
          items={menuItems}
          selectedKeys={[activeMenu]}
          onClick={({ key }) => setActiveMenu(key as MenuKey)}
        />
      </Flexbox>
      <div style={{ flex: 1 }} />
      <Flexbox padding={12}>
        <NetworkSwitcher currentNetwork={currentNetwork} onNetworkChange={onNetworkChange} />
      </Flexbox>
    </Flexbox>
  );
});

interface NetworkSwitcherProps {
  currentNetwork: NetworkInfo | null;
  onNetworkChange: (network: NetworkInfo) => void;
}

const NetworkSwitcher = memo<NetworkSwitcherProps>(({ currentNetwork, onNetworkChange }) => {
  const theme = useTheme();
  const [networks, setNetworks] = useState<NetworkInfo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadNetworks();
  }, []);

  const loadNetworks = async () => {
    try {
      const supportedNetworks = await JarvisService.GetSupportedNetworks();
      // Sort: testnets first (for development), then mainnets
      const sorted = supportedNetworks.sort((a, b) => {
        if (a.isTestnet !== b.isTestnet) return a.isTestnet ? -1 : 1;
        return a.name.localeCompare(b.name);
      });
      setNetworks(sorted);
    } catch (e) {
      console.error('Failed to load networks', e);
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
          <div style={{ padding: '4px 8px', fontSize: 11, color: theme.colorTextTertiary, textTransform: 'uppercase', position: 'sticky', top: 0, background: theme.colorBgElevated }}>
            Switch Network
          </div>
          {loading ? (
            <Flexbox align="center" justify="center" style={{ padding: 20 }}>
              <Spin size="small" />
            </Flexbox>
          ) : (
            <>
              {/* Testnets Section */}
              <div style={{ padding: '4px 8px', fontSize: 10, color: theme.colorTextQuaternary, marginTop: 8 }}>TESTNETS</div>
              {networks.filter(n => n.isTestnet).map((network) => (
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
                    <span style={{ fontSize: 13, fontWeight: currentNetwork?.id === network.id ? 600 : 400 }}>{network.name}</span>
                    <span style={{ fontSize: 10, color: theme.colorTextTertiary }}>{network.nativeTokenSymbol}</span>
                  </Flexbox>
                  {currentNetwork?.id === network.id && <Check size={14} color={theme.colorSuccess} />}
                </Flexbox>
              ))}

              {/* Mainnets Section */}
              <div style={{ padding: '4px 8px', fontSize: 10, color: theme.colorTextQuaternary, marginTop: 8 }}>MAINNETS</div>
              {networks.filter(n => !n.isTestnet).map((network) => (
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
                    <span style={{ fontSize: 13, fontWeight: currentNetwork?.id === network.id ? 600 : 400 }}>{network.name}</span>
                    <span style={{ fontSize: 10, color: theme.colorTextTertiary }}>{network.nativeTokenSymbol}</span>
                  </Flexbox>
                  {currentNetwork?.id === network.id && <Check size={14} color={theme.colorSuccess} />}
                </Flexbox>
              ))}
            </>
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
          background: '#22c55e',
          boxShadow: '0 0 8px rgba(34, 197, 94, 0.5)'
        }} />
        <Flexbox flex={1}>
          <div style={{ fontSize: 10, color: theme.colorTextTertiary, lineHeight: 1 }}>Network</div>
          <div style={{ fontSize: 12, fontWeight: 600, display: 'flex', alignItems: 'center', gap: 4 }}>
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

// Default Monad Testnet chain ID

const DEFAULT_CHAIN_ID = 10143;
const KAWAI_TOKEN_ADDRESS = "0x3EC7A3b85f9658120490d5a76705d4d304f4068D";

const DesktopWalletLayout = memo(() => {
  const { styles, theme } = useStyles();
  const [activeMenu, setActiveMenu] = useState<MenuKey>('home');
  // Use global wallet address from store
  const address = useUserStore((s) => s.walletAddress);
  const [balance, setBalance] = useState('0.00'); // USDT Balance (Vault)
  const [nativeBalance, setNativeBalance] = useState('0.00');
  const [kawaiBalance, setKawaiBalance] = useState('0.00');
  const [balanceVisible, setBalanceVisible] = useState(true);
  const [loading, setLoading] = useState(false);
  const [modalType, setModalType] = useState<'send' | 'receive' | 'swap' | 'deposit' | 'addAccount' | 'createWallet' | 'importWallet' | 'addToken' | null>(null);
  const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
  const [pendingSwitch, setPendingSwitch] = useState<string | null>(null);
  const { message } = App.useApp();
  const { isWalletLoaded, refreshWalletStatus } = useUserStore();

  // Network state
  const [currentNetwork, setCurrentNetwork] = useState<NetworkInfo | null>(null);
  const [gasEstimate, setGasEstimate] = useState<{ maxGasPriceGwei: number; maxTipGwei: number } | null>(null);
  const [currentBlock, setCurrentBlock] = useState<number>(0);

  // Loading states
  const [balancesLoading, setBalancesLoading] = useState(false);

  // Initialize default network
  useEffect(() => {
    initializeNetwork();
  }, []);

  const initializeNetwork = async () => {
    try {
      const network = await JarvisService.GetNetworkByID(DEFAULT_CHAIN_ID);
      if (network) {
        setCurrentNetwork(network);
      }
    } catch (e) {
      console.error('Failed to initialize network', e);
    }
  };

  const handleNetworkChange = useCallback(async (network: NetworkInfo) => {
    setCurrentNetwork(network);
    message.info(`Switched to ${network.name}`);
    // Reload balances for new network
    if (address) {
      loadNativeBalance(network.id);
      loadKawaiBalance(network.id);
      loadGasEstimate(network.id);
      loadCurrentBlock(network.id);
    }
  }, [address, message]);

  useEffect(() => {
    // Only load data if address is available
    if (address && currentNetwork) {
      setBalancesLoading(true);
      Promise.all([
        loadBalance(),
        loadHistory(),
        loadNativeBalance(currentNetwork.id),
        loadKawaiBalance(currentNetwork.id),
        loadGasEstimate(currentNetwork.id),
        loadCurrentBlock(currentNetwork.id),
      ]).finally(() => setBalancesLoading(false));
    }
  }, [address, isWalletLoaded, currentNetwork]);

  const loadNativeBalance = async (networkId: number) => {
    if (!address) return;
    try {
      const result = await JarvisService.GetNativeBalance(address, networkId);
      if (result) {
        setNativeBalance(result.formatted);
      }
    } catch (e) {
      console.error('Failed to load native balance', e);
    }
  };

  const loadKawaiBalance = async (networkId: number) => {
    if (!address) return;
    try {
      const result = await JarvisService.GetTokenBalance(KAWAI_TOKEN_ADDRESS, address, networkId);
      if (result) {
        setKawaiBalance(result.formatted);
      }
    } catch (e) {
      console.error('Failed to load KAWAI balance', e);
      setKawaiBalance('0.00');
    }
  };

  const loadGasEstimate = async (networkId: number) => {
    try {
      const result = await JarvisService.EstimateGas(networkId);
      if (result) {
        setGasEstimate(result);
      }
    } catch (e) {
      console.error('Failed to load gas estimate', e);
    }
  };

  const loadCurrentBlock = async (networkId: number) => {
    try {
      const block = await JarvisService.GetCurrentBlock(networkId);
      setCurrentBlock(block);
    } catch (e) {
      console.error('Failed to load current block', e);
    }
  };

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

  const handleSend = async (to: string, amount: number, assetType: string, customTokenAddress?: string) => {
    setLoading(true);
    const hide = message.loading(`Sending ${assetType.toUpperCase()}...`, 0);
    try {
      let tx = '';
      if (assetType === 'native') {
        tx = await DeAIService.TransferNative(to, amount.toString());
      } else if (assetType === 'usdt') {
        // USDT has 6 decimals
        const rawAmount = Math.floor(amount * 1_000_000).toString();
        tx = await DeAIService.TransferUSDT(to, rawAmount);
      } else if (assetType === 'kawai') {
        // Bug Fix #4: Use high-precision arithmetic for 18 decimals
        // Convert amount to string with full precision, then to BigInt
        // This avoids JavaScript Number precision limits (53 bits)
        const amountStr = amount.toFixed(18); // Ensure 18 decimal places
        const [intPart, decPart = '0'] = amountStr.split('.');
        const paddedDec = decPart.padEnd(18, '0').substring(0, 18);
        const rawAmount = (BigInt(intPart) * BigInt(10 ** 18) + BigInt(paddedDec)).toString();
        tx = await DeAIService.TransferToken(KAWAI_TOKEN_ADDRESS, to, rawAmount);
      } else if (customTokenAddress) {
        // Custom token transfer (assume 18 decimals)
        const amountStr = amount.toFixed(18);
        const [intPart, decPart = '0'] = amountStr.split('.');
        const paddedDec = decPart.padEnd(18, '0').substring(0, 18);
        const rawAmount = (BigInt(intPart) * BigInt(10 ** 18) + BigInt(paddedDec)).toString();
        tx = await DeAIService.TransferToken(customTokenAddress, to, rawAmount);
      } else {
        throw new Error("Unknown asset type");
      }

      message.success(`Transfer Successful! TX: ${tx.substring(0, 10)}...`);
      loadBalance();
      if (currentNetwork) {
        loadNativeBalance(currentNetwork.id);
        loadKawaiBalance(currentNetwork.id);
      }
      loadHistory();
      setModalType(null);
    } catch (e: any) {
      // Bug Fix #5: Enhanced error handling with specific error types
      console.error('Transfer error:', e);

      let errorMessage = 'Transfer Failed';

      // Check for specific error types
      if (e.message) {
        const msg = e.message.toLowerCase();

        if (msg.includes('insufficient funds') || msg.includes('insufficient balance')) {
          errorMessage = 'Insufficient balance for this transfer';
        } else if (msg.includes('invalid address') || msg.includes('invalid recipient')) {
          errorMessage = 'Invalid recipient address';
        } else if (msg.includes('gas') && msg.includes('required exceeds allowance')) {
          errorMessage = 'Insufficient gas. Please add more native tokens';
        } else if (msg.includes('user rejected') || msg.includes('user denied')) {
          errorMessage = 'Transaction cancelled by user';
        } else if (msg.includes('nonce')) {
          errorMessage = 'Transaction nonce error. Please try again';
        } else if (msg.includes('timeout') || msg.includes('deadline')) {
          errorMessage = 'Transaction timeout. Please try again';
        } else if (msg.includes('network') || msg.includes('connection')) {
          errorMessage = 'Network error. Please check your connection';
        } else {
          // Use the original error message if it's descriptive
          errorMessage = `Transfer Failed: ${e.message}`;
        }
      } else {
        errorMessage = `Transfer Failed: ${e.toString()}`;
      }

      message.error(errorMessage, 5); // Show for 5 seconds
    } finally {
      hide();
      setLoading(false);
    }
  };

  const handleSwitchAccount = async (password: string) => {
    if (!pendingSwitch) return;
    setLoading(true);
    try {
      await WalletService.SwitchWallet(pendingSwitch, password);
      message.success("Switched account successfully");
      setPendingSwitch(null);
      await refreshWalletStatus();
      loadBalance();
      loadHistory();
    } catch (e: any) {
      message.error(e.message || "Failed to switch account");
    } finally {
      setLoading(false);
    }
  };

  const renderContent = () => {
    switch (activeMenu) {
      case 'home':
        return <HomeContent
          address={address}
          balance={balance}
          nativeBalance={nativeBalance}
          kawaiBalance={kawaiBalance}
          balanceVisible={balanceVisible}
          setBalanceVisible={setBalanceVisible}
          setModalType={setModalType}
          transactions={transactions}
          styles={styles}
          theme={theme}
          currentNetwork={currentNetwork}
          gasEstimate={gasEstimate}
          currentBlock={currentBlock}
          balancesLoading={balancesLoading}
        />;
      case 'otc':
        return <OTCContent styles={styles} theme={theme} />;
      case 'rewards':
        return <RewardsContent styles={styles} theme={theme} currentNetwork={currentNetwork} transactions={transactions} />;
      case 'settings':
        return <SettingsContent address={address} styles={styles} theme={theme} currentNetwork={currentNetwork} />;
      default:
        return null;
    }
  };

  return (
    <div className={styles.container}>
      {/* Sidebar */}
      <WalletSidePanel>
        <MenuContent
          activeMenu={activeMenu}
          setActiveMenu={setActiveMenu}
          styles={styles}
          currentNetwork={currentNetwork}
          onNetworkChange={handleNetworkChange}
        />
      </WalletSidePanel>

      {/* Content */}
      <div className={styles.content}>
        {renderContent()}
      </div>

      {/* Right Sidebar: Account Management */}
      <WalletAccountPanel>
        <AccountList
          activeAddress={address}
          onAccountSwitch={(addr) => setPendingSwitch(addr)}
          onAddAccount={() => setModalType('addAccount')}
        />
      </WalletAccountPanel>

      {/* Modals */}
      <Modal
        title="Confirm Switch Account"
        open={!!pendingSwitch}
        onCancel={() => setPendingSwitch(null)}
        footer={null}
      >
        <Form layout="vertical" onFinish={(v) => handleSwitchAccount(v.password)}>
          <p style={{ color: theme.colorTextSecondary, marginBottom: 16 }}>
            Please enter your wallet password to unlock: <br />
            <code style={{ fontSize: 11 }}>{pendingSwitch}</code>
          </p>
          <Form.Item label="Password" name="password" rules={[{ required: true }]}>
            <Input.Password autoFocus placeholder="Enter password" size="large" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block size="large" loading={loading}>
            Unlock & Switch
          </Button>
        </Form>
      </Modal>

      <Modal
        title="Add New Account"
        open={modalType === 'addAccount'}
        onCancel={() => setModalType(null)}
        footer={null}
      >
        <Flexbox gap={12}>
          <Button block size="large" onClick={() => { setModalType('createWallet'); }}>Create New Wallet</Button>
          <Button block size="large" onClick={() => { setModalType('importWallet'); }}>Import Existing Wallet (Keystore)</Button>
        </Flexbox>
      </Modal>

      <Modal title="Create Wallet" open={modalType === 'createWallet'} onCancel={() => setModalType(null)} footer={null}>
        <SetupForm type="create" onSuccess={() => { setModalType(null); refreshWalletStatus(); }} />
      </Modal>

      <Modal title="Import Wallet" open={modalType === 'importWallet'} onCancel={() => setModalType(null)} footer={null}>
        <SetupForm type="import" onSuccess={() => { setModalType(null); refreshWalletStatus(); }} />
      </Modal>
      <Modal title="Smart Deposit" open={modalType === 'deposit'} onCancel={() => setModalType(null)} footer={null} destroyOnHidden>
        <SmartDepositForm onDeposit={handleDeposit} loading={loading} />
      </Modal>

      <Modal title="Send Assets" open={modalType === 'send'} onCancel={() => setModalType(null)} footer={null} destroyOnHidden>
        <SendForm onSend={handleSend} loading={loading} currentNetwork={currentNetwork} />
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

      <Modal
        title="Add Token"
        open={modalType === 'addToken'}
        onCancel={() => setModalType(null)}
        footer={null}
        width={450}
      >
        <AddTokenModal
          currentNetwork={currentNetwork}
          onClose={() => setModalType(null)}
        />
      </Modal>
    </div>
  );
});

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

const SendForm = ({ onSend, loading, currentNetwork }: { onSend: (to: string, val: number, asset: string, customAddr?: string) => void, loading: boolean, currentNetwork?: NetworkInfo | null }) => {
  const [form] = Form.useForm();
  const [selectedAsset, setSelectedAsset] = useState('usdt');
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [pendingTx, setPendingTx] = useState<{ to: string; amount: number; assetType: string; customAddr?: string } | null>(null);
  const theme = useTheme();

  const assetOptions = [
    { label: `Native Token (${currentNetwork?.nativeTokenSymbol || 'ETH'})`, value: 'native' },
    { label: 'USDT (Tether)', value: 'usdt' },
    { label: 'KAWAI (Kawai Token)', value: 'kawai' },
  ];

  const getAssetLabel = (assetType: string) => {
    if (assetType === 'native') return currentNetwork?.nativeTokenSymbol || 'ETH';
    if (assetType === 'usdt') return 'USDT';
    if (assetType === 'kawai') return 'KAWAI';
    return assetType.toUpperCase();
  };

  const handleFinish = (values: any) => {
    let assetType = selectedAsset;
    let customAddr: string | undefined = undefined;

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
              <span style={{ fontWeight: 600 }}>{pendingTx.amount} {getAssetLabel(selectedAsset)}</span>
            </Flexbox>
            <Flexbox horizontal justify="space-between">
              <span style={{ color: theme.colorTextSecondary }}>Network</span>
              <span>{currentNetwork?.name || 'Unknown'}</span>
            </Flexbox>
          </Flexbox>
        </div>

        <div style={{ fontSize: 12, color: theme.colorTextTertiary, textAlign: 'center' }}>
          Please review the transaction details before confirming.
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
          options={assetOptions}
          value={selectedAsset}
          onChange={setSelectedAsset}
          size="large"
        />
      </Form.Item>
      <Form.Item label="Recipient Address" name="to" rules={[{ required: true, message: 'Address is required' }, { pattern: /^0x[a-fA-F0-9]{40}$/, message: 'Invalid EVM address' }]}>
        <Input placeholder="0x..." size="large" />
      </Form.Item>
      <Form.Item label="Amount" name="amount" rules={[{ required: true, type: 'number', min: 0.000001 }]}>
        <InputNumber style={{ width: '100%' }} size="large" min={0.000001} />
      </Form.Item>
      <Button type="primary" htmlType="submit" block size="large" style={{ marginTop: 16 }}>
        Review Transaction
      </Button>
    </Form>
  );
};

// ============ ADD TOKEN MODAL ============
const AddTokenModal = memo<{ currentNetwork: NetworkInfo | null; onClose: () => void }>(({ currentNetwork, onClose }) => {
  const theme = useTheme();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(true);
  const [projectTokens, setProjectTokens] = useState<Array<{ address: string; name: string; symbol: string }>>([]);

  useEffect(() => {
    loadProjectTokens();
  }, []);

  const loadProjectTokens = async () => {
    try {
      // Check if method exists (bindings may need regeneration after Go update)
      if (typeof (JarvisService as any).GetProjectTokens === 'function') {
        const tokens = await (JarvisService as any).GetProjectTokens();
        setProjectTokens(tokens);
      } else {
        // Fallback: hardcoded project tokens (Monad Testnet)
        setProjectTokens([
          { address: '0xa6Fc4FaF4CD7a4E3f300D164a37CB45d35bf28eD', name: 'MockUSDT', symbol: 'USDT' },
          { address: '0x3EC7A3b85f9658120490d5a76705d4d304f4068D', name: 'KawaiToken', symbol: 'KAWAI' },
        ]);
      }
    } catch (e) {
      console.error('Failed to load project tokens', e);
      // Fallback on error
      setProjectTokens([
        { address: '0xa6Fc4FaF4CD7a4E3f300D164a37CB45d35bf28eD', name: 'MockUSDT', symbol: 'USDT' },
        { address: '0x3EC7A3b85f9658120490d5a76705d4d304f4068D', name: 'KawaiToken', symbol: 'KAWAI' },
      ]);
    } finally {
      setLoading(false);
    }
  };

  const handleCopyAddress = (address: string) => {
    navigator.clipboard.writeText(address);
    message.success('Address copied!');
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
                  background: token.symbol === 'USDT' ? '#26a17b' : 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontWeight: 700,
                  fontSize: 12,
                }}>
                  {token.symbol === 'USDT' ? <TokenUSDT size={24} variant="branded" /> : token.symbol.substring(0, 2)}
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
        <strong>Note:</strong> These are the official project tokens deployed on Monad Testnet.
        Custom token import is not supported yet.
      </div>

      <Button block onClick={onClose}>Close</Button>
    </Flexbox>
  );
});

export default DesktopWalletLayout;
const SetupForm = memo<{ type: 'create' | 'import'; onSuccess: () => void }>(({ type, onSuccess }) => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [mnemonic, setMnemonic] = useState('');
  const [description, setDescription] = useState('');
  const [step, setStep] = useState<'form' | 'mnemonic'>(type === 'create' ? 'mnemonic' : 'form');
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();

  const { generateMnemonic, createWallet } = useUserStore();

  useEffect(() => {
    if (type === 'create' && step === 'mnemonic') {
      generateMnemonic().then(setMnemonic);
    }
  }, [type, step]);

  const handleFinish = async () => {
    if (password !== confirmPassword) {
      return message.error("Passwords do not match");
    }
    if (password.length < 8) {
      return message.error("Password too short");
    }
    setLoading(true);
    try {
      await createWallet(password, mnemonic, description);
      message.success("Wallet created successfully");
      onSuccess();
    } catch (e: any) {
      message.error(e.message || "Failed to create wallet");
    } finally {
      setLoading(false);
    }
  };

  if (type === 'create' && step === 'mnemonic') {
    return (
      <Flexbox gap={12}>
        <p>Save these 12 words securely:</p>
        <div style={{ background: 'rgba(0,0,0,0.05)', padding: 16, borderRadius: 8, textAlign: 'center' }}>
          <code style={{ fontSize: 16, fontWeight: 700 }}>{mnemonic}</code>
          <div style={{ marginTop: 8 }}>
            <CopyButton text={mnemonic} />
          </div>
        </div>
        <Button type="primary" block onClick={() => setStep('form')}>I have written it down</Button>
      </Flexbox>
    );
  }

  return (
    <Form layout="vertical" onFinish={handleFinish}>
      {type === 'import' && (
        <Form.Item label="Mnemonic Phrase" required>
          <Input.TextArea rows={3} value={mnemonic} onChange={e => setMnemonic(e.target.value)} placeholder="word1 word2 ..." />
        </Form.Item>
      )}
      <Form.Item label="Account Description" name="description">
        <Input placeholder="Main account, Savings, etc." value={description} onChange={e => setDescription(e.target.value)} />
      </Form.Item>
      <Form.Item label="Lock Password" required>
        <Input.Password value={password} onChange={e => setPassword(e.target.value)} placeholder="At least 8 characters" />
      </Form.Item>
      <Form.Item label="Confirm Password" required>
        <Input.Password value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} placeholder="Repeat password" />
      </Form.Item>
      <Button type="primary" htmlType="submit" block size="large" loading={loading}>
        {type === 'create' ? 'Create Account' : 'Import Account'}
      </Button>
    </Form>
  );
});
