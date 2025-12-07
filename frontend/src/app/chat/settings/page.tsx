'use client';

import { Tabs } from '@lobehub/ui';
import { useTheme } from 'antd-style';
import isEqual from 'fast-deep-equal';
import { memo, useState } from 'react';

import { useCategory } from '@/features/AgentSetting/AgentCategory/useCategory';
import AgentSettings from '@/features/AgentSetting/AgentSettings';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/selectors';
import { ChatSettingsTabs } from '@/store/global/initialState';
import { useSessionStore } from '@/store/session';
import { sessionMetaSelectors } from '@/store/session/selectors';

const EditPage = memo(() => {
  const [tab, setTab] = useState(ChatSettingsTabs.Prompt);
  const theme = useTheme();
  const cateItems = useCategory();
  const [id, updateAgentMeta] = useSessionStore((s) => [
    s.activeId,
    s.updateSessionMeta,
    sessionMetaSelectors.currentAgentTitle(s),
  ]);

  const [updateAgentConfig] = useAgentStore((s) => [s.updateAgentConfig]);

  const config = useAgentStore(agentSelectors.currentAgentConfig, isEqual);
  const meta = useSessionStore(sessionMetaSelectors.currentAgentMeta, isEqual);

  // Check if config is loaded from agentConfigInitMap
  const isConfigLoaded = useAgentStore((s) => !!s.agentConfigInitMap[id]);
  const isLoading = !isConfigLoaded;

  return (
    <>
      <Tabs
        activeKey={tab}
        compact
        items={cateItems as any}
        onChange={(value) => setTab(value as ChatSettingsTabs)}
        style={{
          borderBottom: `1px solid ${theme.colorBorderSecondary}`,
        }}
      />
      <AgentSettings
        config={config}
        id={id}
        loading={isLoading}
        meta={meta}
        onConfigChange={updateAgentConfig}
        onMetaChange={updateAgentMeta}
        tab={tab}
      />
    </>
  );
});

export default EditPage;
