import { GenerationTopicItem } from '@/types/database-legacy';
import { ImageGenerationTopic } from '@/types/generation';
import { UpdateTopicValue } from '@/types/generation-types';

export interface IGenerationTopicService {
  getAllGenerationTopics(): Promise<ImageGenerationTopic[]>;
  createTopic(): Promise<string>;
  updateTopic(id: string, data: UpdateTopicValue): Promise<GenerationTopicItem | undefined>;
  updateTopicCover(id: string, coverUrl: string): Promise<GenerationTopicItem | undefined>;
  deleteTopic(id: string): Promise<GenerationTopicItem | undefined>;
}

