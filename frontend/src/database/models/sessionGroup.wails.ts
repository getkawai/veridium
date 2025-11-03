import { nanoid } from 'nanoid';

import { SessionGroupItem } from '../schemas';
import {
  DB,
  toNullString,
  toNullInt,
  getNullableString,
  currentTimestampMs,
} from '@/types/database';

export class SessionGroupModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: { name: string; sort?: number }) => {
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

    return this.mapSessionGroup(result);
  };

  delete = async (id: string) => {
    await DB.DeleteSessionGroup({
      id,
      userId: this.userId,
    });
  };

  deleteAll = async () => {
    await DB.DeleteAllSessionGroups(this.userId);
  };

  query = async () => {
    const groups = await DB.ListSessionGroups(this.userId);
    return groups.map((g) => this.mapSessionGroup(g));
  };

  findById = async (id: string) => {
    try {
      const group = await DB.GetSessionGroup({
        id,
        userId: this.userId,
      });
      return this.mapSessionGroup(group);
    } catch {
      return undefined;
    }
  };

  update = async (id: string, value: Partial<SessionGroupItem>) => {
    const now = currentTimestampMs();

    await DB.UpdateSessionGroup({
      id,
      userId: this.userId,
      name: value.name || '',
      sort: toNullInt(value.sort as any),
      updatedAt: now,
    });
  };

  updateOrder = async (sortMap: { id: string; sort: number }[]) => {
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
  };

  // **************** Helper *************** //

  private genId = () => nanoid();

  private mapSessionGroup = (group: any): SessionGroupItem => {
    return {
      id: group.id,
      name: group.name,
      sort: group.sort,
      userId: group.userId,
      clientId: getNullableString(group.clientId as any),
      createdAt: new Date(group.createdAt),
      updatedAt: new Date(group.updatedAt),
    } as SessionGroupItem;
  };
}

