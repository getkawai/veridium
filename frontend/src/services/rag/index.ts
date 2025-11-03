import { getClientDBConfig } from '@/database/client/db';

import { ClientService } from './client';

export const ragService = getClientDBConfig().mode === 'client' ? new ClientService() : null;

if (!ragService) {
  throw new Error('RAG service not initialized - client mode required');
}

export { type IRAGService } from './type';

