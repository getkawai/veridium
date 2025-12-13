// import { clientDB } from '@/database/client/db';
import { GenerationTopicModel } from '@/database/models/generationTopic';
import { BaseClientService } from '@/services/baseClientService';
import { ImageGenerationTopic } from '@/types/generation';
import { UpdateTopicValue } from '@/types/generation-types';
import { GenerationTopic } from '@/database';
import { IGenerationTopicService } from './type';

export class ClientService extends BaseClientService implements IGenerationTopicService {
  private get generationTopicModel(): GenerationTopicModel {
    return new GenerationTopicModel(this.userId);
  }

  getAllGenerationTopics = async (): Promise<ImageGenerationTopic[]> => {
    return this.generationTopicModel.queryAll() as unknown as Promise<ImageGenerationTopic[]>;
  };

  createTopic = async (): Promise<string> => {
    const topic = await this.generationTopicModel.create('Untitled');
    return topic.id;
  };

  updateTopic = async (id: string, data: UpdateTopicValue): Promise<GenerationTopic | undefined> => {
    return this.generationTopicModel.update(id, data);
  };

  updateTopicCover = async (id: string, coverUrl: string): Promise<GenerationTopic | undefined> => {
    return this.generationTopicModel.update(id, { coverUrl });
  };

  deleteTopic = async (id: string): Promise<GenerationTopic | undefined> => {
    const result = await this.generationTopicModel.delete(id);
    if (result) {
      return result.deletedTopic;
    }
    return undefined;
  };
}

