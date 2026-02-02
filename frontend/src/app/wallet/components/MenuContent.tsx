import { memo } from 'react';
import { Home, ShoppingCart, Gift, Settings } from 'lucide-react';
import { Icon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import type { NetworkInfo, BackendConfig } from '@@/github.com/kawai-network/veridium/internal/services/models';
import Menu from '@/components/Menu';
import PanelTitle from '@/components/PanelTitle';
import { NetworkSwitcher } from './NetworkSwitcher';
import type { MenuKey } from '../types';

interface MenuContentProps {
  activeMenu: MenuKey;
  setActiveMenu: (key: MenuKey) => void;
  currentNetwork: NetworkInfo | null;
  onNetworkChange: (network: NetworkInfo) => void;
  backendConfig: BackendConfig | null;
}

export const MenuContent = memo<MenuContentProps>(({
  activeMenu,
  setActiveMenu,
  currentNetwork,
  onNetworkChange,
  backendConfig,
}) => {
  const menuItems = [
    { key: 'home', icon: <Icon icon={Home} />, label: 'Home' },
    { key: 'otc', icon: <Icon icon={ShoppingCart} />, label: 'OTC Market' },
    { key: 'rewards', icon: <Icon icon={Gift} />, label: 'Rewards' },
    { key: 'settings', icon: <Icon icon={Settings} />, label: 'Settings' },
  ];

  return (
    <Flexbox gap={16} height={'100%'}>
      <Flexbox paddingInline={8}>
        <PanelTitle desc="Manage your digital assets and transactions" title="Wallet" />
        <Menu
          compact
          selectable
          items={menuItems}
          selectedKeys={[activeMenu]}
          onClick={({ key }) => setActiveMenu(key as MenuKey)}
        />
      </Flexbox>
      <div style={{ flex: 1 }} />
      <Flexbox padding={12}>
        <NetworkSwitcher
          currentNetwork={currentNetwork}
          onNetworkChange={onNetworkChange}
          backendConfig={backendConfig}
        />
      </Flexbox>
    </Flexbox>
  );
});

MenuContent.displayName = 'MenuContent';
