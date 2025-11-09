export interface LobeChatGroupMetaConfig {
  description: string;
  title: string;
}

export interface LobeChatGroupChatConfig {
  allowDM: boolean;
  enableSupervisor: boolean;
  maxResponseInRow: number;
  orchestratorModel: string;
  orchestratorProvider: string;
  responseOrder: 'sequential' | 'natural';
  responseSpeed: 'slow' | 'medium' | 'fast';
  revealDM: boolean;
  scene: 'casual' | 'productive';
  systemPrompt?: string;
}

// Database config type (flat structure)
export type LobeChatGroupConfig = LobeChatGroupChatConfig;

// Full group type with nested structure for UI components
export interface LobeChatGroupFullConfig {
  chat: LobeChatGroupChatConfig;
  meta: LobeChatGroupMetaConfig;
}

// Chat Group Agent types (independent from schema)
export interface ChatGroupAgent {
  agentId: string;
  chatGroupId: string;
  createdAt?: Date | number;
  enabled?: boolean;
  order?: number;
  role?: string | null;
  updatedAt?: Date | number;
  userId: string;
}

export interface NewChatGroupAgent {
  agentId: string;
  chatGroupId: string;
  enabled?: boolean;
  order?: number;
  role?: string;
  userId: string;
}

// Chat Group Item types (for database operations)
export interface ChatGroupItem {
  id: string;
  title: string | null;
  description: string | null;
  config: LobeChatGroupConfig | null;
  pinned: boolean;
  userId: string;
  createdAt: number;
  updatedAt: number;
}

export interface NewChatGroup {
  id?: string;
  title?: string | null;
  description?: string | null;
  config?: LobeChatGroupConfig | null;
  pinned?: boolean;
  userId: string;
  createdAt?: number;
  updatedAt?: number;
}

// ChatGroupAgentItem is the same as ChatGroupAgent (reuse)
export type ChatGroupAgentItem = ChatGroupAgent;
