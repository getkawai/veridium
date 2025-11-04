import { getNullableInt, SessionGroup } from '@/database';
import { LobeSessions } from './agentSession';

export type SessionGroupId = string;

export enum SessionDefaultGroup {
  Default = 'default',
  Pinned = 'pinned',
}

export interface SessionGroupItem {
  createdAt: Date;
  id: string;
  name: string;
  sort?: number | null;
  updatedAt: Date;
}

export type SessionGroups = SessionGroupItem[];

export interface CustomSessionGroup extends SessionGroupItem {
  children: LobeSessions;
}

export type LobeSessionGroups = SessionGroupItem[];

export const mapSessionGroup = (group: SessionGroup): SessionGroupItem => {
  return {
    createdAt: new Date(group.createdAt),
    id: group.id,
    name: group.name,
    sort: getNullableInt(group.sort as any) ?? null,
    updatedAt: new Date(group.updatedAt),
  };
};