import { Session } from '@/types/database';

export enum LobeSessionType {
  Agent = 'agent',
  Group = 'group',
}

// Alias Session to LobeSession to match user request
export type LobeSession = Session;
export type LobeSessions = LobeSession[];

