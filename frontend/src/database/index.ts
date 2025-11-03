/**
 * Database Module
 * 
 * Central export point for all database-related functionality.
 * This replaces the old Drizzle ORM setup with direct Wails bindings.
 */

// Export all types and utilities
export * from '@/types/database';

// Export queries with a clean namespace
export { DB } from '@/types/database';

/**
 * Usage examples:
 * 
 * ```typescript
 * import { DB, Agent, User, toNullString, parseNullableJSON } from '@/database';
 * 
 * // Use queries
 * const user = await DB.GetUser('user-123');
 * const agents = await DB.ListAgents({ userID: 'user-123', limit: 10, offset: 0 });
 * 
 * // Work with nullable fields
 * const agent = await DB.GetAgent({ id: 'agent-123', userID: 'user-123' });
 * const title = parseNullableJSON(agent.title); // Extract value
 * 
 * // Create params
 * const params = {
 *   id: 'new-agent',
 *   userId: 'user-123',
 *   title: toNullString('My Agent'),
 *   // ...
 * };
 * const newAgent = await DB.CreateAgent(params);
 * ```
 */
