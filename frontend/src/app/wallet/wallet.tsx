import { Card, Modal, QRCode, App, Table, Tag, Button, InputNumber, Form, Input, Empty, Popover, Spin, Tooltip } from 'antd';
import { memo, useEffect, useState, useCallback } from 'react';
import { DeAIService, WalletService, JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { NetworkInfo, TokenInfo, BalanceInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { ListWalletTransactions } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import type { WalletTransaction } from '@@/github.com/kawai-network/veridium/internal/database/generated/models';
import { useUserStore } from '@/store/user';
import { ArrowDownToLine, Copy, Send, Eye, EyeOff, Repeat2, History, Home, ShoppingCart, Gift, Settings, Coins, ExternalLink, Globe, Plus, Check, X, Loader2, Fuel } from 'lucide-react';
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

// Network icons mapping
// Network icons mapping removed in favor of dynamic NetworkIcon

// getNetworkIconComponent removed

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
              {networks.filter(n => !n.isTestnet).slice(0, 10).map((network) => (
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

// Custom token interface
interface CustomToken {
  address: string;
  symbol: string;
  name: string;
  decimals: number;
  balance: string;
  isCustom: boolean;
}

// Default Monad Testnet chain ID
const DEFAULT_CHAIN_ID = 10143;

const DesktopWalletLayout = memo(() => {
  const { styles, theme } = useStyles();
  const [activeMenu, setActiveMenu] = useState<MenuKey>('home');
  // Use global wallet address from store
  const address = useUserStore((s) => s.walletAddress);
  const [balance, setBalance] = useState('0.00');
  const [nativeBalance, setNativeBalance] = useState('0.00');
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

  // Custom tokens state
  const [customTokens, setCustomTokens] = useState<CustomToken[]>([]);

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
      loadGasEstimate(network.id);
      loadCurrentBlock(network.id);
      // Reload custom token balances
      loadCustomTokenBalances(network.id);
    }
  }, [address, message]);

  useEffect(() => {
    // Only load data if address is available
    if (address && currentNetwork) {
      loadBalance();
      loadHistory();
      loadNativeBalance(currentNetwork.id);
      loadGasEstimate(currentNetwork.id);
      loadCurrentBlock(currentNetwork.id);
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

  const loadCustomTokenBalances = async (networkId: number) => {
    if (!address || customTokens.length === 0) return;

    const updatedTokens = await Promise.all(
      customTokens.map(async (token) => {
        try {
          const result = await JarvisService.GetTokenBalance(token.address, address, networkId);
          return { ...token, balance: result?.formatted || '0' };
        } catch {
          return { ...token, balance: '0' };
        }
      })
    );
    setCustomTokens(updatedTokens);
  };

  const handleAddToken = async (tokenAddress: string) => {
    if (!currentNetwork) {
      message.error('No network selected');
      return;
    }

    setLoading(true);
    try {
      // Fetch token info from blockchain
      const tokenInfo = await JarvisService.GetTokenInfo(tokenAddress, currentNetwork.id);
      if (!tokenInfo) {
        message.error('Could not fetch token info. Is this a valid ERC20 token?');
        return;
      }

      // Check if already added
      if (customTokens.some(t => t.address.toLowerCase() === tokenAddress.toLowerCase())) {
        message.warning('Token already added');
        return;
      }

      // Get balance
      let tokenBalance = '0';
      if (address) {
        const balanceResult = await JarvisService.GetTokenBalance(tokenAddress, address, currentNetwork.id);
        tokenBalance = balanceResult?.formatted || '0';
      }

      const newToken: CustomToken = {
        address: tokenAddress,
        symbol: tokenInfo.symbol,
        name: tokenInfo.name,
        decimals: Number(tokenInfo.decimals),
        balance: tokenBalance,
        isCustom: true,
      };

      setCustomTokens(prev => [...prev, newToken]);
      message.success(`Added ${tokenInfo.symbol} to your wallet`);
      setModalType(null);
    } catch (e: any) {
      message.error(`Failed to add token: ${e.message || e}`);
    } finally {
      setLoading(false);
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
          balanceVisible={balanceVisible}
          setBalanceVisible={setBalanceVisible}
          setModalType={setModalType}
          transactions={transactions}
          styles={styles}
          theme={theme}
          currentNetwork={currentNetwork}
          gasEstimate={gasEstimate}
          currentBlock={currentBlock}
          customTokens={customTokens}
        />;
      case 'otc':
        return <OTCContent styles={styles} theme={theme} />;
      case 'rewards':
        return <RewardsContent styles={styles} theme={theme} />;
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
            <Input.Password placeholder="Enter password" size="large" />
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

      <Modal title="Add Custom Token" open={modalType === 'addToken'} onCancel={() => setModalType(null)} footer={null} destroyOnClose>
        <AddTokenForm onAddToken={handleAddToken} loading={loading} currentNetwork={currentNetwork} />
      </Modal>
    </div>
  );
});

// ============ HOME CONTENT ============
const HomeContent = ({ address, balance, nativeBalance, balanceVisible, setBalanceVisible, setModalType, transactions, styles, theme, currentNetwork, gasEstimate, currentBlock, customTokens }: any) => {
  return (
    <Flexbox style={{ maxWidth: 900, width: '100%' }} gap={20}>
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
            {/* Native token balance */}
            {currentNetwork && (
              <div style={{ fontSize: 13, color: theme.colorTextSecondary, marginTop: 4 }}>
                {balanceVisible ? nativeBalance : '••••'} {currentNetwork.nativeTokenSymbol}
              </div>
            )}
          </Flexbox>
          {/* Network & Gas Info */}
          <Flexbox gap={8} align="flex-end">
            {gasEstimate && (
              <Tooltip title={`Max Tip: ${gasEstimate.maxTipGwei.toFixed(2)} Gwei`}>
                <Flexbox horizontal align="center" gap={4} style={{
                  padding: '4px 8px',
                  background: 'rgba(255,255,255,0.1)',
                  borderRadius: 8,
                  fontSize: 11
                }}>
                  <Fuel size={12} />
                  <span>{gasEstimate.maxGasPriceGwei.toFixed(1)} Gwei</span>
                </Flexbox>
              </Tooltip>
            )}
            {currentBlock > 0 && (
              <div style={{ fontSize: 10, color: theme.colorTextTertiary }}>
                Block #{currentBlock.toLocaleString()}
              </div>
            )}
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Quick Actions */}
      <Flexbox horizontal gap={12} style={{ marginTop: 4 }}>
        {[
          { label: 'Send', icon: Send, color: '#06b6d4', action: () => setModalType('send') },
          { label: 'Receive', icon: ArrowDownToLine, color: '#22c55e', action: () => setModalType('receive') },
          { label: 'Swap', icon: Repeat2, color: '#eab308', action: () => setModalType('swap') },
        ].map((item) => (
          <div key={item.label} className={styles.actionButton} onClick={item.action}>
            <div className={styles.actionCircle} style={{ background: `${item.color}20`, color: item.color }}>
              <item.icon size={24} />
            </div>
            <span style={{ fontWeight: 600, fontSize: 13 }}>{item.label}</span>
          </div>
        ))}
      </Flexbox>

      {/* Token List */}
      <Card
        title={<Flexbox horizontal align="center" gap={8}><Coins size={16} /> Tokens</Flexbox>}
        size="small"
        extra={
          <Button
            type="text"
            icon={<Plus size={14} />}
            size="small"
            onClick={() => setModalType('addToken')}
          >
            Add Token
          </Button>
        }
      >
        <Flexbox gap={8}>
          {/* Native Token */}
          {currentNetwork && (
            <div className={styles.tokenRow}>
              <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
                <div style={{
                  width: 36,
                  height: 36,
                  borderRadius: '50%',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontWeight: 800,
                  fontSize: 14,
                }}>
                  {currentNetwork && (
                    <NetworkIcon
                      name={currentNetwork.icon || 'ethereum'}
                      size={24}
                      variant="mono"
                      color="#fff"
                    />
                  )}
                </div>
                <div>
                  <div style={{ fontWeight: 600 }}>{currentNetwork.nativeTokenSymbol}</div>
                  <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Native Token</div>
                </div>
              </Flexbox>
              <div style={{ textAlign: 'right', minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>{balanceVisible ? nativeBalance : '••••'}</div>
              </div>
            </div>
          )}

          {/* USDT */}
          <div className={styles.tokenRow}>
            <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
              <div style={{
                width: 36,
                height: 36,
                borderRadius: '50%',
                background: '#26a17b',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                fontWeight: 800,
                fontSize: 16,
                fontFamily: 'Arial, sans-serif',
                textShadow: '0 1px 2px rgba(0,0,0,0.2)'
              }}>
                <TokenUSDT size={36} variant="branded" />
              </div>
              <div>
                <div style={{ fontWeight: 600 }}>USDT</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Tether USD</div>
              </div>
            </Flexbox>
            <Flexbox horizontal align="center" gap={16}>
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>$1.00</span>
              <div style={{ textAlign: 'right', minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>{balanceVisible ? balance : '••••'}</div>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>${balanceVisible ? balance : '••••'}</div>
              </div>
            </Flexbox>
          </div>

          {/* KAWAI */}
          <div className={styles.tokenRow}>
            <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
              <div style={{
                width: 36,
                height: 36,
                borderRadius: '50%',
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                fontWeight: 800,
                fontSize: 14,
                fontFamily: 'Arial, sans-serif',
                boxShadow: '0 2px 8px rgba(102, 126, 234, 0.3)'
              }}>🐱</div>
              <div>
                <div style={{ fontWeight: 600 }}>KAWAI</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Kawai Token</div>
              </div>
            </Flexbox>
            <Flexbox horizontal align="center" gap={16}>
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>$0.001</span>
              <div style={{ textAlign: 'right', minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>0</div>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>$0.00</div>
              </div>
            </Flexbox>
          </div>

          {/* Custom Tokens */}
          {customTokens && customTokens.map((token: CustomToken) => (
            <div key={token.address} className={styles.tokenRow}>
              <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
                <div style={{
                  width: 36,
                  height: 36,
                  borderRadius: '50%',
                  background: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontWeight: 700,
                  fontSize: 12,
                }}>{token.symbol.substring(0, 2).toUpperCase()}</div>
                <div>
                  <div style={{ fontWeight: 600 }}>{token.symbol}</div>
                  <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>{token.name}</div>
                </div>
              </Flexbox>
              <Flexbox horizontal align="center" gap={8}>
                <Tag color="blue" style={{ margin: 0, fontSize: 10 }}>Custom</Tag>
                <div style={{ textAlign: 'right', minWidth: 70 }}>
                  <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>{balanceVisible ? token.balance : '••••'}</div>
                </div>
              </Flexbox>
            </div>
          ))}
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
              {
                title: 'TX',
                dataIndex: 'txHash',
                key: 'txHash',
                render: (txHash) => txHash ? (
                  <TransactionLink txHash={txHash} networkId={currentNetwork?.id} />
                ) : '-'
              },
            ]}
          />
        ) : (
          <Flexbox align="center" gap={16} style={{ padding: '24px 0' }}>
            <Empty description={false} image={Empty.PRESENTED_IMAGE_SIMPLE} />
            <span style={{ color: theme.colorTextSecondary }}>No transactions yet</span>
            <Button
              type="primary"
              size="small"
              onClick={() => window.open('https://testnet.monad.xyz/faucet', '_blank')}
            >
              Get Test Tokens (Faucet)
            </Button>
          </Flexbox>
        )}
      </Card>
    </Flexbox>
  );
};

// Transaction Link with analysis popup
const TransactionLink = memo<{ txHash: string; networkId?: number }>(({ txHash, networkId }) => {
  const theme = useTheme();
  const [analyzing, setAnalyzing] = useState(false);
  const [analysis, setAnalysis] = useState<any>(null);

  const handleAnalyze = async () => {
    if (!networkId || analyzing) return;
    setAnalyzing(true);
    try {
      const result = await JarvisService.AnalyzeTransaction(txHash, networkId);
      setAnalysis(result);
    } catch (e) {
      console.error('Failed to analyze transaction', e);
    } finally {
      setAnalyzing(false);
    }
  };

  const shortHash = `${txHash.substring(0, 6)}...${txHash.substring(txHash.length - 4)}`;

  return (
    <Popover
      trigger="click"
      onOpenChange={(open) => open && handleAnalyze()}
      content={
        <div style={{ width: 300, maxHeight: 400, overflowY: 'auto' }}>
          {analyzing ? (
            <Flexbox align="center" justify="center" style={{ padding: 20 }}>
              <Spin size="small" />
              <span style={{ marginLeft: 8 }}>Analyzing...</span>
            </Flexbox>
          ) : analysis ? (
            <Flexbox gap={12}>
              <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>TRANSACTION ANALYSIS</div>

              <Flexbox gap={8}>
                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Status</span>
                  <Tag color={analysis.status === 'done' ? 'green' : analysis.status === 'reverted' ? 'red' : 'orange'}>
                    {analysis.status}
                  </Tag>
                </Flexbox>

                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Type</span>
                  <span style={{ fontSize: 12, fontWeight: 600 }}>{analysis.txType || 'Unknown'}</span>
                </Flexbox>

                {analysis.method && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Method</span>
                    <Tag color="blue" style={{ fontFamily: 'monospace' }}>{analysis.method}</Tag>
                  </Flexbox>
                )}

                {analysis.value && analysis.value !== '0' && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Value</span>
                    <span style={{ fontSize: 12 }}>{analysis.value}</span>
                  </Flexbox>
                )}

                {analysis.gasUsed && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Gas Used</span>
                    <span style={{ fontSize: 12 }}>{parseInt(analysis.gasUsed).toLocaleString()}</span>
                  </Flexbox>
                )}

                {analysis.gasCost && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Gas Cost</span>
                    <span style={{ fontSize: 12 }}>{analysis.gasCost}</span>
                  </Flexbox>
                )}

                {analysis.blockNumber > 0 && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Block</span>
                    <span style={{ fontSize: 12, fontFamily: 'monospace' }}>#{analysis.blockNumber.toLocaleString()}</span>
                  </Flexbox>
                )}
              </Flexbox>

              {/* Decoded Parameters */}
              {analysis.params && analysis.params.length > 0 && (
                <>
                  <div style={{ fontSize: 11, color: theme.colorTextTertiary, marginTop: 8 }}>PARAMETERS</div>
                  <Flexbox gap={4}>
                    {analysis.params.map((param: any, i: number) => (
                      <div key={i} style={{
                        padding: '4px 8px',
                        background: theme.colorFillTertiary,
                        borderRadius: 4,
                        fontSize: 11
                      }}>
                        <span style={{ color: theme.colorTextSecondary }}>{param.name}</span>
                        <span style={{ color: theme.colorTextTertiary }}> ({param.type})</span>
                        <div style={{ fontFamily: 'monospace', wordBreak: 'break-all', marginTop: 2 }}>
                          {param.value?.substring(0, 50)}{param.value?.length > 50 ? '...' : ''}
                        </div>
                      </div>
                    ))}
                  </Flexbox>
                </>
              )}

              {/* Event Logs */}
              {analysis.logs && analysis.logs.length > 0 && (
                <>
                  <div style={{ fontSize: 11, color: theme.colorTextTertiary, marginTop: 8 }}>EVENTS ({analysis.logs.length})</div>
                  <Flexbox gap={4}>
                    {analysis.logs.slice(0, 3).map((log: any, i: number) => (
                      <Tag key={i} color="purple">{log.name || 'Unknown Event'}</Tag>
                    ))}
                    {analysis.logs.length > 3 && (
                      <span style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                        +{analysis.logs.length - 3} more
                      </span>
                    )}
                  </Flexbox>
                </>
              )}

              {analysis.error && (
                <div style={{ color: theme.colorError, fontSize: 12, marginTop: 8 }}>
                  Error: {analysis.error}
                </div>
              )}
            </Flexbox>
          ) : (
            <div style={{ padding: 16, textAlign: 'center', color: theme.colorTextSecondary }}>
              Click to analyze transaction
            </div>
          )}
        </div>
      }
    >
      <span
        style={{
          fontFamily: 'monospace',
          fontSize: 11,
          cursor: 'pointer',
          color: theme.colorPrimary,
          textDecoration: 'underline'
        }}
      >
        {shortHash}
      </span>
    </Popover>
  );
});

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
const SettingsContent = ({ address, styles, theme, currentNetwork }: any) => {
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

const AddTokenForm = ({ onAddToken, loading, currentNetwork }: { onAddToken: (address: string) => void, loading: boolean, currentNetwork: NetworkInfo | null }) => {
  const [form] = Form.useForm();
  const [tokenPreview, setTokenPreview] = useState<TokenInfo | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const theme = useTheme();

  const handleAddressChange = async (address: string) => {
    if (!address || address.length !== 42 || !address.startsWith('0x') || !currentNetwork) {
      setTokenPreview(null);
      return;
    }

    setPreviewLoading(true);
    try {
      const isValid = await JarvisService.ValidateAddress(address);
      if (!isValid) {
        setTokenPreview(null);
        return;
      }

      const info = await JarvisService.GetTokenInfo(address, currentNetwork.id);
      setTokenPreview(info);
    } catch {
      setTokenPreview(null);
    } finally {
      setPreviewLoading(false);
    }
  };

  return (
    <Form form={form} layout="vertical" onFinish={(v) => onAddToken(v.address)}>
      <Form.Item
        label="Token Contract Address"
        name="address"
        rules={[
          { required: true, message: 'Token address is required' },
          { pattern: /^0x[a-fA-F0-9]{40}$/, message: 'Invalid contract address' }
        ]}
      >
        <Input
          placeholder="0x..."
          size="large"
          onChange={(e) => handleAddressChange(e.target.value)}
        />
      </Form.Item>

      {/* Token Preview */}
      {previewLoading && (
        <Flexbox align="center" justify="center" style={{ padding: 16 }}>
          <Loader2 size={20} className="animate-spin" />
          <span style={{ marginLeft: 8, color: theme.colorTextSecondary }}>Fetching token info...</span>
        </Flexbox>
      )}

      {tokenPreview && !previewLoading && (
        <div style={{
          padding: 16,
          background: theme.colorFillTertiary,
          borderRadius: 12,
          marginBottom: 16
        }}>
          <Flexbox horizontal align="center" gap={12}>
            <div style={{
              width: 40,
              height: 40,
              borderRadius: '50%',
              background: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#fff',
              fontWeight: 700,
              fontSize: 14,
            }}>
              {tokenPreview.symbol.substring(0, 2).toUpperCase()}
            </div>
            <Flexbox flex={1}>
              <div style={{ fontWeight: 600, fontSize: 16 }}>{tokenPreview.symbol}</div>
              <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>{tokenPreview.name}</div>
            </Flexbox>
            <Flexbox align="flex-end">
              <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>Decimals: {tokenPreview.decimals}</div>
              {tokenPreview.isKnown && (
                <Tag color="green" style={{ margin: 0, marginTop: 4 }}>
                  <Check size={10} /> Verified
                </Tag>
              )}
            </Flexbox>
          </Flexbox>
        </div>
      )}

      <div style={{ fontSize: 12, color: theme.colorTextSecondary, marginBottom: 16 }}>
        Adding token on: <strong>{currentNetwork?.name || 'Unknown Network'}</strong>
      </div>

      <Button
        type="primary"
        htmlType="submit"
        block
        size="large"
        loading={loading}
        disabled={!tokenPreview || previewLoading}
      >
        {tokenPreview ? `Add ${tokenPreview.symbol}` : 'Add Token'}
      </Button>
    </Form>
  );
};

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
