/**
 * User Store Helper Utilities
 * 
 * Utility functions for user-related operations
 */

import { getResolvedUserId } from '@/utils/userId';

/**
 * Get the default user ID
 */
export const getUserId = () => getResolvedUserId();
