import { TokenUSDT, TokenUSDC } from '@web3icons/react';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';

interface StablecoinIconProps {
  currentNetwork: NetworkInfo | null;
  size?: number;
  variant?: 'branded' | 'mono';
}

/**
 * Dynamic stablecoin icon that shows:
 * - TokenUSDT on testnet (MockUSDT)
 * - TokenUSDC on mainnet (USDC)
 */
export const StablecoinIcon = ({ currentNetwork, size = 24, variant = 'branded' }: StablecoinIconProps) => {
  // Show USDC icon on mainnet, USDT icon on testnet (or fallback)
  const isMainnet = currentNetwork && !currentNetwork.isTestnet;
  
  if (isMainnet) {
    return <TokenUSDC size={size} variant={variant} />;
  }
  
  return <TokenUSDT size={size} variant={variant} />;
};
