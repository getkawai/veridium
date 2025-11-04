import { DEFAULT_AGENT_CONFIG, DEFAULT_INBOX_AVATAR, INBOX_SESSION_ID } from '@/const';
import {
  ChatSessionList,
  LobeAgentConfig,
  LobeAgentSession,
  LobeGroupSession,
  SessionRankItem,
} from '@/types';
import type { PartialDeep } from 'type-fest';
import { nanoid } from 'nanoid';

import { merge } from '@/utils/merge';
import { createModelLogger } from '@/utils/logger';

import {
  DB,
  type Session,
  type Agent,
  type CreateSessionParams,
  type CreateAgentParams,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
} from '@/types/database';

// Type aliases for compatibility
type NewSession = Partial<CreateSessionParams>;
type NewAgent = Partial<CreateAgentParams>;

export class SessionModel {
  private userId: string;
  private logger = createModelLogger('Session', 'SessionModel', 'database/models/session');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // **************** Query *************** //

  query = async ({ current = 0, pageSize = 9999 } = {}) => {
    await this.logger.methodEntry('query', { current, pageSize, userId: this.userId });
    
    const offset = current * pageSize;

    // Get sessions with agents
    const sessions = await DB.ListSessions({
      userId: this.userId,
      limit: pageSize,
      offset,
    });
    
    await this.logger.debug(`Retrieved ${sessions.length} sessions from DB`);

    // Filter out inbox session
    const filtered = sessions.filter(
      (s) => s.slug !== INBOX_SESSION_ID,
    );
    
    await this.logger.debug(`Filtered to ${filtered.length} sessions (excluding inbox)`);

    // Enrich with agents and groups
    const enriched = await Promise.all(
      filtered.map(async (session) => {
        const agents = await DB.GetSessionAgents({
          sessionId: session.id,
          userId: this.userId,
        });

        let group: any = undefined;
        if (session.groupId.Valid && session.groupId.String) {
          try {
            group = await DB.GetSessionGroup({
              id: session.groupId.String,
              userId: this.userId,
            });
          } catch {
            // Group not found
          }
        }

        return {
          ...session,
          agentsToSessions: agents.map((agent) => ({ agent })),
          group,
        };
      }),
    );
    
    await this.logger.methodExit('query', { count: enriched.length });
    return enriched;
  };

  queryWithGroups = async (): Promise<ChatSessionList> => {
    const result = await this.query();
    const groups = await DB.ListSessionGroups(this.userId);

    return {
      sessionGroups: groups as unknown as ChatSessionList['sessionGroups'],
      sessions: result.map((item) => this.mapSession(item as any)),
    };
  };

  queryByKeyword = async (keyword: string) => {
    if (!keyword) return [];

    const keywordLowerCase = keyword.toLowerCase();
    const data = await this.findSessionsByKeywords({ keyword: keywordLowerCase });

    return data.map((item) => this.mapSession(item as any));
  };

  findByIdOrSlug = async (
    idOrSlug: string,
  ): Promise<(Session & { agent: Agent }) | undefined> => {
    await this.logger.methodEntry('findByIdOrSlug', { idOrSlug, userId: this.userId });
    
    // Use single query to find by ID or slug
    let session: Session | undefined;
    
    try {
      session = await DB.GetSessionByIdOrSlug({
        id: idOrSlug,
        slug: idOrSlug,
        userId: this.userId,
      });
    } catch (error) {
      await this.logger.methodError('findByIdOrSlug', error, { idOrSlug });
      return undefined;
    }

    if (!session) {
      await this.logger.debug(`Session not found: ${idOrSlug}`);
      return undefined;
    }
    
    await this.logger.debug(`Found session: ${session.id}`);

    // Get agents
    const agents = await DB.GetSessionAgents({
      sessionId: session.id,
      userId: this.userId,
    });

    // Get group if exists
    let group: any = undefined;
    if (session.groupId.Valid && session.groupId.String) {
      try {
        group = await DB.GetSessionGroup({
          id: session.groupId.String,
          userId: this.userId,
        });
      } catch {
        // Group not found
      }
    }

    // If no agents found, return the session with a default agent structure
    // to prevent undefined errors downstream
    if (agents.length === 0) {
      console.warn(`Session ${idOrSlug} has no associated agent. Using default agent configuration.`);
      return {
        ...session,
        agent: {
          id: '', // Empty string for no actual agent ID
          slug: toNullString(undefined),
          title: toNullString(undefined),
          description: toNullString(undefined),
          tags: toNullJSON([]),
          avatar: toNullString(undefined),
          backgroundColor: toNullString(undefined),
          plugins: toNullJSON([]),
          clientId: toNullString(undefined),
          userId: this.userId,
          chatConfig: toNullJSON({}),
          fewShots: toNullJSON(undefined),
          model: toNullString(undefined),
          params: toNullJSON(undefined),
          provider: toNullString(undefined),
          systemRole: toNullString(undefined),
          tts: toNullJSON(undefined),
          virtual: 0,
          openingMessage: toNullString(undefined),
          openingQuestions: toNullJSON([]),
          createdAt: session.createdAt,
          updatedAt: session.updatedAt,
        } as Agent,
        agentsToSessions: [],
        group,
      } as any;
    }

    return {
      ...session,
      agent: agents[0],
      agentsToSessions: agents.map((agent) => ({ agent })),
      group,
    } as any;
  };

  count = async (params?: {
    endDate?: string;
    range?: [string, string];
    startDate?: string;
  }): Promise<number> => {
    // Use database query for counting
    if (!params) {
      return await DB.CountSessions(this.userId);
    }

    // Determine date range
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

    return await DB.CountSessionsByDateRange({
      userId: this.userId,
      createdAt: startTime,
      createdAt2: endTime,
    });
  };

  rank = async (limit: number = 10): Promise<SessionRankItem[]> => {
    // Get inbox count separately
    const inboxCount = await DB.CountTopicsBySession({
      sessionId: toNullString(''), // Empty for inbox (null session_id)
      userId: this.userId,
    });

    // Get ranked sessions
    const ranked = await DB.GetSessionRank({
      userId: this.userId,
      limit: inboxCount > 0 ? limit - 1 : limit,
    });

    const result = ranked.map((item) => ({
      id: item.id,
      title: getNullableString(item.title as any) || null,
      avatar: getNullableString(item.avatar as any) || null,
      backgroundColor: getNullableString(item.backgroundColor as any) || null,
      count: Number(item.topicCount) || 0,
    }));

    // Add inbox if it has topics
    if (inboxCount > 0) {
      return [
        {
          id: INBOX_SESSION_ID,
          title: 'inbox.title',
          avatar: DEFAULT_INBOX_AVATAR,
          backgroundColor: null,
          count: inboxCount,
        },
        ...result,
      ].sort((a, b) => b.count - a.count);
    }

    return result;
  };

  hasMoreThanN = async (n: number): Promise<boolean> => {
    const sessions = await DB.ListSessions({
      userId: this.userId,
      limit: n + 1,
      offset: 0,
    });

    return sessions.length > n;
  };

  // **************** Create *************** //

  create = async ({
    id = nanoid(),
    type = 'agent',
    session = {},
    config = {},
    slug,
  }: {
    config?: Partial<NewAgent>;
    id?: string;
    session?: Partial<NewSession>;
    slug?: string;
    type: 'agent' | 'group';
  }): Promise<Session> => {
    await this.logger.methodEntry('create', { id, type, slug, userId: this.userId });
    
    // Check if slug exists
    if (slug) {
      try {
        const existing = await DB.GetSessionBySlug({
          slug,
          userId: this.userId,
        });
        if (existing) {
          await this.logger.debug(`Session with slug ${slug} already exists, returning existing`);
          return existing;
        }
      } catch {
        // Doesn't exist, continue
      }
    }

    const now = currentTimestampMs();

    // Create session
    const newSession = await DB.CreateSession({
      id,
      userId: this.userId,
      slug: slug || id,
      title: toNullString(session.title as any),
      description: toNullString(session.description as any),
      avatar: toNullString(session.avatar as any),
      backgroundColor: toNullString(session.backgroundColor as any),
      type: toNullString(type),
      groupId: toNullString(session.groupId as any),
      clientId: toNullString(session.clientId as any),
      pinned: boolToInt(false),
      createdAt: now,
      updatedAt: now,
    });
    
    await this.logger.debug(`Created session: ${newSession.id}`);

    // If agent type, create agent and link
    if (type === 'agent') {
      const agentId = nanoid();
      
      await DB.CreateAgent({
        id: agentId,
        userId: this.userId,
        slug: toNullString(undefined),
        title: toNullString(config.title as any),
        description: toNullString(config.description as any),
        tags: toNullJSON(config.tags || []),
        avatar: toNullString(config.avatar as any),
        backgroundColor: toNullString(config.backgroundColor as any),
        plugins: toNullJSON(config.plugins || []),
        clientId: toNullString(config.clientId as any),
        chatConfig: toNullJSON(config.chatConfig),
        fewShots: toNullJSON(config.fewShots),
        model: toNullString(config.model as any),
        params: toNullJSON(config.params),
        provider: toNullString(config.provider as any),
        systemRole: toNullString(config.systemRole as any),
        tts: toNullJSON(config.tts),
        virtual: boolToInt(false),
        openingMessage: toNullString(config.openingMessage as any),
        openingQuestions: toNullJSON(config.openingQuestions || []),
        createdAt: now,
        updatedAt: now,
      });
      
      await this.logger.debug(`Created agent: ${agentId}`);

      // Link agent to session
      await DB.LinkAgentToSession({
        agentId,
        sessionId: id,
        userId: this.userId,
      });
      
      await this.logger.debug(`Linked agent ${agentId} to session ${id}`);
    }
    
    await this.logger.methodExit('create', { sessionId: newSession.id, type });
    return newSession;
  };

  createInbox = async (defaultAgentConfig: PartialDeep<LobeAgentConfig>) => {
    try {
      const existing = await DB.GetSessionBySlug({
        slug: INBOX_SESSION_ID,
        userId: this.userId,
      });
      if (existing) return existing;
    } catch {
      // Doesn't exist, create it
    }

    // Use try-catch to handle race condition where multiple
    // processes try to create inbox simultaneously
    try {
      return await this.create({
        config: merge(DEFAULT_AGENT_CONFIG, defaultAgentConfig) as any,
        slug: INBOX_SESSION_ID,
        type: 'agent',
      });
    } catch (error: any) {
      // If UNIQUE constraint error, another process created it first
      // Fetch and return the existing inbox
      if (error?.message?.includes('UNIQUE constraint') || error?.message?.includes('2067')) {
        const existing = await DB.GetSessionBySlug({
          slug: INBOX_SESSION_ID,
          userId: this.userId,
        });
        return existing;
      }
      // Re-throw other errors
      throw error;
    }
  };

  batchCreate = async (newSessions: NewSession[]) => {
    // Create sessions one by one (no batch insert in current bindings)
    const results = await Promise.all(
      newSessions.map((session) =>
        this.create({
          id: nanoid(),
          session,
          type: 'agent',
        }),
      ),
    );

    return results;
  };

  duplicate = async (id: string, newTitle?: string) => {
    const result = await this.findByIdOrSlug(id);
    if (!result) return;

    const { agent, clientId, ...session } = result;
    const sessionId = nanoid();

    return this.create({
      config: {
        ...agent,
        id: undefined,
        slug: undefined,
      } as any,
      id: sessionId,
      session: {
        ...session,
        title: newTitle || getNullableString(session.title),
      } as any,
      type: 'agent',
    });
  };

  // **************** Delete *************** //

  /**
   * Delete a session and its associated agent data if no longer referenced.
   */
  delete = async (id: string) => {
    await this.logger.methodEntry('delete', { id, userId: this.userId });
    
    // Get agents linked to this session
    const agents = await DB.GetSessionAgents({
      sessionId: id,
      userId: this.userId,
    });
    
    await this.logger.debug(`Found ${agents.length} agents linked to session ${id}`);

    // Unlink agents
    for (const agent of agents) {
      await DB.UnlinkAgentFromSession({
        agentId: agent.id,
        sessionId: id,
        userId: this.userId,
      });
      await this.logger.debug(`Unlinked agent ${agent.id} from session ${id}`);
    }

    // Delete session
    await DB.DeleteSession({
      id,
      userId: this.userId,
    });
    
    await this.logger.debug(`Deleted session ${id}`);

    // Delete orphaned agents - check if they're still linked to other sessions
    for (const agent of agents) {
      const agentSessions = await DB.GetAgentSessions({
        agentId: agent.id,
        userId: this.userId,
      });

      if (agentSessions.length === 0) {
        await DB.DeleteAgent({
          id: agent.id,
          userId: this.userId,
        });
        await this.logger.debug(`Deleted orphaned agent ${agent.id}`);
      }
    }
    
    await this.logger.methodExit('delete', { id, deletedAgents: agents.length });
  };

  /**
   * Batch delete sessions and their associated agent data if no longer referenced.
   */
  batchDelete = async (ids: string[]) => {
    if (ids.length === 0) return { count: 0 };

    // Get all agents linked to these sessions
    const allAgents = await Promise.all(
      ids.map(async (id) => {
        try {
          return await DB.GetSessionAgents({
            sessionId: id,
            userId: this.userId,
          });
        } catch {
          return [];
        }
      })
    );

    const agentIds = [...new Set(allAgents.flat().map((a) => a.id))];

    // Unlink all agents from these sessions
    await Promise.all(
      ids.flatMap((sessionId) =>
        agentIds.map((agentId) =>
          DB.UnlinkAgentFromSession({
            agentId,
            sessionId,
            userId: this.userId,
          }).catch(() => {})
        )
      )
    );

    // Batch delete sessions
    await DB.BatchDeleteSessions({
      userId: this.userId,
      ids,
    });

    // Delete orphaned agents
    const orphanedAgents = await DB.GetOrphanedAgents(this.userId);
    await Promise.all(
      orphanedAgents.map((agent) =>
        DB.DeleteAgent({
          id: agent.id,
          userId: this.userId,
        })
      )
    );

    return { count: ids.length };
  };

  /**
   * Delete all sessions and their associated agent data for this user.
   */
  deleteAll = async () => {
    const sessions = await DB.ListSessions({
      userId: this.userId,
      limit: 10000,
      offset: 0,
    });

    await this.batchDelete(sessions.map((s) => s.id));
  };

  // **************** Update *************** //

  update = async (id: string, data: Partial<Session>) => {
    await this.logger.methodEntry('update', { id, data, userId: this.userId });
    
    // Resolve slug to actual ID if needed
    const session = await this.findByIdOrSlug(id);
    if (!session) {
      await this.logger.error(`Session not found for update: ${id}`);
      throw new Error(`Session not found: ${id}`);
    }
    const actualId = session.id;

    const updated = await DB.UpdateSession({
      id: actualId,
      userId: this.userId,
      title: data.title !== undefined ? toNullString(getNullableString(data.title as any)) : toNullString(""),
      description: data.description !== undefined ? toNullString(getNullableString(data.description as any)) : toNullString(""),
      avatar: data.avatar !== undefined ? toNullString(getNullableString(data.avatar as any)) : toNullString(""),
      backgroundColor: data.backgroundColor !== undefined ? toNullString(getNullableString(data.backgroundColor as any)) : toNullString(""),
      groupId: data.groupId !== undefined ? toNullString(getNullableString(data.groupId as any)) : toNullString(""),
      pinned: data.pinned !== undefined ? data.pinned : 0,
      updatedAt: currentTimestampMs(),
    });
    
    await this.logger.methodExit('update', { id: updated.id });
    return [updated];
  };

  updateConfig = async (sessionId: string, data: PartialDeep<Agent> | undefined | null) => {
    await this.logger.methodEntry('updateConfig', { sessionId, hasData: !!data, userId: this.userId });
    
    if (!data || Object.keys(data).length === 0) {
      await this.logger.debug('No config data to update');
      return;
    }

    const session = await this.findByIdOrSlug(sessionId);
    if (!session || !session.agent || !session.agent.id) {
      await this.logger.error(`Session ${sessionId} not assigned with an agent`);
      throw new Error(
        'this session is not assigned with an agent, please contact with admin to fix this issue.',
      );
    }

    // Handle params field - undefined means delete, null means disabled
    const existingParams = parseNullableJSON(session.agent.params) ?? {};
    const updatedParams: Record<string, any> = { ...existingParams };

    if (data.params) {
      const incomingParams = data.params as Record<string, any>;
      Object.keys(incomingParams).forEach((key) => {
        const incomingValue = incomingParams[key];

        // undefined means explicitly delete the field
        if (incomingValue === undefined) {
          delete updatedParams[key];
          return;
        }

        // Other values (including null) are directly overwritten
        updatedParams[key] = incomingValue;
      });
    }

    // Merge data, excluding params (handled separately)
    const { params: _params, ...restData } = data;
    const mergedValue = merge(
      {
        ...session.agent,
        // Parse JSON fields for merging
        tags: parseNullableJSON(session.agent.tags),
        plugins: parseNullableJSON(session.agent.plugins),
        chatConfig: parseNullableJSON(session.agent.chatConfig),
        fewShots: parseNullableJSON(session.agent.fewShots),
        params: existingParams,
        tts: parseNullableJSON(session.agent.tts),
        openingQuestions: parseNullableJSON(session.agent.openingQuestions),
      },
      restData,
    );

    // Apply processed params
    mergedValue.params = Object.keys(updatedParams).length > 0 ? updatedParams : undefined;

    // Clean undefined values
    if (mergedValue.params) {
      Object.keys(mergedValue.params).forEach((key) => {
        if (mergedValue.params[key] === undefined) {
          delete mergedValue.params[key];
        }
      });
      if (Object.keys(mergedValue.params).length === 0) {
        mergedValue.params = undefined;
      }
    }

    // Update agent
    await DB.UpdateAgent({
      id: session.agent.id,
      userId: this.userId,
      title: toNullString(mergedValue.title as any),
      description: toNullString(mergedValue.description as any),
      tags: toNullJSON(mergedValue.tags),
      avatar: toNullString(mergedValue.avatar as any),
      backgroundColor: toNullString(mergedValue.backgroundColor as any),
      plugins: toNullJSON(mergedValue.plugins),
      chatConfig: toNullJSON(mergedValue.chatConfig),
      fewShots: toNullJSON(mergedValue.fewShots),
      model: toNullString(mergedValue.model as any),
      params: toNullJSON(mergedValue.params),
      provider: toNullString(mergedValue.provider as any),
      systemRole: toNullString(mergedValue.systemRole as any),
      tts: toNullJSON(mergedValue.tts),
      openingMessage: toNullString(mergedValue.openingMessage as any),
      openingQuestions: toNullJSON(mergedValue.openingQuestions),
      updatedAt: currentTimestampMs(),
    });
    
    await this.logger.methodExit('updateConfig', { agentId: session.agent.id });
  };

  // **************** Helper *************** //

  private mapSession = ({
    agentsToSessions,
    title,
    backgroundColor,
    description,
    avatar,
    groupId,
    type,
    ...res
  }: Session & { agentsToSessions?: { agent: Agent }[] }):
    | LobeAgentSession
    | LobeGroupSession => {
    const meta = {
      avatar: getNullableString(avatar as any) ?? undefined,
      backgroundColor: getNullableString(backgroundColor as any) ?? undefined,
      description: getNullableString(description as any) ?? undefined,
      tags: undefined,
      title: getNullableString(title as any) ?? undefined,
    };

    const typeStr = getNullableString(type as any);

    if (typeStr === 'group') {
      // For group sessions, transform agentsToSessions to members
      const members =
        agentsToSessions?.map((item, index) => ({
          ...item.agent,
          agentId: item.agent.id,
          chatGroupId: res.id,
          enabled: true,
          order: index,
          role: 'participant',
        })) || [];

      return {
        ...res,
        createdAt: new Date(res.createdAt),
        updatedAt: new Date(res.updatedAt),
        group: getNullableString(groupId as any),
        members,
        meta,
        type: 'group',
      } as unknown as LobeGroupSession;
    }

    // For agent sessions, include agent-specific fields
    const agent = agentsToSessions?.[0]?.agent;
    return {
      ...res,
      createdAt: new Date(res.createdAt),
      updatedAt: new Date(res.updatedAt),
      config: agent
        ? ({
            ...agent,
            model: getNullableString(agent.model as any) || '',
            plugins: parseNullableJSON(agent.plugins) || [],
            chatConfig: parseNullableJSON(agent.chatConfig) || {},
            params: parseNullableJSON(agent.params) || {},
            systemRole: getNullableString(agent.systemRole as any) || '',
            tts: parseNullableJSON(agent.tts) || {},
          } as any)
        : { model: '', plugins: [], chatConfig: {}, params: {}, systemRole: '', tts: {} },
      group: getNullableString(groupId as any),
      meta: {
        avatar: getNullableString(agent?.avatar as any) ?? getNullableString(avatar as any) ?? undefined,
        backgroundColor:
          getNullableString(agent?.backgroundColor as any) ??
          getNullableString(backgroundColor as any) ??
          undefined,
        description:
          getNullableString(agent?.description as any) ??
          getNullableString(description as any) ??
          undefined,
        tags: parseNullableJSON(agent?.tags as any) ?? undefined,
        title: getNullableString(agent?.title as any) ?? getNullableString(title as any) ?? undefined,
      },
      model: getNullableString(agent?.model as any) || '',
      type: 'agent',
    } as unknown as LobeAgentSession;
  };

  findSessionsByKeywords = async (params: {
    current?: number;
    keyword: string;
    pageSize?: number;
  }) => {
    const { keyword, pageSize = 9999 } = params;

    try {
      // Search agents by keyword
      const agents = await DB.SearchAgents({
        userId: this.userId,
        title: toNullString(`%${keyword}%`),
        description: toNullString(`%${keyword}%`),
        limit: pageSize,
      });

      // Get sessions for these agents
      const sessions = await Promise.all(
        agents.map(async (agent) => {
          try {
            const agentSessions = await DB.GetAgentSessions({
              agentId: agent.id,
              userId: this.userId,
            });
            // Return first session
            return agentSessions.length > 0 ? agentSessions[0] : null;
          } catch {
            return null;
          }
        }),
      );

      return sessions.filter((s) => s !== null);
    } catch (e) {
      console.error('findSessionsByKeywords error:', e, { keyword });
      return [];
    }
  };
}
