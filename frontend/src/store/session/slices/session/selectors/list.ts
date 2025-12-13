import { DEFAULT_AGENT_LOBE_SESSION, INBOX_SESSION_ID } from '@/const/session';
import { sessionHelpers } from '@/store/session/slices/session/helpers';
import { MetaData } from '@/types/meta';
import {
  CustomSessionGroup,
  LobeSession,
  LobeSessions,
} from '@/types/session';

import { SessionStore } from '../../../store';

const defaultSessions = (s: SessionStore): LobeSessions => s.defaultSessions;
const pinnedSessions = (s: SessionStore): LobeSessions => s.pinnedSessions;
const customSessionGroups = (s: SessionStore): CustomSessionGroup[] => s.customSessionGroups;

const allSessions = (s: SessionStore): LobeSessions => s.sessions;

const getSessionById =
  (id: string) =>
    (s: SessionStore): LobeSession =>
      sessionHelpers.getSessionById(id, allSessions(s));

const getSessionMetaById =
  (id: string) =>
    (s: SessionStore): MetaData => {
      const session = getSessionById(id)(s);

      if (!session) return {};
      // Session type has flat properties, not nested meta
      return {
        title: typeof session.title === 'string' ? session.title : (session.title as any)?.String,
        description: typeof session.description === 'string' ? session.description : (session.description as any)?.String,
        avatar: typeof session.avatar === 'string' ? session.avatar : (session.avatar as any)?.String,
        backgroundColor: typeof session.backgroundColor === 'string' ? session.backgroundColor : (session.backgroundColor as any)?.String,
      };
    };

const currentSession = (s: SessionStore): LobeSession | undefined => {
  if (!s.activeId) return;

  return allSessions(s).find((i) => i.id === s.activeId);
};

const currentSessionSafe = (s: SessionStore): LobeSession => {
  return currentSession(s) || DEFAULT_AGENT_LOBE_SESSION;
};

const hasCustomAgents = (s: SessionStore) => defaultSessions(s).length > 0;

const isInboxSession = (s: SessionStore) => s.activeId === INBOX_SESSION_ID;

const isCurrentSessionGroupSession = (s: SessionStore): boolean => {
  const session = currentSession(s);
  const sessionType = typeof session?.type === 'string' ? session.type : (session?.type as any)?.String;
  return sessionType === 'group';
};

const currentGroupAgents = (s: SessionStore): any[] => {
  const session = currentSession(s);
  if (!session) return [];

  const sessionType = typeof session.type === 'string' ? session.type : (session.type as any)?.String;
  if (sessionType !== 'group') return [];

  // Group sessions don't have members in the Session table anymore
  // This would need to be fetched from a separate table
  return [];
};

const isSessionListInit = (s: SessionStore) => s.isSessionsFirstFetchFinished;

const hasMoreSessions = (s: SessionStore) => s.sessionsHasMore;

// use to judge whether a session is fully activated
const isSomeSessionActive = (s: SessionStore) => !!s.activeId && isSessionListInit(s);

export const sessionSelectors = {
  currentGroupAgents,
  currentSession,
  currentSessionSafe,
  customSessionGroups,
  defaultSessions,
  getSessionById,
  getSessionMetaById,
  hasCustomAgents,
  hasMoreSessions,
  isCurrentSessionGroupSession,
  isInboxSession,
  isSessionListInit,
  isSomeSessionActive,
  pinnedSessions,
};
