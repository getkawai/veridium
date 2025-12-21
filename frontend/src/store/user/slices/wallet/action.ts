import { StateCreator } from 'zustand/vanilla';
import * as WalletService from '@@/github.com/kawai-network/veridium/internal/services/walletservice';
import type { UserStore } from '../../store';

export interface WalletAction {
  refreshWalletStatus: () => Promise<void>;
  unlockWallet: (password: string) => Promise<boolean>;
  setupWallet: (password: string, mnemonic: string) => Promise<string>;
  generateMnemonic: () => Promise<string>;
  lockWallet: () => Promise<void>;
}

export const createWalletSlice: StateCreator<
  UserStore,
  [['zustand/devtools', never]],
  [],
  WalletAction
> = (set, get) => ({
  refreshWalletStatus: async () => {
    try {
      const status = await WalletService.GetStatus();
      set({
        hasWallet: status.hasWallet,
        isWalletLocked: status.isLocked,
        walletAddress: status.address,
        isWalletLoaded: true
      });
    } catch (error) {
      console.error('Failed to refresh wallet status:', error);
    }
  },

  unlockWallet: async (password: string) => {
    try {
      const address = await WalletService.UnlockWallet(password);
      set({
        walletAddress: address,
        isWalletLocked: false
      });
      return true;
    } catch (error) {
      console.error('Failed to unlock wallet:', error);
      return false;
    }
  },

  setupWallet: async (password: string, mnemonic: string) => {
    try {
      const address = await WalletService.SetupWallet(password, mnemonic);
      set({
        walletAddress: address,
        isWalletLocked: false,
        hasWallet: true
      });
      return address;
    } catch (error) {
      console.error('Failed to setup wallet:', error);
      throw error;
    }
  },

  generateMnemonic: async () => {
    return await WalletService.GenerateMnemonic();
  },

  lockWallet: async () => {
    // WalletService.LockWallet() is not implemented yet in backend but follows the plan
    // We'll just reset state for now
    set({ isWalletLocked: true, walletAddress: '' });
  },
});
