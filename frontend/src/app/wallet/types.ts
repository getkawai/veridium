import type { NetworkInfo, ClaimableRewardsResponse, ClaimableReward } from '@@/github.com/kawai-network/veridium/internal/services/models';
import type { WalletTransaction } from '@@/github.com/getkawai/database/db/models';
import type { UserBalanceInfo } from '@@/github.com/kawai-network/x/jarvis/models';

// Export NetworkInfo from the generated models to ensure consistency with Go backend
export type { ClaimableRewardsResponse, ClaimableReward, NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';

export type MenuKey = 'home' | 'otc' | 'rewards' | 'settings';

export interface ContentStyles {
  container: string;
  content: string;
  balanceCard: string;
  actionButton: string;
  actionCircle: string;
  eyeButton: string;
  statValue: string;
  tokenRow: string;
  placeholderCard: string;
}

export interface HomeContentProps {
  address: string;
  onChainBalance: string;
  trackedBalance: UserBalanceInfo | null;
  nativeBalance: string;
  kawaiBalance: string;
  nativePrice: number;
  kawaiPrice: number;
  balanceVisible: boolean;
  setBalanceVisible: (visible: boolean) => void;
  setModalType: (type: 'send' | 'receive' | 'swap' | 'deposit' | 'addAccount' | 'createWallet' | 'importWallet' | 'addToken' | null) => void;
  transactions: WalletTransaction[];
  styles: ContentStyles;
  theme: any;
  currentNetwork: NetworkInfo | null;
  gasEstimate: { maxGasPriceGwei: number; maxTipGwei: number } | null;
  currentBlock: number;
  balancesLoading: boolean;
}

export interface OTCContentProps {
  styles: ContentStyles;
  theme: any;
  currentNetwork: NetworkInfo | null;
}

export interface RewardsContentProps {
  styles: ContentStyles;
  theme: any;
  currentNetwork: NetworkInfo | null;
  transactions: WalletTransaction[];
  setModalType?: (type: 'send' | 'receive' | 'swap' | 'deposit' | 'addAccount' | 'createWallet' | 'importWallet' | 'addToken' | null) => void;
}

export interface SettingsContentProps {
  address: string;
  styles: ContentStyles;
  theme: any;
  currentNetwork: NetworkInfo | null;
}
