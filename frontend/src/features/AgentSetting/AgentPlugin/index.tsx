'use client';

import { Form, type FormGroupItemType } from '@lobehub/ui';
import { Empty } from 'antd';
import isEqual from 'fast-deep-equal';

import { memo, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { Center, Flexbox } from 'react-layout-kit';

import PluginAvatar from '@/components/Plugins/PluginAvatar';
import PluginTag from '@/components/Plugins/PluginTag';
import { FORM_STYLE } from '@/const/layoutTokens';
import { useFetchInstalledPlugins } from '@/hooks/useFetchInstalledPlugins';
import { featureFlagsSelectors, useServerConfigStore } from '@/store/serverConfig';
import { pluginHelpers, useToolStore } from '@/store/tool';
import { toolSelectors } from '@/store/tool/selectors';

import { useStore } from '../store';
import LoadingList from './LoadingList';
import LocalPluginItem from './LocalPluginItem';
import PluginAction from './PluginAction';

const AgentPlugin = memo(() => {
  const { t } = useTranslation('setting');

  const [showStore, setShowStore] = useState(false);

  const [userEnabledPlugins, toggleAgentPlugin] = useStore((s) => [
    s.config.plugins || [],
    s.toggleAgentPlugin,
  ]);

  const { showDalle } = useServerConfigStore(featureFlagsSelectors);
  const installedPlugins = useToolStore(toolSelectors.metaList(showDalle), isEqual);

  const { isLoading } = useFetchInstalledPlugins();

  const isEmpty = installedPlugins.length === 0 && userEnabledPlugins.length === 0;

  //  =========== Plugin List =========== //

  const list = installedPlugins.map(({ identifier, type, meta, author }) => {
    const isCustomPlugin = type === 'customPlugin';

    return {
      avatar: <PluginAvatar avatar={pluginHelpers.getPluginAvatar(meta)} size={40} />,
      children: isCustomPlugin ? (
        <LocalPluginItem id={identifier} />
      ) : (
        <PluginAction identifier={identifier} />
      ),
      desc: pluginHelpers.getPluginDesc(meta),
      label: (
        <Flexbox align={'center'} gap={8} horizontal>
          {pluginHelpers.getPluginTitle(meta)}
          <PluginTag author={author} type={type} />
        </Flexbox>
      ),
      layout: 'horizontal',
      minWidth: undefined,
    };
  });

  const loadingSkeleton = LoadingList();

  const empty = (
    <Center padding={40}>
      <Empty
        description={
          <Trans i18nKey={'plugin.empty'} ns={'setting'}>
            暂无安装插件，
            <span
              onClick={() => setShowStore(true)}
              style={{ cursor: 'pointer', color: 'inherit' }}
            >
              前往插件市场
            </span>
            安装
          </Trans>
        }
        image={Empty.PRESENTED_IMAGE_SIMPLE}
      />
    </Center>
  );

  const plugin: FormGroupItemType = {
    children: isLoading ? loadingSkeleton : isEmpty ? empty : [...list],
    title: t('settingPlugin.title'),
  };

  return <Form items={[plugin]} itemsType={'group'} variant={'borderless'} {...FORM_STYLE} />;
});

export default AgentPlugin;
