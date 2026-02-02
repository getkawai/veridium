import { Modal, QRCode, App, Button, Form, Input } from 'antd';
import { memo, useEffect, useState, useCallback } from 'react';
import { DeAIService, WalletService, DepositSyncService, JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { NetworkInfo, BackendConfig, GasEstimate } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { ListWalletTransactions } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';
import type { WalletTransaction } from '@@/github.com/kawai-network/veridium/internal/database/generated/models';
import { useUserStore } from '@/store/user';
import { Repeat2 } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { createStyles } from 'antd-style';
import { getBackendNetworkConfig, DEFAULT_CHAIN_ID } from '@/config/network';
import WalletSidePanel from './WalletSidePanel';
import WalletAccountPanel from './WalletAccountPanel';
import AccountList from './AccountList';

// Import content components from separate files
import HomeContent from './HomeContent';
import OTCContent from './OTCContent';
import RewardsContent from './RewardsContent';
import SettingsContent from './SettingsContent';

// Import refactored components
import {
  MenuContent,
  CopyButton,
  SmartDepositForm,
  SendForm,
  AddTokenModal,
  SetupForm,
} from './components';

import type { MenuKey } from './types';

const useStyles = createStyles(({ css, token, appearance }) => {
  const isDark = appearance === 'dark';
  return {
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
      background: ${token.colorBgContainer};
      border: 1px solid ${token.colorBorderSecondary};
      border-radius: 20px;
      position: relative;
      overflow: hidden;

      /* Dynamic gradient based on theme */
      background-image:
        linear-gradient(135deg, ${token.colorFillContent} 0%, ${token.colorFillQuaternary} 50%, ${token.colorFillSecondary} 100%),
        repeating-linear-gradient(
          45deg,
          transparent,
          transparent 20px,
          rgba(255,255,255,0.05) 20px,
          rgba(255,255,255,0.05) 40px
        );

      /* Glow effects using theme colors */
      &::before {
        content: '';
        position: absolute;
        top: -50%;
        right: -20%;
        width: 300px;
        height: 300px;
        background: radial-gradient(circle, ${isDark ? 'rgba(102, 126, 234, 0.4)' : token.colorPrimaryBg} 0%, transparent 60%);
        border-radius: 50%;
      }

      &::after {
        content: '';
        position: absolute;
        bottom: -30%;
        left: -10%;
        width: 200px;
        height: 200px;
        background: radial-gradient(circle, ${isDark ? 'rgba(118, 75, 162, 0.25)' : token.colorInfoBg} 0%, transparent 60%);
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
  };
});

const DesktopWalletLayout = memo(() => {
  const { styles, theme } = useStyles();
  const [activeMenu, setActiveMenu] = useState<MenuKey>('home');
  const address = useUserStore((s) => s.walletAddress);
  const [balance, setBalance] = useState('0.00');
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
  const [gasEstimate, setGasEstimate] = useState<GasEstimate | null>(null);
  const [currentBlock, setCurrentBlock] = useState<number>(0);

  // Loading states
  const [balancesLoading, setBalancesLoading] = useState(false);

  // Price state
  const [nativePrice, setNativePrice] = useState(0);
  const [kawaiPrice, setKawaiPrice] = useState(0);

  // Backend config state
  const [backendConfig, setBackendConfig] = useState<BackendConfig | null>(null);

  // Helper to get KAWAI token address from backend config
  const getKawaiTokenAddress = useCallback((networkId?: number): string => {
    if (!backendConfig) return '';

    if (networkId === 10143 || backendConfig.environment === 'testnet') {
      return backendConfig.contracts.kawai || '';
    }

    if (networkId === 143 && backendConfig.environment === 'mainnet') {
      return backendConfig.contracts.kawai || '';
    }

    return backendConfig.contracts.kawai || '';
  }, [backendConfig]);

  // Initialize default network and backend config
  useEffect(() => {
    loadBackendConfig();
  }, []);

  // Re-initialize network when backend config is loaded
  useEffect(() => {
    if (backendConfig) {
      initializeNetwork();
    }
  }, [backendConfig]);

  const loadBackendConfig = async () => {
    try {
      const config = await getBackendNetworkConfig();
      setBackendConfig(config);
      console.log('Backend config loaded:', config);
    } catch (e) {
      console.error('Failed to load backend config', e);
      try {
        const network = await JarvisService.GetNetworkByID(DEFAULT_CHAIN_ID);
        if (network) {
          setCurrentNetwork(network);
          console.log('Initialized with fallback network:', network.name);
        }
      } catch (fallbackError) {
        console.error('Failed to initialize fallback network', fallbackError);
      }
    }
  };

  const initializeNetwork = async () => {
    try {
      const chainId = backendConfig?.network.chainId || DEFAULT_CHAIN_ID;
      const network = await JarvisService.GetNetworkByID(chainId);
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
    if (address) {
      loadNativeBalance(network.id);
      loadKawaiBalance(network.id);
      loadGasEstimate(network.id);
      loadCurrentBlock(network.id);
      loadPrices(network.id);
    }
  }, [address, message]);

  useEffect(() => {
    if (address && currentNetwork) {
      setBalancesLoading(true);
      Promise.all([
        loadBalance(),
        loadHistory(),
        loadNativeBalance(currentNetwork.id),
        loadKawaiBalance(currentNetwork.id),
        loadGasEstimate(currentNetwork.id),
        loadCurrentBlock(currentNetwork.id),
        loadPrices(currentNetwork.id),
      ]).finally(() => setBalancesLoading(false));
    }
  }, [address, isWalletLoaded, currentNetwork]);

  const loadPrices = async (networkId: number) => {
    try {
      const nPrice = await JarvisService.GetTokenPrice("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", networkId);
      setNativePrice(nPrice);

      const kawaiAddress = getKawaiTokenAddress(networkId);
      if (!kawaiAddress) {
        console.warn('KAWAI address not available, skipping price load');
        setKawaiPrice(0);
        return;
      }
      const kPrice = await JarvisService.GetTokenPrice(kawaiAddress, networkId);
      setKawaiPrice(kPrice);
    } catch (e) {
      console.error('Failed to load prices', e);
    }
  };

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
      const kawaiAddress = getKawaiTokenAddress(networkId);
      if (!kawaiAddress) {
        console.warn('KAWAI address not available, skipping balance load');
        setKawaiBalance('0.00');
        return;
      }
      const result = await JarvisService.GetTokenBalance(kawaiAddress, address, networkId);
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

      hide();
      const syncHide = message.loading("Syncing balance (waiting for block confirmation)...", 0);

      let synced = false;
      let attempts = 0;
      const maxAttempts = 15;

      while (!synced && attempts < maxAttempts) {
        await new Promise(resolve => setTimeout(resolve, 2000));
        attempts++;

        const userAddress = await WalletService.GetCurrentAddress();
        const syncResult = await DepositSyncService.SyncDeposit({
          txHash: txHash,
          userAddress: userAddress
        });

        if (syncResult?.success) {
          syncHide();
          message.success(`Balance synced! New balance: ${(parseFloat(syncResult.newBalance || '0') / 1_000_000).toFixed(2)} ${currentNetwork?.stablecoinSymbol || 'USDT'}`);
          synced = true;
          break;
        }
      }

      if (!synced) {
        syncHide();
        message.warning(`Deposit transaction sent, but sync timed out. Please click "Refresh" in a few moments.`);
      }

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
        const rawAmount = Math.floor(amount * 1_000_000).toString();
        tx = await DeAIService.TransferUSDT(to, rawAmount);
      } else if (assetType === 'kawai') {
        const amountStr = amount.toFixed(18);
        const [intPart, decPart = '0'] = amountStr.split('.');
        const paddedDec = decPart.padEnd(18, '0').substring(0, 18);
        const rawAmount = (BigInt(intPart) * BigInt(10 ** 18) + BigInt(paddedDec)).toString();
        const kawaiAddress = currentNetwork ? getKawaiTokenAddress(currentNetwork.id) : (backendConfig?.contracts.kawai || '');
        if (!kawaiAddress) {
          message.error('KAWAI contract not available for this network');
          return;
        }
        tx = await DeAIService.TransferToken(kawaiAddress, to, rawAmount);
      } else if (customTokenAddress) {
        const networkId = currentNetwork?.id || DEFAULT_CHAIN_ID;
        const tokenInfo = await JarvisService.GetTokenInfo(customTokenAddress, networkId);
        const decimals = tokenInfo?.decimals || 18;

        console.log(`Sending custom token: ${customTokenAddress}, Decimals: ${decimals}`);

        const amountStr = amount.toFixed(decimals);
        const [intPart, decPart = '0'] = amountStr.split('.');
        const paddedDec = decPart.padEnd(decimals, '0').substring(0, decimals);
        const rawAmount = (BigInt(intPart) * (BigInt(10) ** BigInt(decimals)) + BigInt(paddedDec)).toString();

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
      console.error('Transfer error:', e);

      let errorMessage = 'Transfer Failed';

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
          errorMessage = `Transfer Failed: ${e.message}`;
        }
      } else {
        errorMessage = `Transfer Failed: ${e.toString()}`;
      }

      message.error(errorMessage, 5);
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
          nativePrice={nativePrice}
          kawaiPrice={kawaiPrice}
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
        return <OTCContent styles={styles} theme={theme} currentNetwork={currentNetwork} />;
      case 'rewards':
        return <RewardsContent styles={styles} theme={theme} currentNetwork={currentNetwork} transactions={transactions} setModalType={setModalType} />;
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
          currentNetwork={currentNetwork}
          onNetworkChange={handleNetworkChange}
          backendConfig={backendConfig}
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
        <SmartDepositForm onDeposit={handleDeposit} loading={loading} currentNetwork={currentNetwork} />
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
              Your Wallet Address ({currentNetwork?.name || 'Monad Testnet'})
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

DesktopWalletLayout.displayName = 'DesktopWalletLayout';

export default DesktopWalletLayout;
