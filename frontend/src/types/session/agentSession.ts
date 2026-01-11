import { Session, parseNullableJSON } from "@/types/database";
import { LobeAgentConfig } from "@/types/agent/item";
import { MetaData } from "@/types/meta";

export enum LobeSessionType {
  Agent = "agent",
  Group = "group",
}

// Alias Session to LobeSession to match user request
export type LobeSession = Session;
export type LobeSessions = LobeSession[];

// Extended interface for agent sessions with config and meta properties
export interface ExtendedLobeAgentSession extends LobeSession {
  config?: Partial<LobeAgentConfig>;
  meta?: MetaData;
}

// Extended interface for group sessions with members and meta properties
export interface GroupMember {
  id: string;
  title?: string;
  avatar?: string;
  backgroundColor?: string;
  name?: string;
}

export interface ExtendedLobeGroupSession extends LobeSession {
  members?: GroupMember[];
  meta?: MetaData;
}

// Type aliases for backward compatibility
export type LobeAgentSession = ExtendedLobeAgentSession;
export type LobeGroupSession = ExtendedLobeGroupSession;

/**
 * Get meta data from Session object
 * Session has flat properties (title, description, avatar, backgroundColor)
 * instead of nested meta object
 */
export const getSessionMeta = (session: LobeSession | undefined): MetaData => {
  if (!session) return {};

  const getNullable = (value: any): string | undefined => {
    if (!value) return undefined;
    if (typeof value === "object" && "Valid" in value && !value.Valid)
      return undefined;
    if (typeof value === "string") return value;
    if (typeof value === "object" && "String" in value && value.Valid) {
      return value.String;
    }
    return undefined;
  };

  return {
    title: getNullable(session.title),
    description: getNullable(session.description),
    avatar: getNullable(session.avatar),
    backgroundColor: getNullable(session.backgroundColor),
    tags: parseNullableJSON<string[]>(session.tags),
  };
};

/**
 * Get config from Session object
 * Config is stored separately and not available directly in Session
 * Must be fetched from AgentStore via agentMap[sessionId]
 */
export const getSessionConfig = (
  sessionId: string | undefined,
  agentMap: Record<string, Partial<LobeAgentConfig>>,
): Partial<LobeAgentConfig> | undefined => {
  if (!sessionId) return undefined;
  return agentMap[sessionId];
};

/**
 * Get members from Group Session object
 * Members are stored in separate table (ChatGroupAgents)
 * Must be fetched from ChatGroup store
 */
export const getSessionMembers = (
  sessionId: string | undefined,
  members: GroupMember[],
): GroupMember[] => {
  // Members passed from ChatGroup store
  return members || [];
};
