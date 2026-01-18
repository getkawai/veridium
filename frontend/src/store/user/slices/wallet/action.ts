import { StateCreator } from "zustand/vanilla";
import * as WalletService from "@@/github.com/kawai-network/veridium/internal/services/walletservice";
import type { WalletInfo } from "@@/github.com/kawai-network/veridium/internal/services/models";
import type { UserStore } from "../../store";
import { message } from "antd";

export interface WalletAction {
  refreshWalletStatus: () => Promise<void>;
  unlockWallet: (password: string) => Promise<boolean>;
  setupWallet: (
    password: string,
    mnemonic: string,
    name?: string,
  ) => Promise<string>;
  generateMnemonic: () => Promise<string>;
  lockWallet: () => Promise<void>;
  // Multi-wallet actions
  createWallet: (
    password: string,
    mnemonic: string,
    description?: string,
  ) => Promise<string>;
  switchWallet: (address: string, password: string) => Promise<boolean>;
  deleteWallet: (address: string) => Promise<boolean>;
  exportKeystore: (address: string) => Promise<string>;
  importKeystore: (
    keystoreJSON: string,
    password: string,
    description?: string,
  ) => Promise<string>;
  importPrivateKey: (
    privateKey: string,
    password: string,
    description?: string,
  ) => Promise<string>;
  updateWalletDescription: (
    address: string,
    description: string,
  ) => Promise<boolean>;
  getWallets: () => Promise<WalletInfo[]>;
  getAPIKey: () => Promise<string>;
}

export const createWalletSlice: StateCreator<
  UserStore,
  [["zustand/devtools", never]],
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
        wallets: status.wallets || [],
      });
    } catch (error) {
      const errorMessage =
        error instanceof Error
          ? error.message
          : "Failed to refresh wallet status";
      console.error("Failed to refresh wallet status:", error);
      message.error(errorMessage);
      set({ isWalletLoaded: true });
    }
  },

  unlockWallet: async (password: string) => {
    try {
      await WalletService.UnlockWallet(password);
      await get().refreshWalletStatus();

      // Auto-claim trial if needed (with referral code if available)
      try {
        const pendingReferralCode =
          localStorage.getItem("pendingReferralCode") || "";
        const [claimed, usdtAmount, kawaiAmount] =
          await WalletService.AutoClaimTrialIfNeeded(pendingReferralCode);

        if (claimed) {
          // Clear pending referral code after successful claim
          localStorage.removeItem("pendingReferralCode");

          // Show success message
          const usdtFormatted = usdtAmount.toFixed(2);
          const kawaiFormatted = (parseFloat(kawaiAmount) / 1e18).toFixed(0);
          message.success(
            `🎉 Free trial claimed: ${usdtFormatted} USDT + ${kawaiFormatted} KAWAI`,
          );
        }
      } catch (claimError) {
        // Log but show warning to user without failing unlock
        const warningMessage =
          claimError instanceof Error
            ? claimError.message
            : "Failed to claim trial bonus";
        console.warn("Failed to auto-claim trial:", claimError);
        message.warning(warningMessage);
      }

      return true;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to unlock wallet";
      console.error("Failed to unlock wallet:", error);
      message.error(errorMessage);
      return false;
    }
  },

  setupWallet: async (password: string, mnemonic: string, name?: string) => {
    try {
      const address = await WalletService.SetupWallet(
        password,
        mnemonic,
        name || "My Wallet",
      );
      await get().refreshWalletStatus();
      message.success("Wallet created successfully");
      return address;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to setup wallet";
      console.error("Failed to setup wallet:", error);
      message.error(errorMessage);
      throw error;
    }
  },

  generateMnemonic: async () => {
    try {
      const mnemonic = await WalletService.GenerateMnemonic();
      return mnemonic;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to generate mnemonic";
      console.error("Failed to generate mnemonic:", error);
      message.error(errorMessage);
      throw error;
    }
  },

  lockWallet: async () => {
    try {
      await WalletService.LockWallet();
      set({ isWalletLocked: true, walletAddress: "" });
      message.success("Wallet locked successfully");
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to lock wallet";
      console.error("Failed to lock wallet:", error);
      message.error(errorMessage);
    }
  },

  // Multi-wallet actions
  createWallet: async (
    password: string,
    mnemonic: string,
    description?: string,
  ) => {
    try {
      const address = await WalletService.CreateWallet(
        password,
        mnemonic,
        description || "",
      );
      await get().refreshWalletStatus();
      message.success("Wallet created successfully");
      return address;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to create wallet";
      console.error("Failed to create wallet:", error);
      message.error(errorMessage);
      throw error;
    }
  },

  switchWallet: async (address: string, password: string) => {
    try {
      await WalletService.SwitchWallet(address, password);
      await get().refreshWalletStatus();
      message.success("Wallet switched successfully");
      return true;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to switch wallet";
      console.error("Failed to switch wallet:", error);
      message.error(errorMessage);
      return false;
    }
  },

  deleteWallet: async (address: string) => {
    try {
      await WalletService.DeleteWallet(address);
      await get().refreshWalletStatus();
      message.success("Wallet deleted successfully");
      return true;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to delete wallet";
      console.error("Failed to delete wallet:", error);
      message.error(errorMessage);
      return false;
    }
  },

  exportKeystore: async (address: string) => {
    try {
      const keystore = await WalletService.ExportKeystore(address);
      message.success("Keystore exported successfully");
      return keystore;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to export keystore";
      console.error("Failed to export keystore:", error);
      message.error(errorMessage);
      throw error;
    }
  },

  importKeystore: async (
    keystoreJSON: string,
    password: string,
    description?: string,
  ) => {
    try {
      const address = await WalletService.ImportKeystore(
        keystoreJSON,
        password,
        description || "",
      );
      await get().refreshWalletStatus();
      message.success("Keystore imported successfully");
      return address;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to import keystore";
      console.error("Failed to import keystore:", error);
      message.error(errorMessage);
      throw error;
    }
  },

  importPrivateKey: async (
    privateKey: string,
    password: string,
    description?: string,
  ) => {
    try {
      const address = await WalletService.ImportPrivateKey(
        privateKey,
        password,
        description || "",
      );
      await get().refreshWalletStatus();
      message.success("Private key imported successfully");
      return address;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to import private key";
      console.error("Failed to import private key:", error);
      message.error(errorMessage);
      throw error;
    }
  },

  updateWalletDescription: async (address: string, description: string) => {
    try {
      await WalletService.UpdateWalletDescription(address, description);
      await get().refreshWalletStatus();
      message.success("Wallet description updated successfully");
      return true;
    } catch (error) {
      const errorMessage =
        error instanceof Error
          ? error.message
          : "Failed to update wallet description";
      console.error("Failed to update wallet description:", error);
      message.error(errorMessage);
      return false;
    }
  },

  getWallets: async () => {
    try {
      const wallets = await WalletService.GetWallets();
      set({ wallets });
      return wallets;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to get wallets";
      console.error("Failed to get wallets:", error);
      message.error(errorMessage);
      return [];
    }
  },

  getAPIKey: async () => {
    try {
      return await WalletService.GetAPIKey();
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to get API key";
      console.error("Failed to get API key:", error);
      message.error(errorMessage);
      return "";
    }
  },
});
