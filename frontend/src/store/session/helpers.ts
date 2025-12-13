
/**
 * Session Store Helper Utilities
 * 
 * Utility functions for mapping database results to frontend types
 * and converting frontend types to database params.
 */

import { LobeAgentConfig, LobeSessionType } from '@/types';
import { DEFAULT_AGENT_CHAT_CONFIG } from '@/const/settings';
import {
  getNullableString,
  getNullableInt,
  toNullString,
  toNullJSON,
  toNullInt,
  boolToInt,
  parseNullableJSON,
} from '@/types/database';

const DEFAULT_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

/**
 * Get the default user ID
 */
export const getUserId = () => DEFAULT_USER_ID;

/**
 * Map Session from database to frontend type
 */
// Map function removed as we use Session type directly
// export const mapSessionFromDB = ...

/**
 * Map Agent config from database to LobeAgentConfig
 */
export const mapAgentConfigFromDB = (agent: any): LobeAgentConfig => {
  return {
    id: agent.id,
    model: getNullableString(agent.model) || '',
    systemRole: getNullableString(agent.systemRole) || '',
    plugins: parseNullableJSON(agent.plugins) || [],
    chatConfig: parseNullableJSON(agent.chatConfig) || DEFAULT_AGENT_CHAT_CONFIG,
    params: parseNullableJSON(agent.params) || {},
    openingMessage: getNullableString(agent.openingMessage),
    openingQuestions: parseNullableJSON(agent.openingQuestions) || [],
    fewShots: parseNullableJSON(agent.fewShots),
    virtual: Boolean(agent.virtual),
    provider: getNullableString(agent.provider),
  };
};

/**
 * Map SessionGroup from database to frontend type
 */
export const mapSessionGroupFromDB = (dbGroup: any) => {
  return {
    id: dbGroup.id,
    name: getNullableString(dbGroup.name) || '',
    sort: getNullableInt(dbGroup.sort) || 0,
    createdAt: new Date(dbGroup.createdAt),
    updatedAt: new Date(dbGroup.updatedAt),
    userId: dbGroup.userId,
  };
};

/**
 * Convert frontend session data to database params for create
 */
export const sessionToCreateParams = (type: LobeSessionType, data: any, userId: string) => {
  const now = Date.now();

  return {
    session: {
      id: data.id || crypto.randomUUID(),
      type,
      slug: toNullString(data.slug),
      title: toNullString(data.meta?.title || data.title),
      description: toNullString(data.meta?.description || data.description),
      avatar: toNullString(data.meta?.avatar || data.avatar),
      backgroundColor: toNullString(data.meta?.backgroundColor || data.backgroundColor),
      groupId: toNullString(data.group === 'default' ? null : data.group),
      pinned: toNullInt(data.pinned ? 1 : 0),
      userId,
      createdAt: data.createdAt ? data.createdAt.getTime() : now,
      updatedAt: data.updatedAt ? data.updatedAt.getTime() : now,
    },
    agent: {
      id: data.config?.id || crypto.randomUUID(),
      sessionId: data.id || '',
      model: toNullString(data.config?.model),
      systemRole: toNullString(data.config?.systemRole),
      plugins: toNullJSON(data.config?.plugins),
      chatConfig: toNullJSON(data.config?.chatConfig),
      params: toNullJSON(data.config?.params),
      openingMessage: toNullString(data.config?.openingMessage),
      openingQuestions: toNullJSON(data.config?.openingQuestions),
      fewShots: toNullJSON(data.config?.fewShots),
      virtual: toNullInt(data.config?.virtual ? 1 : 0),
      provider: toNullString(data.config?.provider),
      userId,
      createdAt: now,
      updatedAt: now,
    },
  };
};

/**
 * Convert frontend session update data to database params
 */
export const sessionToUpdateParams = (id: string, data: any, userId: string) => {
  const now = Date.now();

  return {
    id,
    userId,
    title: data.title ? toNullString(data.title) : undefined,
    description: data.description ? toNullString(data.description) : undefined,
    avatar: data.avatar ? toNullString(data.avatar) : undefined,
    backgroundColor: data.backgroundColor ? toNullString(data.backgroundColor) : undefined,
    groupId: data.group !== undefined ? toNullString(data.group === 'default' ? null : data.group) : undefined,
    pinned: data.pinned !== undefined ? toNullInt(boolToInt(data.pinned)) : undefined,
    updatedAt: now,
  };
};

/**
 * Get session pinned status
 */
export const getSessionPinned = (session: any) => session.pinned;

export const sessionHelpers = {
  getSessionPinned,
};
