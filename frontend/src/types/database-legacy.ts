/**
 * Legacy type compatibility layer
 * Provides Drizzle-style types for services/stores that still reference @/database/schemas
 * 
 * These map directly to the Go-generated database models.
 */

import type { 
  ChatGroup, 
  ChatGroupsAgent, 
  GenerationTopic, 
  User,
  File,
  Chunk,
  UnstructuredChunk,
} from '@@/github.com/kawai-network/veridium/internal/database/generated/models';

// Chat Group types
export type ChatGroupItem = ChatGroup;
export type ChatGroupAgentItem = ChatGroupsAgent;

// Generation types
export type GenerationTopicItem = GenerationTopic;

// User types  
export type UserItem = User;
export type NewUser = Omit<User, 'createdAt' | 'updatedAt'> & {
  createdAt?: number;
  updatedAt?: number;
};

// File types
export type FileItem = File;

// Chunk types
export type ChunkItem = Chunk;
export type NewChunkItem = Omit<Chunk, 'createdAt' | 'updatedAt' | 'id'> & {
  id?: string;
  createdAt?: number;
  updatedAt?: number;
};

export type UnstructuredChunkItem = UnstructuredChunk;
export type NewUnstructuredChunkItem = Omit<UnstructuredChunk, 'createdAt' | 'updatedAt' | 'id'> & {
  id?: string;
  createdAt?: number;
  updatedAt?: number;
};

// New insert types (for backward compatibility)
export type NewChatGroup = Omit<ChatGroup, 'id' | 'createdAt' | 'updatedAt'> & {
  id?: string;
  createdAt?: number;
  updatedAt?: number;
};

