import { nanoid } from 'nanoid';

import { SessionGroupItem } from '@/types/session/sessionGroup';
import {
  DB,
  toNullString,
  toNullInt,
  getNullableInt,
  currentTimestampMs,
  SessionGroup,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';

export class SessionGroupModel {
  private userId: string;
  private logger = createModelLogger('SessionGroup', 'SessionGroupModel', 'database/models/sessionGroup');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  /**
   * Show error notification to user
   */
  private async showErrorNotification(title: string, message: string) {
    try {
      await NotificationService.SendNotification(
        new NotificationOptions({
          id: `session-group-error-${Date.now()}`,
          title: `Session Group Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  create = async (params: { name: string; sort?: number }) => {
    await this.logger.methodEntry('create', { name: params.name, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      const result = await DB.CreateSessionGroup({
        id: this.genId(),
        name: params.name,
        sort: toNullInt(params.sort as any),
        userId: this.userId,
        clientId: toNullString(''),
        createdAt: now,
        updatedAt: now,
      });

      await this.logger.methodExit('create', { groupId: result.id });
      return this.mapSessionGroup(result);
    } catch (error) {
      await this.logger.error('Failed to create session group', { error, params });
      await this.showErrorNotification(
        'Create Failed',
        `Failed to create session group "${params.name}". Please try again.`
      );
      throw error;
    }
  };

  delete = async (id: string) => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    try {
      await DB.DeleteSessionGroup({
        id,
        userId: this.userId,
      });
      
      await this.logger.methodExit('delete', { id });
    } catch (error) {
      await this.logger.error('Failed to delete session group', { error, id });
      await this.showErrorNotification(
        'Delete Failed',
        `Failed to delete session group. Please try again.`
      );
      throw error;
    }
  };

  deleteAll = async () => {
    try {
      await DB.DeleteAllSessionGroups(this.userId);
    } catch (error) {
      await this.logger.error('Failed to delete all session groups', { error });
      await this.showErrorNotification(
        'Delete All Failed',
        `Failed to delete all session groups. Please try again.`
      );
      throw error;
    }
  };

  query = async () => {
    try {
      const groups = await DB.ListSessionGroups(this.userId);
      return groups.map((g) => this.mapSessionGroup(g));
    } catch (error) {
      await this.logger.error('Failed to query session groups', { error });
      throw error;
    }
  };

  findById = async (id: string) => {
    try {
      const group = await DB.GetSessionGroup({
        id,
        userId: this.userId,
      });
      return this.mapSessionGroup(group);
    } catch (error) {
      await this.logger.warn('Session group not found', { id, error });
      return undefined;
    }
  };

  update = async (id: string, value: Partial<SessionGroupItem>) => {
    await this.logger.methodEntry('update', { id, value, userId: this.userId });
    
    try {
      const now = currentTimestampMs();

      await DB.UpdateSessionGroup({
        id,
        userId: this.userId,
        name: value.name || '',
        sort: toNullInt(value.sort as any),
        updatedAt: now,
      });
      
      await this.logger.methodExit('update', { id });
    } catch (error) {
      await this.logger.error('Failed to update session group', { error, id, value });
      await this.showErrorNotification(
        'Update Failed',
        `Failed to update session group "${value.name || ''}". Please try again.`
      );
      throw error;
    }
  };

  updateOrder = async (sortMap: { id: string; sort: number }[]) => {
    await this.logger.methodEntry('updateOrder', { count: sortMap.length, userId: this.userId });
    
    try {
      // Note: No transaction support in Wails!
      // This is a potential data consistency issue
      
      const now = currentTimestampMs();

      await Promise.all(
        sortMap.map(({ id, sort }) =>
          DB.UpdateSessionGroupOrder({
            id,
            userId: this.userId,
            sort: toNullInt(sort as any),
            updatedAt: now,
          }),
        ),
      );
      
      await this.logger.methodExit('updateOrder', { count: sortMap.length });
    } catch (error) {
      await this.logger.error('Failed to update session group order', { error, count: sortMap.length });
      await this.showErrorNotification(
        'Update Order Failed',
        `Failed to update session group order. Please try again.`
      );
      throw error;
    }
  };

  // **************** Helper *************** //

  private genId = () => nanoid();

  private mapSessionGroup = (group: SessionGroup): SessionGroupItem => {
    return {
      id: group.id,
      name: group.name,
      sort: getNullableInt(group.sort as any) ?? null,
      createdAt: new Date(group.createdAt),
      updatedAt: new Date(group.updatedAt),
    };
  };
}

