import { z } from 'zod';

import { ChatSemanticSearchChunk } from './chunk';

export const SemanticSearchSchema = z.object({
  fileIds: z.array(z.string()).optional(),
  knowledgeIds: z.array(z.string()).optional(),
  messageId: z.string(),
  model: z.string().optional(),
  rewriteQuery: z.string(),
  userQuery: z.string(),
});

export type SemanticSearchSchemaType = z.infer<typeof SemanticSearchSchema>;

export type MessageSemanticSearchChunk = Pick<ChatSemanticSearchChunk, 'id' | 'similarity'>;

export type SemanticSearchResult = {
  id: string;
  similarity: number;
  text: string | null;
  fileId: string;
  fileName: string;
  type: string | null;
  index: number;
  metadata: any;
};
