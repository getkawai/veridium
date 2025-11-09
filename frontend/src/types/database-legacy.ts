/**
 * Legacy type compatibility layer
 * Provides Drizzle-style types for services/stores that still reference @/database/schemas
 * 
 * These map directly to the Go-generated database models.
 */

import type { 
  User,
  File,
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
// Convert Go-generated Chunk to TypeScript-friendly NewChunkItem
export interface NewChunkItem {
  id?: string;
  text?: string | null;
  abstract?: string | null;
  metadata?: string | null;
  index?: number;
  type?: string | null;
  clientId?: string | null;
  userId?: string;
  createdAt?: number;
  updatedAt?: number;
}

// Convert Go-generated UnstructuredChunk to TypeScript-friendly NewUnstructuredChunkItem
export interface NewUnstructuredChunkItem {
  id?: string;
  text?: string | null;
  metadata?: string | null;
  index?: number;
  type?: string | null;
  parentId?: string | null;
  compositeId?: string | null;
  clientId?: string | null;
  userId?: string;
  fileId?: string | null;
  createdAt?: number;
  updatedAt?: number;
}

