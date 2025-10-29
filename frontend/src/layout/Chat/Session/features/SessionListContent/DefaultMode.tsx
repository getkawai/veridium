import { CollapseProps } from 'antd';
import isEqual from 'fast-deep-equal';
import { memo, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  LobeAgentSession,
  LobeSessionType,
  LobeSessions,
  SessionDefaultGroup,
} from '@/types/session';

import CollapseGroup from './CollapseGroup';
import Actions from './CollapseGroup/Actions';
import Inbox from './Inbox';
import SessionList from './List';
import ConfigGroupModal from './Modals/ConfigGroupModal';
import RenameGroupModal from './Modals/RenameGroupModal';

// Dummy implementations for development
const useFetchSessions = () => {
  console.log('Mock useFetchSessions called');
  // No-op for dummy implementation
};

const useServerConfigStore = (selector: any) => {
  if (selector) {
    return selector({ isMobile: false });
  }
  return { isMobile: false };
};

const serverConfigSelectors = {
  isMobile: (state: any) => state.isMobile,
};

const useSessionStore = (selector?: any, comparator?: any) => {
  if (selector) {
    return selector({
      defaultSessions: [],
      customSessionGroups: [],
      pinnedSessions: [],
      isSearching: false,
      sessionSearchKeywords: '',
      useSearchSessions: () => ({ data: [], isLoading: false }),
    });
  }
  return {
    defaultSessions: [],
    customSessionGroups: [],
    pinnedSessions: [],
    isSearching: false,
    sessionSearchKeywords: '',
    useSearchSessions: () => ({ data: [], isLoading: false }),
  };
};

const sessionSelectors = {
  defaultSessions: (state: any) => state.defaultSessions,
  customSessionGroups: (state: any) => state.customSessionGroups,
  pinnedSessions: (state: any) => state.pinnedSessions,
};

const useGlobalStore = (selector: any) => {
  if (selector) {
    return selector({
      sessionGroupKeys: [SessionDefaultGroup.Default],
      updateSystemStatus: (status: any) => {
        console.log('Mock updateSystemStatus called with:', status);
      },
    });
  }
  return {
    sessionGroupKeys: [SessionDefaultGroup.Default],
    updateSystemStatus: (status: any) => {
      console.log('Mock updateSystemStatus called with:', status);
    },
  };
};

const systemStatusSelectors = {
  sessionGroupKeys: (state: any) => state.sessionGroupKeys,
};

const DefaultMode = memo(() => {
  const { t } = useTranslation('chat');

  const [activeGroupId, setActiveGroupId] = useState<string>();
  const [renameGroupModalOpen, setRenameGroupModalOpen] = useState(false);
  const [configGroupModalOpen, setConfigGroupModalOpen] = useState(false);

  useFetchSessions();

  const isMobile = useServerConfigStore(serverConfigSelectors.isMobile);

  const defaultSessions = useSessionStore(sessionSelectors.defaultSessions, isEqual);
  const customSessionGroups = useSessionStore(sessionSelectors.customSessionGroups, isEqual);
  const pinnedSessions = useSessionStore(sessionSelectors.pinnedSessions, isEqual);

  const shouldHideSession = (session: LobeSessions[0]) =>
    !isMobile &&
    session.type === LobeSessionType.Agent &&
    Boolean((session as LobeAgentSession).config?.virtual);

  const filterSessionsForView = (sessions: LobeSessions): LobeSessions => {
    const filteredForDevice = isMobile
      ? sessions.filter((session) => session.type !== LobeSessionType.Group)
      : sessions;

    if (isMobile) return filteredForDevice;

    return filteredForDevice.filter((session) => !shouldHideSession(session));
  };

  const filteredDefaultSessions = filterSessionsForView(defaultSessions);
  const filteredPinnedSessions = filterSessionsForView(pinnedSessions);
  const filteredCustomSessionGroups = customSessionGroups?.map((group) => ({
    ...group,
    children: filterSessionsForView(group.children),
  }));

  const [sessionGroupKeys, updateSystemStatus] = useGlobalStore((s) => [
    systemStatusSelectors.sessionGroupKeys(s),
    s.updateSystemStatus,
  ]);

  const items = useMemo(
    () =>
      [
        filteredPinnedSessions &&
          filteredPinnedSessions.length > 0 && {
            children: <SessionList dataSource={filteredPinnedSessions} />,
            extra: <Actions isPinned openConfigModal={() => setConfigGroupModalOpen(true)} />,
            key: SessionDefaultGroup.Pinned,
            label: t('pin'),
          },
        ...(filteredCustomSessionGroups || []).map(({ id, name, children }) => ({
          children: <SessionList dataSource={children} groupId={id} />,
          extra: (
            <Actions
              id={id}
              isCustomGroup
              onOpenChange={(isOpen) => {
                if (isOpen) setActiveGroupId(id);
              }}
              openConfigModal={() => setConfigGroupModalOpen(true)}
              openRenameModal={() => setRenameGroupModalOpen(true)}
            />
          ),
          key: id,
          label: name,
        })),
        {
          children: <SessionList dataSource={filteredDefaultSessions || []} />,
          extra: <Actions openConfigModal={() => setConfigGroupModalOpen(true)} />,
          key: SessionDefaultGroup.Default,
          label: t('defaultList'),
        },
      ].filter(Boolean) as CollapseProps['items'],
    [t, filteredCustomSessionGroups, filteredPinnedSessions, filteredDefaultSessions],
  );

  return (
    <>
      <Inbox />
      <CollapseGroup
        activeKey={sessionGroupKeys}
        items={items}
        onChange={(keys) => {
          const expandSessionGroupKeys = typeof keys === 'string' ? [keys] : keys;

          updateSystemStatus({ expandSessionGroupKeys });
        }}
      />
      {activeGroupId && (
        <RenameGroupModal
          id={activeGroupId}
          onCancel={() => setRenameGroupModalOpen(false)}
          open={renameGroupModalOpen}
        />
      )}
      <ConfigGroupModal
        onCancel={() => setConfigGroupModalOpen(false)}
        open={configGroupModalOpen}
      />
    </>
  );
});

DefaultMode.displayName = 'SessionDefaultMode';

export default DefaultMode;
