import { Avatar, ItemType } from '@lobehub/ui';
import isEqual from 'fast-deep-equal';
import { useTranslation } from 'react-i18next';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/selectors';
import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
import { useToolStore } from '@/store/tool';
import { builtinToolSelectors } from '@/store/tool/selectors';

import ToolItem from './ToolItem';

export const useControls = ({
  setUpdating,
}: {
  setUpdating: (updating: boolean) => void;
}) => {
  const { t } = useTranslation('setting');
  const { showDalle } = useServerConfigStore(featureFlagsSelectors);
  const [checked, togglePlugin] = useAgentStore((s) => [
    agentSelectors.currentAgentPlugins(s),
    s.togglePlugin,
  ]);
  const builtinList = useToolStore(builtinToolSelectors.metaList(showDalle), isEqual);

  const items: ItemType[] = [
    {
      children: builtinList.map((item) => ({
        icon: <Avatar avatar={item.meta.avatar} size={20} style={{ flex: 'none' }} />,
        key: item.identifier,
        label: (
          <ToolItem
            checked={checked.includes(item.identifier)}
            id={item.identifier}
            label={item.meta?.title}
            onUpdate={async () => {
              setUpdating(true);
              await togglePlugin(item.identifier);
              setUpdating(false);
            }}
          />
        ),
      })),

      key: 'builtins',
      label: t('tools.builtins.groupName'),
      type: 'group',
    },
  ];

  return items;
};
