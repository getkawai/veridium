import { GenerationTopic } from '@/database';
import { ImageGenerationTopic } from '@/types/generation';
import { UpdateTopicValue } from '@/types/generation-types';

export interface IGenerationTopicService {
  getAllGenerationTopics(): Promise<ImageGenerationTopic[]>;
  createTopic(): Promise<string>;
  updateTopic(id: string, data: UpdateTopicValue): Promise<GenerationTopic | undefined>;
  updateTopicCover(id: string, coverUrl: string): Promise<GenerationTopic | undefined>;
  deleteTopic(id: string): Promise<GenerationTopic | undefined>;
}

