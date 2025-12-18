'use client';

import { Drawer } from '@lobehub/ui';
import isEqual from 'fast-deep-equal';
import { memo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import BrandWatermark from '@/components/BrandWatermark';
import PanelTitle from '@/components/PanelTitle';
import { INBOX_SESSION_ID } from '@/const/session';
import { AgentCategory, AgentSettings as Settings, HeaderContent } from '@/features/AgentSetting';
import { AgentSettingsProvider } from '@/features/AgentSetting/AgentSettingsProvider';
import Footer from '@/features/Setting/Footer';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/slices/chat';
import { ChatSettingsTabs } from '@/store/global/initialState';
import { useSessionStore } from '@/store/session';
import { sessionMetaSelectors } from '@/store/session/selectors';

export interface AgentSettingsProps {
  agentId?: string;
  onClose?: () => void;
  open?: boolean;
}

const AgentSettings = memo<AgentSettingsProps>(({ agentId, onClose, open }) => {
  const { t } = useTranslation('setting');

  const activeId = useSessionStore((s) => s.activeId);
  const id = agentId || activeId;

  const config = useAgentStore((s) => {
    if (agentId) {
      return agentSelectors.getAgentConfigByAgentId(agentId)(s);
    } else if (id) {
      return agentSelectors.getAgentConfigById(id)(s);
    } else {
      return agentSelectors.currentAgentConfig(s);
    }
  }, isEqual);
  const meta = useSessionStore((s) => {
    if (agentId) {
      return sessionMetaSelectors.getAgentMetaByAgentId(agentId)(s);
    } else {
      return sessionMetaSelectors.currentAgentMeta(s);
    }
  }, isEqual);

  const isConfigLoaded = useAgentStore((s) =>
    agentId ? !!s.agentConfigInitMap[agentId] : !!s.agentConfigInitMap[id],
  );
  const isLoading = !isConfigLoaded;

  const [showAgentSetting, globalUpdateAgentConfig] = useAgentStore((s) => [
    s.showAgentSetting,
    s.updateAgentConfig,
  ]);
  const [globalUpdateAgentMeta] = useSessionStore((s) => [
    s.updateSessionMeta,
    sessionMetaSelectors.currentAgentTitle(s),
  ]);

  const updateAgentConfig = async (config: any) => {
    if (agentId) {
      // If agentId is provided, we assume it refers to the session ID for now
      // as that's how the store handles it in this context
      await useAgentStore.getState().internal_updateAgentConfig(agentId, config);
    } else {
      await globalUpdateAgentConfig(config);
    }
  };

  const updateAgentMeta = async (meta: any) => {
    if (agentId) {
      const currentActiveId = useSessionStore.getState().activeId;
      useSessionStore.getState().switchSession(agentId);
      await useSessionStore.getState().updateSessionMeta(meta);
      if (currentActiveId !== agentId) {
        useSessionStore.getState().switchSession(currentActiveId);
      }
    } else {
      await globalUpdateAgentMeta(meta);
    }
  };

  const isOpen = open !== undefined ? open : showAgentSetting;
  const handleClose = onClose || (() => useAgentStore.setState({ showAgentSetting: false }));

  const isInbox = id === INBOX_SESSION_ID;
  const [tab, setTab] = useState(isInbox ? ChatSettingsTabs.Prompt : ChatSettingsTabs.Meta);

  return (
    <AgentSettingsProvider
      config={config}
      id={id}
      loading={isLoading}
      meta={meta}
      onConfigChange={updateAgentConfig}
      onMetaChange={updateAgentMeta}
    >
      <Drawer
        containerMaxWidth={1280}
        height={'100vh'}
        noHeader
        onClose={handleClose}
        open={isOpen}
        placement={'bottom'}
        sidebar={
          <Flexbox
            gap={20}
            style={{
              height: 'calc(100vh - 28px)',
            }}
          >
            <PanelTitle desc={t('header.sessionDesc')} title={t('header.session')} />
            <Flexbox flex={1} width={'100%'}>
              <AgentCategory setTab={setTab} tab={tab} />
            </Flexbox>
            <Flexbox align={'center'} gap={8} paddingInline={8} width={'100%'}>
              <HeaderContent modal />
            </Flexbox>
            <BrandWatermark paddingInline={12} />
          </Flexbox>
        }
        sidebarWidth={280}
        styles={{
          sidebarContent: {
            gap: 48,
            justifyContent: 'space-between',
            minHeight: '100%',
            paddingBlock: 24,
            paddingInline: 48,
          },
        }}
      >
        <Settings
          config={config}
          id={id}
          loading={isLoading}
          meta={meta}
          onConfigChange={updateAgentConfig}
          onMetaChange={updateAgentMeta}
          tab={tab}
        />
        <Footer />
      </Drawer>
    </AgentSettingsProvider>
  );
});

export default AgentSettings;
