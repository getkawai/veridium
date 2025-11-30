import {
  ChatFileItem,
  ChatImageItem,
  ChatToolPayload,
  ChatTranslate,
  ChatTTS,
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
} from '@/types';
import type { HeatmapsProps } from '@lobehub/charts';
import dayjs from 'dayjs';
import { nanoid } from 'nanoid';

import { merge } from '@/utils/merge';
import { today } from '@/utils/time';
import { createModelLogger } from '@/utils/logger';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
} from '@/types/database';

// Import transaction methods (Optimization 4A)
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';

export class MessageModel {
  private userId: string;
  private logger = createModelLogger('Message', 'MessageModel', 'database/models/message');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // **************** Query *************** //

  /**
   * OPTIMIZED: Uses server-side filtering (2A) and JOIN queries (3A)
   * Much faster than previous client-side filtering approach
   */
  query = async (
    { current = 0, pageSize = 1000, sessionId, topicId, groupId }: QueryMessageParams = {},
    options: {
      postProcessUrl?: (path: string | null, file: { fileType: string }) => Promise<string>;
    } = {},
  ) => {
    await this.logger.methodEntry('query', {
      current,
      pageSize,
      sessionId,
      topicId,
      groupId,
      userId: this.userId
    });

    const offset = current * pageSize;

    // OPTIMIZATION 2A: Server-side filtering
    // Use specific query based on filter type
    let messages;
    if (sessionId !== undefined && sessionId !== null) {
      await this.logger.debug(`Querying messages by session: ${sessionId}`);
      messages = await DB.ListMessagesBySession({
        userId: this.userId,
        sessionId: toNullString(sessionId) as any,
        limit: pageSize,
        offset,
      });
    } else if (topicId !== undefined && topicId !== null) {
      await this.logger.debug(`Querying messages by topic: ${topicId}`);
      messages = await DB.ListMessagesByTopic({
        userId: this.userId,
        topicId: toNullString(topicId) as any,
        limit: pageSize,
        offset,
      });
    } else if (groupId !== undefined && groupId !== null) {
      await this.logger.debug(`Querying messages by group: ${groupId}`);
      messages = await DB.ListMessagesByGroup({
        userId: this.userId,
        groupId: toNullString(groupId) as any,
        limit: pageSize,
        offset,
      });
    } else {
      await this.logger.debug('Querying all messages');
      messages = await DB.ListMessages({
        userId: this.userId,
        limit: pageSize,
        offset,
      });
    }

    await this.logger.debug(`Retrieved ${messages.length} messages from DB`);

    if (messages.length === 0) {
      await this.logger.methodExit('query', { count: 0 });
      return [];
    }

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
    const fileIds = processedFileList.map((file) => file.id).filter(Boolean);
    let documentsMap: Record<string, string> = {};

    // Note: GetDocumentsByFileIds doesn't support IN clause with sqlc.slice
    // Fetching documents one by one (N+1 query)
    if (fileIds.length > 0) {
      await Promise.all(
        fileIds.map(async (fileId) => {
          try {
            const doc = await DB.GetDocumentByFileId({
              fileId: toNullString(fileId),
              userId: this.userId,
            });
            if (doc && doc.fileId) {
              const fileIdStr = getNullableString(doc.fileId as any);
              if (fileIdStr) {
                documentsMap[fileIdStr] = getNullableString(doc.content as any) || '';
              }
            }
          } catch {
            // Document not found
          }
        }),
      );
    }

    const imageList = processedFileList.filter((i) => (i.fileType || '').startsWith('image'));
    const videoList = processedFileList.filter((i) => (i.fileType || '').startsWith('video'));
    const fileList = processedFileList.filter(
      (i) => !(i.fileType || '').startsWith('image') && !(i.fileType || '').startsWith('video'),
    );

    // 3. Get relative file chunks
    // Note: GetMessageQueryChunks expects messageIds as NullString[]
    let chunksList: any[] = [];
    await Promise.all(
      messageIds.map(async (msgId) => {
        try {
          const chunks = await DB.GetMessageQueryChunks({
            messageIds: [toNullString(msgId)],
            userId: this.userId,
          });
          chunksList.push(...chunks.map((c: any) => ({ ...c, messageId: msgId })));
        } catch {
          // Chunks query might fail
        }
      }),
    );

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
        } catch { }

        try {
          const translate = await DB.GetMessageTranslate({
            id: messageId,
            userId: this.userId,
          });
          if (translate) translatesMap.set(messageId, translate);
        } catch { }

        try {
          const tts = await DB.GetMessageTTS({
            id: messageId,
            userId: this.userId,
          });
          if (tts) ttsMap.set(messageId, tts);
        } catch { }
      }),
    );

    // Map results
    const result = messages.map((message) => {
      const plugin = pluginsMap.get(message.id);
      const translate = translatesMap.get(message.id);
      const tts = ttsMap.get(message.id);
      const messageQuery = flatMessageQueries.find((q) => q.messageId === message.id);

      return {
        id: message.id,
        role: message.role,
        content: getNullableString(message.content as any),
        reasoning: parseNullableJSON(message.reasoning as any),
        search: parseNullableJSON(message.search as any),
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
              content: getNullableString(translate.content as any) || '',
              from: getNullableString(translate.from as any) || '',
              to: getNullableString(translate.to as any) || '',
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
            id: getNullableString(c.id as any) || c.id,
            fileId: getNullableString(c.fileId as any) || getNullableString(c.fileID as any) || c.fileId,
            filename: getNullableString(c.filename as any) || getNullableString(c.Filename as any) || c.filename,
            fileType: getNullableString(c.fileType as any) || getNullableString(c.FileType as any) || c.fileType,
            fileUrl: getNullableString(c.fileUrl as any) || getNullableString(c.FileUrl as any) || c.fileUrl,
            text: getNullableString(c.text as any) || c.text,
            similarity: typeof c.similarity === 'number' 
              ? c.similarity 
              : (c.similarity?.Int64 ?? c.similarity?.Valid ? Number(c.similarity.Int64) : undefined),
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

    await this.logger.methodExit('query', { count: result.length });
    return result;
  };

  findById = async (id: string) => {
    await this.logger.methodEntry('findById', { id, userId: this.userId });

    try {
      const message = await DB.GetMessage({
        id,
        userId: this.userId,
      });
      await this.logger.methodExit('findById', { found: !!message });
      return message;
    } catch (error) {
      await this.logger.methodError('findById', error, { id });
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

      const query = queries[0];
      let embeddings: any = null;

      // Get embeddings if embeddingsId exists
      const embeddingsIdStr = getNullableString(query.embeddingsId as any);
      if (embeddingsIdStr) {
        try {
          const emb = await DB.GetEmbeddingsItem({
            id: embeddingsIdStr,
            userId: toNullString(this.userId),
          });
          if (emb && emb.embeddings) {
            embeddings = emb.embeddings;
          }
        } catch {
          // Embeddings not found
        }
      }

      return {
        id: query.id,
        query: getNullableString(query.rewriteQuery as any),
        rewriteQuery: getNullableString(query.rewriteQuery as any),
        userQuery: getNullableString(query.userQuery as any),
        embeddings,
      };
    } catch {
      return undefined;
    }
  };

  queryAll = async () => {
    const messages = await DB.ListMessages({
      userId: this.userId,
      limit: 10000,
      offset: 0,
    });

    return messages as unknown as DBMessageItem[];
  };

  queryBySessionId = async (sessionId?: string | null) => {
    const allMessages = await DB.ListMessages({
      userId: this.userId,
      limit: 10000,
      offset: 0,
    });

    // Filter by sessionId
    const messages = allMessages.filter((m) => {
      const msgSessionId = getNullableString(m.sessionId as any);
      return sessionId === null || sessionId === undefined
        ? !msgSessionId
        : msgSessionId === sessionId;
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
    await this.logger.methodEntry('count', { params, userId: this.userId });

    if (!params) {
      const result = await DB.CountMessages(this.userId);
      const count = Number(result) || 0;
      await this.logger.methodExit('count', { count });
      return count;
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

    await this.logger.debug(`Counting messages in date range: ${new Date(startTime).toISOString()} to ${new Date(endTime).toISOString()}`);

    const result = await DB.CountMessagesByDateRange({
      userId: this.userId,
      createdAt: startTime,
      createdAt2: endTime,
    });
    const count = Number(result) || 0;
    await this.logger.methodExit('count', { count, dateRange: true });
    return count;
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

    // Use optimized query with GROUP BY date
    const result = await DB.GetMessageHeatmaps({
      userId: this.userId,
      createdAt: startDate.valueOf(),
      createdAt2: endDate.valueOf(),
    });

    const dateCountMap = new Map<string, number>();

    for (const item of result) {
      if (item.date) {
        const dateStr = dayjs(item.date as any).format('YYYY-MM-DD');
        dateCountMap.set(dateStr, Number(item.count) || 0);
      }
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
      limit: n + 1,
      offset: 0,
    });

    return messages.length > n;
  };

  // **************** Create *************** //

  /**
   * OPTIMIZED: Uses atomic transaction (4A)
   * All operations succeed or fail together - no partial writes!
   */
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
    await this.logger.methodEntry('create', {
      id,
      role: message.role,
      hasPlugin: !!plugin,
      filesCount: files?.length || 0,
      userId: this.userId
    });

    const normalizedMessage = message.groupId ? { ...message, sessionId: null } : message;
    const now = currentTimestampMs();

    // OPTIMIZATION 4A: Use atomic transaction
    // Ensure fromModel and fromProvider are plain strings (not already wrapped NullStrings)
    const modelStr = typeof fromModel === 'string' ? fromModel : (fromModel as any)?.String || undefined;
    const providerStr = typeof fromProvider === 'string' ? fromProvider : (fromProvider as any)?.String || undefined;

    const item = await DBService.CreateMessageWithRelations({
      Message: {
        id,
        role: normalizedMessage.role,
        content: toNullString(normalizedMessage.content as any),
        reasoning: toNullJSON((normalizedMessage as any).reasoning),  // reasoning is object, use toNullJSON
        search: toNullJSON((normalizedMessage as any).search),  // search is object, use toNullJSON
        metadata: toNullJSON((normalizedMessage as any).metadata),
        model: toNullString(modelStr as any),
        provider: toNullString(providerStr as any),
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
      },
      Plugin: message.role === 'tool' ? {
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
      } : null,
      FileIds: files || [],
      FileChunks: (fileChunks && ragQueryId) ? fileChunks.map((chunk) => ({
        ChunkId: chunk.id,
        QueryId: ragQueryId,
        Similarity: { Int64: chunk.similarity || 0, Valid: !!chunk.similarity } as any,
      })) : [],
    }, this.userId);

    await this.logger.methodExit('create', { messageId: item.id });
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
    await this.logger.methodEntry('batchCreate', {
      count: newMessages.length,
      userId: this.userId
    });

    // No batch insert support - create one by one
    await Promise.all(
      newMessages.map((message) =>
        this.create(message as any, message.id),
      ),
    );

    await this.logger.methodExit('batchCreate', { createdCount: newMessages.length });
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

  /**
   * OPTIMIZED: Uses atomic transaction (4A) when updating with images
   */
  update = async (id: string, { imageList, ...message }: Partial<UpdateMessageParams>) => {
    await this.logger.methodEntry('update', {
      id,
      hasImages: !!(imageList && imageList.length > 0),
      imageCount: imageList?.length || 0,
      userId: this.userId
    });

    const updateParams = {
      id,
      userId: this.userId,
      content: toNullString(message.content as any),
      reasoning: toNullJSON(message.reasoning),  // reasoning is object, use toNullJSON
      metadata: toNullJSON(message.metadata),
      favorite: boolToInt(false), // favorite not in UpdateMessageParams
      updatedAt: currentTimestampMs(),
    };

    // OPTIMIZATION 4A: Use transaction if updating with images
    if (imageList && imageList.length > 0) {
      await this.logger.debug(`Updating message with ${imageList.length} images`);
      const result = await DBService.UpdateMessageWithImages({
        MessageId: id,
        Message: updateParams,
        ImageIds: imageList.map((file) => file.id),
      }, this.userId);
      await this.logger.methodExit('update', { id, withImages: true });
      return result;
    }

    // Simple update without images
    await this.logger.debug('Updating message without images');
    const result = await DB.UpdateMessage(updateParams);
    await this.logger.methodExit('update', { id, withImages: false });
    return result;
  };

  updateMetadata = async (id: string, metadata: Record<string, any>) => {
    await this.logger.methodEntry('updateMetadata', {
      id,
      metadataKeys: Object.keys(metadata),
      userId: this.userId
    });

    const item = await this.findById(id);
    if (!item) {
      await this.logger.warn('Message not found for metadata update', { id });
      return;
    }

    const currentMetadata = parseNullableJSON(item.metadata as any) || {};
    const mergedMetadata = merge(currentMetadata, metadata);

    await this.logger.debug('Merging metadata', {
      currentKeys: Object.keys(currentMetadata).length,
      newKeys: Object.keys(metadata).length,
      mergedKeys: Object.keys(mergedMetadata).length
    });

    const result = await DB.UpdateMessage({
      id,
      userId: this.userId,
      content: toNullString(item.content as any),
      reasoning: toNullJSON(item.reasoning),  // reasoning is object, use toNullJSON
      metadata: toNullJSON(mergedMetadata),
      favorite: item.favorite,
      updatedAt: currentTimestampMs(),
    });

    await this.logger.methodExit('updateMetadata', { id });
    return result;
  };

  updatePluginState = async (id: string, state: Record<string, any>) => {
    await this.logger.methodEntry('updatePluginState', {
      id,
      stateKeys: Object.keys(state),
      userId: this.userId
    });

    const item = await DB.GetMessagePlugin({
      id,
      userId: this.userId,
    });

    if (!item) {
      await this.logger.error('Plugin not found', null, { id });
      throw new Error('Plugin not found');
    }

    const currentState = parseNullableJSON(item.state as any) || {};
    const mergedState = merge(currentState, state);

    await this.logger.debug('Merging plugin state', {
      currentKeys: Object.keys(currentState).length,
      newKeys: Object.keys(state).length,
      mergedKeys: Object.keys(mergedState).length
    });

    const result = await DB.UpdateMessagePlugin({
      id,
      userId: this.userId,
      state: toNullJSON(mergedState),
      error: item.error,
    });

    await this.logger.methodExit('updatePluginState', { id });
    return result;
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

  // TTS functionality moved to backend (native OS TTS)
  // This method kept for backward compatibility but no longer used
  updateTTS = async (id: string, tts: Partial<ChatTTS>) => {
    await this.logger.warn('updateTTS called but TTS now handled by backend', { id });
    return;
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

  /**
   * OPTIMIZED: Uses atomic transaction (4A) and batch query (1A)
   * Deletes message and all related tool messages atomically
   */
  deleteMessage = async (id: string) => {
    await this.logger.methodEntry('deleteMessage', { id, userId: this.userId });

    const message = await this.findById(id);
    if (!message) {
      await this.logger.warn('Message not found for deletion', { id });
      return;
    }

    const tools = parseNullableJSON(message.tools as any) as ChatToolPayload[] | null;
    const toolCallIds = tools?.map((tool) => tool.id).filter(Boolean) || [];

    await this.logger.debug(`Deleting message with ${toolCallIds.length} tool calls`);

    // OPTIMIZATION 4A: Use atomic transaction
    // OPTIMIZATION 1A: Batch query for tool call IDs
    await DBService.DeleteMessageWithRelated(
      JSON.stringify(toolCallIds),
      [id],
      this.userId
    );

    await this.logger.methodExit('deleteMessage', { id, deletedToolCalls: toolCallIds.length });
  };

  deleteMessages = async (ids: string[]) => {
    await this.logger.methodEntry('deleteMessages', {
      count: ids.length,
      userId: this.userId
    });

    await DB.BatchDeleteMessages({
      userId: this.userId,
      ids,
    });

    await this.logger.methodExit('deleteMessages', { deletedCount: ids.length });
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


  deleteMessagesBySession = async (
    sessionId?: string | null,
    topicId?: string | null,
    groupId?: string | null,
  ) => {
    await this.logger.methodEntry('deleteMessagesBySession', {
      sessionId,
      topicId,
      groupId,
      userId: this.userId
    });

    if (sessionId) {
      await this.logger.debug(`Deleting messages by session: ${sessionId}`);
      await DB.DeleteMessagesBySession({
        sessionId: toNullString(sessionId),
        userId: this.userId,
      });
    } else if (topicId) {
      await this.logger.debug(`Deleting messages by topic: ${topicId}`);
      await DB.DeleteMessagesByTopic({
        topicId: toNullString(topicId),
        userId: this.userId,
      });
    } else if (groupId) {
      await this.logger.debug(`Deleting messages by group: ${groupId}`);
      await DB.DeleteMessagesByGroup({
        groupId: toNullString(groupId),
        userId: this.userId,
      });
    }

    await this.logger.methodExit('deleteMessagesBySession', { sessionId, topicId, groupId });
  };

  deleteAllMessages = async () => {
    return await DB.DeleteAllMessages(this.userId);
  };
}

