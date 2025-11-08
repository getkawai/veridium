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
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';

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

  /**
   * Show error notification to user
   */
  private async showErrorNotification(title: string, message: string) {
    try {
      await NotificationService.SendNotification(
        new NotificationOptions({
          id: `topic-error-${Date.now()}`,
          title: `Topic Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  // **************** Query *************** //

  query = async ({ current = 0, pageSize = 9999, containerId }: QueryTopicParams = {}) => {
    try {
      const offset = current * pageSize;

      // Note: containerId can be null for inbox session (converted by toDbSessionId)
      // We need to query with the null value to match topics with session_id IS NULL
      if (containerId !== undefined) {
        // Try session first (containerId can be null for inbox)
        const sessionTopics = await DB.ListTopics({
          userId: this.userId,
          sessionId: toNullString(containerId),
          limit: pageSize,
          offset,
        });

        if (sessionTopics.length > 0) {
          return mapTopicsToChatTopics(sessionTopics);
        }

        // Try group if no session topics found
        // TODO: Add ListTopicsByGroup query
        // For now, return empty
        return [];
      }

      // If containerId is undefined, return topics with no session/group
      // TODO: Add ListTopicsWithoutContainer query
      return [];
    } catch (error) {
      await this.logger.error('Failed to query topics', { error, containerId, current, pageSize });
      throw error;
    }
  };

  findById = async (id: string) => {
    try {
      const topic = await DB.GetTopic({
        id,
        userId: this.userId,
      });
      return mapTopicToChatTopic(topic);
    } catch (error) {
      await this.logger.warn('Topic not found', { id, error });
      return undefined;
    }
  };

  queryAll = async (): Promise<ChatTopic[]> => {
    try {
      const topics = await DB.ListAllTopics(this.userId);
      return mapTopicsToChatTopics(topics);
    } catch (error) {
      await this.logger.error('Failed to query all topics', { error });
      throw error;
    }
  };

  queryByKeyword = async (keyword: string, containerId?: string | null): Promise<ChatTopic[]> => {
    try {
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
    } catch (error) {
      await this.logger.error('Failed to query topics by keyword', { error, keyword, containerId });
      throw error;
    }
  };

  count = async (params?: {
    endDate?: string;
    range?: [string, string];
    startDate?: string;
  }): Promise<number> => {
    try {
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
    } catch (error) {
      await this.logger.error('Failed to count topics', { error, params });
      throw error;
    }
  };

  rank = async (limit: number = 10): Promise<TopicRankItem[]> => {
    try {
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
    } catch (error) {
      await this.logger.error('Failed to rank topics', { error, limit });
      throw error;
    }
  };

  // **************** Create *************** //

  create = async (
    { messages: messageIds, ...params }: CreateTopicParams,
    id: string = this.genId(),
  ): Promise<ChatTopic> => {
    await this.logger.methodEntry('create', { id, title: params.title, messagesCount: messageIds?.length || 0, userId: this.userId });
    
    try {
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
    } catch (error) {
      await this.logger.error('Failed to create topic', { error, id, params });
      await this.showErrorNotification(
        'Create Failed',
        `Failed to create topic "${params.title || 'Untitled'}". Please try again.`
      );
      throw error;
    }
  };

  batchCreate = async (topicParams: (CreateTopicParams & { id?: string })[]) => {
    try {
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
    } catch (error) {
      await this.logger.error('Failed to batch create topics', { error, count: topicParams.length });
      await this.showErrorNotification(
        'Batch Create Failed',
        `Failed to create ${topicParams.length} topics. Please try again.`
      );
      throw error;
    }
  };

  duplicate = async (topicId: string, newTitle?: string) => {
    try {
      // No transaction support!

      // Find original topic
      const originalTopic = await this.findById(topicId);
      if (!originalTopic) {
        throw new Error(`Topic with id ${topicId} not found`);
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
    } catch (error) {
      await this.logger.error('Failed to duplicate topic', { error, topicId, newTitle });
      await this.showErrorNotification(
        'Duplicate Failed',
        `Failed to duplicate topic. Please try again.`
      );
      throw error;
    }
  };

  // **************** Delete *************** //

  delete = async (id: string): Promise<void> => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    try {
      await DB.DeleteTopic({
        id,
        userId: this.userId,
      });
      
      await this.logger.methodExit('delete', { id });
    } catch (error) {
      await this.logger.error('Failed to delete topic', { error, id });
      await this.showErrorNotification(
        'Delete Failed',
        `Failed to delete topic. Please try again.`
      );
      throw error;
    }
  };

  batchDeleteBySessionId = async (sessionId?: string | null) => {
    if (!sessionId) return;

    try {
      await DB.DeleteTopicsBySession({
        sessionId: toNullString(sessionId),
        userId: this.userId,
      });
    } catch (error) {
      await this.logger.error('Failed to batch delete topics by session', { error, sessionId });
      await this.showErrorNotification(
        'Batch Delete Failed',
        `Failed to delete topics for session. Please try again.`
      );
      throw error;
    }
  };

  batchDeleteByGroupId = async (groupId?: string | null) => {
    if (!groupId) return;

    try {
      await DB.DeleteTopicsByGroup({
        groupId: toNullString(groupId),
        userId: this.userId,
      });
    } catch (error) {
      await this.logger.error('Failed to batch delete topics by group', { error, groupId });
      await this.showErrorNotification(
        'Batch Delete Failed',
        `Failed to delete topics for group. Please try again.`
      );
      throw error;
    }
  };

  batchDelete = async (ids: string[]): Promise<void> => {
    try {
      await DB.BatchDeleteTopics({
        userId: this.userId,
        ids,
      });
    } catch (error) {
      await this.logger.error('Failed to batch delete topics', { error, count: ids.length });
      await this.showErrorNotification(
        'Batch Delete Failed',
        `Failed to delete ${ids.length} topics. Please try again.`
      );
      throw error;
    }
  };

  deleteAll = async (): Promise<void> => {
    try {
      await DB.DeleteAllTopics(this.userId);
    } catch (error) {
      await this.logger.error('Failed to delete all topics', { error });
      await this.showErrorNotification(
        'Delete All Failed',
        `Failed to delete all topics. Please try again.`
      );
      throw error;
    }
  };

  // **************** Update *************** //

  update = async (id: string, data: Partial<ChatTopic>): Promise<ChatTopic> => {
    await this.logger.methodEntry('update', { id, data, userId: this.userId });
    
    try {
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
    } catch (error) {
      await this.logger.error('Failed to update topic', { error, id, data });
      await this.showErrorNotification(
        'Update Failed',
        `Failed to update topic "${data.title || ''}". Please try again.`
      );
      throw error;
    }
  };

  // **************** Helper *************** //

  private genId = () => nanoid();
}

