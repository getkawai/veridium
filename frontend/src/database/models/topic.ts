import { ChatTopic, DBMessageItem, mapTopicsToChatTopics, mapTopicToChatTopic, TopicRankItem } from '@/types';
import { nanoid } from 'nanoid';
import { createModelLogger } from '@/utils/logger';
import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
  Topic,
} from '@/types/database';

export interface CreateTopicParams {
  favorite?: boolean;
  groupId?: string | null;
  messages?: string[];
  sessionId?: string | null;
  title?: string;
}

interface QueryTopicParams {
  containerId?: string | null; // sessionId or groupId
  current?: number;
  pageSize?: number;
}

export class TopicModel {
  private userId: string;
  private logger = createModelLogger('Topic', 'TopicModel', 'database/models/topic');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // **************** Query *************** //

  query = async ({ current = 0, pageSize = 9999, containerId }: QueryTopicParams = {}) => {
    const offset = current * pageSize;

    // Note: Drizzle uses complex WHERE with OR for containerId
    // Wails requires separate queries or custom SQL
    if (containerId) {
      // Try session first
      const sessionTopics = await DB.ListTopics({
        userId: this.userId,
        sessionId: toNullString(containerId),
        limit: pageSize,
        offset,
      });

      if (sessionTopics.length > 0) {
        return sessionTopics;
      }

      // Try group if no session topics found
      // TODO: Add ListTopicsByGroup query
      // For now, return empty
      return [];
    }

    // If no containerId, return topics with no session/group
    // TODO: Add ListTopicsWithoutContainer query
    return [];
  };

  findById = async (id: string) => {
    try {
      const topic = await DB.GetTopic({
        id,
        userId: this.userId,
      });
      return topic;
    } catch {
      return undefined;
    }
  };

  queryAll = async (): Promise<ChatTopic[]> => {
    const topics = await DB.ListAllTopics(this.userId);
    return mapTopicsToChatTopics(topics);
  };

  queryByKeyword = async (keyword: string, containerId?: string | null): Promise<ChatTopic[]> => {
    if (!keyword) return [];

    const keywordLowerCase = keyword.toLowerCase();
    const containerParam = containerId || '';

    // Search by title
    const topicsByTitle = await DB.SearchTopicsByTitle({
      userId: this.userId,
      title: toNullString(`%${keywordLowerCase}%`) as any,
      column3: containerParam, // containerId check
      sessionId: toNullString(containerParam) as any,
      groupId: toNullString(containerParam) as any,
    });

    // Search by message content
    const topicsByMessages = await DB.SearchTopicsByMessageContent({
      userId: this.userId,
      content: toNullString(`%${keywordLowerCase}%`) as any,
      column3: containerParam, // containerId check
      sessionId: toNullString(containerParam) as any,
      groupId: toNullString(containerParam) as any,
    });

    // If no message results, return title results
    if (topicsByMessages.length === 0) {
      return mapTopicsToChatTopics(topicsByTitle);
    }

    // Merge and deduplicate
    const allTopics = [...topicsByTitle];
    const existingIds = new Set(topicsByTitle.map((t) => t.id));

    for (const topic of topicsByMessages) {
      if (!existingIds.has(topic.id)) {
        allTopics.push(topic);
      }
    }

    // Sort by updated_at
    return mapTopicsToChatTopics(allTopics)
      .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
  };

  count = async (params?: {
    endDate?: string;
    range?: [string, string];
    startDate?: string;
  }): Promise<number> => {
    if (!params) {
      const result = await DB.CountTopics(this.userId);
      return Number(result) || 0;
    }

    let startTime: number;
    let endTime: number;

    if (params.range) {
      const [start, end] = params.range;
      startTime = new Date(start).getTime();
      endTime = new Date(end).getTime();
    } else {
      startTime = params.startDate ? new Date(params.startDate).getTime() : 0;
      endTime = params.endDate ? new Date(params.endDate).getTime() : Date.now();
    }

    const result = await DB.CountTopicsByDateRange({
      userId: this.userId,
      createdAt: startTime,
      createdAt2: endTime,
    });
    return Number(result) || 0;
  };

  rank = async (limit: number = 10): Promise<TopicRankItem[]> => {
    const result = await DB.RankTopics({
      userId: this.userId,
      limit,
    });

    return result.map((r) => ({
      id: r.id,
      title: getNullableString(r.title as any) || null,
      sessionId: getNullableString(r.sessionId as any) || null,
      count: Number(r.count) || 0,
    }));
  };

  // **************** Create *************** //

  create = async (
    { messages: messageIds, ...params }: CreateTopicParams,
    id: string = this.genId(),
  ): Promise<ChatTopic> => {
    await this.logger.methodEntry('create', { id, title: params.title, messagesCount: messageIds?.length || 0, userId: this.userId });
    
    // Note: No transaction support in Wails!
    // This is a potential data consistency issue

    const now = currentTimestampMs();

    const topic: Topic = await DB.CreateTopic({
      id,
      title: toNullString(params.title as any),
      favorite: boolToInt(params.favorite || false),
      sessionId: toNullString(params.sessionId as any),
      groupId: toNullString(params.groupId as any),
      userId: this.userId,
      clientId: toNullString(''),
      historySummary: toNullString(''),
      metadata: toNullJSON(null),
      createdAt: now,
      updatedAt: now,
    });

    // Update associated messages' topicId
    if (messageIds && messageIds.length > 0) {
      await this.logger.debug(`Updating ${messageIds.length} messages with topic ID`);
      await DB.UpdateMessagesTopicId({
        topicId: toNullString(topic.id) as any,
        userId: this.userId,
        ids: messageIds,
      });
    }
    
    await this.logger.methodExit('create', { topicId: topic.id });
    return mapTopicToChatTopic(topic);
  };

  batchCreate = async (topicParams: (CreateTopicParams & { id?: string })[]) => {
    // No transaction support - create one by one
    // No batch insert support

    const createdTopics: ChatTopic[] = await Promise.all(
      topicParams.map(async (params) => {
        const now = currentTimestampMs();
        const topic: Topic = await DB.CreateTopic({
          id: params.id || this.genId(),
          title: toNullString(params.title as any),
          favorite: boolToInt(params.favorite || false),
          sessionId: toNullString(params.sessionId as any),
          groupId: toNullString(params.groupId as any),
          userId: this.userId,
          clientId: toNullString(''),
          historySummary: toNullString(''),
          metadata: toNullJSON(null),
          createdAt: now,
          updatedAt: now,
        });

        // Update messages
        if (params.messages && params.messages.length > 0) {
          await DB.UpdateMessagesTopicId({
            topicId: toNullString(topic.id) as any,
            userId: this.userId,
            ids: params.messages,
          });
        }

        return mapTopicToChatTopic(topic);
      }),
    );

    return createdTopics;
  };

  duplicate = async (topicId: string, newTitle?: string) => {
    // No transaction support!

    // Find original topic
    const originalTopic = await this.findById(topicId);
    if (!originalTopic) {
      throw new Error(`ChatTopic with id ${topicId} not found`);
    }

    // Copy topic
    const now = currentTimestampMs();
    const duplicatedTopic = await DB.CreateTopic({
      id: this.genId(),
      title: toNullString(newTitle || originalTopic.title as any),
      favorite: originalTopic.favorite || 0,
      sessionId: toNullString(originalTopic.sessionId as any),
      groupId: toNullString(originalTopic.groupId as any),
      userId: this.userId,
      clientId: toNullString(''),
      historySummary: toNullString(originalTopic.historySummary as any),
      metadata: toNullJSON(parseNullableJSON(originalTopic.metadata as any)),
      createdAt: now,
      updatedAt: now,
    });

    // Get original messages
    const originalMessages = await DB.GetMessagesByTopicId({
      topicId: toNullString(topicId) as any,
      userId: this.userId,
    });

    // Copy messages
    const duplicatedMessages = await Promise.all(
      originalMessages.map(async (message) => {
        const msgNow = currentTimestampMs();
        return await DB.CreateMessage({
          id: nanoid(14),
          role: message.role,
          content: message.content,
          reasoning: message.reasoning,
          search: message.search,
          metadata: message.metadata,
          model: message.model,
          provider: message.provider,
          favorite: message.favorite,
          error: message.error,
          tools: message.tools,
          traceId: message.traceId,
          observationId: message.observationId,
          clientId: toNullString(''),
          userId: this.userId,
          sessionId: message.sessionId,
          topicId: toNullString(duplicatedTopic.id),
          threadId: message.threadId,
          parentId: message.parentId,
          quotaId: message.quotaId,
          agentId: message.agentId,
          groupId: message.groupId,
          targetId: message.targetId,
          messageGroupId: message.messageGroupId,
          createdAt: msgNow,
          updatedAt: msgNow,
        }) as unknown as DBMessageItem;
      }),
    );

    return {
      topic: duplicatedTopic,
      messages: duplicatedMessages,
    };
  };

  // **************** Delete *************** //

  delete = async (id: string): Promise<void> => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    await DB.DeleteTopic({
      id,
      userId: this.userId,
    });
    
    await this.logger.methodExit('delete', { id });
  };

  batchDeleteBySessionId = async (sessionId?: string | null) => {
    if (!sessionId) return;

    await DB.DeleteTopicsBySession({
      sessionId: toNullString(sessionId),
      userId: this.userId,
    });
  };

  batchDeleteByGroupId = async (groupId?: string | null) => {
    if (!groupId) return;

    await DB.DeleteTopicsByGroup({
      groupId: toNullString(groupId),
      userId: this.userId,
    });
  };

  batchDelete = async (ids: string[]): Promise<void> => {
    await DB.BatchDeleteTopics({
      userId: this.userId,
      ids,
    });
  };

  deleteAll = async (): Promise<void> => {
    await DB.DeleteAllTopics(this.userId);
  };

  // **************** Update *************** //

  update = async (id: string, data: Partial<ChatTopic>): Promise<ChatTopic> => {
    await this.logger.methodEntry('update', { id, data, userId: this.userId });
    
    const result: Topic = await DB.UpdateTopic({
      id,
      userId: this.userId,
      title: toNullString(data.title as any),
      historySummary: toNullString(data.historySummary as any),
      metadata: toNullJSON(data.metadata),
      updatedAt: currentTimestampMs(),
    });
    
    await this.logger.methodExit('update', { id });
    return mapTopicToChatTopic(result);
  };

  // **************** Helper *************** //

  private genId = () => nanoid();
}

