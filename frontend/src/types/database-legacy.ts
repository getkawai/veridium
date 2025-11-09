/**
 * Legacy type compatibility layer
 * Provides Drizzle-style types for services/stores that still reference @/database/schemas
 * 
 * These map directly to the Go-generated database models.
 */

import type { 
  User,
  File,
  Chunk,
  UnstructuredChunk,
} from '@@/github.com/kawai-network/veridium/internal/database/generated/models';

// Generation types
// Map GenerationTopic from Go-generated models to TypeScript-friendly types
export interface GenerationTopicItem {
  id: string;
  userId: string;
  title: string | null;
  coverUrl: string | null;
  createdAt: number;
  updatedAt: number;
}

// User types  
export type UserItem = User;
export type NewUser = Omit<User, 'createdAt' | 'updatedAt'> & {
  createdAt?: number;
  updatedAt?: number;
};

// File types
export type FileItem = File;

// Chunk types
export type NewChunkItem = Omit<Chunk, 'createdAt' | 'updatedAt' | 'id'> & {
  id?: string;
  createdAt?: number;
  updatedAt?: number;
};

export type NewUnstructuredChunkItem = Omit<UnstructuredChunk, 'createdAt' | 'updatedAt' | 'id'> & {
  id?: string;
  createdAt?: number;
  updatedAt?: number;
};

