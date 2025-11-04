import { BaseDataModel } from '@/types/meta';
import { Topic, getNullableString, parseNullableJSON, intToBool } from '@/types/database';

// 类型定义
export type TimeGroupId =
  | 'today'
  | 'yesterday'
  | 'week'
  | 'month'
  | `${number}-${string}`
  | `${number}`;

/* eslint-disable typescript-sort-keys/string-enum */
export enum TopicDisplayMode {
  ByTime = 'byTime',
  Flat = 'flat',
  // AscMessages = 'ascMessages',
  // DescMessages = 'descMessages',
}
/* eslint-enable */

export interface GroupedTopic {
  children: ChatTopic[];
  id: string;
  title?: string;
}

export interface ChatTopicMetadata {
  model?: string;
  provider?: string;
}

export interface ChatTopicSummary {
  content: string;
  model: string;
  provider: string;
}

export interface ChatTopic extends Omit<BaseDataModel, 'meta'> {
  favorite?: boolean;
  historySummary?: string;
  metadata?: ChatTopicMetadata;
  sessionId?: string;
  title: string;
}

export type ChatTopicMap = Record<string, ChatTopic>;

export interface TopicRankItem {
  count: number;
  id: string;
  sessionId: string | null;
  title: string | null;
}

/**
 * Map database Topic type to frontend ChatTopic type
 */
export function mapTopicToChatTopic(topic: Topic): ChatTopic {
  return {
    id: topic.id,
    title: getNullableString(topic.title) || 'Untitled',
    favorite: topic.favorite ? intToBool(topic.favorite) : false,
    sessionId: getNullableString(topic.sessionId),
    historySummary: getNullableString(topic.historySummary),
    metadata: parseNullableJSON<ChatTopicMetadata>(topic.metadata),
    createdAt: topic.createdAt,
    updatedAt: topic.updatedAt,
  };
}

/**
 * Map array of database Topics to ChatTopics
 */
export function mapTopicsToChatTopics(topics: Topic[]): ChatTopic[] {
  return topics.map(mapTopicToChatTopic);
}
