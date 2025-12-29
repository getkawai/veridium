import {
  NetworkEthereum,
  NetworkBinanceSmartChain,
  NetworkPolygon,
  NetworkPolygonZkevm,
  NetworkAvalanche,
  NetworkFantom,
  NetworkArbitrumOne,
  NetworkOptimism,
  NetworkBase,
  NetworkScroll,
  NetworkLinea,
  NetworkMonad,
} from '@web3icons/react';

interface NetworkIconProps {
  name: string;
  size?: number;
  variant?: 'branded' | 'mono';
  color?: string;
}

const NETWORK_COMPONENT_MAP: Record<string, React.ElementType> = {
  'ethereum': NetworkEthereum,
  'binance-smart-chain': NetworkBinanceSmartChain,
  'polygon': NetworkPolygon,
  'polygon-zkevm': NetworkPolygonZkevm,
  'avalanche': NetworkAvalanche,
  'fantom': NetworkFantom,
  'arbitrum': NetworkArbitrumOne,
  'optimism': NetworkOptimism,
  'base': NetworkBase,
  'scroll': NetworkScroll,
  'linea': NetworkLinea,
  'monad': NetworkMonad,
};

export const NetworkIcon = ({ name, size = 24, variant = 'branded' }: NetworkIconProps) => {
  const IconComponent = NETWORK_COMPONENT_MAP[name];

  // Fallback to Ethereum icon for unknown networks
  if (!IconComponent) {
    return <NetworkEthereum size={size} variant={variant} />;
  }

  return <IconComponent size={size} variant={variant} />;
};
