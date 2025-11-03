import { getClientDBConfig } from '@/database/client/db';

import { ClientService } from './client';

export const knowledgeBaseService =
  getClientDBConfig().mode === 'client' ? new ClientService() : null;

if (!knowledgeBaseService) {
  throw new Error('KnowledgeBase service not initialized - client mode required');
}

export { type IKnowledgeBaseService } from './type';

