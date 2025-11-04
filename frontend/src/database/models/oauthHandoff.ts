import { NewOAuthHandoff, OAuthHandoffItem } from '../schemas';
import {
  DB,
  parseNullableJSON,
  currentTimestampMs,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

export class OAuthHandoffModel {
  private logger = createModelLogger('OAuthHandoff', 'OAuthHandoffModel', 'database/models/oauthHandoff');
  
  constructor(_db: any) {}

  /**
   * Create a new OAuth handoff record
   * @param params Credential data
   * @returns The created record
   */
  create = async (params: NewOAuthHandoff): Promise<OAuthHandoffItem> => {
    const now = currentTimestampMs();

    try {
      const result = await DB.CreateOAuthHandoff({
        id: params.id,
        client: params.client,
        payload: JSON.stringify(params.payload),
        createdAt: now,
        updatedAt: now,
      });

      return this.mapOAuthHandoff(result);
    } catch {
      // ON CONFLICT DO NOTHING simulation - if fails, just return empty
      return null as any;
    }
  };

  /**
   * Fetch and consume OAuth credentials
   * This method first queries the record, and if found, immediately deletes it to ensure the credential can only be used once
   * @param id Credential ID
   * @param client Client type
   * @returns Credential data, returns null if not exists or expired
   */
  fetchAndConsume = async (id: string, client: string): Promise<OAuthHandoffItem | null> => {
    // First find the record, and check if expired (5 minutes TTL)
    const fiveMinutesAgo = Date.now() - 5 * 60 * 1000;

    try {
      const handoff = await DB.GetOAuthHandoffByClient({
        id,
        client,
        createdAt: fiveMinutesAgo,
      });

      if (!handoff) {
        return null;
      }

      // Immediately delete the record to ensure one-time use
      await DB.DeleteOAuthHandoff(id);

      return this.mapOAuthHandoff(handoff);
    } catch {
      return null;
    }
  };

  /**
   * Clean up expired OAuth handoff records
   * This method should be called periodically (e.g., via cron job) to clean up expired records
   * @returns Number of records cleaned up
   */
  cleanupExpired = async (): Promise<number> => {
    const fiveMinutesAgo = Date.now() - 5 * 60 * 1000;

    try {
      await DB.CleanupExpiredOAuthHandoffs(fiveMinutesAgo);
      // Note: SQLite exec queries don't return rowCount easily
      // Would need a separate COUNT query for accurate results
      return 0;
    } catch {
      return 0;
    }
  };

  /**
   * Check if credential exists (without consuming)
   * Mainly used for testing and debugging
   * @param id Credential ID
   * @param client Client type
   * @returns Whether exists and not expired
   */
  exists = async (id: string, client: string): Promise<boolean> => {
    const fiveMinutesAgo = Date.now() - 5 * 60 * 1000;

    try {
      const handoff = await DB.GetOAuthHandoffByClient({
        id,
        client,
        createdAt: fiveMinutesAgo,
      });
      return !!handoff;
    } catch {
      return false;
    }
  };

  // **************** Helper *************** //

  private mapOAuthHandoff = (item: any): OAuthHandoffItem => {
    return {
      id: item.id,
      client: item.client,
      payload: parseNullableJSON(item.payload as any) || {},
      createdAt: new Date(item.createdAt),
      updatedAt: new Date(item.updatedAt),
    } as OAuthHandoffItem;
  };
}
