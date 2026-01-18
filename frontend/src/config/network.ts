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
 * Network Environment Types
 */
export type NetworkEnv = 'testnet' | 'mainnet';

/**
 * Get current network environment from backend
 */
export async function detectNetworkEnv(): Promise<NetworkEnv> {
  try {
    const config = await ConfigService.GetConfig();
    const env = config.environment;
    
    // Validate environment value
    if (env !== 'testnet' && env !== 'mainnet') {
      console.error(`Invalid network environment from backend: ${env}`);
      throw new Error(`Unsupported network environment: ${env}`);
    }
    
    return env as NetworkEnv;
  } catch (e) {
    console.error('CRITICAL: Failed to detect network environment from backend:', e);
    // Re-throw to allow UI to handle critical configuration failure
    throw e;
  }
}

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
        name: 'Tether USD',
        symbol: 'USDT',
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
 * Check if current environment is mainnet
 */
export function isMainnet(config: BackendConfig): boolean {
  return !config.network.isTestnet;
}
