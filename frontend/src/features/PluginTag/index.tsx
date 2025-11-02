'use client';

import { Dropdown, Icon, type MenuProps, Tag } from '@lobehub/ui';
import isEqual from 'fast-deep-equal';
import { LucideToyBrick } from 'lucide-react';
import { memo } from 'react';
import { Center } from 'react-layout-kit';

import Avatar from '@/components/Plugins/PluginAvatar';
// import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
// import { pluginHelpers, useToolStore } from '@/store/tool';
// import { toolSelectors } from '@/store/tool/selectors';

// Dummy implementations for development - memoized
const mockServerConfig = {
  featureFlags: { showDalle: true }
};

const useServerConfigStore = (selector?: any) => {
  if (selector && typeof selector === 'function') {
    return selector(mockServerConfig);
  }
  return mockServerConfig;
};

const featureFlagsSelectors = (s: any) => s.featureFlags || { showDalle: true };

const mockToolStore = {
  metaList: [
    {
      identifier: 'plugin-1',
      meta: {
        title: 'Sample Plugin 1',
        avatar: '🔌'
      }
    },
    {
      identifier: 'plugin-2',
      meta: {
        title: 'Sample Plugin 2',
        avatar: '⚡'
      }
    }
  ],
  getMetaById: (id: string) => {
    const plugin = mockToolStore.metaList.find(p => p.identifier === id);
    return plugin?.meta || { title: id, avatar: '🔌' };
  }
};

const useToolStore = (selector?: any, comparator?: any) => {
  if (selector) {
    return selector(mockToolStore);
  }
  return mockToolStore;
};

const toolSelectors = {
  metaList: (showDalle?: boolean) => (state: any) => state.metaList,
  getMetaById: (id: string) => (state: any) => state.getMetaById(id)
};

const pluginHelpers = {
  getPluginAvatar: (meta: any) => meta?.avatar || '🔌',
  getPluginTitle: (meta: any) => meta?.title || 'Unknown Plugin'
};

import PluginStatus from './PluginStatus';

export interface PluginTagProps {
  plugins: string[];
}

const PluginTag = memo<PluginTagProps>(({ plugins }) => {
  const { showDalle } = useServerConfigStore(featureFlagsSelectors);
  const list = useToolStore(toolSelectors.metaList(showDalle), isEqual);

  const displayPlugin = useToolStore(toolSelectors.getMetaById(plugins[0]), isEqual);

  if (plugins.length === 0) return null;

  const items: MenuProps['items'] = plugins.map((id) => {
    const item = list.find((i) => i.identifier === id);

    const isDeprecated = !item;
    const avatar = isDeprecated ? '♻️' : pluginHelpers.getPluginAvatar(item.meta || item);

    return {
      icon: (
        <Center style={{ minWidth: 24 }}>
          <Avatar avatar={avatar} size={24} />
        </Center>
      ),
      key: id,
      label: (
        <PluginStatus
          deprecated={isDeprecated}
          id={id}
          title={pluginHelpers.getPluginTitle(item?.meta || item)}
        />
      ),
    };
  });

  const count = plugins.length;

  return (
    <Dropdown menu={{ items }}>
      <div>
        <Tag>
          {<Icon icon={LucideToyBrick} />}
          {pluginHelpers.getPluginTitle(displayPlugin) || plugins[0]}
          {count > 1 && <div>({plugins.length - 1}+)</div>}
        </Tag>
      </div>
    </Dropdown>
  );
});

export default PluginTag;
