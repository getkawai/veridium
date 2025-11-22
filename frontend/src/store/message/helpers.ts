/**
 * Message Store Helper Utilities
 * 
 * Utility functions for message-related operations
 */

import { INBOX_SESSION_ID } from '@/const/session';

const DEFAULT_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

/**
 * Get the default user ID
 */
export const getUserId = () => DEFAULT_USER_ID;

/**
 * Convert session ID to DB session ID
 * Inbox session is stored as empty string in DB
 */
export const toDbSessionId = (sessionId: string | undefined): string => {
  if (!sessionId || sessionId === INBOX_SESSION_ID) return '';
  return sessionId;
};

/**
 * Convert DB session ID back to session ID
 * Empty string in DB means inbox session
 */
export const fromDbSessionId = (dbSessionId: string | null | undefined): string => {
  if (!dbSessionId || dbSessionId === '') return INBOX_SESSION_ID;
  return dbSessionId;
};

