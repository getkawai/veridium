export interface WalletState {
  walletAddress: string;
  isWalletLocked: boolean;
  hasWallet: boolean;
  isWalletLoaded: boolean;
}

export const initialWalletState: WalletState = {
  walletAddress: '',
  isWalletLocked: true,
  hasWallet: false,
  isWalletLoaded: false,
};
