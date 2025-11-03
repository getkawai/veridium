import {
  ChatFileItem,
  ChatImageItem,
  ChatTTS,
  ChatToolPayload,
  ChatTranslate,
  ChatVideoItem,
  CreateMessageParams,
  CreateMessageResult,
  DBMessageItem,
  ModelRankItem,
  NewMessageQueryParams,
  QueryMessageParams,
  UIChatMessage,
  UpdateMessageParams,
  UpdateMessageRAGParams,
} from  '@/types';
import type { HeatmapsProps } from '@lobehub/charts';
import dayjs from 'dayjs';
import { nanoid } from 'nanoid';

import { merge } from '@/utils/merge';
import { today } from '@/utils/time';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
} from '@/types/database';

export class MessageModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // **************** Query *************** //
  
  /**
   * Complex query method - WARNING: This is a simplified version
   * The Drizzle version does multiple JOINs in a single query which is more efficient
   * This Wails version requires multiple queries (N+1 problem)
   * 
   * For production, consider creating a dedicated SQL query with all JOINs
   */
  query = async (
    { current = 0, pageSize = 1000, sessionId, topicId, groupId }: QueryMessageParams = {},
    options: {
      postProcessUrl?: (path: string | null, file: { fileType: string }) => Promise<string>;
    } = {},
  ) => {
    const offset = current * pageSize;

    // 1. Get basic messages
    // TODO: Add query that handles sessionId/topicId/groupId filtering
    const messages = await DB.ListMessages({
      userId: this.userId,
      sessionId: toNullString(sessionId as any),
      limit: pageSize,
      offset,
    });

    if (messages.length === 0) return [];

    const messageIds = messages.map((m) => m.id);

    // 2. Get relative files (N+1 query)
    const relatedFileList = await Promise.all(
      messageIds.map(async (messageId) => {
        const files = await DB.GetMessageFiles({
          messageId,
          userId: this.userId,
        });
        
        return files.map((file) => ({
          id: file.id,
          messageId,
          name: getNullableString(file.name as any),
          size: file.size,
          fileType: getNullableString(file.fileType as any),
          url: getNullableString(file.url as any),
        }));
      }),
    );

    const flatFileList = relatedFileList.flat();

    // Post-process URLs
    const processedFileList = await Promise.all(
      flatFileList.map(async (file) => ({
        ...file,
        url: options.postProcessUrl
          ? await options.postProcessUrl(file.url || null, file as any)
          : file.url,
      })),
    );

    // Get documents content
    const documentsMap: Record<string, string> = {};

    // TODO: Add batch query for documents
    // For now, simplified version

    const imageList = processedFileList.filter((i) => (i.fileType || '').startsWith('image'));
    const videoList = processedFileList.filter((i) => (i.fileType || '').startsWith('video'));
    const fileList = processedFileList.filter(
      (i) => !(i.fileType || '').startsWith('image') && !(i.fileType || '').startsWith('video'),
    );

    // 3. Get relative file chunks
    // TODO: Add GetMessageQueryChunks query
    const chunksList: any[] = [];

    // 4. Get relative message queries
    const messageQueriesList = await Promise.all(
      messageIds.map(async (messageId) => {
        try {
          const queries = await DB.ListMessageQueriesByMessage({
            messageId,
            userId: this.userId,
          });
          return queries.map((q) => ({
            id: q.id,
            messageId: q.messageId,
            rewriteQuery: getNullableString(q.rewriteQuery as any),
            userQuery: getNullableString(q.userQuery as any),
          }));
        } catch {
          return [];
        }
      }),
    );

    const flatMessageQueries = messageQueriesList.flat();

    // 5. Get plugins, translates, TTS for each message
    const pluginsMap = new Map();
    const translatesMap = new Map();
    const ttsMap = new Map();

    await Promise.all(
      messageIds.map(async (messageId) => {
        try {
          const plugin = await DB.GetMessagePlugin({
            id: messageId,
            userId: this.userId,
          });
          if (plugin) pluginsMap.set(messageId, plugin);
        } catch {}

        try {
          const translate = await DB.GetMessageTranslate({
            id: messageId,
            userId: this.userId,
          });
          if (translate) translatesMap.set(messageId, translate);
        } catch {}

        try {
          const tts = await DB.GetMessageTTS({
            id: messageId,
            userId: this.userId,
          });
          if (tts) ttsMap.set(messageId, tts);
        } catch {}
      }),
    );

    // Map results
    return messages.map((message) => {
      const plugin = pluginsMap.get(message.id);
      const translate = translatesMap.get(message.id);
      const tts = ttsMap.get(message.id);
      const messageQuery = flatMessageQueries.find((q) => q.messageId === message.id);

      return {
        id: message.id,
        role: message.role,
        content: getNullableString(message.content as any),
        reasoning: getNullableString(message.reasoning as any),
        search: getNullableString(message.search as any),
        metadata: parseNullableJSON(message.metadata as any),
        error: parseNullableJSON(message.error as any),
        createdAt: message.createdAt,
        updatedAt: message.updatedAt,
        topicId: getNullableString(message.topicId as any),
        parentId: getNullableString(message.parentId as any),
        threadId: getNullableString(message.threadId as any),
        groupId: getNullableString(message.groupId as any),
        agentId: getNullableString(message.agentId as any),
        targetId: getNullableString(message.targetId as any),
        tools: parseNullableJSON(message.tools as any),
        tool_call_id: plugin ? getNullableString(plugin.toolCallId as any) : undefined,
        plugin: plugin
          ? {
              apiName: getNullableString(plugin.apiName as any),
              arguments: getNullableString(plugin.arguments as any),
              identifier: getNullableString(plugin.identifier as any),
              type: getNullableString(plugin.type as any),
            }
          : undefined,
        pluginError: plugin ? parseNullableJSON(plugin.error as any) : undefined,
        pluginState: plugin ? parseNullableJSON(plugin.state as any) : undefined,
        translate: translate
          ? {
              content: getNullableString(translate.content as any),
              from: getNullableString(translate.from as any),
              to: getNullableString(translate.to as any),
            }
          : undefined,
        extra: {
          fromModel: getNullableString(message.model as any),
          fromProvider: getNullableString(message.provider as any),
          translate: translate
            ? {
                content: getNullableString(translate.content as any),
                from: getNullableString(translate.from as any),
                to: getNullableString(translate.to as any),
              }
            : undefined,
          tts: tts
            ? {
                contentMd5: getNullableString(tts.contentMd5 as any),
                file: getNullableString(tts.fileId as any),
                voice: getNullableString(tts.voice as any),
              }
            : undefined,
        },
        chunksList: chunksList
          .filter((c: any) => c.messageId === message.id)
          .map((c: any) => ({
            ...c,
            similarity: Number(c.similarity) ?? undefined,
          })),
        fileList: fileList
          .filter((f) => f.messageId === message.id)
          .map<ChatFileItem>((f) => ({
            content: documentsMap[f.id],
            fileType: f.fileType!,
            id: f.id,
            name: f.name!,
            size: f.size!,
            url: f.url || '',
          })),
        imageList: imageList
          .filter((f) => f.messageId === message.id)
          .map<ChatImageItem>((f) => ({ alt: f.name!, id: f.id, url: f.url || '' })),
        videoList: videoList
          .filter((f) => f.messageId === message.id)
          .map<ChatVideoItem>((f) => ({ alt: f.name!, id: f.id, url: f.url || '' })),
        meta: {},
        ragQuery: messageQuery?.rewriteQuery,
        ragQueryId: messageQuery?.id,
        ragRawQuery: messageQuery?.userQuery,
      } as unknown as UIChatMessage;
    });
  };

  findById = async (id: string) => {
    try {
      return await DB.GetMessage({
        id,
        userId: this.userId,
      });
    } catch {
      return undefined;
    }
  };

  findMessageQueriesById = async (messageId: string) => {
    try {
      const queries = await DB.ListMessageQueriesByMessage({
        messageId,
        userId: this.userId,
      });

      if (queries.length === 0) return undefined;

      return {
        id: queries[0].id,
        query: getNullableString(queries[0].rewriteQuery as any),
        rewriteQuery: getNullableString(queries[0].rewriteQuery as any),
        userQuery: getNullableString(queries[0].userQuery as any),
        embeddings: null, // TODO: Join with embeddings table
      };
    } catch {
      return undefined;
    }
  };

  queryAll = async () => {
    // TODO: Add ListAllMessages query
    const messages = await DB.ListMessages({
      userId: this.userId,
      sessionId: toNullString(''),
      limit: 10000,
      offset: 0,
    });

    return messages as unknown as DBMessageItem[];
  };

  queryBySessionId = async (sessionId?: string | null) => {
    const messages = await DB.ListMessages({
      userId: this.userId,
      sessionId: toNullString(sessionId as any),
      limit: 10000,
      offset: 0,
    });

    return messages as unknown as DBMessageItem[];
  };

  queryByKeyword = async (keyword: string) => {
    if (!keyword) return [];

    const messages = await DB.SearchMessagesByKeyword({
      userId: this.userId,
      content: toNullString(`%${keyword}%`),
      limit: 1000,
    });

    return messages as unknown as DBMessageItem[];
  };

  count = async (params?: {
    endDate?: string;
    range?: [string, string];
    startDate?: string;
  }): Promise<number> => {
    if (!params) {
      const result = await DB.CountMessages(this.userId);
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

    const result = await DB.CountMessagesByDateRange({
      userId: this.userId,
      createdAt: startTime,
      createdAt2: endTime,
    });
    return Number(result) || 0;
  };

  countWords = async (params?: {
    endDate?: string;
    range?: [string, string];
    startDate?: string;
  }): Promise<number> => {
    if (!params) {
      const result = await DB.CountMessageWords(this.userId);
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

    const result = await DB.CountMessageWordsByDateRange({
      userId: this.userId,
      createdAt: startTime,
      createdAt2: endTime,
    });
    return Number(result) || 0;
  };

  rankModels = async (limit: number = 10): Promise<ModelRankItem[]> => {
    const result = await DB.RankModels({
      userId: this.userId,
      limit,
    });

    return result.map((r) => ({
      id: getNullableString(r.id as any) || '',
      count: Number(r.count) || 0,
    }));
  };

  getHeatmaps = async (): Promise<HeatmapsProps['data']> => {
    const startDate = today().subtract(1, 'year').startOf('day');
    const endDate = today().endOf('day');

    // TODO: Add GetMessageHeatmaps query with GROUP BY date
    // For now, simplified version - get all messages and group in memory
    const messages = await DB.ListMessages({
      userId: this.userId,
      sessionId: toNullString(''),
      limit: 100000,
      offset: 0,
    });

    const dateCountMap = new Map<string, number>();
    
    for (const message of messages) {
      const date = dayjs(message.createdAt).format('YYYY-MM-DD');
      dateCountMap.set(date, (dateCountMap.get(date) || 0) + 1);
    }

    const heatmapData: HeatmapsProps['data'] = [];
    let currentDate = startDate.clone();

    while (currentDate.isBefore(endDate) || currentDate.isSame(endDate, 'day')) {
      const formattedDate = currentDate.format('YYYY-MM-DD');
      const count = dateCountMap.get(formattedDate) || 0;

      const levelCount = count > 0 ? Math.ceil(count / 5) : 0;
      const level = levelCount > 4 ? 4 : levelCount;

      heatmapData.push({
        count,
        date: formattedDate,
        level,
      });

      currentDate = currentDate.add(1, 'day');
    }

    return heatmapData;
  };

  hasMoreThanN = async (n: number): Promise<boolean> => {
    const messages = await DB.ListMessages({
      userId: this.userId,
      sessionId: toNullString(''),
      limit: n + 1,
      offset: 0,
    });

    return messages.length > n;
  };

  // **************** Create *************** //

  create = async (
    {
      fromModel,
      fromProvider,
      files,
      plugin,
      pluginState,
      fileChunks,
      ragQueryId,
      updatedAt,
      createdAt,
      ...message
    }: CreateMessageParams,
    id: string = nanoid(14),
  ): Promise<DBMessageItem> => {
    // Note: No transaction support in Wails!
    // This is a potential data consistency issue

    const normalizedMessage = message.groupId ? { ...message, sessionId: null } : message;
    const now = currentTimestampMs();

    const item = await DB.CreateMessage({
      id,
      role: normalizedMessage.role,
      content: toNullString(normalizedMessage.content as any),
      reasoning: toNullString((normalizedMessage as any).reasoning as any),
      search: toNullString((normalizedMessage as any).search as any),
      metadata: toNullJSON((normalizedMessage as any).metadata),
      model: toNullString(fromModel as any),
      provider: toNullString(fromProvider as any),
      favorite: boolToInt(false),
      error: toNullJSON(normalizedMessage.error),
      tools: toNullJSON((normalizedMessage as any).tools),
      traceId: toNullString(normalizedMessage.traceId as any),
      observationId: toNullString((normalizedMessage as any).observationId as any),
      clientId: toNullString((normalizedMessage as any).clientId as any),
      userId: this.userId,
      sessionId: toNullString(normalizedMessage.sessionId as any),
      topicId: toNullString(normalizedMessage.topicId as any),
      threadId: toNullString(normalizedMessage.threadId as any),
      parentId: toNullString((normalizedMessage as any).parentId as any),
      quotaId: toNullString((normalizedMessage as any).quotaId as any),
      agentId: toNullString((normalizedMessage as any).agentId as any),
      groupId: toNullString(normalizedMessage.groupId as any),
      targetId: toNullString(normalizedMessage.targetId as any),
      messageGroupId: toNullString((normalizedMessage as any).messageGroupId as any),
      createdAt: createdAt || now,
      updatedAt: updatedAt || now,
    });

    // Insert plugin data if message is a tool
    if (message.role === 'tool') {
      await DB.CreateMessagePlugin({
        id,
        toolCallId: toNullString(message.tool_call_id as any),
        type: toNullString(plugin?.type as any),
        apiName: toNullString(plugin?.apiName as any),
        arguments: toNullString(plugin?.arguments as any),
        identifier: toNullString(plugin?.identifier as any),
        state: toNullJSON(pluginState),
        error: toNullJSON(null),
        clientId: toNullString(null),
        userId: this.userId,
      });
    }

    // Link files
    if (files && files.length > 0) {
      await Promise.all(
        files.map((fileId) =>
          DB.LinkMessageToFile({
            fileId,
            messageId: id,
            userId: this.userId,
          }),
        ),
      );
    }

    // Link file chunks
    if (fileChunks && fileChunks.length > 0 && ragQueryId) {
      await Promise.all(
        fileChunks.map((chunk) =>
          DB.LinkMessageQueryToChunk({
            messageId: toNullString(id),
            queryId: toNullString(ragQueryId),
            chunkId: toNullString(chunk.id),
            similarity: { Int64: chunk.similarity || 0, Valid: !!chunk.similarity } as any,
            userId: this.userId,
          }),
        ),
      );
    }

    return item as unknown as DBMessageItem;
  };

  createNewMessage = async (
    params: CreateMessageParams,
    options: {
      postProcessUrl?: (path: string | null, file: { fileType: string }) => Promise<string>;
    } = {},
  ): Promise<CreateMessageResult> => {
    const item = await this.create(params);

    const messages = await this.query(
      {
        current: 0,
        groupId: params.groupId,
        pageSize: 9999,
        sessionId: params.sessionId,
        topicId: params.topicId,
      },
      options,
    );

    return {
      id: item.id,
      messages,
    };
  };

  batchCreate = async (newMessages: DBMessageItem[]) => {
    // No batch insert support - create one by one
    await Promise.all(
      newMessages.map((message) =>
        this.create(message as any, message.id),
      ),
    );
  };

  createMessageQuery = async (params: NewMessageQueryParams) => {
    return await DB.CreateMessageQuery({
      id: nanoid(),
      messageId: params.messageId,
      rewriteQuery: toNullString(params.rewriteQuery as any),
      userQuery: toNullString(params.userQuery as any),
      clientId: toNullString(''),
      userId: this.userId,
      embeddingsId: toNullString(params.embeddingsId as any),
    });
  };

  // **************** Update *************** //

  update = async (id: string, { imageList, ...message }: Partial<UpdateMessageParams>) => {
    // No transaction support!
    
    if (imageList && imageList.length > 0) {
      await Promise.all(
        imageList.map((file) =>
          DB.LinkMessageToFile({
            fileId: file.id,
            messageId: id,
            userId: this.userId,
          }),
        ),
      );
    }

    return await DB.UpdateMessage({
      id,
      userId: this.userId,
      content: toNullString(message.content as any),
      reasoning: toNullString(message.reasoning as any),
      metadata: toNullJSON(message.metadata),
      favorite: boolToInt(false), // favorite not in UpdateMessageParams
      updatedAt: currentTimestampMs(),
    });
  };

  updateMetadata = async (id: string, metadata: Record<string, any>) => {
    const item = await this.findById(id);
    if (!item) return;

    const currentMetadata = parseNullableJSON(item.metadata as any) || {};
    const mergedMetadata = merge(currentMetadata, metadata);

    return await DB.UpdateMessage({
      id,
      userId: this.userId,
      content: toNullString(item.content as any),
      reasoning: toNullString(item.reasoning as any),
      metadata: toNullJSON(mergedMetadata),
      favorite: item.favorite,
      updatedAt: currentTimestampMs(),
    });
  };

  updatePluginState = async (id: string, state: Record<string, any>) => {
    const item = await DB.GetMessagePlugin({
      id,
      userId: this.userId,
    });
    
    if (!item) throw new Error('Plugin not found');

    const currentState = parseNullableJSON(item.state as any) || {};
    const mergedState = merge(currentState, state);

    return await DB.UpdateMessagePlugin({
      id,
      userId: this.userId,
      state: toNullJSON(mergedState),
      error: item.error,
    });
  };

  updateMessagePlugin = async (id: string, value: { state?: any; error?: any }) => {
    const item = await DB.GetMessagePlugin({
      id,
      userId: this.userId,
    });
    
    if (!item) throw new Error('Plugin not found');

    return await DB.UpdateMessagePlugin({
      id,
      userId: this.userId,
      state: value.state !== undefined ? toNullJSON(value.state) : item.state,
      error: value.error !== undefined ? toNullJSON(value.error) : item.error,
    });
  };

  updateTranslate = async (id: string, translate: Partial<ChatTranslate>) => {
    return await DB.UpsertMessageTranslate({
      id,
      content: toNullString(translate.content as any),
      from: toNullString(translate.from as any),
      to: toNullString(translate.to as any),
      clientId: toNullString(null),
      userId: this.userId,
    });
  };

  updateTTS = async (id: string, tts: Partial<ChatTTS>) => {
    return await DB.UpsertMessageTTS({
      id,
      contentMd5: toNullString(tts.contentMd5 as any),
      fileId: toNullString(tts.file as any),
      voice: toNullString(tts.voice as any),
      clientId: toNullString(null),
      userId: this.userId,
    });
  };

  async updateMessageRAG(id: string, { ragQueryId, fileChunks }: UpdateMessageRAGParams) {
    await Promise.all(
      fileChunks.map((chunk) =>
        DB.LinkMessageQueryToChunk({
          messageId: toNullString(id),
          queryId: toNullString(ragQueryId),
          chunkId: toNullString(chunk.id),
          similarity: { Int64: chunk.similarity || 0, Valid: !!chunk.similarity } as any,
          userId: this.userId,
        }),
      ),
    );
  }

  // **************** Delete *************** //

  deleteMessage = async (id: string) => {
    // No transaction support!
    // This is a potential data consistency issue

    const message = await this.findById(id);
    if (!message) return;

    const tools = parseNullableJSON(message.tools as any) as ChatToolPayload[] | null;
    const toolCallIds = tools?.map((tool) => tool.id).filter(Boolean) || [];

    let relatedMessageIds: string[] = [];

    if (toolCallIds.length > 0) {
      // TODO: Add query to get messages by tool_call_id
      // For now, simplified
    }

    const messageIdsToDelete = [id, ...relatedMessageIds];

    await DB.BatchDeleteMessages({
      userId: this.userId,
      ids: messageIdsToDelete,
    });
  };

  deleteMessages = async (ids: string[]) => {
    await DB.BatchDeleteMessages({
      userId: this.userId,
      ids,
    });
  };

  deleteMessageTranslate = async (id: string) => {
    await DB.DeleteMessageTranslate({
      id,
      userId: this.userId,
    });
  };

  deleteMessageTTS = async (id: string) => {
    await DB.DeleteMessageTTS({
      id,
      userId: this.userId,
    });
  };

  deleteMessageQuery = async (id: string) => {
    await DB.DeleteMessageQuery({
      id,
      userId: this.userId,
    });
  };

  deleteMessagesBySession = async (
    sessionId?: string | null,
    topicId?: string | null,
    groupId?: string | null,
  ) => {
    if (sessionId) {
      await DB.DeleteMessagesBySession({
        sessionId: toNullString(sessionId),
        userId: this.userId,
      });
    } else if (topicId) {
      await DB.DeleteMessagesByTopic({
        topicId: toNullString(topicId),
        userId: this.userId,
      });
    }
    // TODO: Add DeleteMessagesByGroup query for groupId
  };

  deleteAllMessages = async () => {
    return await DB.DeleteAllMessages(this.userId);
  };
}

