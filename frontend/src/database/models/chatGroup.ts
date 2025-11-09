import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';
import {
  ChatGroupItem,
  ChatGroupAgentItem,
  NewChatGroup,
  NewChatGroupAgent,
} from '@/types/chatGroup';

export class ChatGroupModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // ******* Query Methods ******* //

  async findById(id: string): Promise<ChatGroupItem | undefined> {
    const result = await DB.GetChatGroup({
      id,
      userId: this.userId,
    });

    if (!result) return undefined;

    return {
      id: result.id,
      title: getNullableString(result.title as any) ?? null,
      description: getNullableString(result.description as any) ?? null,
      config: parseNullableJSON(result.config as any) ?? null,
      pinned: intToBool(result.pinned || 0),
      userId: result.userId,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  }

  async query(): Promise<ChatGroupItem[]> {
    const results = await DB.ListChatGroups(this.userId);

    return results.map((r) => ({
      id: r.id,
      title: getNullableString(r.title as any) ?? null,
      description: getNullableString(r.description as any) ?? null,
      config: parseNullableJSON(r.config as any) ?? null,
      pinned: intToBool(r.pinned || 0),
      userId: r.userId,
      createdAt: r.createdAt,
      updatedAt: r.updatedAt,
    }));
  }

  /**
   * OPTIMIZED: Uses single JOIN query (3A) to fetch groups with agents
   * Much faster than N+1 approach
   */
  async queryWithMemberDetails(): Promise<(ChatGroupItem & { members?: any[] })[]> {
    // Single query with JOINs
    const results = await DB.ListChatGroupsWithAgents(this.userId);

    // Group by group_id
    const groupsMap = new Map<string, ChatGroupItem & { members?: any[] }>();

    for (const row of results) {
      const groupId = row.groupId;

      if (!groupsMap.has(groupId)) {
        groupsMap.set(groupId, {
          id: groupId,
          title: getNullableString(row.groupTitle as any) ?? null,
          description: getNullableString(row.groupDescription as any) ?? null,
          config: parseNullableJSON(row.groupConfig as any) ?? null,
          pinned: intToBool(row.groupPinned || 0),
          createdAt: row.groupCreatedAt,
          updatedAt: row.groupUpdatedAt,
          userId: this.userId,
          members: [],
        });
      }

      // Add agent if exists
      if (row.agentId) {
        const group = groupsMap.get(groupId)!;
        group.members!.push({
          id: row.agentId,
          title: getNullableString(row.agentTitle as any),
          description: getNullableString(row.agentDescription as any),
          avatar: getNullableString(row.agentAvatar as any),
          backgroundColor: getNullableString(row.agentBgColor as any),
          chatConfig: parseNullableJSON(row.agentChatConfig as any),
          params: parseNullableJSON(row.agentParams as any),
          systemRole: getNullableString(row.agentSystemRole as any),
          tts: parseNullableJSON(row.agentTts as any),
          model: getNullableString(row.agentModel as any),
          provider: getNullableString(row.agentProvider as any),
          createdAt: row.agentCreatedAt,
          updatedAt: row.agentUpdatedAt,
        });
      }
    }

    return Array.from(groupsMap.values());
  }

  /**
   * OPTIMIZED: Uses JOIN query to fetch group with agents
   */
  async findGroupWithAgents(groupId: string): Promise<{
    agents: ChatGroupAgentItem[];
    group: ChatGroupItem;
  } | null> {
    const results = await DB.GetChatGroupWithAgents({
      id: groupId,
      userId: this.userId,
    });

    if (results.length === 0) return null;

    // First row contains group data
    const firstRow = results[0];
    const group: ChatGroupItem = {
      id: firstRow.groupId,
      title: getNullableString(firstRow.groupTitle as any) ?? null,
      description: getNullableString(firstRow.groupDescription as any) ?? null,
      config: parseNullableJSON(firstRow.groupConfig as any) ?? null,
      pinned: intToBool(firstRow.groupPinned || 0),
      createdAt: firstRow.groupCreatedAt,
      updatedAt: firstRow.groupUpdatedAt,
      userId: this.userId,
    };

    // Extract agents
    const agents: ChatGroupAgentItem[] = results
      .filter((row) => row.agentId)
      .map((row) => ({
        agentId: String(row.agentId),
        chatGroupId: groupId,
        order: row.agentSortOrder?.Int64 || 0,
        role: getNullableString(row.agentRole as any) ?? null,
        enabled: intToBool(row.agentEnabled?.Int64 || 1),
        userId: this.userId,
      }));

    return { agents, group };
  }

  // ******* Create Methods ******* //

  async create(params: Partial<NewChatGroup>): Promise<ChatGroupItem> {
    const id = params.id || nanoid();
    const now = currentTimestampMs();

    const result = await DB.CreateChatGroup({
      id,
      title: toNullString(params.title),
      description: toNullString(params.description),
      config: toNullJSON(params.config),
      clientId: toNullString((params as any).clientId),
      userId: this.userId,
      pinned: boolToInt(params.pinned || false),
      createdAt: now,
      updatedAt: now,
    });

    return {
      id: result.id,
      title: getNullableString(result.title as any) ?? null,
      description: getNullableString(result.description as any) ?? null,
      config: parseNullableJSON(result.config as any) ?? null,
      pinned: intToBool(result.pinned || 0),
      userId: result.userId,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  }

  async createWithAgents(
    groupParams: Partial<NewChatGroup>,
    agentIds: string[],
  ): Promise<{ agents: ChatGroupAgentItem[]; group: ChatGroupItem }> {
    const group = await this.create(groupParams);

    if (agentIds.length === 0) {
      return { agents: [], group };
    }

    const agents: any[] = [];
    const now = currentTimestampMs();

    for (let i = 0; i < agentIds.length; i++) {
      await DB.LinkChatGroupToAgent({
        chatGroupId: group.id,
        agentId: agentIds[i],
        userId: this.userId,
        enabled: boolToInt(true),
        sortOrder: { Int64: i, Valid: true },
        role: toNullString('assistant'),
        createdAt: now,
        updatedAt: now,
      });

      agents.push({
        agentId: agentIds[i],
        chatGroupId: group.id,
        order: i,
        role: 'assistant',
        userId: this.userId,
      });
    }

    return { agents, group };
  }

  // ******* Update Methods ******* //

  async update(id: string, value: Partial<ChatGroupItem>): Promise<ChatGroupItem> {
    const result = await DB.UpdateChatGroup({
      id,
      userId: this.userId,
      title: toNullString(value.title),
      description: toNullString(value.description),
      config: toNullJSON(value.config),
      pinned: boolToInt(value.pinned ?? false),
      updatedAt: currentTimestampMs(),
    });

    if (!result) {
      throw new Error('Chat group not found or access denied');
    }

    return {
      id: result.id,
      title: getNullableString(result.title as any) ?? null,
      description: getNullableString(result.description as any) ?? null,
      config: parseNullableJSON(result.config as any) ?? null,
      pinned: intToBool(result.pinned || 0),
      userId: result.userId,
      createdAt: result.createdAt,
      updatedAt: result.updatedAt,
    };
  }

  async addAgentToGroup(
    groupId: string,
    agentId: string,
    options?: { order?: number; role?: string },
  ): Promise<ChatGroupAgentItem> {
    const now = currentTimestampMs();

    await DB.LinkChatGroupToAgent({
      chatGroupId: groupId,
      agentId: agentId,
      userId: this.userId,
      enabled: boolToInt(true),
      sortOrder: { Int64: options?.order || 0, Valid: true },
      role: toNullString(options?.role || 'assistant'),
      createdAt: now,
      updatedAt: now,
    });

    return {
      agentId,
      chatGroupId: groupId,
      order: options?.order || 0,
      role: options?.role || 'assistant',
      userId: this.userId,
    };
  }

  async addAgentsToGroup(groupId: string, agentIds: string[]): Promise<ChatGroupAgentItem[]> {
    const group = await this.findById(groupId);
    if (!group) throw new Error('Group not found');

    const existingAgents = await this.getGroupAgents(groupId);
    const existingAgentIds = new Set(existingAgents.map((a) => a.agentId));

    const newAgentIds = agentIds.filter((id) => !existingAgentIds.has(id));

    if (newAgentIds.length === 0) {
      return [];
    }

    const newAgents: ChatGroupAgentItem[] = [];
    const now = currentTimestampMs();

    for (const agentId of newAgentIds) {
      await DB.LinkChatGroupToAgent({
        chatGroupId: groupId,
        agentId: agentId,
        userId: this.userId,
        enabled: boolToInt(true),
        sortOrder: { Int64: 0, Valid: true },
        role: toNullString('assistant'),
        createdAt: now,
        updatedAt: now,
      });

      newAgents.push({
        agentId,
        chatGroupId: groupId,
        enabled: true,
        userId: this.userId,
      });
    }

    return newAgents;
  }

  async removeAgentFromGroup(groupId: string, agentId: string): Promise<void> {
    await DB.UnlinkChatGroupFromAgent({
      chatGroupId: groupId,
      agentId: agentId,
      userId: this.userId,
    });
  }

  async updateAgentInGroup(
    groupId: string,
    agentId: string,
    updates: Partial<Pick<NewChatGroupAgent, 'order' | 'role' | 'enabled'>>,
  ): Promise<ChatGroupAgentItem> {
    const result = await DB.UpdateChatGroupAgentLink({
      chatGroupId: groupId,
      agentId: agentId,
      userId: this.userId,
      sortOrder: { Int64: updates.order ?? 0, Valid: true },
      role: toNullString(updates.role || 'assistant'),
      enabled: boolToInt(updates.enabled ?? true),
      updatedAt: currentTimestampMs(),
    });

    return {
      chatGroupId: result.chatGroupId,
      agentId: result.agentId,
      userId: result.userId,
      role: getNullableString(result.role as any),
      enabled: intToBool(result.enabled || 0),
      order: result.sortOrder?.Int64 || 0,
    };
  }

  // ******* Delete Methods ******* //

  async delete(id: string): Promise<ChatGroupItem> {
    // Get group first to return it
    const group = await this.findById(id);
    if (!group) {
      throw new Error('Chat group not found or access denied');
    }

    // Delete (agents are automatically deleted due to CASCADE)
    await DB.DeleteChatGroup({
      id,
      userId: this.userId,
    });

    return group;
  }

  async deleteAll(): Promise<void> {
    await DB.DeleteAllChatGroups(this.userId);
  }

  // ******* Agent Query Methods ******* //

  async getGroupAgents(groupId: string): Promise<ChatGroupAgentItem[]> {
    const results = await DB.GetChatGroupAgentLinks({
      chatGroupId: groupId,
      userId: this.userId,
    });

    return results.map((r) => ({
      chatGroupId: r.chatGroupId,
      agentId: r.agentId,
      userId: r.userId,
      role: getNullableString(r.role as any),
      enabled: intToBool(r.enabled || 0),
      order: r.sortOrder?.Int64 || 0,
    }));
  }

  async getEnabledGroupAgents(groupId: string): Promise<ChatGroupAgentItem[]> {
    const results = await DB.GetEnabledChatGroupAgentLinks({
      chatGroupId: groupId,
      userId: this.userId,
    });

    return results.map((r) => ({
      chatGroupId: r.chatGroupId,
      agentId: r.agentId,
      userId: r.userId,
      role: getNullableString(r.role as any),
      enabled: intToBool(r.enabled || 0),
      order: r.sortOrder?.Int64 || 0,
    }));
  }

  async getGroupsWithAgents(agentIds?: string[]): Promise<(ChatGroupItem & { members?: any[] })[]> {
    if (!agentIds || agentIds.length === 0) {
      return this.query();
    }

    // Get all groups, then filter by agents
    // Note: This is not fully optimized - would need specific query for this
    const allGroups = await this.queryWithMemberDetails();
    
    return allGroups.filter((group) =>
      group.members?.some((member: any) => agentIds.includes(member.id)),
    );
  }
}

