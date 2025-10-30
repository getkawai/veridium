import { Icon, Tag, Tooltip } from '@lobehub/ui';
import { Users } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

// import { useSessionStore } from '@/store/session';
// import { sessionSelectors } from '@/store/session/selectors';
import { LobeGroupSession } from '@/types/session';

// Dummy implementations for development - memoized
const mockSessionStore = {
  currentSession: {
    type: 'group',
    members: [
      { id: 'member-1', name: 'User 1' },
      { id: 'member-2', name: 'User 2' },
      { id: 'member-3', name: 'User 3' }
    ]
  }
};

const useSessionStore = (selector?: any) => {
  if (selector) {
    return selector(mockSessionStore);
  }
  return mockSessionStore;
};

const sessionSelectors = {
  currentSession: (state: any) => state.currentSession
};

const MemberCountTag = memo(() => {
  const { t } = useTranslation('chat');
  const currentSession = useSessionStore(sessionSelectors.currentSession);

  const memberCount = (currentSession as LobeGroupSession).members?.length ?? 0 + 1;

  if (memberCount < 0) return null;

  return (
    <Tooltip title={t('group.memberTooltip', { count: memberCount })}>
      <Flexbox height={22}>
        <Tag>
          <Icon icon={Users} />
          <span>{memberCount}</span>
        </Tag>
      </Flexbox>
    </Tooltip>
  );
});

export default MemberCountTag;
