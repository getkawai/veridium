import { clientDB } from '@/database/client/db';
import { GenerationTopicItem } from '@/database/schemas';
import { GenerationTopicModel } from '@/database/models/generationTopic';
import { BaseClientService } from '@/services/baseClientService';
import { ImageGenerationTopic } from '@/types/generation';
import { UpdateTopicValue } from '@/types/generation-types';

import { IGenerationTopicService } from './type';

export class ClientService extends BaseClientService implements IGenerationTopicService {
  private get generationTopicModel(): GenerationTopicModel {
    return new GenerationTopicModel(clientDB as any, this.userId);
  }

  getAllGenerationTopics = async (): Promise<ImageGenerationTopic[]> => {
    return this.generationTopicModel.queryAll() as Promise<ImageGenerationTopic[]>;
  };

  createTopic = async (): Promise<string> => {
    const topic = await this.generationTopicModel.create('Untitled');
    return topic.id;
  };

  updateTopic = async (id: string, data: UpdateTopicValue): Promise<GenerationTopicItem | undefined> => {
    return this.generationTopicModel.update(id, data);
  };

  updateTopicCover = async (id: string, coverUrl: string): Promise<GenerationTopicItem | undefined> => {
    return this.generationTopicModel.update(id, { coverUrl });
  };

  deleteTopic = async (id: string): Promise<GenerationTopicItem | undefined> => {
    await this.generationTopicModel.delete(id);
    return undefined; // Model delete doesn't return value
  };
}

