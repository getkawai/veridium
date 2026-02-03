/**
 * Network Configuration
 *
 * This config automatically fetches the active network environment from the backend.
 * The Single Source of Truth is the backend config (internal/services/config_service.go),
 * which is generated from the server-side .env file.
 *
 * NO HARDCODED ADDRESSES OR CONSTANTS ALLOWED HERE.
 */

import { ConfigService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { BackendConfig } from '@@/github.com/kawai-network/veridium/internal/services/models';

/**
 * Default Chain ID (Monad Testnet)
 * Used as fallback when backend config fails to load
 */
export const DEFAULT_CHAIN_ID = 10143;

/**
 * Network Environment Types
 */
export type NetworkEnv = 'testnet' | 'mainnet';


/**
 * Get full network configuration from backend.
 * This object contains Chain ID, Names, and all contract addresses.
 */
export async function getBackendNetworkConfig(): Promise<BackendConfig> {
  try {
    return await ConfigService.GetConfig();
  } catch (e) {
    console.error('CRITICAL: Failed to load backend network config:', e);
    throw new Error("Failed to load network configuration. Please check backend connection.");
  }
}

/**
 * Get formatted token list for UI usage (e.g. Add Token Modal).
 * Fetched dynamically from backend config to ensure addresses are correct for current env.
 */
export async function getTokenListFromBackend(): Promise<Array<{
  address: string;
  name: string;
  symbol: string;
  decimals: number;
}>> {
  try {
    const config = await getBackendNetworkConfig();
    return [
      {
        address: config.contracts.usdt,
        name: getStablecoinDisplayName(config),
        symbol: getStablecoinSymbol(config),
        decimals: 6,
      },
      {
        address: config.contracts.kawai,
        name: 'Kawai AI Token',
        symbol: 'KAWAI',
        decimals: 18,
      },
    ];
  } catch (e) {
    console.error('Failed to load token list from backend:', e);
    return [];
  }
}

/**
 * Check if current environment is testnet
 */
export function isTestnet(config: BackendConfig): boolean {
  return config.network.isTestnet;
}


/**
 * Get stablecoin symbol based on network environment
 * Returns "MockUSDT" for testnet, "USDC" for mainnet
 */
export function getStablecoinSymbol(config: BackendConfig): string {
  return config.network.isTestnet ? 'MockUSDT' : 'USDC';
}

/**
 * Get stablecoin display name based on network environment
 * Returns full name for UI display
 */
export function getStablecoinDisplayName(config: BackendConfig): string {
  return config.network.isTestnet
    ? 'Mock Tether USD (Testnet)'
    : 'USD Coin';
}

