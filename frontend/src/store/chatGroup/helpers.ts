/**
 * Chat Group Store Helper Utilities
 * 
 * Utility functions for mapping database results to frontend types
 * and converting frontend types to database params.
 */

import { ChatGroupItem, ChatGroupAgentItem } from '@/types/chatGroup';
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
 * Map Chat Group from database to frontend type
 */
export const mapChatGroupFromDB = (dbGroup: any): ChatGroupItem => {
  return {
    id: dbGroup.id,
    title: getNullableString(dbGroup.title) || 'Untitled Group',
    description: getNullableString(dbGroup.description) || '',
    config: parseNullableJSON(dbGroup.config) || {},
    groupId: getNullableString(dbGroup.groupId) || null,
    pinned: Boolean(dbGroup.pinned),
    createdAt: dbGroup.createdAt,
    updatedAt: dbGroup.updatedAt,
    userId: dbGroup.userId,
    clientId: getNullableString(dbGroup.clientId) || null,
    // Members will be populated separately
    members: [],
  };
};

/**
 * Map Chat Group Agent from database to frontend type
 */
export const mapChatGroupAgentFromDB = (dbAgent: any): ChatGroupAgentItem => {
  return {
    id: dbAgent.id,
    title: getNullableString(dbAgent.title) || '',
    description: getNullableString(dbAgent.description) || '',
    avatar: getNullableString(dbAgent.avatar) || '',
    backgroundColor: getNullableString(dbAgent.backgroundColor) || '',
    chatConfig: parseNullableJSON(dbAgent.chatConfig) || {},
    params: parseNullableJSON(dbAgent.params) || {},
    systemRole: getNullableString(dbAgent.systemRole) || '',
    tts: parseNullableJSON(dbAgent.tts) || null,
    model: getNullableString(dbAgent.model) || '',
    provider: getNullableString(dbAgent.provider) || '',
    createdAt: dbAgent.createdAt,
    updatedAt: dbGroup.updatedAt,
    // Junction table fields
    sortOrder: getNullableInt(dbAgent.sortOrder) || 0,
    enabled: Boolean(dbAgent.enabled),
    role: getNullableString(dbAgent.role) || 'member',
  };
};
