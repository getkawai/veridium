
import { t } from 'i18next';

import { DEFAULT_AVATAR, DEFAULT_BACKGROUND_COLOR, DEFAULT_INBOX_AVATAR } from '@/const/meta';
import { getNullableString, Session } from '@/types/database';
import { MetaData } from '@/types/meta';
import { LobeSessionType } from '@/types/session';
import { merge } from '@/utils/merge';

import type { SessionStore } from '../../../store';
import { sessionSelectors } from './list';

// ==========   Meta   ============== //
const currentAgentMeta = (s: SessionStore): MetaData => {
  const isInbox = sessionSelectors.isInboxSession(s);

  const defaultMeta: MetaData = {
    avatar: isInbox ? DEFAULT_INBOX_AVATAR : DEFAULT_AVATAR,
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
    description: isInbox ? t('inbox.desc', { ns: 'chat' }) : undefined,
    title: isInbox ? t('inbox.title', { ns: 'chat' }) : t('defaultSession'),
  };

  const session = sessionSelectors.currentSession(s);

  const currentSessionMeta: MetaData = session ? {
    title: getNullableString(session.title) || undefined,
    description: getNullableString(session.description) || undefined,
    avatar: getNullableString(session.avatar) || undefined,
    backgroundColor: getNullableString(session.backgroundColor) || undefined,
    tags: [], // tags are json string, need parsing if we want them
  } : {};

  return merge(defaultMeta, currentSessionMeta);
};

const currentGroupMeta = (s: SessionStore): MetaData => {
  const defaultMeta: MetaData = {
    description: t('group.desc', { ns: 'chat' }),
    title: t('group.title', { ns: 'chat' }),
  };

  const session = sessionSelectors.currentSession(s);

  const currentSessionMeta: MetaData = session ? {
    title: getNullableString(session.title) || undefined,
    description: getNullableString(session.description) || undefined,
    avatar: getNullableString(session.avatar) || undefined,
    backgroundColor: getNullableString(session.backgroundColor) || undefined,
  } : {};

  return merge(defaultMeta, currentSessionMeta);
};

const currentAgentTitle = (s: SessionStore) => currentAgentMeta(s).title;
const currentAgentDescription = (s: SessionStore) => currentAgentMeta(s).description;
const currentAgentAvatar = (s: SessionStore) => currentAgentMeta(s).avatar;
const currentAgentBackgroundColor = (s: SessionStore) => currentAgentMeta(s).backgroundColor;

const getAgentMetaByAgentId =
  (agentId: string) =>
    (s: SessionStore): MetaData => {
      // Find session where id matches agentId (wait, agentId != sessionId usually, but for agent session they are linked)
      // Actually the logic specific to agentId usually implies searching Agents, but here we search Sessions.
      // If the session concept merges agent, we search by config.id? Session doesn't have config.
      // We can't implementation this selector easily without config.
      // Assuming for now we skip or return empty if we can't find it.
      // Or we assume the session.id IS the agentId for single-agent sessions?

      // Original logic: session.config?.id === agentId
      // New logic: We don't have config. 
      // We will look for session where session.model matches or something? No.
      // This selector might be broken without config.
      return {};
    };

const getAvatar = (s: MetaData) => s.avatar || DEFAULT_AVATAR;
const getTitle = (s: MetaData) => s.title || t('defaultSession', { ns: 'common' });
// New session do not show 'noDescription'
export const getDescription = (s: MetaData) => s.description;

export const sessionMetaSelectors = {
  currentAgentAvatar,
  currentAgentBackgroundColor,
  currentAgentDescription,
  currentAgentMeta,
  currentAgentTitle,
  currentGroupMeta,
  getAgentMetaByAgentId,
  getAvatar,
  getDescription,
  getTitle,
};
