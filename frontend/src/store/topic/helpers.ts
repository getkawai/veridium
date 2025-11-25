/**
 * Topic Store Helper Utilities
 * 
 * Utility functions for mapping database results to frontend types
 * and converting frontend types to database params.
 */

import { INBOX_SESSION_ID } from '@/const/session';
import {
  getNullableString,
  toNullString,
  toNullJSON,
  toNullInt,
  parseNullableJSON,
} from '@/types/database';

const DEFAULT_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

/**
 * Get the default user ID
 */
export const getUserId = () => DEFAULT_USER_ID;

/**
 * Convert sessionId to DB format (inbox -> null)
 */
export const toDbSessionId = (sessionId?: string | null) => {
  return sessionId;
};

/**
 * Convert DB sessionId to frontend format (null -> inbox)
 */
export const fromDbSessionId = (sessionId?: string | null) => {
  return sessionId || INBOX_SESSION_ID;
};

/**
 * Map Topic from database to frontend type
 */
export const mapTopicFromDB = (dbTopic: any) => {
  return {
    id: dbTopic.id,
    title: getNullableString(dbTopic.title) || 'Untitled',
    favorite: Boolean(dbTopic.favorite),
    sessionId: fromDbSessionId(getNullableString(dbTopic.sessionId)),
    groupId: getNullableString(dbTopic.groupId) || undefined,
    clientId: getNullableString(dbTopic.clientId) || undefined,
    historySummary: getNullableString(dbTopic.historySummary) || undefined,
    metadata: parseNullableJSON(dbTopic.metadata) || {},
    createdAt: Number(dbTopic.createdAt) || Date.now(),
    updatedAt: Number(dbTopic.updatedAt) || Date.now(),
  };
};

/**
 * Convert frontend topic data to database params for create
 */
export const topicToCreateParams = (data: any, userId: string) => {
  const now = Date.now();

  return {
    id: data.id || crypto.randomUUID(),
    title: toNullString(data.title || 'Untitled'),
    favorite: toNullInt(data.favorite ? 1 : 0),
    sessionId: toNullString(toDbSessionId(data.sessionId)),
    groupId: toNullString(data.groupId),
    userId,
    clientId: toNullString(data.clientId),
    historySummary: toNullString(data.historySummary),
    metadata: toNullJSON(data.metadata || {}),
    createdAt: data.createdAt ? data.createdAt.getTime() : now,
    updatedAt: data.updatedAt ? data.updatedAt.getTime() : now,
  };
};

/**
 * Convert frontend topic update data to database params
 */
export const topicToUpdateParams = (id: string, data: any, userId: string) => {
  const now = Date.now();

  return {
    id,
    userId,
    title: data.title ? toNullString(data.title) : undefined,
    historySummary: data.historySummary !== undefined ? toNullString(data.historySummary) : undefined,
    metadata: data.metadata ? toNullJSON(data.metadata) : undefined,
    updatedAt: now,
  };
};

