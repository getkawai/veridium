import { FilesConfigItem } from '../user/settings/filesConfig';

export enum KnowledgeBaseTabs {
  Files = 'files',
  Settings = 'Settings',
  Testing = 'testing',
}

export interface KnowledgeBaseItem {
  avatar: string | null;
  createdAt: Date;
  description?: string | null;
  enabled?: boolean;
  id: string;
  isPublic: boolean | null;
  name: string;
  settings: any;
  // different types of knowledge bases need to be distinguished
  type: string | null;
  updatedAt: Date;
  userId: string;
  clientId?: string | null;
}

/**
 * New knowledge base (for insert operations)
 * Equivalent to: typeof knowledgeBases.$inferInsert
 */
export interface NewKnowledgeBase {
  id?: string;
  name: string;
  description?: string | null;
  avatar?: string | null;
  type?: string | null;
  userId: string;
  clientId?: string | null;
  isPublic?: boolean;
  settings?: any;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Knowledge base files junction table
 * Equivalent to: typeof knowledgeBaseFiles.$inferInsert
 */
export interface KnowledgeBaseFile {
  knowledgeBaseId: string;
  fileId: string;
  userId: string;
  createdAt: Date;
}

/**
 * New knowledge base file link (for insert operations)
 */
export interface NewKnowledgeBaseFile {
  knowledgeBaseId: string;
  fileId: string;
  userId: string;
  createdAt?: Date;
}

export interface CreateKnowledgeBaseParams {
  avatar?: string;
  description?: string;
  name: string;
}

export enum KnowledgeType {
  File = 'file',
  KnowledgeBase = 'knowledgeBase',
}

export interface KnowledgeItem {
  avatar?: string | null;
  description?: string | null;
  enabled?: boolean;
  fileType?: string;
  id: string;
  name: string;
  type: KnowledgeType;
}

export interface SystemEmbeddingConfig {
  embeddingModel: FilesConfigItem;
  queryMode: string;
  rerankerModel: FilesConfigItem;
}
