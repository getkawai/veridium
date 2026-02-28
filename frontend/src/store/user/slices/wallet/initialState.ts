import type { WalletInfo } from '@@/github.com/kawai-network/x/jarvis/models';

export interface WalletState {
  walletAddress: string;
  isWalletLocked: boolean;
  hasWallet: boolean;
  isWalletLoaded: boolean;
  wallets: WalletInfo[];
}

export const initialWalletState: WalletState = {
  walletAddress: '',
  isWalletLocked: true,
  hasWallet: false,
  isWalletLoaded: false,
  wallets: [],
};
