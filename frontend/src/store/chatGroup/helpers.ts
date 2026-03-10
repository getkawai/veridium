/**
 * Chat Group Store Helper Utilities
 * 
 * Utility functions for mapping database results to frontend types
 * and converting frontend types to database params.
 */

import { ChatGroupItem } from '@/types/chatGroup';
import {
  getNullableString,
  getNullableInt,
  parseNullableJSON,
} from '@/types/database';
import { getResolvedUserId } from '@/utils/userId';

/**
 * Get the default user ID
 */
export const getUserId = () => getResolvedUserId();

/**
 * Map Chat Group from database to frontend type
 */
export const mapChatGroupFromDB = (dbGroup: any): ChatGroupItem => {
  const config = parseNullableJSON(dbGroup.config);
  
  return {
    id: dbGroup.id,
    title: getNullableString(dbGroup.title) || 'Untitled Group',
    description: getNullableString(dbGroup.description) || '',
    config: config || null,
    groupId: getNullableString(dbGroup.groupId) || null,
    pinned: Boolean(dbGroup.pinned),
    createdAt: dbGroup.createdAt,
    updatedAt: dbGroup.updatedAt,
    userId: dbGroup.userId,
    clientId: getNullableString(dbGroup.clientId) || null,
  } as ChatGroupItem;
};

/**
 * Map Chat Group Agent from database to frontend type
 * Note: ChatGroupAgentItem is actually ChatGroupAgent (not the full agent)
 */
export const mapChatGroupAgentFromDB = (dbAgent: any) => {
  return {
    agentId: dbAgent.agentId || dbAgent.id,
    chatGroupId: dbAgent.chatGroupId,
    sortOrder: getNullableInt(dbAgent.sortOrder) || 0,
    enabled: Boolean(dbAgent.enabled),
    role: getNullableString(dbAgent.role) || 'member',
    createdAt: dbAgent.createdAt,
    updatedAt: dbAgent.updatedAt,
  };
};
