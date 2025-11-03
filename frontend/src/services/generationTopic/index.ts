import { getClientDBConfig } from '@/database/client/db';

import { ClientService } from './client';

export const generationTopicService =
  getClientDBConfig().mode === 'client' ? new ClientService() : null;

if (!generationTopicService) {
  throw new Error('GenerationTopic service not initialized - client mode required');
}

export { type IGenerationTopicService } from './type';

