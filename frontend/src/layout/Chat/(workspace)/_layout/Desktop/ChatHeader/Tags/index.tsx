import { ModelTag } from '@lobehub/icons';
import { Skeleton } from 'antd';
import isEqual from 'fast-deep-equal';
import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

import ModelSwitchPanel from '@/features/ModelSwitchPanel';
import PluginTag from '@/features/PluginTag';
// import { useAgentEnableSearch } from '@/hooks/useAgentEnableSearch';
// import { useModelSupportToolUse } from '@/hooks/useModelSupportToolUse';
// import { useAgentStore } from '@/store/agent';
// import { agentChatConfigSelectors, agentSelectors } from '@/store/agent/selectors';
// import { useSessionStore } from '@/store/session';
// import { sessionSelectors } from '@/store/session/selectors';
// import { useUserStore } from '@/store/user';
// import { authSelectors } from '@/store/user/selectors';

// Dummy implementations for development - memoized
const useAgentEnableSearch = () => true;

const useModelSupportToolUse = (model: string, provider: string) => true;

const mockAgentStore = {
  currentAgentModel: 'gpt-4',
  currentAgentModelProvider: 'openai',
  hasKnowledge: false,
  isAgentConfigLoading: false,
  displayableAgentPlugins: ['plugin-1'],
  currentEnabledKnowledge: [],
  enableHistoryCount: true
};

const useAgentStore = (selector?: any, comparator?: any) => {
  if (selector) {
    return selector(mockAgentStore);
  }
  return mockAgentStore;
};

const agentSelectors = {
  currentAgentModel: (state: any) => state.currentAgentModel,
  currentAgentModelProvider: (state: any) => state.currentAgentModelProvider,
  hasKnowledge: (state: any) => state.hasKnowledge,
  isAgentConfigLoading: (state: any) => state.isAgentConfigLoading,
  displayableAgentPlugins: (state: any) => state.displayableAgentPlugins,
  currentEnabledKnowledge: (state: any) => state.currentEnabledKnowledge
};

const agentChatConfigSelectors = {
  enableHistoryCount: (state: any) => state.enableHistoryCount
};

const mockUserStore = {
  isLogin: true
};

const useUserStore = (selector?: any) => {
  if (selector) {
    return selector(mockUserStore);
  }
  return mockUserStore;
};

const authSelectors = {
  isLogin: (state: any) => state.isLogin
};

const mockSessionStore = {
  isCurrentSessionGroupSession: false
};

const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector(mockSessionStore);
  }
  return mockSessionStore;
};

const sessionSelectors = {
  isCurrentSessionGroupSession: (state: any) => state.isCurrentSessionGroupSession
};

import HistoryLimitTags from './HistoryLimitTags';
import KnowledgeTag from './KnowledgeTag';
import MemberCountTag from './MemberCountTag';
import SearchTags from './SearchTags';

const TitleTags = memo(() => {
  const [model, provider, hasKnowledge, isLoading] = useAgentStore((s) => [
    agentSelectors.currentAgentModel(s),
    agentSelectors.currentAgentModelProvider(s),
    agentSelectors.hasKnowledge(s),
    agentSelectors.isAgentConfigLoading(s),
  ]);

  const plugins = useAgentStore(agentSelectors.displayableAgentPlugins, isEqual);
  const enabledKnowledge = useAgentStore(agentSelectors.currentEnabledKnowledge, isEqual);
  const enableHistoryCount = useAgentStore(agentChatConfigSelectors.enableHistoryCount);

  const showPlugin = useModelSupportToolUse(model, provider);
  const isLogin = useUserStore(authSelectors.isLogin);
  const isGroupSession = useSessionStore(sessionSelectors.isCurrentSessionGroupSession);

  const isAgentEnableSearch = useAgentEnableSearch();

  if (isGroupSession) {
    return (
      <Flexbox align={'center'} gap={12} horizontal>
        <MemberCountTag />
      </Flexbox>
    );
  }

  return isLoading && isLogin ? (
    <Skeleton.Button active size={'small'} style={{ height: 20 }} />
  ) : (
    <Flexbox align={'center'} gap={4} horizontal>
      <ModelSwitchPanel>
        <ModelTag model={model} />
      </ModelSwitchPanel>
      {isAgentEnableSearch && <SearchTags />}
      {showPlugin && plugins?.length > 0 && <PluginTag plugins={plugins} />}
      {hasKnowledge && <KnowledgeTag data={enabledKnowledge} />}
      {enableHistoryCount && <HistoryLimitTags />}
    </Flexbox>
  );
});

export default TitleTags;
