import { StateCreator } from 'zustand/vanilla';
import * as WalletService from '@@/github.com/kawai-network/veridium/internal/services/walletservice';
import type { WalletInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';
import type { UserStore } from '../../store';

export interface WalletAction {
  refreshWalletStatus: () => Promise<void>;
  unlockWallet: (password: string) => Promise<boolean>;
  setupWallet: (password: string, mnemonic: string) => Promise<string>;
  generateMnemonic: () => Promise<string>;
  lockWallet: () => Promise<void>;
  // Multi-wallet actions
  createWallet: (password: string, mnemonic: string, description?: string) => Promise<string>;
  switchWallet: (address: string, password: string) => Promise<boolean>;
  deleteWallet: (address: string) => Promise<boolean>;
  exportKeystore: (address: string) => Promise<string>;
  importKeystore: (keystoreJSON: string, password: string, description?: string) => Promise<string>;
  updateWalletDescription: (address: string, description: string) => Promise<boolean>;
  getWallets: () => Promise<WalletInfo[]>;
  getAPIKey: () => Promise<string>;
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
        isWalletLoaded: true,
        wallets: status.wallets || []
      });
    } catch (error) {
      console.error('Failed to refresh wallet status:', error);
    }
  },

  unlockWallet: async (password: string) => {
    try {
      await WalletService.UnlockWallet(password);
      await get().refreshWalletStatus();
      return true;
    } catch (error) {
      console.error('Failed to unlock wallet:', error);
      return false;
    }
  },

  setupWallet: async (password: string, mnemonic: string) => {
    try {
      const address = await WalletService.SetupWallet(password, mnemonic);
      await get().refreshWalletStatus();
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
    try {
      await WalletService.LockWallet();
      set({ isWalletLocked: true, walletAddress: '' });
    } catch (error) {
      console.error('Failed to lock wallet:', error);
    }
  },

  // Multi-wallet actions
  createWallet: async (password: string, mnemonic: string, description?: string) => {
    try {
      const address = await WalletService.CreateWallet(password, mnemonic, description || '');
      await get().refreshWalletStatus();
      return address;
    } catch (error) {
      console.error('Failed to create wallet:', error);
      throw error;
    }
  },

  switchWallet: async (address: string, password: string) => {
    try {
      await WalletService.SwitchWallet(address, password);
      await get().refreshWalletStatus();
      return true;
    } catch (error) {
      console.error('Failed to switch wallet:', error);
      return false;
    }
  },

  deleteWallet: async (address: string) => {
    try {
      await WalletService.DeleteWallet(address);
      await get().refreshWalletStatus();
      return true;
    } catch (error) {
      console.error('Failed to delete wallet:', error);
      return false;
    }
  },

  exportKeystore: async (address: string) => {
    try {
      return await WalletService.ExportKeystore(address);
    } catch (error) {
      console.error('Failed to export keystore:', error);
      throw error;
    }
  },

  importKeystore: async (keystoreJSON: string, password: string, description?: string) => {
    try {
      const address = await WalletService.ImportKeystore(keystoreJSON, password, description || '');
      await get().refreshWalletStatus();
      return address;
    } catch (error) {
      console.error('Failed to import keystore:', error);
      throw error;
    }
  },

  updateWalletDescription: async (address: string, description: string) => {
    try {
      await WalletService.UpdateWalletDescription(address, description);
      await get().refreshWalletStatus();
      return true;
    } catch (error) {
      console.error('Failed to update wallet description:', error);
      return false;
    }
  },

  getWallets: async () => {
    try {
      const wallets = await WalletService.GetWallets();
      set({ wallets });
      return wallets;
    } catch (error) {
      console.error('Failed to get wallets:', error);
      return [];
    }
  },

  getAPIKey: async () => {
    try {
      return await WalletService.GetAPIKey();
    } catch (error) {
      console.error('Failed to get API key:', error);
      return '';
    }
  },
});
